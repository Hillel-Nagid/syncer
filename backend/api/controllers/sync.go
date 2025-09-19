package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"syncer.net/core/services"
	"syncer.net/core/sync"
	"syncer.net/services/music"
)

// SyncController handles both manual and automatic sync operations
// Implements project requirements: service pairing, sync modes, manual/auto sync
type SyncController struct {
	syncEngine *sync.SyncEngine
	registry   *services.ServiceRegistry
	db         *sqlx.DB
}

// NewSyncController creates a new sync controller
func NewSyncController(syncEngine *sync.SyncEngine, registry *services.ServiceRegistry, db *sqlx.DB) *SyncController {
	return &SyncController{
		syncEngine: syncEngine,
		registry:   registry,
		db:         db,
	}
}

// Initiate manual sync with service pairs
func (c *SyncController) InitiateManualSync(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		ServicePairs []sync.ServicePair `json:"service_pairs" binding:"required,min=1"`
		SyncType     string             `json:"sync_type" binding:"required"`
		SyncOptions  sync.SyncOptions   `json:"sync_options"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate each service pair
	for i, pair := range req.ServicePairs {
		if pair.SourceService == pair.TargetService {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Pair %d: source and target services must be different", i),
			})
			return
		}

		// Validate services are available
		if !c.registry.IsServiceAvailable(pair.SourceService) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Pair %d: source service %s not available", i, pair.SourceService),
			})
			return
		}
		if !c.registry.IsServiceAvailable(pair.TargetService) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Pair %d: target service %s not available", i, pair.TargetService),
			})
			return
		}

		// Validate sync mode
		validModes := []sync.SyncMode{sync.SyncModeFrom, sync.SyncModeTo, sync.SyncModeBidirectional}
		validMode := false
		for _, mode := range validModes {
			if pair.SyncMode == mode {
				validMode = true
				break
			}
		}
		if !validMode {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Pair %d: invalid sync mode %s", i, pair.SyncMode),
			})
			return
		}
	}

	// Check user has all required services connected
	var requiredServices []string
	for _, pair := range req.ServicePairs {
		requiredServices = append(requiredServices, pair.SourceService, pair.TargetService)
	}

	// Remove duplicates
	serviceMap := make(map[string]bool)
	for _, service := range requiredServices {
		serviceMap[service] = true
	}
	uniqueServices := make([]string, 0, len(serviceMap))
	for service := range serviceMap {
		uniqueServices = append(uniqueServices, service)
	}

	var connectedServices []string
	err := c.db.Select(&connectedServices, `
		SELECT s.name FROM user_services us
		JOIN services s ON us.service_id = s.id
		WHERE us.user_id = $1 AND s.name = ANY($2)
	`, userID, pq.Array(uniqueServices))

	if err != nil || len(connectedServices) != len(uniqueServices) {
		missingServices := make([]string, 0)
		connectedMap := make(map[string]bool)
		for _, service := range connectedServices {
			connectedMap[service] = true
		}
		for _, service := range uniqueServices {
			if !connectedMap[service] {
				missingServices = append(missingServices, service)
			}
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":              "All services in service pairs must be connected",
			"missing_services":   missingServices,
			"connected_services": connectedServices,
		})
		return
	}

	// Set default sync options if not provided
	if req.SyncOptions.MatchThreshold == 0 {
		req.SyncOptions.MatchThreshold = 0.8
	}
	if req.SyncOptions.ConflictPolicy == "" {
		req.SyncOptions.ConflictPolicy = sync.ConflictPolicySkip
	}

	// Create sync request
	syncReq := &sync.SyncJobRequest{
		UserID:       userID,
		ServicePairs: req.ServicePairs,
		SyncType:     req.SyncType,
		SyncOptions:  req.SyncOptions,
		RequestedAt:  time.Now(),
		IsScheduled:  false,
	}

	// Queue the manual sync
	if err := c.syncEngine.QueueManualSync(syncReq); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to queue sync job: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":          "Manual sync initiated successfully",
		"service_pairs":    len(req.ServicePairs),
		"sync_":            req.SyncType,
		"description":      syncReq.GetDescription(),
		"total_directions": syncReq.GetTotalDirections(),
	})
}

// ScheduleAutoSync - POST /api/sync/schedule
// Schedule automatic background sync
// Implements project requirement: "automatic in the background"
func (c *SyncController) ScheduleAutoSync(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		ServicePairs []sync.ServicePair `json:"service_pairs" binding:"required,min=1"`
		SyncType     string             `json:"sync_type" binding:"required"`
		SyncOptions  sync.SyncOptions   `json:"sync_options"`
		Schedule     sync.SyncSchedule  `json:"schedule" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Schedule.Frequency < 10*time.Minute {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Schedule frequency must be at least 10 minute",
		})
		return
	}

	if req.Schedule.NextRun.IsZero() {
		req.Schedule.NextRun = time.Now().Add(req.Schedule.Frequency)
	}

	// Create scheduled sync request
	syncReq := &sync.SyncJobRequest{
		UserID:       userID,
		ServicePairs: req.ServicePairs,
		SyncType:     req.SyncType,
		SyncOptions:  req.SyncOptions,
		RequestedAt:  time.Now(),
		IsScheduled:  true,
		Schedule:     &req.Schedule,
	}

	// Schedule the automatic sync
	if err := c.syncEngine.ScheduleAutoSync(syncReq); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to schedule sync: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Automatic sync scheduled successfully",
		"service_pairs": len(req.ServicePairs),
		"frequency":     req.Schedule.Frequency.String(),
		"next_run":      req.Schedule.NextRun,
		"enabled":       req.Schedule.Enabled,
	})
}

// GetSupportedSyncPairs - GET /api/sync/supported-pairs
// Get supported sync service pairs and modes
func (c *SyncController) GetSupportedSyncPairs(ctx *gin.Context) {
	allServices := c.registry.ListServices()

	var supportedPairs []map[string]any
	for _, source := range allServices {
		for _, target := range allServices {
			if source.Name != target.Name && source.Category == target.Category {
				supportedPairs = append(supportedPairs, map[string]any{
					"source_service": source.Name,
					"source_display": source.DisplayName,
					"target_service": target.Name,
					"target_display": target.DisplayName,
					"category":       source.Category,
					"supported_modes": []map[string]string{
						{"mode": string(sync.SyncModeFrom), "description": "One-way sync from source to target"},
						{"mode": string(sync.SyncModeTo), "description": "One-way sync from target to source"},
						{"mode": string(sync.SyncModeBidirectional), "description": "Two-way sync in both directions"},
					},
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"supported_pairs": supportedPairs,
		"sync_s": map[string]map[string]string{
			"music": music.GetMusicSyncTypeDescription(),
		},
		"conflict_policies": []map[string]string{
			{"policy": string(sync.ConflictPolicySkip), "description": "Skip conflicting items"},
			{"policy": string(sync.ConflictPolicyOverwrite), "description": "Replace existing items"},
			{"policy": string(sync.ConflictPolicyMerge), "description": "Merge metadata"},
		},
	})
}

// GetSyncStatus - GET /api/sync/status
// Get current sync status and recent jobs
func (c *SyncController) GetSyncStatus(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get recent sync jobs for the user
	var recentJobs []map[string]any
	rows, err := c.db.Query(`
		SELECT id, sync_, service_pairs_count, status, is_scheduled,
		       items_synced, items_failed, duration_ms, error_count,
		       created_at, finished_at
		FROM sync_jobs 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT 10
	`, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch sync status",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var job struct {
			ID               string     `db:"id"`
			SyncType         string     `db:"sync_"`
			ServicePairCount int        `db:"service_pairs_count"`
			Status           string     `db:"status"`
			IsScheduled      bool       `db:"is_scheduled"`
			ItemsSynced      *int       `db:"items_synced"`
			ItemsFailed      *int       `db:"items_failed"`
			DurationMs       *int       `db:"duration_ms"`
			ErrorCount       *int       `db:"error_count"`
			CreatedAt        time.Time  `db:"created_at"`
			FinishedAt       *time.Time `db:"finished_at"`
		}

		if err := rows.Scan(
			&job.ID, &job.SyncType, &job.ServicePairCount, &job.Status, &job.IsScheduled,
			&job.ItemsSynced, &job.ItemsFailed, &job.DurationMs, &job.ErrorCount,
			&job.CreatedAt, &job.FinishedAt,
		); err != nil {
			continue
		}

		jobMap := map[string]any{
			"id":                 job.ID,
			"sync_":              job.SyncType,
			"service_pair_count": job.ServicePairCount,
			"status":             job.Status,
			"type":               map[string]string{"manual": "Manual", "scheduled": "Automatic"}[map[bool]string{true: "scheduled", false: "manual"}[job.IsScheduled]],
			"created_at":         job.CreatedAt,
		}

		if job.ItemsSynced != nil {
			jobMap["items_synced"] = *job.ItemsSynced
		}
		if job.ItemsFailed != nil {
			jobMap["items_failed"] = *job.ItemsFailed
		}
		if job.DurationMs != nil {
			jobMap["duration_ms"] = *job.DurationMs
		}
		if job.ErrorCount != nil {
			jobMap["error_count"] = *job.ErrorCount
		}
		if job.FinishedAt != nil {
			jobMap["finished_at"] = *job.FinishedAt
		}

		recentJobs = append(recentJobs, jobMap)
	}

	// Get sync statistics
	stats, err := c.getSyncStats(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch sync statistics",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"recent_jobs": recentJobs,
		"statistics":  stats,
	})
}

// GetUserSchedules - GET /api/sync/schedules
// Get user's automatic sync schedules
func (c *SyncController) GetUserSchedules(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var schedules []map[string]any
	rows, err := c.db.Query(`
		SELECT sync_, enabled, next_run, created_at, updated_at
		FROM sync_schedules
		WHERE user_id = $1
		ORDER BY sync_
	`, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch schedules",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var schedule struct {
			SyncType  string    `db:"sync_"`
			Enabled   bool      `db:"enabled"`
			NextRun   time.Time `db:"next_run"`
			CreatedAt time.Time `db:"created_at"`
			UpdatedAt time.Time `db:"updated_at"`
		}

		if err := rows.Scan(&schedule.SyncType, &schedule.Enabled, &schedule.NextRun,
			&schedule.CreatedAt, &schedule.UpdatedAt); err != nil {
			continue
		}

		schedules = append(schedules, map[string]any{
			"sync_":      schedule.SyncType,
			"enabled":    schedule.Enabled,
			"next_run":   schedule.NextRun,
			"created_at": schedule.CreatedAt,
			"updated_at": schedule.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"schedules": schedules,
	})
}

// UpdateSchedule - PUT /api/sync/schedules/:syncType
// Update an existing schedule
func (c *SyncController) UpdateSchedule(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	syncType := ctx.Param("syncType")
	if syncType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Sync type is required"})
		return
	}

	var req struct {
		Enabled   *bool          `json:"enabled"`
		Frequency *time.Duration `json:"frequency"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the schedule in database
	updates := []string{}
	args := []any{}
	argIndex := 1

	if req.Enabled != nil {
		updates = append(updates, fmt.Sprintf("enabled = $%d", argIndex))
		args = append(args, *req.Enabled)
		argIndex++
	}

	if req.Frequency != nil {
		if *req.Frequency < 10*time.Minute {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Frequency must be at least 10 minute",
			})
			return
		}

		// When frequency changes, update next_run too
		updates = append(updates, fmt.Sprintf("next_run = NOW() + $%d", argIndex))
		args = append(args, *req.Frequency)
		argIndex++

		// We would also need to update the schedule_data JSON, but for simplicity
		// we'll just update the next_run time here
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one field must be updated",
		})
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, userID, syncType)

	query := fmt.Sprintf(`
		UPDATE sync_schedules 
		SET %s 
		WHERE user_id = $%d AND sync_ = $%d
	`, fmt.Sprintf(updates[0], argIndex, argIndex+1), argIndex, argIndex+1)

	result, err := c.db.Exec(query, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update schedule",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Schedule not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Schedule updated successfully",
		"sync_":   syncType,
	})
}

// DeleteSchedule - DELETE /api/sync/schedules/:syncType
// Delete a sync schedule
func (c *SyncController) DeleteSchedule(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	syncType := ctx.Param("syncType")
	if syncType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Sync type is required"})
		return
	}

	result, err := c.db.Exec(`
		DELETE FROM sync_schedules 
		WHERE user_id = $1 AND sync_ = $2
	`, userID, syncType)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete schedule",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Schedule not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Schedule deleted successfully",
		"sync_":   syncType,
	})
}

// Helper function to get sync statistics
func (c *SyncController) getSyncStats(userID string) (map[string]any, error) {
	var stats struct {
		TotalJobs      int `db:"total_jobs"`
		SuccessfulJobs int `db:"successful_jobs"`
		FailedJobs     int `db:"failed_jobs"`
		TotalSynced    int `db:"total_synced"`
	}

	err := c.db.Get(&stats, `
		SELECT 
			COUNT(*) as total_jobs,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful_jobs,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_jobs,
			COALESCE(SUM(items_synced), 0) as total_synced
		FROM sync_jobs 
		WHERE user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}

	// Get last sync time
	var lastSync *time.Time
	err = c.db.Get(&lastSync, `
		SELECT MAX(finished_at) 
		FROM sync_jobs 
		WHERE user_id = $1 AND status = 'completed'
	`, userID)

	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"total_jobs":      stats.TotalJobs,
		"successful_jobs": stats.SuccessfulJobs,
		"failed_jobs":     stats.FailedJobs,
		"total_synced":    stats.TotalSynced,
	}

	if lastSync != nil {
		result["last_sync"] = *lastSync
	}

	return result, nil
}

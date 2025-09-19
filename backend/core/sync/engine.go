package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"syncer.net/core/services"
)

// SyncEngine handles real-time synchronization between paired services
type SyncEngine struct {
	oauth       *services.OAuthManager
	transformer DataTransformer
	adder       CrossServiceAdder
	scheduler   *SyncScheduler
	db          *sqlx.DB
	manualQueue chan *CrossServiceSyncRequest
	autoQueue   chan *CrossServiceSyncRequest
	workers     int
	logger      *log.Logger
	metrics     *SyncMetrics
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

// NewSyncEngine creates a new sync engine with generic interfaces
func NewSyncEngine(
	oauth *services.OAuthManager,
	transformer DataTransformer,
	adder CrossServiceAdder,
	db *sqlx.DB,
	workers int,
) *SyncEngine {
	return &SyncEngine{
		oauth:       oauth,
		transformer: transformer,
		adder:       adder,
		scheduler:   NewSyncScheduler(db),
		db:          db,
		manualQueue: make(chan *CrossServiceSyncRequest, 500),
		autoQueue:   make(chan *CrossServiceSyncRequest, 200),
		workers:     workers,
		logger:      log.New(log.Writer(), "[SyncEngine] ", log.LstdFlags),
		metrics:     NewSyncMetrics(),
		stopChan:    make(chan struct{}),
	}
}

// Start initializes worker goroutines and automatic scheduler
func (e *SyncEngine) Start(ctx context.Context) error {
	e.logger.Printf("Starting sync engine with %d workers", e.workers)

	for i := range e.workers / 2 {
		e.wg.Add(1)
		go e.manualWorker(ctx, i)
		e.wg.Add(1)
		go e.autoWorker(ctx, i)
	}

	e.wg.Add(1)
	go e.scheduler.Start(ctx, e.autoQueue)

	return nil
}

// Stop gracefully shuts down the sync engine
func (e *SyncEngine) Stop() error {
	e.logger.Printf("Stopping sync engine")
	close(e.stopChan)
	e.wg.Wait()
	return nil
}

// QueueManualSync queues a user-initiated sync job
func (e *SyncEngine) QueueManualSync(req *SyncJobRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid sync request: %w", err)
	}

	if err := e.validateServicesAvailability(req); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	crossServiceReq := &CrossServiceSyncRequest{
		SyncJobRequest: req,
		Priority:       PriorityMedium,
		RequestedBy:    req.UserID,
	}
	crossServiceReq.IsScheduled = false

	select {
	case e.manualQueue <- crossServiceReq:
		e.logger.Printf("Queued manual sync for user %s with %d service pairs",
			req.UserID, len(req.ServicePairs))
		return nil
	default:
		return fmt.Errorf("manual sync queue is full")
	}
}

// ScheduleAutoSync sets up automatic background sync
func (e *SyncEngine) ScheduleAutoSync(req *SyncJobRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid sync request: %w", err)
	}
	if req.Schedule == nil {
		return fmt.Errorf("schedule configuration required for automatic sync")
	}

	req.IsScheduled = true
	return e.scheduler.Schedule(req)
}

// validateServicesAvailability ensures all required services are registered
func (e *SyncEngine) validateServicesAvailability(req *SyncJobRequest) error {
	for _, pair := range req.ServicePairs {
		if !e.oauth.Registry.IsServiceAvailable(pair.SourceService) {
			return fmt.Errorf("source service %s not available", pair.SourceService)
		}
		if !e.oauth.Registry.IsServiceAvailable(pair.TargetService) {
			return fmt.Errorf("target service %s not available", pair.TargetService)
		}
	}
	return nil
}

// manualWorker processes user-initiated sync requests
func (e *SyncEngine) manualWorker(ctx context.Context, workerID int) {
	defer e.wg.Done()
	logger := log.New(log.Writer(), fmt.Sprintf("[ManualWorker-%d] ", workerID), log.LstdFlags)
	logger.Printf("Manual sync worker started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-e.stopChan:
			return
		case req := <-e.manualQueue:
			e.processSyncJob(ctx, req, logger)
		}
	}
}

// autoWorker processes scheduled background sync requests
func (e *SyncEngine) autoWorker(ctx context.Context, workerID int) {
	defer e.wg.Done()
	logger := log.New(log.Writer(), fmt.Sprintf("[AutoWorker-%d] ", workerID), log.LstdFlags)
	logger.Printf("Automatic sync worker started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-e.stopChan:
			return
		case req := <-e.autoQueue:
			e.processSyncJob(ctx, req, logger)
		}
	}
}

// processSyncJob performs synchronization for all service pairs in the request
func (e *SyncEngine) processSyncJob(ctx context.Context, req *CrossServiceSyncRequest, logger *log.Logger) {
	jobID := uuid.New().String()
	startTime := time.Now()
	syncType := "manual"
	if req.IsScheduled {
		syncType = "automatic"
	}

	logger.Printf("Processing %s sync job %s for user %s with %d service pairs (type: %s)",
		syncType, jobID, req.UserID, len(req.ServicePairs), req.SyncType)

	if err := e.createSyncJobRecord(jobID, req); err != nil {
		logger.Printf("Failed to create sync job record: %v", err)
	}

	var servicePairResults []ServicePairResult
	totalSynced := []UniversalItem{}
	totalFailed := []UniversalItem{}
	var allErrors []services.SyncError

	for i, pair := range req.ServicePairs {
		pairLogger := log.New(log.Writer(), fmt.Sprintf("[Job-%s-Pair-%d] ", jobID[:8], i), log.LstdFlags)

		result := e.processServicePair(ctx, req.UserID, pair, req.SyncType, req.SyncOptions, pairLogger)
		servicePairResults = append(servicePairResults, result)

		totalSynced = append(totalSynced, result.ItemsSynced...)
		totalFailed = append(totalFailed, result.ItemsFailed...)
		allErrors = append(allErrors, result.Errors...)

		if result.Success {
			pairLogger.Printf("Service pair sync completed: %s ↔ %s - %d items synced in %v",
				pair.SourceService, pair.TargetService, result.ItemsSynced, result.Duration)
		} else {
			pairLogger.Printf("Service pair sync failed: %s ↔ %s - %d errors",
				pair.SourceService, pair.TargetService, len(result.Errors))
		}
	}

	duration := time.Since(startTime)

	syncResult := &CrossServiceSyncResult{
		JobID:        jobID,
		Success:      len(totalFailed) == 0,
		ServicePairs: servicePairResults,
		TotalSynced:  len(totalSynced),
		TotalFailed:  len(totalFailed),
		Duration:     duration,
		Errors:       allErrors,
		Metadata: map[string]any{
			"sync_type":     req.SyncType,
			"request_type":  syncType,
			"service_pairs": len(req.ServicePairs),
			"user_id":       req.UserID,
			"requested_by":  req.RequestedBy,
			"timestamp":     startTime,
		},
	}

	if err := e.updateSyncJobRecord(jobID, syncResult); err != nil {
		logger.Printf("Failed to update sync job record: %v", err)
	}

	if err := e.storeSyncResult(syncResult); err != nil {
		logger.Printf("Failed to store sync result: %v", err)
	} else {
		logger.Printf("Stored complete sync result for job %s", jobID)
	}

	if syncResult.Success {
		e.metrics.RecordSyncJobSuccess(req.UserID, syncType, len(req.ServicePairs), len(totalSynced), duration)
		logger.Printf("Sync job %s completed successfully in %v: %d total items synced across %d pairs",
			jobID, duration, len(totalSynced), len(req.ServicePairs))
	} else {
		e.metrics.RecordSyncJobFailure(req.UserID, syncType, len(req.ServicePairs), len(totalFailed))
		logger.Printf("Sync job %s completed with errors after %v: %d total items failed",
			jobID, duration, len(totalFailed))
	}
}

// processServicePair handles sync for a single service pair based on sync mode
func (e *SyncEngine) processServicePair(ctx context.Context, userID string, pair ServicePair, syncType string, options SyncOptions, logger *log.Logger) ServicePairResult {
	startTime := time.Now()

	result := ServicePairResult{
		SourceService: pair.SourceService,
		TargetService: pair.TargetService,
		SyncMode:      pair.SyncMode,
		Success:       false,
		ItemsSynced:   []UniversalItem{},
		ItemsFailed:   []UniversalItem{},
		Errors:        []services.SyncError{},
	}

	logger.Printf("Processing service pair: %s ↔ %s (mode: %s)", pair.SourceService, pair.TargetService, pair.SyncMode)

	sourceService, err := e.oauth.Registry.GetService(pair.SourceService)
	if err != nil {
		result.Errors = append(result.Errors, services.SyncError{
			Type:    "service_error",
			Error:   fmt.Sprintf("failed to get source service: %v", err),
			Context: "service_resolution",
		})
		result.Duration = time.Since(startTime)
		return result
	}

	targetService, err := e.oauth.Registry.GetService(pair.TargetService)
	if err != nil {
		result.Errors = append(result.Errors, services.SyncError{
			Type:    "service_error",
			Error:   fmt.Sprintf("failed to get target service: %v", err),
			Context: "service_resolution",
		})
		result.Duration = time.Since(startTime)
		return result
	}

	sourceTokens, err := e.getUserTokens(userID, pair.SourceService)
	if err != nil {
		result.Errors = append(result.Errors, services.SyncError{
			Type:    "auth_error",
			Error:   fmt.Sprintf("failed to get source tokens: %v", err),
			Context: "token_retrieval",
		})
		result.Duration = time.Since(startTime)
		return result
	}

	targetTokens, err := e.getUserTokens(userID, pair.TargetService)
	if err != nil {
		result.Errors = append(result.Errors, services.SyncError{
			Type:    "auth_error",
			Error:   fmt.Sprintf("failed to get target tokens: %v", err),
			Context: "token_retrieval",
		})
		result.Duration = time.Since(startTime)
		return result
	}

	switch pair.SyncMode {
	case SyncModeFrom:
		synced, syncErrors := e.performDirectionalSync(ctx, sourceService, targetService, sourceTokens, targetTokens, syncType, options, logger)
		result.ItemsSynced = synced.Items
		result.ItemsFailed = synced.Failed
		result.Errors = syncErrors
		result.Success = len(syncErrors) == 0

	case SyncModeTo:
		synced, syncErrors := e.performDirectionalSync(ctx, targetService, sourceService, targetTokens, sourceTokens, syncType, options, logger)
		result.ItemsSynced = synced.Items
		result.ItemsFailed = synced.Failed
		result.Errors = syncErrors
		result.Success = len(syncErrors) == 0

	case SyncModeBidirectional:
		synced1, errors1 := e.performDirectionalSync(ctx, sourceService, targetService, sourceTokens, targetTokens, syncType, options, logger)
		synced2, errors2 := e.performDirectionalSync(ctx, targetService, sourceService, targetTokens, sourceTokens, syncType, options, logger)

		result.ItemsSynced = append(synced1.Items, synced2.Items...)
		result.ItemsFailed = append(synced1.Failed, synced2.Failed...)
		result.Errors = append(errors1, errors2...)
		result.Success = len(result.Errors) == 0

	default:
		result.Errors = append(result.Errors, services.SyncError{
			Type:    "config_error",
			Error:   fmt.Sprintf("unsupported sync mode: %s", pair.SyncMode),
			Context: "sync_mode_validation",
		})
	}

	result.Duration = time.Since(startTime)
	return result
}

// performDirectionalSync performs one-way sync from source to target
func (e *SyncEngine) performDirectionalSync(
	ctx context.Context,
	sourceService, targetService services.ServiceProvider,
	sourceTokens, targetTokens *services.OAuthTokens,
	syncType string,
	options SyncOptions,
	logger *log.Logger,
) (*SyncResult, []services.SyncError) {
	logger.Printf("Starting directional sync: %s → %s (type: %s)", sourceService.Name(), targetService.Name(), syncType)

	sourceResult, err := sourceService.GetUserData(ctx, sourceTokens, time.Time{})
	if err != nil {
		return nil, []services.SyncError{{
			Type:    "sync_error",
			Error:   fmt.Sprintf("failed to fetch source data: %v", err),
			Context: "source_data_fetch",
		}}
	}

	if !sourceResult.Success || len(sourceResult.Items) == 0 {
		logger.Printf("No data found in source service %s", sourceService.Name())
		return nil, nil
	}

	var universalItems []UniversalItem
	var transformErrors []services.SyncError

	for _, item := range sourceResult.Items {
		if e.itemMatchesSyncType(item.ItemType, syncType) {
			universalItem, err := e.transformer.TransformToUniversal(sourceService.Name(), item.Data)
			if err != nil {
				transformErrors = append(transformErrors, services.SyncError{
					Type:    "transform_error",
					Error:   fmt.Sprintf("failed to transform item: %v", err),
					ItemID:  item.ExternalID,
					Context: "item_transformation",
				})
				continue
			}
			universalItems = append(universalItems, universalItem)
		}
	}

	logger.Printf("Transformed %d items to universal format", len(universalItems))

	if options.DryRun {
		logger.Printf("DRY RUN: Would sync %d items", len(universalItems))
		return nil, transformErrors
	}

	syncedItems := make([]UniversalItem, 0, len(universalItems))
	var syncErrors []services.SyncError

	for _, universalItem := range universalItems {
		err := e.adder.AddItemToService(ctx, targetService, targetTokens, universalItem, options)
		if err != nil {
			syncErrors = append(syncErrors, services.SyncError{
				Type:    "add_error",
				Error:   fmt.Sprintf("failed to add item to %s: %v", targetService.Name(), err),
				Context: fmt.Sprintf("adding_to_%s", targetService.Name()),
			})
			continue
		}

		syncedItems = append(syncedItems, universalItem)

		logger.Printf("Successfully synced item to %s", targetService.Name())

		time.Sleep(100 * time.Millisecond)
	}

	allErrors := append(transformErrors, syncErrors...)

	logger.Printf("Generic sync completed: %d/%d items synced successfully", len(syncedItems), len(universalItems))
	return &SyncResult{
		Items:    syncedItems,
		Errors:   allErrors,
		Metadata: map[string]any{},
	}, allErrors
}

// itemMatchesSyncType checks if an item type matches the requested sync type
func (e *SyncEngine) itemMatchesSyncType(itemType string, syncType string) bool {
	return e.transformer.MatchesSyncType(itemType, syncType)
}

// Helper method to get user tokens
func (e *SyncEngine) getUserTokens(userID, serviceName string) (*services.OAuthTokens, error) {
	var userServiceID string
	err := e.db.Get(&userServiceID, `
		SELECT us.id
		FROM user_services us
		JOIN services s ON us.service_id = s.id
		WHERE us.user_id = $1 AND s.name = $2
	`, userID, serviceName)

	if err != nil {
		return nil, fmt.Errorf("user service not found: %w", err)
	}

	return e.oauth.GetUserTokens(userServiceID)
}

// Database operations for sync job tracking (metadata only)
func (e *SyncEngine) createSyncJobRecord(jobID string, req *CrossServiceSyncRequest) error {
	_, err := e.db.Exec(`
		INSERT INTO sync_jobs (
			id, user_id, status, sync_type, service_pairs_count, 
			is_scheduled, priority, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`, jobID, req.UserID, "running", req.SyncType, len(req.ServicePairs),
		req.IsScheduled, req.Priority)

	return err
}

func (e *SyncEngine) updateSyncJobRecord(jobID string, result *CrossServiceSyncResult) error {
	status := "completed"
	if !result.Success {
		status = "failed"
	}

	_, err := e.db.Exec(`
		UPDATE sync_jobs SET 
			status = $1, 
			items_synced = $2, 
			items_failed = $3, 
			duration_ms = $4,
			error_count = $5,
			finished_at = NOW()
		WHERE id = $6
	`, status, result.TotalSynced, result.TotalFailed,
		result.Duration.Milliseconds(), len(result.Errors), jobID)

	return err
}

// storeSyncResult stores the complete sync result as JSON in the database
func (e *SyncEngine) storeSyncResult(result *CrossServiceSyncResult) error {

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal sync result: %w", err)
	}

	_, err = e.db.Exec(`
		INSERT INTO sync_results (job_id, result_data)
		VALUES ($1, $2)
		ON CONFLICT (job_id) DO UPDATE SET
			result_data = EXCLUDED.result_data,
			updated_at = NOW()
	`, result.JobID, resultJSON)

	if err != nil {
		return fmt.Errorf("failed to store sync result: %w", err)
	}

	return nil
}

// GetSyncResult retrieves a sync result by job ID
func (e *SyncEngine) GetSyncResult(jobID string) (*CrossServiceSyncResult, error) {
	var resultJSON []byte

	err := e.db.Get(&resultJSON, `
		SELECT result_data 
		FROM sync_results 
		WHERE job_id = $1
	`, jobID)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sync result: %w", err)
	}

	var result CrossServiceSyncResult
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sync result: %w", err)
	}

	return &result, nil
}

// GetUserSyncResults retrieves sync results for a specific user with optional pagination and filtering
func (e *SyncEngine) GetUserSyncResults(userID string, limit int, offset int, successOnly *bool) ([]*CrossServiceSyncResult, error) {
	//TODO: properly implement this
	query := `
		SELECT result_data 
		FROM sync_results 
		WHERE user_id = $1
	`

	args := []any{}
	argIndex := 2

	if successOnly != nil {
		query += fmt.Sprintf(" AND (result_data->>'success')::boolean = $%d", argIndex)
		args = append(args, *successOnly)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++

		if offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, offset)
		}
	}

	rows, err := e.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user sync results: %w", err)
	}
	defer rows.Close()

	var results []*CrossServiceSyncResult

	for rows.Next() {
		var resultJSON []byte
		if err := rows.Scan(&resultJSON); err != nil {
			return nil, fmt.Errorf("failed to scan sync result: %w", err)
		}

		var result CrossServiceSyncResult
		if err := json.Unmarshal(resultJSON, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sync result: %w", err)
		}

		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sync results: %w", err)
	}

	return results, nil
}

// GetSyncResultsByDateRange retrieves sync results for a user within a specific date range
func (e *SyncEngine) GetSyncResultsByDateRange(userID string, startTime, endTime time.Time) ([]*CrossServiceSyncResult, error) {
	//TODO: properly implement this
	query := `
		SELECT result_data 
		FROM sync_results 
		WHERE user_id = $1 
		AND created_at >= $2 
		AND created_at <= $3
		ORDER BY created_at DESC
	`

	rows, err := e.db.Query(query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sync results by date range: %w", err)
	}
	defer rows.Close()

	var results []*CrossServiceSyncResult

	for rows.Next() {
		var resultJSON []byte
		if err := rows.Scan(&resultJSON); err != nil {
			return nil, fmt.Errorf("failed to scan sync result: %w", err)
		}

		var result CrossServiceSyncResult
		if err := json.Unmarshal(resultJSON, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sync result: %w", err)
		}

		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sync results: %w", err)
	}

	return results, nil
}

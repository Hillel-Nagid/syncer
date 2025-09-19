package sync

import (
	"context"
	"fmt"
	"slices"
	"time"

	"syncer.net/core/services"
)

// DataTransformer defines the interface for transforming service-specific data to universal format
type DataTransformer interface {
	TransformToUniversal(serviceName string, data any) (UniversalItem, error)
	FindBestMatch(sourceItem UniversalItem, candidates []UniversalItem, threshold float64) UniversalMatch
	AnalyzeMatches(matches []UniversalMatch) map[string]any
	MatchesSyncType(itemType string, syncType string) bool
}

// CrossServiceAdder defines the interface for adding items to services
type CrossServiceAdder interface {
	AddItemToService(ctx context.Context, targetService services.ServiceProvider, tokens *services.OAuthTokens, universalItem UniversalItem, options any) error
}

// SyncJobRequest defines a sync operation between paired services
type SyncJobRequest struct {
	UserID       string        `json:"user_id"`
	ServicePairs []ServicePair `json:"service_pairs" binding:"required,min=1"`
	SyncType     string        `json:"sync_type"`
	SyncOptions  SyncOptions   `json:"sync_options"`
	RequestedAt  time.Time     `json:"requested_at"`
	IsScheduled  bool          `json:"is_scheduled"`
	Schedule     *SyncSchedule `json:"schedule,omitempty"`
}

// ServicePair defines a sync relationship between two services with direction
type ServicePair struct {
	SourceService string   `json:"source_service" binding:"required"`
	TargetService string   `json:"target_service" binding:"required"`
	SyncMode      SyncMode `json:"sync_mode" binding:"required"`
}

// SyncMode defines the direction of synchronization
type SyncMode string

const (
	SyncModeFrom          SyncMode = "sync-from"
	SyncModeTo            SyncMode = "sync-to"
	SyncModeBidirectional SyncMode = "bidirectional"
)

// SyncSchedule defines automatic background sync configuration
type SyncSchedule struct {
	Enabled   bool          `json:"enabled"`
	Frequency time.Duration `json:"frequency"`
	NextRun   time.Time     `json:"next_run"`
}

// SyncOptions defines options for sync operations
type SyncOptions struct {
	ConflictPolicy ConflictPolicy `json:"conflict_policy"`
	MatchThreshold float64        `json:"match_threshold"`
	DryRun         bool           `json:"dry_run"`
}

// ConflictPolicy defines how to handle sync conflicts
type ConflictPolicy string

const (
	ConflictPolicySkip      ConflictPolicy = "skip"
	ConflictPolicyOverwrite ConflictPolicy = "overwrite"
	ConflictPolicyMerge     ConflictPolicy = "merge"
)

type SyncResult struct {
	Items    []UniversalItem      `json:"items"`
	Failed   []UniversalItem      `json:"failed"`
	Errors   []services.SyncError `json:"errors"`
	Metadata map[string]any       `json:"metadata,omitempty"`
}

// SyncPriority defines the priority level for sync jobs
type SyncPriority int

const (
	PriorityLow    SyncPriority = 0
	PriorityMedium SyncPriority = 1
	PriorityHigh   SyncPriority = 2
	PriorityUrgent SyncPriority = 3
)

type UniversalItem interface {
	GetItemType() string
	GetItemIdentifier() string
	GetItemAction() services.SyncAction
}

type UniversalMatch struct {
	Target     UniversalItem `json:"target"`
	Confidence float64       `json:"confidence"` // Match confidence score
}

// SyncJobStatus defines the current status of a sync job
type SyncJobStatus string

const (
	SyncStatusPending   SyncJobStatus = "pending"
	SyncStatusRunning   SyncJobStatus = "running"
	SyncStatusCompleted SyncJobStatus = "completed"
	SyncStatusFailed    SyncJobStatus = "failed"
	SyncStatusCancelled SyncJobStatus = "cancelled"
)

// CrossServiceSyncRequest wraps the enhanced sync job request
type CrossServiceSyncRequest struct {
	*SyncJobRequest
	Priority    SyncPriority `json:"priority"`
	RequestedBy string       `json:"requested_by"`
}

// SyncResult represents the result of a cross-service sync operation
type CrossServiceSyncResult struct {
	JobID        string               `json:"job_id"`
	Success      bool                 `json:"success"`
	ServicePairs []ServicePairResult  `json:"service_pairs"`
	TotalSynced  int                  `json:"total_synced"`
	TotalFailed  int                  `json:"total_failed"`
	Duration     time.Duration        `json:"duration"`
	Errors       []services.SyncError `json:"errors"`
	Metadata     map[string]any       `json:"metadata"`
}

// ServicePairResult represents the result for a single service pair
type ServicePairResult struct {
	SourceService string               `json:"source_service"`
	TargetService string               `json:"target_service"`
	SyncMode      SyncMode             `json:"sync_mode"`
	Success       bool                 `json:"success"`
	ItemsSynced   []UniversalItem      `json:"items_synced"`
	ItemsFailed   []UniversalItem      `json:"items_failed"`
	Errors        []services.SyncError `json:"errors"`
	Duration      time.Duration        `json:"duration"`
}

// SyncStats represents statistics about sync operations
type SyncStats struct {
	TotalJobs        int           `json:"total_jobs"`
	SuccessfulJobs   int           `json:"successful_jobs"`
	FailedJobs       int           `json:"failed_jobs"`
	TotalItemsSynced int           `json:"total_items_synced"`
	AverageDuration  time.Duration `json:"average_duration"`
	LastSyncAt       time.Time     `json:"last_sync_at"`
}

// ValidateSyncJobRequest validates a sync job request
func (r *SyncJobRequest) Validate() error {
	if len(r.ServicePairs) == 0 {
		return fmt.Errorf("at least one service pair is required - cannot sync a service alone")
	}

	if r.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	if r.SyncType == "" {
		return fmt.Errorf("sync type is required")
	}

	for i, pair := range r.ServicePairs {
		if pair.SourceService == "" {
			return fmt.Errorf("service pair %d: source service is required", i)
		}
		if pair.TargetService == "" {
			return fmt.Errorf("service pair %d: target service is required", i)
		}
		if pair.SourceService == pair.TargetService {
			return fmt.Errorf("service pair %d: source and target services must be different", i)
		}
		if pair.SyncMode == "" {
			return fmt.Errorf("service pair %d: sync mode is required", i)
		}

		validModes := []SyncMode{SyncModeFrom, SyncModeTo, SyncModeBidirectional}
		if !slices.Contains(validModes, pair.SyncMode) {
			return fmt.Errorf("service pair %d: invalid sync mode %s", i, pair.SyncMode)
		}
	}

	if r.SyncOptions.MatchThreshold < 0 || r.SyncOptions.MatchThreshold > 1 {
		return fmt.Errorf("match threshold must be between 0 and 1")
	}

	if r.SyncOptions.MatchThreshold == 0 {
		r.SyncOptions.MatchThreshold = 0.8
	}

	if r.SyncOptions.ConflictPolicy == "" {
		r.SyncOptions.ConflictPolicy = ConflictPolicySkip
	}

	return nil
}

// GetDescription returns a human-readable description of the sync job
func (r *SyncJobRequest) GetDescription() string {
	description := fmt.Sprintf("Sync %s between %d service pairs", r.SyncType, len(r.ServicePairs))

	if r.IsScheduled {
		description += " (scheduled)"
	} else {
		description += " (manual)"
	}

	return description
}

// GetTotalDirections returns the total number of sync directions
func (r *SyncJobRequest) GetTotalDirections() int {
	total := 0
	for _, pair := range r.ServicePairs {
		if pair.SyncMode == SyncModeBidirectional {
			total += 2
		} else {
			total += 1
		}
	}
	return total
}

// GetUniqueServices returns all unique services involved in the sync
func (r *SyncJobRequest) GetUniqueServices() []string {
	serviceMap := make(map[string]bool)

	for _, pair := range r.ServicePairs {
		serviceMap[pair.SourceService] = true
		serviceMap[pair.TargetService] = true
	}

	services := make([]string, 0, len(serviceMap))
	for service := range serviceMap {
		services = append(services, service)
	}

	return services
}

package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

// SyncScheduler handles automatic background sync scheduling
type SyncScheduler struct {
	db        *sqlx.DB
	schedules map[string]*SyncJobRequest
	ticker    *time.Ticker
	logger    *log.Logger
	mu        sync.RWMutex
}

// NewSyncScheduler creates a new sync scheduler
func NewSyncScheduler(db *sqlx.DB) *SyncScheduler {
	return &SyncScheduler{
		db:        db,
		schedules: make(map[string]*SyncJobRequest),
		ticker:    time.NewTicker(10 * time.Minute),
		logger:    log.New(log.Writer(), "[SyncScheduler] ", log.LstdFlags),
	}
}

// Start begins the automatic sync scheduler
func (s *SyncScheduler) Start(ctx context.Context, autoQueue chan *CrossServiceSyncRequest) {
	s.logger.Printf("Starting automatic sync scheduler")

	s.loadSchedules()

	for {
		select {
		case <-ctx.Done():
			s.logger.Printf("Sync scheduler stopping due to context cancellation")
			return
		case <-s.ticker.C:
			s.checkScheduledSyncs(autoQueue)
		}
	}
}

// Schedule adds or updates an automatic sync schedule
func (s *SyncScheduler) Schedule(req *SyncJobRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Schedule == nil {
		return fmt.Errorf("schedule configuration is required")
	}

	if req.Schedule.NextRun.IsZero() {
		req.Schedule.NextRun = time.Now().Add(req.Schedule.Frequency)
	}

	scheduleID, err := s.saveScheduleToDatabase(req)
	if err != nil {
		return fmt.Errorf("failed to save schedule: %w", err)
	}
	s.schedules[scheduleID] = req

	s.logger.Printf("Scheduled automatic sync for user %s: type=%s, frequency=%v, next_run=%v",
		req.UserID, req.SyncType, req.Schedule.Frequency, req.Schedule.NextRun)

	return nil
}

// Unschedule removes an automatic sync schedule
func (s *SyncScheduler) Unschedule(scheduleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.schedules, scheduleID)

	_, err := s.db.Exec(`
		DELETE FROM sync_schedules 
		WHERE id = $1
	`, scheduleID)

	if err != nil {
		return fmt.Errorf("failed to remove schedule: %w", err)
	}

	s.logger.Printf("Unscheduled automatic sync for schedule %s", scheduleID)
	return nil
}

// checkScheduledSyncs looks for sync jobs that are due to run
func (s *SyncScheduler) checkScheduledSyncs(autoQueue chan *CrossServiceSyncRequest) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	scheduled := 0

	for scheduleID, req := range s.schedules {
		if !req.Schedule.Enabled {
			continue
		}

		if now.After(req.Schedule.NextRun) {
			crossServiceReq := &CrossServiceSyncRequest{
				SyncJobRequest: req,
				Priority:       PriorityLow,
				RequestedBy:    "system",
			}

			select {
			case autoQueue <- crossServiceReq:
				s.logger.Printf("Queued automatic sync for user %s, type %s", req.UserID, req.SyncType)
				scheduled++

				req.Schedule.NextRun = now.Add(req.Schedule.Frequency)

				if _, err := s.saveScheduleToDatabase(req); err != nil {
					s.logger.Printf("Failed to update next run time for schedule %s: %v", scheduleID, err)
				}

			default:
				s.logger.Printf("Failed to queue automatic sync - queue full")
			}
		}
	}

	if scheduled > 0 {
		s.logger.Printf("Scheduled %d automatic sync jobs", scheduled)
	}
}

// GetSchedules returns all active schedules for a user
func (s *SyncScheduler) GetSchedules(userID string) []SyncJobRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userSchedules []SyncJobRequest

	for _, req := range s.schedules {
		if req.UserID == userID && req.Schedule != nil && req.Schedule.Enabled {
			userSchedules = append(userSchedules, *req)
		}
	}

	return userSchedules
}

// GetAllSchedules returns information about all schedules
func (s *SyncScheduler) GetAllSchedules() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	enabled := 0
	disabled := 0
	byType := make(map[string]int)
	nextRuns := make([]time.Time, 0, len(s.schedules))

	for _, req := range s.schedules {
		if req.Schedule == nil {
			continue
		}

		if req.Schedule.Enabled {
			enabled++
			nextRuns = append(nextRuns, req.Schedule.NextRun)
		} else {
			disabled++
		}

		byType[string(req.SyncType)]++
	}

	var nextRun time.Time
	for _, runTime := range nextRuns {
		if nextRun.IsZero() || runTime.Before(nextRun) {
			nextRun = runTime
		}
	}

	return map[string]any{
		"total_schedules":    len(s.schedules),
		"enabled_schedules":  enabled,
		"disabled_schedules": disabled,
		"by_sync_type":       byType,
		"next_run":           nextRun,
	}
}

// saveScheduleToDatabase persists a sync schedule
func (s *SyncScheduler) saveScheduleToDatabase(req *SyncJobRequest) (string, error) {
	var scheduleID string
	scheduleData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schedule: %w", err)
	}

	err = s.db.QueryRow(`
		INSERT INTO sync_schedules (user_id, sync_type, schedule_data, next_run, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (user_id, sync_type) DO UPDATE SET
			schedule_data = $3,
			next_run = $4,
			enabled = $5,
			updated_at = NOW()
		RETURNING id
	`, req.UserID, req.SyncType, scheduleData, req.Schedule.NextRun, req.Schedule.Enabled).Scan(&scheduleID)

	if err != nil {
		return "", fmt.Errorf("failed to save schedule: %w", err)
	}

	return scheduleID, nil
}

// loadSchedules loads existing sync schedules from database
func (s *SyncScheduler) loadSchedules() {
	s.logger.Printf("Loading existing sync schedules from database")

	rows, err := s.db.Query(`
		SELECT id, schedule_data FROM sync_schedules
		WHERE enabled = true AND next_run > NOW()
	`)
	if err != nil {
		s.logger.Printf("Failed to load schedules: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var scheduleID string
		var scheduleData []byte
		if err := rows.Scan(&scheduleID, &scheduleData); err != nil {
			s.logger.Printf("Failed to scan schedule data: %v", err)
			continue
		}

		var req SyncJobRequest
		if err := json.Unmarshal(scheduleData, &req); err != nil {
			s.logger.Printf("Failed to unmarshal schedule data: %v", err)
			continue
		}

		s.schedules[scheduleID] = &req
		count++
	}

	s.logger.Printf("Loaded %d existing sync schedules", count)
}

// Stop stops the scheduler
func (s *SyncScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.logger.Printf("Sync scheduler stopped")
	}
}

// EnableSchedule enables a disabled schedule
func (s *SyncScheduler) EnableSchedule(scheduleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, exists := s.schedules[scheduleID]
	if !exists {
		return fmt.Errorf("schedule not found")
	}

	if req.Schedule == nil {
		return fmt.Errorf("invalid schedule configuration")
	}

	req.Schedule.Enabled = true
	req.Schedule.NextRun = time.Now().Add(req.Schedule.Frequency)

	_, err := s.saveScheduleToDatabase(req)
	if err != nil {
		return fmt.Errorf("failed to enable schedule: %w", err)
	}

	s.logger.Printf("Enabled schedule for schedule %s", scheduleID)
	return nil
}

// DisableSchedule disables an active schedule
func (s *SyncScheduler) DisableSchedule(scheduleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, exists := s.schedules[scheduleID]
	if !exists {
		return fmt.Errorf("schedule not found")
	}

	if req.Schedule == nil {
		return fmt.Errorf("invalid schedule configuration")
	}

	req.Schedule.Enabled = false

	_, err := s.saveScheduleToDatabase(req)
	if err != nil {
		return fmt.Errorf("failed to disable schedule: %w", err)
	}

	s.logger.Printf("Disabled schedule for schedule %s", scheduleID)
	return nil
}

// UpdateScheduleFrequency updates the frequency of an existing schedule
func (s *SyncScheduler) UpdateScheduleFrequency(scheduleID string, frequency time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if frequency < time.Minute {
		return fmt.Errorf("minimum frequency is 1 minute")
	}

	req, exists := s.schedules[scheduleID]
	if !exists {
		return fmt.Errorf("schedule not found")
	}

	if req.Schedule == nil {
		return fmt.Errorf("invalid schedule configuration")
	}

	req.Schedule.Frequency = frequency
	req.Schedule.NextRun = time.Now().Add(frequency)

	_, err := s.saveScheduleToDatabase(req)
	if err != nil {
		return fmt.Errorf("failed to update schedule frequency: %w", err)
	}

	s.logger.Printf("Updated schedule frequency for schedule %s to %v", scheduleID, frequency)
	return nil
}

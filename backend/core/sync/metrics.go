package sync

import (
	"log"
	"sync"
	"time"
)

// SyncMetrics collects and tracks synchronization metrics
type SyncMetrics struct {
	mu           sync.RWMutex
	jobsTotal    map[string]int
	jobsSuccess  map[string]int
	jobsFailure  map[string]int
	itemsSynced  map[string]int
	avgDurations map[string]time.Duration
	lastSync     map[string]time.Time
	logger       *log.Logger
}

// NewSyncMetrics creates a new sync metrics collector
func NewSyncMetrics() *SyncMetrics {
	return &SyncMetrics{
		jobsTotal:    make(map[string]int),
		jobsSuccess:  make(map[string]int),
		jobsFailure:  make(map[string]int),
		itemsSynced:  make(map[string]int),
		avgDurations: make(map[string]time.Duration),
		lastSync:     make(map[string]time.Time),
		logger:       log.New(log.Writer(), "[SyncMetrics] ", log.LstdFlags),
	}
}

// RecordSyncJobSuccess records a successful sync job
func (m *SyncMetrics) RecordSyncJobSuccess(userID, syncType string, pairCount, itemsSynced int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := syncType

	m.jobsTotal[key]++
	m.jobsSuccess[key]++
	m.itemsSynced[key] += itemsSynced
	m.lastSync[key] = time.Now()

	if prev, exists := m.avgDurations[key]; exists {
		m.avgDurations[key] = (prev + duration) / 2
	} else {
		m.avgDurations[key] = duration
	}

	m.logger.Printf("Recorded successful sync: type=%s, items=%d, duration=%v",
		syncType, itemsSynced, duration)
}

// RecordSyncJobFailure records a failed sync job
func (m *SyncMetrics) RecordSyncJobFailure(userID, syncType string, pairCount, itemsFailed int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := syncType

	m.jobsTotal[key]++
	m.jobsFailure[key]++

	m.logger.Printf("Recorded failed sync: type=%s, failures=%d", syncType, itemsFailed)
}

// GetStats returns current sync statistics
func (m *SyncMetrics) GetStats() SyncStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalJobs := 0
	successfulJobs := 0
	failedJobs := 0
	totalItemsSynced := 0
	var avgDuration time.Duration
	var lastSyncAt time.Time

	for _, count := range m.jobsTotal {
		totalJobs += count
	}

	for _, count := range m.jobsSuccess {
		successfulJobs += count
	}

	for _, count := range m.jobsFailure {
		failedJobs += count
	}

	for _, count := range m.itemsSynced {
		totalItemsSynced += count
	}

	var totalDuration time.Duration
	durationCount := 0
	for _, duration := range m.avgDurations {
		totalDuration += duration
		durationCount++
	}
	if durationCount > 0 {
		avgDuration = totalDuration / time.Duration(durationCount)
	}

	for _, syncTime := range m.lastSync {
		if syncTime.After(lastSyncAt) {
			lastSyncAt = syncTime
		}
	}

	return SyncStats{
		TotalJobs:        totalJobs,
		SuccessfulJobs:   successfulJobs,
		FailedJobs:       failedJobs,
		TotalItemsSynced: totalItemsSynced,
		AverageDuration:  avgDuration,
		LastSyncAt:       lastSyncAt,
	}
}

// GetStatsForSyncType returns stats for a specific sync type
func (m *SyncMetrics) GetStatsForSyncType(syncType string) map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]any{
		"total_jobs":      m.jobsTotal[syncType],
		"successful_jobs": m.jobsSuccess[syncType],
		"failed_jobs":     m.jobsFailure[syncType],
		"items_synced":    m.itemsSynced[syncType],
		"avg_duration":    m.avgDurations[syncType],
		"last_sync":       m.lastSync[syncType],
	}
}

// Reset clears all metrics
func (m *SyncMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobsTotal = make(map[string]int)
	m.jobsSuccess = make(map[string]int)
	m.jobsFailure = make(map[string]int)
	m.itemsSynced = make(map[string]int)
	m.avgDurations = make(map[string]time.Duration)
	m.lastSync = make(map[string]time.Time)

	m.logger.Printf("Metrics reset")
}

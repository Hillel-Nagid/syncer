CREATE TABLE IF NOT EXISTS sync_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_service_id UUID NOT NULL REFERENCES user_services(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    error TEXT,
    -- Enhanced fields for Phase 1
    priority INTEGER DEFAULT 0,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    next_retry_at TIMESTAMP,
    items_processed INTEGER DEFAULT 0,
    items_added INTEGER DEFAULT 0,
    items_updated INTEGER DEFAULT 0,
    items_deleted INTEGER DEFAULT 0,
    sync_metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Index for job processing (status and priority)
CREATE INDEX IF NOT EXISTS idx_sync_jobs_status_priority ON sync_jobs(status, priority DESC, created_at);
-- Index for retry processing
CREATE INDEX IF NOT EXISTS idx_sync_jobs_retry ON sync_jobs(status, next_retry_at)
WHERE status = 'failed'
    AND next_retry_at IS NOT NULL;
-- Index for job statistics
CREATE INDEX IF NOT EXISTS idx_sync_jobs_stats ON sync_jobs(created_at, status);
-- GIN index for metadata queries
CREATE INDEX IF NOT EXISTS idx_sync_jobs_metadata ON sync_jobs USING GIN (sync_metadata)
WHERE sync_metadata IS NOT NULL;
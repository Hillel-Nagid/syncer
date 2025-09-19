-- Migration: Create sync infrastructure tables
-- Streamlined schema with only actively used columns
CREATE TABLE IF NOT EXISTS sync_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    sync_type TEXT,
    service_pairs_count INTEGER DEFAULT 1 CHECK (service_pairs_count > 0),
    is_scheduled BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 1 CHECK (priority IN (0, 1, 2, 3)),
    items_synced INTEGER DEFAULT 0 CHECK (items_synced >= 0),
    items_failed INTEGER DEFAULT 0 CHECK (items_failed >= 0),
    duration_ms BIGINT CHECK (duration_ms >= 0),
    error_count INTEGER DEFAULT 0 CHECK (error_count >= 0),
    finished_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Create sync_schedules table for automatic background sync
CREATE TABLE IF NOT EXISTS sync_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sync_type TEXT NOT NULL,
    schedule_data JSONB NOT NULL,
    -- Full SyncJobRequest serialized with service pairs and options
    next_run TIMESTAMP NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (next_run > created_at)
);
-- sync_logs enhancements are handled in migration 005
-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_sync_jobs_status_priority ON sync_jobs(status, priority DESC, created_at);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_user_type ON sync_jobs(user_id, sync_type);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_scheduled ON sync_jobs(is_scheduled, created_at);
-- sync_schedules indexes
CREATE INDEX IF NOT EXISTS idx_sync_schedules_user ON sync_schedules(user_id);
CREATE INDEX IF NOT EXISTS idx_sync_schedules_next_run ON sync_schedules(enabled, next_run)
WHERE enabled = true;
CREATE INDEX IF NOT EXISTS idx_sync_schedules_type ON sync_schedules(sync_type);
-- sync_logs indexes are handled in migration 005
-- Add trigger to update sync_schedules updated_at timestamp
CREATE OR REPLACE FUNCTION update_sync_schedules_updated_at() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS trigger_sync_schedules_updated_at ON sync_schedules;
CREATE TRIGGER trigger_sync_schedules_updated_at BEFORE
UPDATE ON sync_schedules FOR EACH ROW EXECUTE FUNCTION update_sync_schedules_updated_at();
-- Create function to clean up old sync jobs (older than 30 days)
CREATE OR REPLACE FUNCTION cleanup_old_sync_jobs() RETURNS VOID AS $$ BEGIN
DELETE FROM sync_jobs
WHERE finished_at < NOW() - INTERVAL '30 days'
    AND status IN ('completed', 'failed', 'cancelled');
END;
$$ LANGUAGE plpgsql;
-- Add some helpful views for sync analytics
-- View for sync job statistics per user
CREATE OR REPLACE VIEW user_sync_stats AS
SELECT user_id,
    COUNT(*) as total_jobs,
    COUNT(
        CASE
            WHEN status = 'completed' THEN 1
        END
    ) as successful_jobs,
    COUNT(
        CASE
            WHEN status = 'failed' THEN 1
        END
    ) as failed_jobs,
    COUNT(
        CASE
            WHEN is_scheduled = true THEN 1
        END
    ) as scheduled_jobs,
    COUNT(
        CASE
            WHEN is_scheduled = false THEN 1
        END
    ) as manual_jobs,
    COALESCE(SUM(items_synced), 0) as total_items_synced,
    COALESCE(AVG(duration_ms), 0) as avg_duration_ms,
    MAX(finished_at) as last_sync_at
FROM sync_jobs
GROUP BY user_id;
-- View for service pair usage statistics
CREATE OR REPLACE VIEW service_pair_stats AS
SELECT sl.source_service,
    sl.target_service,
    sl.sync_direction,
    COUNT(*) as total_syncs,
    COUNT(
        CASE
            WHEN sj.status = 'completed' THEN 1
        END
    ) as successful_syncs,
    COALESCE(AVG(sj.items_synced), 0) as avg_items_synced
FROM sync_logs sl
    JOIN sync_jobs sj ON sl.sync_job_id = sj.id
WHERE sl.source_service IS NOT NULL
    AND sl.target_service IS NOT NULL
GROUP BY sl.source_service,
    sl.target_service,
    sl.sync_direction;
-- View for active schedules summary
CREATE OR REPLACE VIEW active_schedules_summary AS
SELECT sync_type,
    COUNT(*) as total_schedules,
    COUNT(
        CASE
            WHEN enabled = true THEN 1
        END
    ) as enabled_schedules,
    MIN(next_run) as next_scheduled_run,
    COUNT(
        CASE
            WHEN next_run < NOW() THEN 1
        END
    ) as overdue_schedules
FROM sync_schedules
GROUP BY sync_type;
-- Service records are handled in migration 002
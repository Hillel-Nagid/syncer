-- Migration: Create sync results table for storing complete sync result JSON
-- This table stores the full CrossServiceSyncResult as JSON for on-demand retrieval
CREATE TABLE IF NOT EXISTS sync_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL UNIQUE REFERENCES sync_jobs(id) ON DELETE CASCADE,
    result_data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_sync_results_user_id ON sync_results(user_id);
CREATE INDEX IF NOT EXISTS idx_sync_results_job_id ON sync_results(job_id);
CREATE INDEX IF NOT EXISTS idx_sync_results_created_at ON sync_results(created_at DESC);
-- Index for querying by success status within the JSON
CREATE INDEX IF NOT EXISTS idx_sync_results_success ON sync_results USING gin((result_data->'success'));
-- Index for querying by sync type within the JSON metadata
CREATE INDEX IF NOT EXISTS idx_sync_results_sync_type ON sync_results USING gin((result_data->'metadata'->'sync_type'));
-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_sync_results_updated_at() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trigger_sync_results_updated_at BEFORE
UPDATE ON sync_results FOR EACH ROW EXECUTE FUNCTION update_sync_results_updated_at();
-- Add cleanup function to remove old sync results (older than 90 days by default)
CREATE OR REPLACE FUNCTION cleanup_old_sync_results(retention_days INTEGER DEFAULT 90) RETURNS INTEGER AS $$
DECLARE deleted_count INTEGER;
BEGIN
DELETE FROM sync_results
WHERE created_at < NOW() - (retention_days || ' days')::INTERVAL;
GET DIAGNOSTICS deleted_count = ROW_COUNT;
RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
-- Create view for sync results summary 
CREATE OR REPLACE VIEW sync_results_summary AS
SELECT sr.user_id,
    sr.job_id,
    sr.result_data->>'job_id' as result_job_id,
    (sr.result_data->>'success')::boolean as success,
    (sr.result_data->>'total_synced')::integer as total_synced,
    (sr.result_data->>'total_failed')::integer as total_failed,
    sr.result_data->'metadata'->>'sync_type' as sync_type,
    sr.result_data->'metadata'->>'request_type' as request_type,
    (sr.result_data->'metadata'->>'service_pairs')::integer as service_pairs_count,
    to_timestamp(
        (sr.result_data->'metadata'->>'timestamp')::double precision
    ) as sync_timestamp,
    sr.created_at,
    sr.updated_at
FROM sync_results sr;
-- Drop indexes
DROP INDEX IF EXISTS idx_sync_jobs_status_priority;
DROP INDEX IF EXISTS idx_sync_jobs_retry;
DROP INDEX IF EXISTS idx_sync_jobs_stats;
DROP INDEX IF EXISTS idx_sync_jobs_metadata;
-- Drop table
DROP TABLE IF EXISTS sync_jobs;
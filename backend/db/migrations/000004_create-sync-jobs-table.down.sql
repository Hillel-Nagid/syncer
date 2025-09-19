-- Migration rollback: Remove sync infrastructure tables
-- Drop views
DROP VIEW IF EXISTS service_pair_stats;
DROP VIEW IF EXISTS active_schedules_summary;
DROP VIEW IF EXISTS user_sync_stats;
-- Drop function and trigger
DROP TRIGGER IF EXISTS trigger_sync_schedules_updated_at ON sync_schedules;
DROP FUNCTION IF EXISTS update_sync_schedules_updated_at();
DROP FUNCTION IF EXISTS cleanup_old_sync_jobs();
-- Drop indexes
DROP INDEX IF EXISTS idx_sync_schedules_type;
DROP INDEX IF EXISTS idx_sync_schedules_next_run;
DROP INDEX IF EXISTS idx_sync_schedules_user;
DROP INDEX IF EXISTS idx_sync_jobs_scheduled;
DROP INDEX IF EXISTS idx_sync_jobs_user_type;
DROP INDEX IF EXISTS idx_sync_jobs_status_priority;
-- Drop tables
DROP TABLE IF EXISTS sync_schedules;
DROP TABLE IF EXISTS sync_jobs;
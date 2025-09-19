-- Migration: Drop sync results table and related objects
-- Drop the view first
DROP VIEW IF EXISTS sync_results_summary;
-- Drop the trigger
DROP TRIGGER IF EXISTS trigger_sync_results_updated_at ON sync_results;
-- Drop the trigger function
DROP FUNCTION IF EXISTS update_sync_results_updated_at();
-- Drop the cleanup function
DROP FUNCTION IF EXISTS cleanup_old_sync_results(INTEGER);
-- Drop the table (this will cascade to drop the indexes)
DROP TABLE IF EXISTS sync_results;
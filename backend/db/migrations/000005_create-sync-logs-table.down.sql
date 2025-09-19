-- Drop functions
DROP FUNCTION IF EXISTS cleanup_old_sync_logs();
DROP FUNCTION IF EXISTS cleanup_expired_oauth_auth();
-- Drop indexes
DROP INDEX IF EXISTS idx_sync_logs_services;
-- Drop pending OAuth auth table
DROP TABLE IF EXISTS pending_oauth_auth;
-- Drop sync metadata table and trigger
DROP TRIGGER IF EXISTS update_sync_metadata_updated_at ON sync_metadata;
DROP TABLE IF EXISTS sync_metadata;
-- Drop sync logs table
DROP TABLE IF EXISTS sync_logs;
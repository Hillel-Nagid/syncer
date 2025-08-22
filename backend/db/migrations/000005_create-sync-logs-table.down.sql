-- Drop pending OAuth auth table and cleanup function
DROP FUNCTION IF EXISTS cleanup_expired_oauth_auth();
DROP TABLE IF EXISTS pending_oauth_auth;
DROP TRIGGER IF EXISTS update_sync_metadata_updated_at ON sync_metadata;
DROP TABLE IF EXISTS sync_metadata;
DROP TABLE IF EXISTS sync_logs;
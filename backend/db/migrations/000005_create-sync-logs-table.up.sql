CREATE TABLE IF NOT EXISTS sync_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sync_job_id UUID NOT NULL REFERENCES sync_jobs(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Table for storing ONLY sync metadata - NO user data content (privacy-first)
CREATE TABLE IF NOT EXISTS sync_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_service_id UUID NOT NULL REFERENCES user_services(id) ON DELETE CASCADE,
    external_id TEXT NOT NULL,
    item_type TEXT NOT NULL,
    checksum TEXT NOT NULL,
    last_modified TIMESTAMP NOT NULL,
    last_sync_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sync_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_service_id, external_id, item_type)
);
CREATE INDEX IF NOT EXISTS idx_sync_metadata_user_service ON sync_metadata(user_service_id);
CREATE INDEX IF NOT EXISTS idx_sync_metadata_type ON sync_metadata(item_type);
CREATE INDEX IF NOT EXISTS idx_sync_metadata_modified ON sync_metadata(last_modified);
CREATE INDEX IF NOT EXISTS idx_sync_metadata_synced ON sync_metadata(last_sync_at);
CREATE INDEX IF NOT EXISTS idx_sync_metadata_external_id ON sync_metadata(external_id);
CREATE INDEX IF NOT EXISTS idx_sync_metadata_checksum ON sync_metadata(checksum);
-- Trigger to update updated_at for sync_metadata
CREATE TRIGGER update_sync_metadata_updated_at BEFORE
UPDATE ON sync_metadata FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Table for managing pending OAuth authorizations
CREATE TABLE IF NOT EXISTS pending_oauth_auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_name TEXT NOT NULL,
    state TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, service_name)
);
-- Index for efficient state token lookups
CREATE INDEX IF NOT EXISTS idx_pending_oauth_state ON pending_oauth_auth(state);
-- Index for cleanup of expired entries
CREATE INDEX IF NOT EXISTS idx_pending_oauth_expires ON pending_oauth_auth(expires_at);
-- Index for user lookups
CREATE INDEX IF NOT EXISTS idx_pending_oauth_user ON pending_oauth_auth(user_id);
-- Cleanup function to remove expired entries
CREATE OR REPLACE FUNCTION cleanup_expired_oauth_auth() RETURNS INTEGER AS $$
DECLARE deleted_count INTEGER;
BEGIN
DELETE FROM pending_oauth_auth
WHERE expires_at < NOW();
GET DIAGNOSTICS deleted_count = ROW_COUNT;
RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
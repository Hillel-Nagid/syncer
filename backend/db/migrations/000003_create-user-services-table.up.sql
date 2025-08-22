CREATE TABLE IF NOT EXISTS user_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    metadata JSONB,
    -- Enhanced fields for Phase 1
    encrypted_access_token BYTEA,
    encrypted_refresh_token BYTEA,
    token_type TEXT DEFAULT 'Bearer',
    token_expires_at TIMESTAMP,
    last_sync_at TIMESTAMP,
    sync_frequency INTERVAL DEFAULT '1 day',
    sync_enabled BOOLEAN DEFAULT true,
    service_user_id TEXT,
    service_username TEXT,
    scopes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, service_id)
);
-- Create indexes for efficient sync scheduling
CREATE INDEX IF NOT EXISTS idx_user_services_next_sync ON user_services (sync_enabled, last_sync_at, sync_frequency)
WHERE sync_enabled = true;
-- Create index for service user lookups
CREATE INDEX IF NOT EXISTS idx_user_services_service_user ON user_services(service_user_id)
WHERE service_user_id IS NOT NULL;
-- Create index for updated_at for efficient queries
CREATE INDEX IF NOT EXISTS idx_user_services_updated_at ON user_services(updated_at);
-- Update trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';
CREATE TRIGGER update_user_services_updated_at BEFORE
UPDATE ON user_services FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
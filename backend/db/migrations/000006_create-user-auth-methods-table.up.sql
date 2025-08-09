CREATE TABLE IF NOT EXISTS user_auth_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_id TEXT,
    provider_email TEXT,
    password_hash TEXT,
    email_verified BOOLEAN DEFAULT FALSE,
    is_primary BOOLEAN DEFAULT FALSE,
    verification_token TEXT,
    verification_expires TIMESTAMP,
    reset_token TEXT,
    reset_expires TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, provider)
);
-- Create unique constraint for provider_id when it's not null
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auth_methods_provider_id_unique ON user_auth_methods(provider, provider_id)
WHERE provider_id IS NOT NULL;
-- Create unique constraint for email provider
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auth_methods_email_unique ON user_auth_methods(provider, provider_email)
WHERE provider = 'email';
-- Create regular indexes
CREATE INDEX IF NOT EXISTS idx_user_auth_methods_user_id ON user_auth_methods(user_id);
CREATE INDEX IF NOT EXISTS idx_user_auth_methods_provider ON user_auth_methods(provider);
CREATE INDEX IF NOT EXISTS idx_user_auth_methods_provider_email ON user_auth_methods(provider_email);
CREATE INDEX IF NOT EXISTS idx_user_auth_methods_verification_token ON user_auth_methods(verification_token)
WHERE verification_token IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_user_auth_methods_reset_token ON user_auth_methods(reset_token)
WHERE reset_token IS NOT NULL;
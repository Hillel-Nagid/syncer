CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    primary_email TEXT NOT NULL,
    full_name TEXT NOT NULL,
    avatar_url TEXT,
    last_login TIMESTAMP,
    online BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_primary_email ON users(primary_email);
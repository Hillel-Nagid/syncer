CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    category TEXT,
    oauth_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Insert initial service records
INSERT INTO services (
        id,
        name,
        display_name,
        category,
        oauth_required,
        created_at,
        updated_at
    )
VALUES (
        gen_random_uuid(),
        'spotify',
        'Spotify',
        'music',
        TRUE,
        NOW(),
        NOW()
    ),
    (
        gen_random_uuid(),
        'deezer',
        'Deezer',
        'music',
        TRUE,
        NOW(),
        NOW()
    ) ON CONFLICT (name) DO
UPDATE
SET display_name = EXCLUDED.display_name,
    category = EXCLUDED.category,
    oauth_required = EXCLUDED.oauth_required,
    updated_at = NOW();
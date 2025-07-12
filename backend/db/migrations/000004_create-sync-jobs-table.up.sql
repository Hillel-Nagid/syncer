CREATE TABLE IF NOT EXISTS sync_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_service_id UUID NOT NULL REFERENCES user_services(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
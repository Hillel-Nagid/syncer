package sync_jobs

import "time"

type SyncJob struct {
	ID            string    `json:"id" db:"id"`
	UserServiceID string    `json:"user_service_id" db:"user_service_id"`
	Status        string    `json:"status" db:"status"`
	StartedAt     time.Time `json:"started_at" db:"started_at"`
	FinishedAt    time.Time `json:"finished_at" db:"finished_at"`
	Error         string    `json:"error" db:"error"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
type SyncLog struct {
	ID        string    `json:"id" db:"id"`
	SyncJobID string    `json:"sync_job_id" db:"sync_job_id"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

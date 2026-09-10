package entity

import "time"

const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusFailed    = "failed"
	StatusCompleted = "completed"
)

type Job struct {
	ID        string
	Task      string
	Status    string
	Attempts  int32
	MaxRetry  int32
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobStatus struct {
	Pending   int32
	Running   int32
	Failed    int32
	Completed int32
}
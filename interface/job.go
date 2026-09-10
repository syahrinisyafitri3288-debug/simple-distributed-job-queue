package _interface

import (
	"context"
	"jobqueue/entity"
)

type JobService interface {
	Enqueue(ctx context.Context, taskName string) (*entity.Job, error)
	SimultaneousCreateJob(ctx context.Context, tasks []string) ([]*entity.Job, error)
	SimulateUnstableJob(ctx context.Context) (*entity.Job, error)
	GetAllJobs(ctx context.Context) ([]*entity.Job, error)
	GetJobById(ctx context.Context, id string) (*entity.Job, error)
	GetAllJobStatus(ctx context.Context) (*entity.JobStatus, error)
}

type JobRepository interface {
	Save(ctx context.Context, job *entity.Job) error
	FindByID(ctx context.Context, id string) (*entity.Job, error)
	FindAll(ctx context.Context) ([]*entity.Job, error)
}
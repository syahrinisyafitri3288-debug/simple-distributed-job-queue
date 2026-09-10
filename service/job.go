package service

import (
	"context"
	"fmt"
	"jobqueue/entity"
	_interface "jobqueue/interface"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type jobService struct {
	jobRepo _interface.JobRepository
	mu      sync.Mutex
}

type Initiator func(s *jobService) *jobService

func NewJobService() Initiator {
	return func(s *jobService) *jobService { return s }
}

func (i Initiator) SetJobRepository(r _interface.JobRepository) Initiator {
	return func(s *jobService) *jobService {
		i(s).jobRepo = r
		return s
	}
}

func (i Initiator) Build() _interface.JobService {
	return i(&jobService{})
}

func (s *jobService) Enqueue(ctx context.Context, taskName string) (*entity.Job, error) {
	job := &entity.Job{
		ID:        uuid.NewString(),
		Task:      taskName,
		Status:    entity.StatusPending,
		MaxRetry:  3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.jobRepo.Save(ctx, job); err != nil {
		return nil, err
	}

	go s.process(context.Background(), job)

	return job, nil
}

func (s *jobService) SimultaneousCreateJob(ctx context.Context, tasks []string) ([]*entity.Job, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []*entity.Job
		errs    []error
	)

	for _, t := range tasks {
		wg.Add(1)
		go func(taskName string) {
			defer wg.Done()
			job, err := s.Enqueue(ctx, taskName)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			results = append(results, job)
		}(t)
	}
	wg.Wait()

	if len(errs) > 0 {
		return results, fmt.Errorf("some jobs failed to enqueue: %v", errs)
	}
	return results, nil
}

func (s *jobService) SimulateUnstableJob(ctx context.Context) (*entity.Job, error) {
	return s.Enqueue(ctx, "unstable-job")
}

func (s *jobService) GetAllJobs(ctx context.Context) ([]*entity.Job, error) {
	return s.jobRepo.FindAll(ctx)
}

func (s *jobService) GetJobById(ctx context.Context, id string) (*entity.Job, error) {
	return s.jobRepo.FindByID(ctx, id)
}

func (s *jobService) GetAllJobStatus(ctx context.Context) (*entity.JobStatus, error) {
	jobs, err := s.jobRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	stats := &entity.JobStatus{}
	for _, j := range jobs {
		switch j.Status {
		case entity.StatusPending:
			stats.Pending++
		case entity.StatusRunning:
			stats.Running++
		case entity.StatusFailed:
			stats.Failed++
		case entity.StatusCompleted:
			stats.Completed++
		}
	}
	return stats, nil
}

func (s *jobService) process(ctx context.Context, job *entity.Job) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[job %s] recovered from panic: %v", job.ID, r)
			s.mu.Lock()
			job.Status = entity.StatusFailed
			job.Error = fmt.Sprintf("panic: %v", r)
			job.UpdatedAt = time.Now()
			s.jobRepo.Save(ctx, job)
			s.mu.Unlock()
		}
	}()

	var attempt int32
	for attempt = 1; attempt <= job.MaxRetry; attempt++ {
		s.mu.Lock()
		job.Attempts = attempt
		job.Status = entity.StatusRunning
		job.UpdatedAt = time.Now()
		s.jobRepo.Save(ctx, job)
		s.mu.Unlock()

		log.Printf("[job %s] task=%s attempt=%d running", job.ID, job.Task, attempt)

		err := s.execute(job, attempt)

		s.mu.Lock()
		if err == nil {
			job.Status = entity.StatusCompleted
			job.Error = ""
			job.UpdatedAt = time.Now()
			s.jobRepo.Save(ctx, job)
			s.mu.Unlock()
			log.Printf("[job %s] completed", job.ID)
			return
		}

		job.Error = err.Error()
		job.UpdatedAt = time.Now()
		s.jobRepo.Save(ctx, job)
		s.mu.Unlock()

		log.Printf("[job %s] attempt=%d failed: %v", job.ID, attempt, err)

		if attempt < job.MaxRetry {
			time.Sleep(500 * time.Millisecond)
		}
	}

	s.mu.Lock()
	job.Status = entity.StatusFailed
	job.UpdatedAt = time.Now()
	s.jobRepo.Save(ctx, job)
	s.mu.Unlock()
	log.Printf("[job %s] permanently failed after %d attempts", job.ID, job.Attempts)
}

func (s *jobService) execute(job *entity.Job, attempt int32) error {
	time.Sleep(200 * time.Millisecond)

	if job.Task == "unstable-job" && attempt < 3 {
		return fmt.Errorf("simulated failure on attempt %d", attempt)
	}
	return nil
}

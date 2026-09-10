package service

import (
	"context"
	"fmt"
	"jobqueue/entity"
	inmemrepo "jobqueue/repository/inmem"
	"testing"
)

func TestConcurrentLoad(t *testing.T) {
	repo := inmemrepo.NewJobRepository().
		SetInMemConnection(make(map[string]*entity.Job)).
		Build()
	svc := NewJobService().SetJobRepository(repo).Build()

	tasks := make([]string, 100)
	for i := range tasks {
		tasks[i] = fmt.Sprintf("load-task-%d", i)
	}

	jobs, err := svc.SimultaneousCreateJob(context.Background(), tasks)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(jobs) != 100 {
		t.Fatalf("expected 100 jobs, got %d", len(jobs))
	}
}

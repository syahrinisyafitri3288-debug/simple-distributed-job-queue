package mutation

import (
	"context"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/resolver"
	_interface "jobqueue/interface"
)

type JobMutation struct {
	jobService _interface.JobService
	dataloader *_dataloader.GeneralDataloader
}

type enqueueArgs struct {
	Task string
}

func (q JobMutation) Enqueue(ctx context.Context, args enqueueArgs) (*resolver.JobResolver, error) {
	job, err := q.jobService.Enqueue(ctx, args.Task)
	if err != nil {
		return nil, err
	}
	return &resolver.JobResolver{
		Data:       *job,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

type simultaneousCreateJobArgs struct {
	Job1 string
	Job2 string
	Job3 string
}

func (q JobMutation) SimultaneousCreateJob(ctx context.Context, args simultaneousCreateJobArgs) ([]*resolver.JobResolver, error) {
	tasks := []string{args.Job1, args.Job2, args.Job3}
	jobs, err := q.jobService.SimultaneousCreateJob(ctx, tasks)
	if err != nil {
		return nil, err
	}

	resolvers := make([]*resolver.JobResolver, 0, len(jobs))
	for _, j := range jobs {
		resolvers = append(resolvers, &resolver.JobResolver{
			Data:       *j,
			JobService: q.jobService,
			Dataloader: q.dataloader,
		})
	}
	return resolvers, nil
}

func (q JobMutation) SimulateUnstableJob(ctx context.Context) (*resolver.JobResolver, error) {
	job, err := q.jobService.SimulateUnstableJob(ctx)
	if err != nil {
		return nil, err
	}
	return &resolver.JobResolver{
		Data:       *job,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

func NewJobMutation(jobService _interface.JobService, dataloader *_dataloader.GeneralDataloader) JobMutation {
	return JobMutation{
		jobService: jobService,
		dataloader: dataloader,
	}
}

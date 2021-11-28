package service

import (
	context "context"
	"fmt"
	"github.com/beoboo/job-scheduler/library/errors"
	"github.com/beoboo/job-scheduler/library/scheduler"
	"github.com/beoboo/job-scheduler/pkg/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"strings"
)

type GrpcJobService struct {
	scheduler *scheduler.Scheduler

	protocol.UnimplementedJobSchedulerServer
}

func NewGrpcJobService() *GrpcJobService {
	return &GrpcJobService{
		scheduler: scheduler.New(true),
	}
}

func (s GrpcJobService) Start(ctx context.Context, job *protocol.Job) (*protocol.JobStatus, error) {
	executable := strings.TrimSpace(job.Executable)
	if executable == "" {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid executable: \"%s\"", executable))
	}

	id, err := s.scheduler.Start(job.Executable, job.Args)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, fmt.Sprintf("cannot start executable: \"%s\": %s", executable, err))
	}
	sts, err := s.scheduler.Status(id)

	return &protocol.JobStatus{Status: sts, Id: id}, nil
}

func (s GrpcJobService) Stop(ctx context.Context, id *protocol.JobId) (*protocol.JobStatus, error) {
	sts, err := s.scheduler.Stop(id.Id)
	if err != nil {
		return nil, status.Errorf(errorCode(err), err.Error())
	}

	return &protocol.JobStatus{Status: sts, Id: id.Id}, nil
}

func (s GrpcJobService) Status(ctx context.Context, id *protocol.JobId) (*protocol.JobStatus, error) {
	sts, err := s.scheduler.Status(id.Id)
	if err != nil {
		return nil, status.Errorf(errorCode(err), err.Error())
	}

	return &protocol.JobStatus{Status: sts, Id: id.Id}, nil
}

func (s GrpcJobService) Output(id *protocol.JobId, stream protocol.JobScheduler_OutputServer) error {
	output, err := s.scheduler.Output(id.Id)
	if err != nil {
		return status.Errorf(errorCode(err), err.Error())
	}

	for {
		o, err := output.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(errorCode(err), err.Error())
		}
		if err := stream.Send(&protocol.JobOutput{
			Channel: o.Channel,
			Text:    o.Text,
			Time:    fmt.Sprintf("%d", o.Time),
		}); err != nil {
			return err
		}
	}

	return nil
}

func errorCode(error error) codes.Code {
	switch error.(type) {
	case *errors.NotFoundError:
		return codes.NotFound
	default:
		return codes.Internal
	}
}

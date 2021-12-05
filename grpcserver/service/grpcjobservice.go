package service

import (
	context "context"
	"fmt"
	"github.com/beoboo/job-scheduler/library/errors"
	"github.com/beoboo/job-scheduler/library/scheduler"
	"github.com/beoboo/job-scheduler/pkg/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

type GrpcJobService struct {
	scheduler *scheduler.Scheduler

	protocol.UnimplementedJobSchedulerServer
}

func NewGrpcJobService() *GrpcJobService {
	return &GrpcJobService{
		scheduler: scheduler.NewSelf(),
	}
}

func (s GrpcJobService) Start(ctx context.Context, job *protocol.Job) (*protocol.JobStatus, error) {
	executable := strings.TrimSpace(job.Executable)
	if executable == "" {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid executable: \"%s\"", executable))
	}

	id, err := s.scheduler.Start(job.Executable, 0, job.Args)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, fmt.Sprintf("cannot start executable: \"%s\": %s", executable, err))
	}
	sts, err := s.scheduler.Status(id)

	return toProtocol(id, sts), nil
}

func toProtocol(id string, st *scheduler.JobStatus) *protocol.JobStatus {
	return &protocol.JobStatus{
		Id:       id,
		Status:   protocol.JobStatus_Status(st.Type),
		ExitCode: int32(st.ExitCode),
	}
}

func (s GrpcJobService) Stop(ctx context.Context, id *protocol.JobId) (*protocol.JobStatus, error) {
	sts, err := s.scheduler.Stop(id.Id)
	if err != nil {
		return nil, status.Errorf(errorCode(err), err.Error())
	}

	return toProtocol(id.Id, sts), nil
}

func (s GrpcJobService) Status(ctx context.Context, id *protocol.JobId) (*protocol.JobStatus, error) {
	sts, err := s.scheduler.Status(id.Id)
	if err != nil {
		return nil, status.Errorf(errorCode(err), err.Error())
	}

	return toProtocol(id.Id, sts), nil
}

func (s GrpcJobService) Output(id *protocol.JobId, stream protocol.JobScheduler_OutputServer) error {
	output, err := s.scheduler.Output(id.Id)
	if err != nil {
		return status.Errorf(errorCode(err), err.Error())
	}

	for line := range output.Read() {
		if line == nil {
			break
		}

		if err := stream.Send(&protocol.JobOutput{
			Type: protocol.JobOutput_Type(line.Type),
			Text: line.Text,
			Time: timestamppb.New(line.Time),
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

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protocol

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// JobSchedulerClient is the client API for JobScheduler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JobSchedulerClient interface {
	Start(ctx context.Context, in *Job, opts ...grpc.CallOption) (*JobStatus, error)
	Stop(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*JobStatus, error)
	Status(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*JobStatus, error)
	Output(ctx context.Context, in *JobId, opts ...grpc.CallOption) (JobScheduler_OutputClient, error)
}

type jobSchedulerClient struct {
	cc grpc.ClientConnInterface
}

func NewJobSchedulerClient(cc grpc.ClientConnInterface) JobSchedulerClient {
	return &jobSchedulerClient{cc}
}

func (c *jobSchedulerClient) Start(ctx context.Context, in *Job, opts ...grpc.CallOption) (*JobStatus, error) {
	out := new(JobStatus)
	err := c.cc.Invoke(ctx, "/JobScheduler/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobSchedulerClient) Stop(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*JobStatus, error) {
	out := new(JobStatus)
	err := c.cc.Invoke(ctx, "/JobScheduler/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobSchedulerClient) Status(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*JobStatus, error) {
	out := new(JobStatus)
	err := c.cc.Invoke(ctx, "/JobScheduler/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobSchedulerClient) Output(ctx context.Context, in *JobId, opts ...grpc.CallOption) (JobScheduler_OutputClient, error) {
	stream, err := c.cc.NewStream(ctx, &JobScheduler_ServiceDesc.Streams[0], "/JobScheduler/Output", opts...)
	if err != nil {
		return nil, err
	}
	x := &jobSchedulerOutputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type JobScheduler_OutputClient interface {
	Recv() (*JobOutput, error)
	grpc.ClientStream
}

type jobSchedulerOutputClient struct {
	grpc.ClientStream
}

func (x *jobSchedulerOutputClient) Recv() (*JobOutput, error) {
	m := new(JobOutput)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// JobSchedulerServer is the server API for JobScheduler service.
// All implementations must embed UnimplementedJobSchedulerServer
// for forward compatibility
type JobSchedulerServer interface {
	Start(context.Context, *Job) (*JobStatus, error)
	Stop(context.Context, *JobId) (*JobStatus, error)
	Status(context.Context, *JobId) (*JobStatus, error)
	Output(*JobId, JobScheduler_OutputServer) error
	mustEmbedUnimplementedJobSchedulerServer()
}

// UnimplementedJobSchedulerServer must be embedded to have forward compatible implementations.
type UnimplementedJobSchedulerServer struct {
}

func (UnimplementedJobSchedulerServer) Start(context.Context, *Job) (*JobStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedJobSchedulerServer) Stop(context.Context, *JobId) (*JobStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedJobSchedulerServer) Status(context.Context, *JobId) (*JobStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedJobSchedulerServer) Output(*JobId, JobScheduler_OutputServer) error {
	return status.Errorf(codes.Unimplemented, "method Output not implemented")
}
func (UnimplementedJobSchedulerServer) mustEmbedUnimplementedJobSchedulerServer() {}

// UnsafeJobSchedulerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JobSchedulerServer will
// result in compilation errors.
type UnsafeJobSchedulerServer interface {
	mustEmbedUnimplementedJobSchedulerServer()
}

func RegisterJobSchedulerServer(s grpc.ServiceRegistrar, srv JobSchedulerServer) {
	s.RegisterService(&JobScheduler_ServiceDesc, srv)
}

func _JobScheduler_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobSchedulerServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JobScheduler/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobSchedulerServer).Start(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobScheduler_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobSchedulerServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JobScheduler/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobSchedulerServer).Stop(ctx, req.(*JobId))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobScheduler_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobSchedulerServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JobScheduler/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobSchedulerServer).Status(ctx, req.(*JobId))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobScheduler_Output_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(JobId)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JobSchedulerServer).Output(m, &jobSchedulerOutputServer{stream})
}

type JobScheduler_OutputServer interface {
	Send(*JobOutput) error
	grpc.ServerStream
}

type jobSchedulerOutputServer struct {
	grpc.ServerStream
}

func (x *jobSchedulerOutputServer) Send(m *JobOutput) error {
	return x.ServerStream.SendMsg(m)
}

// JobScheduler_ServiceDesc is the grpc.ServiceDesc for JobScheduler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var JobScheduler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "JobScheduler",
	HandlerType: (*JobSchedulerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Start",
			Handler:    _JobScheduler_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _JobScheduler_Stop_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _JobScheduler_Status_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Output",
			Handler:       _JobScheduler_Output_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "job_worker_service.proto",
}
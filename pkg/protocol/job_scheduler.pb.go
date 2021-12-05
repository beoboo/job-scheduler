// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: job_scheduler.proto

package protocol

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type JobStatus_Status int32

const (
	JobStatus_IDLE     JobStatus_Status = 0
	JobStatus_RUNNING  JobStatus_Status = 1
	JobStatus_FINISHED JobStatus_Status = 2
	JobStatus_KILLED   JobStatus_Status = 3
	JobStatus_ERRORED  JobStatus_Status = 4
)

// Enum value maps for JobStatus_Status.
var (
	JobStatus_Status_name = map[int32]string{
		0: "IDLE",
		1: "RUNNING",
		2: "FINISHED",
		3: "KILLED",
		4: "ERRORED",
	}
	JobStatus_Status_value = map[string]int32{
		"IDLE":     0,
		"RUNNING":  1,
		"FINISHED": 2,
		"KILLED":   3,
		"ERRORED":  4,
	}
)

func (x JobStatus_Status) Enum() *JobStatus_Status {
	p := new(JobStatus_Status)
	*p = x
	return p
}

func (x JobStatus_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (JobStatus_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_job_scheduler_proto_enumTypes[0].Descriptor()
}

func (JobStatus_Status) Type() protoreflect.EnumType {
	return &file_job_scheduler_proto_enumTypes[0]
}

func (x JobStatus_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use JobStatus_Status.Descriptor instead.
func (JobStatus_Status) EnumDescriptor() ([]byte, []int) {
	return file_job_scheduler_proto_rawDescGZIP(), []int{2, 0}
}

type JobOutput_Type int32

const (
	JobOutput_OUTPUT JobOutput_Type = 0
	JobOutput_ERROR  JobOutput_Type = 1
)

// Enum value maps for JobOutput_Type.
var (
	JobOutput_Type_name = map[int32]string{
		0: "OUTPUT",
		1: "ERROR",
	}
	JobOutput_Type_value = map[string]int32{
		"OUTPUT": 0,
		"ERROR":  1,
	}
)

func (x JobOutput_Type) Enum() *JobOutput_Type {
	p := new(JobOutput_Type)
	*p = x
	return p
}

func (x JobOutput_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (JobOutput_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_job_scheduler_proto_enumTypes[1].Descriptor()
}

func (JobOutput_Type) Type() protoreflect.EnumType {
	return &file_job_scheduler_proto_enumTypes[1]
}

func (x JobOutput_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use JobOutput_Type.Descriptor instead.
func (JobOutput_Type) EnumDescriptor() ([]byte, []int) {
	return file_job_scheduler_proto_rawDescGZIP(), []int{3, 0}
}

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Executable string  `protobuf:"bytes,1,opt,name=executable,proto3" json:"executable,omitempty"`
	Args       string  `protobuf:"bytes,2,opt,name=args,proto3" json:"args,omitempty"`
	Cpu        float32 `protobuf:"fixed32,3,opt,name=cpu,proto3" json:"cpu,omitempty"`      // 0 = unlimited, # shares per second
	Memory     int32   `protobuf:"varint,4,opt,name=memory,proto3" json:"memory,omitempty"` // 0 = unlimited, in bytes
	Io         int32   `protobuf:"varint,5,opt,name=io,proto3" json:"io,omitempty"`         // 0 = unlimited, 100-1000 otherwise
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_job_scheduler_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_job_scheduler_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_job_scheduler_proto_rawDescGZIP(), []int{0}
}

func (x *Job) GetExecutable() string {
	if x != nil {
		return x.Executable
	}
	return ""
}

func (x *Job) GetArgs() string {
	if x != nil {
		return x.Args
	}
	return ""
}

func (x *Job) GetCpu() float32 {
	if x != nil {
		return x.Cpu
	}
	return 0
}

func (x *Job) GetMemory() int32 {
	if x != nil {
		return x.Memory
	}
	return 0
}

func (x *Job) GetIo() int32 {
	if x != nil {
		return x.Io
	}
	return 0
}

type JobId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *JobId) Reset() {
	*x = JobId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_job_scheduler_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobId) ProtoMessage() {}

func (x *JobId) ProtoReflect() protoreflect.Message {
	mi := &file_job_scheduler_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobId.ProtoReflect.Descriptor instead.
func (*JobId) Descriptor() ([]byte, []int) {
	return file_job_scheduler_proto_rawDescGZIP(), []int{1}
}

func (x *JobId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type JobStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status   JobStatus_Status `protobuf:"varint,2,opt,name=status,proto3,enum=JobStatus_Status" json:"status,omitempty"`
	ExitCode int32            `protobuf:"varint,3,opt,name=exitCode,proto3" json:"exitCode,omitempty"`
}

func (x *JobStatus) Reset() {
	*x = JobStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_job_scheduler_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobStatus) ProtoMessage() {}

func (x *JobStatus) ProtoReflect() protoreflect.Message {
	mi := &file_job_scheduler_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobStatus.ProtoReflect.Descriptor instead.
func (*JobStatus) Descriptor() ([]byte, []int) {
	return file_job_scheduler_proto_rawDescGZIP(), []int{2}
}

func (x *JobStatus) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *JobStatus) GetStatus() JobStatus_Status {
	if x != nil {
		return x.Status
	}
	return JobStatus_IDLE
}

func (x *JobStatus) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

type JobOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type JobOutput_Type         `protobuf:"varint,1,opt,name=type,proto3,enum=JobOutput_Type" json:"type,omitempty"`
	Time *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	Text []byte                 `protobuf:"bytes,3,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *JobOutput) Reset() {
	*x = JobOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_job_scheduler_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobOutput) ProtoMessage() {}

func (x *JobOutput) ProtoReflect() protoreflect.Message {
	mi := &file_job_scheduler_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobOutput.ProtoReflect.Descriptor instead.
func (*JobOutput) Descriptor() ([]byte, []int) {
	return file_job_scheduler_proto_rawDescGZIP(), []int{3}
}

func (x *JobOutput) GetType() JobOutput_Type {
	if x != nil {
		return x.Type
	}
	return JobOutput_OUTPUT
}

func (x *JobOutput) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *JobOutput) GetText() []byte {
	if x != nil {
		return x.Text
	}
	return nil
}

var File_job_scheduler_proto protoreflect.FileDescriptor

var file_job_scheduler_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6a, 0x6f, 0x62, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a, 0x03, 0x4a, 0x6f, 0x62, 0x12, 0x1e, 0x0a,
	0x0a, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x67,
	0x73, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03,
	0x63, 0x70, 0x75, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x6f, 0x22, 0x17, 0x0a, 0x05, 0x4a,
	0x6f, 0x62, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0xaa, 0x01, 0x0a, 0x09, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x29, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x11, 0x2e, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a,
	0x08, 0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x46, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x44, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x0b, 0x0a,
	0x07, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x46, 0x49,
	0x4e, 0x49, 0x53, 0x48, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x4b, 0x49, 0x4c, 0x4c,
	0x45, 0x44, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x45, 0x44, 0x10,
	0x04, 0x22, 0x93, 0x01, 0x0a, 0x09, 0x4a, 0x6f, 0x62, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12,
	0x23, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e,
	0x4a, 0x6f, 0x62, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04,
	0x74, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x1d, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x0a, 0x0a, 0x06, 0x4f, 0x55, 0x54, 0x50, 0x55, 0x54, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x01, 0x32, 0x8b, 0x01, 0x0a, 0x0c, 0x4a, 0x6f, 0x62, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x12, 0x04, 0x2e, 0x4a, 0x6f, 0x62, 0x1a, 0x0a, 0x2e, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x1c, 0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70, 0x12, 0x06, 0x2e,
	0x4a, 0x6f, 0x62, 0x49, 0x64, 0x1a, 0x0a, 0x2e, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x00, 0x12, 0x1e, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x06, 0x2e,
	0x4a, 0x6f, 0x62, 0x49, 0x64, 0x1a, 0x0a, 0x2e, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x00, 0x12, 0x20, 0x0a, 0x06, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x06, 0x2e,
	0x4a, 0x6f, 0x62, 0x49, 0x64, 0x1a, 0x0a, 0x2e, 0x4a, 0x6f, 0x62, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x22, 0x00, 0x30, 0x01, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x65, 0x6f, 0x62, 0x6f, 0x6f, 0x2f, 0x6a, 0x6f, 0x62, 0x2d, 0x73,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_job_scheduler_proto_rawDescOnce sync.Once
	file_job_scheduler_proto_rawDescData = file_job_scheduler_proto_rawDesc
)

func file_job_scheduler_proto_rawDescGZIP() []byte {
	file_job_scheduler_proto_rawDescOnce.Do(func() {
		file_job_scheduler_proto_rawDescData = protoimpl.X.CompressGZIP(file_job_scheduler_proto_rawDescData)
	})
	return file_job_scheduler_proto_rawDescData
}

var file_job_scheduler_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_job_scheduler_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_job_scheduler_proto_goTypes = []interface{}{
	(JobStatus_Status)(0),         // 0: JobStatus.Status
	(JobOutput_Type)(0),           // 1: JobOutput.Type
	(*Job)(nil),                   // 2: Job
	(*JobId)(nil),                 // 3: JobId
	(*JobStatus)(nil),             // 4: JobStatus
	(*JobOutput)(nil),             // 5: JobOutput
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_job_scheduler_proto_depIdxs = []int32{
	0, // 0: JobStatus.status:type_name -> JobStatus.Status
	1, // 1: JobOutput.type:type_name -> JobOutput.Type
	6, // 2: JobOutput.time:type_name -> google.protobuf.Timestamp
	2, // 3: JobScheduler.Start:input_type -> Job
	3, // 4: JobScheduler.Stop:input_type -> JobId
	3, // 5: JobScheduler.Status:input_type -> JobId
	3, // 6: JobScheduler.Output:input_type -> JobId
	4, // 7: JobScheduler.Start:output_type -> JobStatus
	4, // 8: JobScheduler.Stop:output_type -> JobStatus
	4, // 9: JobScheduler.Status:output_type -> JobStatus
	5, // 10: JobScheduler.Output:output_type -> JobOutput
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_job_scheduler_proto_init() }
func file_job_scheduler_proto_init() {
	if File_job_scheduler_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_job_scheduler_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_job_scheduler_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_job_scheduler_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_job_scheduler_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobOutput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_job_scheduler_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_job_scheduler_proto_goTypes,
		DependencyIndexes: file_job_scheduler_proto_depIdxs,
		EnumInfos:         file_job_scheduler_proto_enumTypes,
		MessageInfos:      file_job_scheduler_proto_msgTypes,
	}.Build()
	File_job_scheduler_proto = out.File
	file_job_scheduler_proto_rawDesc = nil
	file_job_scheduler_proto_goTypes = nil
	file_job_scheduler_proto_depIdxs = nil
}

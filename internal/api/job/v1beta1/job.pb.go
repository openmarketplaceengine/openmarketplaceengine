// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: api/job/v1beta1/job.proto

package v1beta1

import (
	v1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/type/v1beta1"
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

type JobAction int32

const (
	JobAction_JOB_ACTION_UNSPECIFIED JobAction = 0
	JobAction_JOB_ACTION_CREATED     JobAction = 1
	JobAction_JOB_ACTION_UPDATED     JobAction = 2
)

// Enum value maps for JobAction.
var (
	JobAction_name = map[int32]string{
		0: "JOB_ACTION_UNSPECIFIED",
		1: "JOB_ACTION_CREATED",
		2: "JOB_ACTION_UPDATED",
	}
	JobAction_value = map[string]int32{
		"JOB_ACTION_UNSPECIFIED": 0,
		"JOB_ACTION_CREATED":     1,
		"JOB_ACTION_UPDATED":     2,
	}
)

func (x JobAction) Enum() *JobAction {
	p := new(JobAction)
	*p = x
	return p
}

func (x JobAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (JobAction) Descriptor() protoreflect.EnumDescriptor {
	return file_api_job_v1beta1_job_proto_enumTypes[0].Descriptor()
}

func (JobAction) Type() protoreflect.EnumType {
	return &file_api_job_v1beta1_job_proto_enumTypes[0]
}

func (x JobAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use JobAction.Descriptor instead.
func (JobAction) EnumDescriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{0}
}

type JobInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	WorkerId    string                 `protobuf:"bytes,2,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
	Created     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created,proto3" json:"created,omitempty"`
	Updated     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated,proto3" json:"updated,omitempty"`
	State       string                 `protobuf:"bytes,5,opt,name=state,proto3" json:"state,omitempty"`
	PickupDate  *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=pickup_date,json=pickupDate,proto3" json:"pickup_date,omitempty"`
	PickupAddr  string                 `protobuf:"bytes,7,opt,name=pickup_addr,json=pickupAddr,proto3" json:"pickup_addr,omitempty"`
	PickupLoc   *v1beta1.Location      `protobuf:"bytes,8,opt,name=pickup_loc,json=pickupLoc,proto3" json:"pickup_loc,omitempty"`
	DropoffAddr string                 `protobuf:"bytes,9,opt,name=dropoff_addr,json=dropoffAddr,proto3" json:"dropoff_addr,omitempty"`
	DropoffLoc  *v1beta1.Location      `protobuf:"bytes,10,opt,name=dropoff_loc,json=dropoffLoc,proto3" json:"dropoff_loc,omitempty"`
	TripType    string                 `protobuf:"bytes,11,opt,name=trip_type,json=tripType,proto3" json:"trip_type,omitempty"`
	Category    string                 `protobuf:"bytes,12,opt,name=category,proto3" json:"category,omitempty"`
}

func (x *JobInfo) Reset() {
	*x = JobInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobInfo) ProtoMessage() {}

func (x *JobInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobInfo.ProtoReflect.Descriptor instead.
func (*JobInfo) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{0}
}

func (x *JobInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *JobInfo) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

func (x *JobInfo) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *JobInfo) GetUpdated() *timestamppb.Timestamp {
	if x != nil {
		return x.Updated
	}
	return nil
}

func (x *JobInfo) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *JobInfo) GetPickupDate() *timestamppb.Timestamp {
	if x != nil {
		return x.PickupDate
	}
	return nil
}

func (x *JobInfo) GetPickupAddr() string {
	if x != nil {
		return x.PickupAddr
	}
	return ""
}

func (x *JobInfo) GetPickupLoc() *v1beta1.Location {
	if x != nil {
		return x.PickupLoc
	}
	return nil
}

func (x *JobInfo) GetDropoffAddr() string {
	if x != nil {
		return x.DropoffAddr
	}
	return ""
}

func (x *JobInfo) GetDropoffLoc() *v1beta1.Location {
	if x != nil {
		return x.DropoffLoc
	}
	return nil
}

func (x *JobInfo) GetTripType() string {
	if x != nil {
		return x.TripType
	}
	return ""
}

func (x *JobInfo) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

type ImportJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job *JobInfo `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
}

func (x *ImportJobRequest) Reset() {
	*x = ImportJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImportJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportJobRequest) ProtoMessage() {}

func (x *ImportJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportJobRequest.ProtoReflect.Descriptor instead.
func (*ImportJobRequest) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{1}
}

func (x *ImportJobRequest) GetJob() *JobInfo {
	if x != nil {
		return x.Job
	}
	return nil
}

type ImportJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action JobAction `protobuf:"varint,1,opt,name=action,proto3,enum=api.job.v1beta1.JobAction" json:"action,omitempty"`
}

func (x *ImportJobResponse) Reset() {
	*x = ImportJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImportJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportJobResponse) ProtoMessage() {}

func (x *ImportJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportJobResponse.ProtoReflect.Descriptor instead.
func (*ImportJobResponse) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{2}
}

func (x *ImportJobResponse) GetAction() JobAction {
	if x != nil {
		return x.Action
	}
	return JobAction_JOB_ACTION_UNSPECIFIED
}

type ExportJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
}

func (x *ExportJobRequest) Reset() {
	*x = ExportJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExportJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExportJobRequest) ProtoMessage() {}

func (x *ExportJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExportJobRequest.ProtoReflect.Descriptor instead.
func (*ExportJobRequest) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{3}
}

func (x *ExportJobRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

type ExportJobItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Job *JobInfo `protobuf:"bytes,2,opt,name=job,proto3" json:"job,omitempty"`
}

func (x *ExportJobItem) Reset() {
	*x = ExportJobItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExportJobItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExportJobItem) ProtoMessage() {}

func (x *ExportJobItem) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExportJobItem.ProtoReflect.Descriptor instead.
func (*ExportJobItem) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{4}
}

func (x *ExportJobItem) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ExportJobItem) GetJob() *JobInfo {
	if x != nil {
		return x.Job
	}
	return nil
}

type ExportJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jobs []*ExportJobItem `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
}

func (x *ExportJobResponse) Reset() {
	*x = ExportJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExportJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExportJobResponse) ProtoMessage() {}

func (x *ExportJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExportJobResponse.ProtoReflect.Descriptor instead.
func (*ExportJobResponse) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{5}
}

func (x *ExportJobResponse) GetJobs() []*ExportJobItem {
	if x != nil {
		return x.Jobs
	}
	return nil
}

type GetAvailableJobsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Any non-empty string representing geographic region
	// i.e. san-fran, new-york, alaska, us, global, test, etc.
	// Should be used consistently across all APIs.
	AreaKey           string `protobuf:"bytes,1,opt,name=area_key,json=areaKey,proto3" json:"area_key,omitempty"`
	WorkerId          string `protobuf:"bytes,2,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
	MaxDistanceMeters int32  `protobuf:"varint,3,opt,name=max_distance_meters,json=maxDistanceMeters,proto3" json:"max_distance_meters,omitempty"`
	// Result set limit.
	Limit int32 `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetAvailableJobsRequest) Reset() {
	*x = GetAvailableJobsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAvailableJobsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAvailableJobsRequest) ProtoMessage() {}

func (x *GetAvailableJobsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAvailableJobsRequest.ProtoReflect.Descriptor instead.
func (*GetAvailableJobsRequest) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{6}
}

func (x *GetAvailableJobsRequest) GetAreaKey() string {
	if x != nil {
		return x.AreaKey
	}
	return ""
}

func (x *GetAvailableJobsRequest) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

func (x *GetAvailableJobsRequest) GetMaxDistanceMeters() int32 {
	if x != nil {
		return x.MaxDistanceMeters
	}
	return 0
}

func (x *GetAvailableJobsRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type AvailableJob struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job            *JobInfo `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	DistanceMeters int32    `protobuf:"varint,2,opt,name=distance_meters,json=distanceMeters,proto3" json:"distance_meters,omitempty"`
}

func (x *AvailableJob) Reset() {
	*x = AvailableJob{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AvailableJob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AvailableJob) ProtoMessage() {}

func (x *AvailableJob) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AvailableJob.ProtoReflect.Descriptor instead.
func (*AvailableJob) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{7}
}

func (x *AvailableJob) GetJob() *JobInfo {
	if x != nil {
		return x.Job
	}
	return nil
}

func (x *AvailableJob) GetDistanceMeters() int32 {
	if x != nil {
		return x.DistanceMeters
	}
	return 0
}

type GetAvailableJobsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jobs []*AvailableJob `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
}

func (x *GetAvailableJobsResponse) Reset() {
	*x = GetAvailableJobsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_job_v1beta1_job_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAvailableJobsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAvailableJobsResponse) ProtoMessage() {}

func (x *GetAvailableJobsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1beta1_job_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAvailableJobsResponse.ProtoReflect.Descriptor instead.
func (*GetAvailableJobsResponse) Descriptor() ([]byte, []int) {
	return file_api_job_v1beta1_job_proto_rawDescGZIP(), []int{8}
}

func (x *GetAvailableJobsResponse) GetJobs() []*AvailableJob {
	if x != nil {
		return x.Jobs
	}
	return nil
}

var File_api_job_v1beta1_job_proto protoreflect.FileDescriptor

var file_api_job_v1beta1_job_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x6a, 0x6f, 0x62, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2f, 0x6a, 0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x61, 0x70, 0x69,
	0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x61,
	0x70, 0x69, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f,
	0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xea,
	0x03, 0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x34, 0x0a,
	0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x70, 0x69, 0x63,
	0x6b, 0x75, 0x70, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x70, 0x69, 0x63, 0x6b,
	0x75, 0x70, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x69, 0x63, 0x6b, 0x75, 0x70,
	0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x69, 0x63,
	0x6b, 0x75, 0x70, 0x41, 0x64, 0x64, 0x72, 0x12, 0x39, 0x0a, 0x0a, 0x70, 0x69, 0x63, 0x6b, 0x75,
	0x70, 0x5f, 0x6c, 0x6f, 0x63, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x70, 0x69, 0x63, 0x6b, 0x75, 0x70, 0x4c,
	0x6f, 0x63, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x72, 0x6f, 0x70, 0x6f, 0x66, 0x66, 0x5f, 0x61, 0x64,
	0x64, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x72, 0x6f, 0x70, 0x6f, 0x66,
	0x66, 0x41, 0x64, 0x64, 0x72, 0x12, 0x3b, 0x0a, 0x0b, 0x64, 0x72, 0x6f, 0x70, 0x6f, 0x66, 0x66,
	0x5f, 0x6c, 0x6f, 0x63, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x64, 0x72, 0x6f, 0x70, 0x6f, 0x66, 0x66, 0x4c,
	0x6f, 0x63, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x72, 0x69, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x72, 0x69, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x22, 0x3e, 0x0a, 0x10, 0x49,
	0x6d, 0x70, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2a, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4a,
	0x6f, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x22, 0x47, 0x0a, 0x11, 0x49,
	0x6d, 0x70, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x32, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x24, 0x0a, 0x10, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x4a, 0x6f,
	0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x22, 0x4b, 0x0a, 0x0d, 0x45, 0x78,
	0x70, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x03, 0x6a,
	0x6f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a,
	0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x22, 0x47, 0x0a, 0x11, 0x45, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x04,
	0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x45, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73,
	0x22, 0x97, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08,
	0x61, 0x72, 0x65, 0x61, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x61, 0x72, 0x65, 0x61, 0x4b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x69, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x11, 0x6d, 0x61, 0x78, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x4d, 0x65,
	0x74, 0x65, 0x72, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x63, 0x0a, 0x0c, 0x41, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4a, 0x6f, 0x62, 0x12, 0x2a, 0x0a, 0x03, 0x6a, 0x6f,
	0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x6f,
	0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x27, 0x0a, 0x0f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0e, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x22,
	0x4d, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4a,
	0x6f, 0x62, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x6a,
	0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x41, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x2a, 0x57,
	0x0a, 0x09, 0x4a, 0x6f, 0x62, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x16, 0x4a,
	0x4f, 0x42, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x4a, 0x4f, 0x42, 0x5f, 0x41,
	0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12,
	0x16, 0x0a, 0x12, 0x4a, 0x4f, 0x42, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x50,
	0x44, 0x41, 0x54, 0x45, 0x44, 0x10, 0x02, 0x32, 0xa3, 0x02, 0x0a, 0x0a, 0x4a, 0x6f, 0x62, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x09, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74,
	0x4a, 0x6f, 0x62, 0x12, 0x21, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x6f, 0x62,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x4a,
	0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x09,
	0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x12, 0x21, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x45, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x45,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x69, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62,
	0x6c, 0x65, 0x4a, 0x6f, 0x62, 0x73, 0x12, 0x28, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x6f, 0x62,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x29, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4a,
	0x6f, 0x62, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x51, 0x5a,
	0x4f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e,
	0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63,
	0x65, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6a, 0x6f, 0x62, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_job_v1beta1_job_proto_rawDescOnce sync.Once
	file_api_job_v1beta1_job_proto_rawDescData = file_api_job_v1beta1_job_proto_rawDesc
)

func file_api_job_v1beta1_job_proto_rawDescGZIP() []byte {
	file_api_job_v1beta1_job_proto_rawDescOnce.Do(func() {
		file_api_job_v1beta1_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_job_v1beta1_job_proto_rawDescData)
	})
	return file_api_job_v1beta1_job_proto_rawDescData
}

var file_api_job_v1beta1_job_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_job_v1beta1_job_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_job_v1beta1_job_proto_goTypes = []interface{}{
	(JobAction)(0),                   // 0: api.job.v1beta1.JobAction
	(*JobInfo)(nil),                  // 1: api.job.v1beta1.JobInfo
	(*ImportJobRequest)(nil),         // 2: api.job.v1beta1.ImportJobRequest
	(*ImportJobResponse)(nil),        // 3: api.job.v1beta1.ImportJobResponse
	(*ExportJobRequest)(nil),         // 4: api.job.v1beta1.ExportJobRequest
	(*ExportJobItem)(nil),            // 5: api.job.v1beta1.ExportJobItem
	(*ExportJobResponse)(nil),        // 6: api.job.v1beta1.ExportJobResponse
	(*GetAvailableJobsRequest)(nil),  // 7: api.job.v1beta1.GetAvailableJobsRequest
	(*AvailableJob)(nil),             // 8: api.job.v1beta1.AvailableJob
	(*GetAvailableJobsResponse)(nil), // 9: api.job.v1beta1.GetAvailableJobsResponse
	(*timestamppb.Timestamp)(nil),    // 10: google.protobuf.Timestamp
	(*v1beta1.Location)(nil),         // 11: api.type.v1beta1.Location
}
var file_api_job_v1beta1_job_proto_depIdxs = []int32{
	10, // 0: api.job.v1beta1.JobInfo.created:type_name -> google.protobuf.Timestamp
	10, // 1: api.job.v1beta1.JobInfo.updated:type_name -> google.protobuf.Timestamp
	10, // 2: api.job.v1beta1.JobInfo.pickup_date:type_name -> google.protobuf.Timestamp
	11, // 3: api.job.v1beta1.JobInfo.pickup_loc:type_name -> api.type.v1beta1.Location
	11, // 4: api.job.v1beta1.JobInfo.dropoff_loc:type_name -> api.type.v1beta1.Location
	1,  // 5: api.job.v1beta1.ImportJobRequest.job:type_name -> api.job.v1beta1.JobInfo
	0,  // 6: api.job.v1beta1.ImportJobResponse.action:type_name -> api.job.v1beta1.JobAction
	1,  // 7: api.job.v1beta1.ExportJobItem.job:type_name -> api.job.v1beta1.JobInfo
	5,  // 8: api.job.v1beta1.ExportJobResponse.jobs:type_name -> api.job.v1beta1.ExportJobItem
	1,  // 9: api.job.v1beta1.AvailableJob.job:type_name -> api.job.v1beta1.JobInfo
	8,  // 10: api.job.v1beta1.GetAvailableJobsResponse.jobs:type_name -> api.job.v1beta1.AvailableJob
	2,  // 11: api.job.v1beta1.JobService.ImportJob:input_type -> api.job.v1beta1.ImportJobRequest
	4,  // 12: api.job.v1beta1.JobService.ExportJob:input_type -> api.job.v1beta1.ExportJobRequest
	7,  // 13: api.job.v1beta1.JobService.GetAvailableJobs:input_type -> api.job.v1beta1.GetAvailableJobsRequest
	3,  // 14: api.job.v1beta1.JobService.ImportJob:output_type -> api.job.v1beta1.ImportJobResponse
	6,  // 15: api.job.v1beta1.JobService.ExportJob:output_type -> api.job.v1beta1.ExportJobResponse
	9,  // 16: api.job.v1beta1.JobService.GetAvailableJobs:output_type -> api.job.v1beta1.GetAvailableJobsResponse
	14, // [14:17] is the sub-list for method output_type
	11, // [11:14] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_api_job_v1beta1_job_proto_init() }
func file_api_job_v1beta1_job_proto_init() {
	if File_api_job_v1beta1_job_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_job_v1beta1_job_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobInfo); i {
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
		file_api_job_v1beta1_job_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImportJobRequest); i {
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
		file_api_job_v1beta1_job_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImportJobResponse); i {
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
		file_api_job_v1beta1_job_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExportJobRequest); i {
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
		file_api_job_v1beta1_job_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExportJobItem); i {
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
		file_api_job_v1beta1_job_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExportJobResponse); i {
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
		file_api_job_v1beta1_job_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAvailableJobsRequest); i {
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
		file_api_job_v1beta1_job_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AvailableJob); i {
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
		file_api_job_v1beta1_job_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAvailableJobsResponse); i {
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
			RawDescriptor: file_api_job_v1beta1_job_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_job_v1beta1_job_proto_goTypes,
		DependencyIndexes: file_api_job_v1beta1_job_proto_depIdxs,
		EnumInfos:         file_api_job_v1beta1_job_proto_enumTypes,
		MessageInfos:      file_api_job_v1beta1_job_proto_msgTypes,
	}.Build()
	File_api_job_v1beta1_job_proto = out.File
	file_api_job_v1beta1_job_proto_rawDesc = nil
	file_api_job_v1beta1_job_proto_goTypes = nil
	file_api_job_v1beta1_job_proto_depIdxs = nil
}

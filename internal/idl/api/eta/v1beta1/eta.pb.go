// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: api/eta/v1beta1/eta.proto

package v1beta1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Type int32

const (
	Type_TYPE_UNSPECIFIED         Type = 0
	Type_TYPE_DRIVER_TO_PICK_UP   Type = 1
	Type_TYPE_PICK_UP_TO_DROP_OFF Type = 2
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "TYPE_DRIVER_TO_PICK_UP",
		2: "TYPE_PICK_UP_TO_DROP_OFF",
	}
	Type_value = map[string]int32{
		"TYPE_UNSPECIFIED":         0,
		"TYPE_DRIVER_TO_PICK_UP":   1,
		"TYPE_PICK_UP_TO_DROP_OFF": 2,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_api_eta_v1beta1_eta_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_api_eta_v1beta1_eta_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_api_eta_v1beta1_eta_proto_rawDescGZIP(), []int{0}
}

type GetEstimatedJobsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Any non-empty string representing geographic region
	// i.e. san-fran, new-york, alaska, us, global, test, etc.
	// Should be used consistently across all APIs.
	AreaKey  string `protobuf:"bytes,1,opt,name=area_key,json=areaKey,proto3" json:"area_key,omitempty"`
	WorkerId string `protobuf:"bytes,2,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
	// Radius in meters to look for jobs relative to current worker position.
	RadiusMeters int32 `protobuf:"varint,3,opt,name=radius_meters,json=radiusMeters,proto3" json:"radius_meters,omitempty"`
	// Result set limit.
	Limit int32 `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	Eta   Type  `protobuf:"varint,5,opt,name=eta,proto3,enum=api.eta.v1beta1.Type" json:"eta,omitempty"`
}

func (x *GetEstimatedJobsRequest) Reset() {
	*x = GetEstimatedJobsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eta_v1beta1_eta_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEstimatedJobsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEstimatedJobsRequest) ProtoMessage() {}

func (x *GetEstimatedJobsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_eta_v1beta1_eta_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEstimatedJobsRequest.ProtoReflect.Descriptor instead.
func (*GetEstimatedJobsRequest) Descriptor() ([]byte, []int) {
	return file_api_eta_v1beta1_eta_proto_rawDescGZIP(), []int{0}
}

func (x *GetEstimatedJobsRequest) GetAreaKey() string {
	if x != nil {
		return x.AreaKey
	}
	return ""
}

func (x *GetEstimatedJobsRequest) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

func (x *GetEstimatedJobsRequest) GetRadiusMeters() int32 {
	if x != nil {
		return x.RadiusMeters
	}
	return 0
}

func (x *GetEstimatedJobsRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetEstimatedJobsRequest) GetEta() Type {
	if x != nil {
		return x.Eta
	}
	return Type_TYPE_UNSPECIFIED
}

type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lat     float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng     float64 `protobuf:"fixed64,2,opt,name=lng,proto3" json:"lng,omitempty"`
	Address string  `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eta_v1beta1_eta_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_api_eta_v1beta1_eta_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_api_eta_v1beta1_eta_proto_rawDescGZIP(), []int{1}
}

func (x *Location) GetLat() float64 {
	if x != nil {
		return x.Lat
	}
	return 0
}

func (x *Location) GetLng() float64 {
	if x != nil {
		return x.Lng
	}
	return 0
}

func (x *Location) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type Estimate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DistanceMeters int32                `protobuf:"varint,1,opt,name=distance_meters,json=distanceMeters,proto3" json:"distance_meters,omitempty"`
	Duration       *durationpb.Duration `protobuf:"bytes,2,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *Estimate) Reset() {
	*x = Estimate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eta_v1beta1_eta_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Estimate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Estimate) ProtoMessage() {}

func (x *Estimate) ProtoReflect() protoreflect.Message {
	mi := &file_api_eta_v1beta1_eta_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Estimate.ProtoReflect.Descriptor instead.
func (*Estimate) Descriptor() ([]byte, []int) {
	return file_api_eta_v1beta1_eta_proto_rawDescGZIP(), []int{2}
}

func (x *Estimate) GetDistanceMeters() int32 {
	if x != nil {
		return x.DistanceMeters
	}
	return 0
}

func (x *Estimate) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

type EstimatedJob struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                      string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	WorkerToPickupEstimate  *Estimate `protobuf:"bytes,2,opt,name=worker_to_pickup_estimate,json=workerToPickupEstimate,proto3" json:"worker_to_pickup_estimate,omitempty"`
	PickupToDropOffEstimate *Estimate `protobuf:"bytes,3,opt,name=pickup_to_drop_off_estimate,json=pickupToDropOffEstimate,proto3" json:"pickup_to_drop_off_estimate,omitempty"`
	WorkerLocation          *Location `protobuf:"bytes,4,opt,name=worker_location,json=workerLocation,proto3" json:"worker_location,omitempty"`
	PickupLocation          *Location `protobuf:"bytes,5,opt,name=pickup_location,json=pickupLocation,proto3" json:"pickup_location,omitempty"`
	DropOffLocation         *Location `protobuf:"bytes,6,opt,name=drop_off_location,json=dropOffLocation,proto3" json:"drop_off_location,omitempty"`
}

func (x *EstimatedJob) Reset() {
	*x = EstimatedJob{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eta_v1beta1_eta_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EstimatedJob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EstimatedJob) ProtoMessage() {}

func (x *EstimatedJob) ProtoReflect() protoreflect.Message {
	mi := &file_api_eta_v1beta1_eta_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EstimatedJob.ProtoReflect.Descriptor instead.
func (*EstimatedJob) Descriptor() ([]byte, []int) {
	return file_api_eta_v1beta1_eta_proto_rawDescGZIP(), []int{3}
}

func (x *EstimatedJob) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EstimatedJob) GetWorkerToPickupEstimate() *Estimate {
	if x != nil {
		return x.WorkerToPickupEstimate
	}
	return nil
}

func (x *EstimatedJob) GetPickupToDropOffEstimate() *Estimate {
	if x != nil {
		return x.PickupToDropOffEstimate
	}
	return nil
}

func (x *EstimatedJob) GetWorkerLocation() *Location {
	if x != nil {
		return x.WorkerLocation
	}
	return nil
}

func (x *EstimatedJob) GetPickupLocation() *Location {
	if x != nil {
		return x.PickupLocation
	}
	return nil
}

func (x *EstimatedJob) GetDropOffLocation() *Location {
	if x != nil {
		return x.DropOffLocation
	}
	return nil
}

type GetEstimatedJobsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jobs []*EstimatedJob `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
}

func (x *GetEstimatedJobsResponse) Reset() {
	*x = GetEstimatedJobsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eta_v1beta1_eta_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEstimatedJobsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEstimatedJobsResponse) ProtoMessage() {}

func (x *GetEstimatedJobsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_eta_v1beta1_eta_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEstimatedJobsResponse.ProtoReflect.Descriptor instead.
func (*GetEstimatedJobsResponse) Descriptor() ([]byte, []int) {
	return file_api_eta_v1beta1_eta_proto_rawDescGZIP(), []int{4}
}

func (x *GetEstimatedJobsResponse) GetJobs() []*EstimatedJob {
	if x != nil {
		return x.Jobs
	}
	return nil
}

var File_api_eta_v1beta1_eta_proto protoreflect.FileDescriptor

var file_api_eta_v1beta1_eta_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x74, 0x61, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2f, 0x65, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x61, 0x70, 0x69,
	0x2e, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe4, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x45, 0x73, 0x74,
	0x69, 0x6d, 0x61, 0x74, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x22, 0x0a, 0x08, 0x61, 0x72, 0x65, 0x61, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x07, 0x61, 0x72,
	0x65, 0x61, 0x4b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0,
	0x01, 0x01, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x0d,
	0x72, 0x61, 0x64, 0x69, 0x75, 0x73, 0x5f, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x0c, 0x72, 0x61,
	0x64, 0x69, 0x75, 0x73, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x1d, 0x0a, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02,
	0x18, 0x19, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x31, 0x0a, 0x03, 0x65, 0x74, 0x61,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x74, 0x61,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x03, 0x65, 0x74, 0x61, 0x22, 0x48, 0x0a, 0x08,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x61, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6e,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x6e, 0x67, 0x12, 0x18, 0x0a, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x6a, 0x0a, 0x08, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61,
	0x74, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x64, 0x69, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x35, 0x0a, 0x08, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x9c, 0x03, 0x0a, 0x0c, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x64,
	0x4a, 0x6f, 0x62, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x54, 0x0a, 0x19, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x74, 0x6f,
	0x5f, 0x70, 0x69, 0x63, 0x6b, 0x75, 0x70, 0x5f, 0x65, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x74, 0x61,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74,
	0x65, 0x52, 0x16, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x54, 0x6f, 0x50, 0x69, 0x63, 0x6b, 0x75,
	0x70, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x12, 0x57, 0x0a, 0x1b, 0x70, 0x69, 0x63,
	0x6b, 0x75, 0x70, 0x5f, 0x74, 0x6f, 0x5f, 0x64, 0x72, 0x6f, 0x70, 0x5f, 0x6f, 0x66, 0x66, 0x5f,
	0x65, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2e, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x52, 0x17, 0x70, 0x69, 0x63, 0x6b, 0x75,
	0x70, 0x54, 0x6f, 0x44, 0x72, 0x6f, 0x70, 0x4f, 0x66, 0x66, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61,
	0x74, 0x65, 0x12, 0x42, 0x0a, 0x0f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x42, 0x0a, 0x0f, 0x70, 0x69, 0x63, 0x6b, 0x75, 0x70,
	0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x70, 0x69, 0x63, 0x6b,
	0x75, 0x70, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x45, 0x0a, 0x11, 0x64, 0x72,
	0x6f, 0x70, 0x5f, 0x6f, 0x66, 0x66, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x74, 0x61, 0x2e,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0f, 0x64, 0x72, 0x6f, 0x70, 0x4f, 0x66, 0x66, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x4d, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65,
	0x64, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a,
	0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x45, 0x73,
	0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73,
	0x2a, 0x56, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a,
	0x0a, 0x16, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x52, 0x49, 0x56, 0x45, 0x52, 0x5f, 0x54, 0x4f,
	0x5f, 0x50, 0x49, 0x43, 0x4b, 0x5f, 0x55, 0x50, 0x10, 0x01, 0x12, 0x1c, 0x0a, 0x18, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x50, 0x49, 0x43, 0x4b, 0x5f, 0x55, 0x50, 0x5f, 0x54, 0x4f, 0x5f, 0x44, 0x52,
	0x4f, 0x50, 0x5f, 0x4f, 0x46, 0x46, 0x10, 0x02, 0x32, 0x80, 0x01, 0x0a, 0x13, 0x45, 0x73, 0x74,
	0x69, 0x6d, 0x61, 0x74, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x69, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x64,
	0x4a, 0x6f, 0x62, 0x73, 0x12, 0x28, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x74, 0x61, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61,
	0x74, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x45, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x64, 0x4a, 0x6f, 0x62,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x55, 0x5a, 0x53, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x6d, 0x61,
	0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f,
	0x6f, 0x70, 0x65, 0x6e, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69,
	0x64, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x74, 0x61, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_eta_v1beta1_eta_proto_rawDescOnce sync.Once
	file_api_eta_v1beta1_eta_proto_rawDescData = file_api_eta_v1beta1_eta_proto_rawDesc
)

func file_api_eta_v1beta1_eta_proto_rawDescGZIP() []byte {
	file_api_eta_v1beta1_eta_proto_rawDescOnce.Do(func() {
		file_api_eta_v1beta1_eta_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_eta_v1beta1_eta_proto_rawDescData)
	})
	return file_api_eta_v1beta1_eta_proto_rawDescData
}

var file_api_eta_v1beta1_eta_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_eta_v1beta1_eta_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_eta_v1beta1_eta_proto_goTypes = []interface{}{
	(Type)(0),                        // 0: api.eta.v1beta1.Type
	(*GetEstimatedJobsRequest)(nil),  // 1: api.eta.v1beta1.GetEstimatedJobsRequest
	(*Location)(nil),                 // 2: api.eta.v1beta1.Location
	(*Estimate)(nil),                 // 3: api.eta.v1beta1.Estimate
	(*EstimatedJob)(nil),             // 4: api.eta.v1beta1.EstimatedJob
	(*GetEstimatedJobsResponse)(nil), // 5: api.eta.v1beta1.GetEstimatedJobsResponse
	(*durationpb.Duration)(nil),      // 6: google.protobuf.Duration
}
var file_api_eta_v1beta1_eta_proto_depIdxs = []int32{
	0, // 0: api.eta.v1beta1.GetEstimatedJobsRequest.eta:type_name -> api.eta.v1beta1.Type
	6, // 1: api.eta.v1beta1.Estimate.duration:type_name -> google.protobuf.Duration
	3, // 2: api.eta.v1beta1.EstimatedJob.worker_to_pickup_estimate:type_name -> api.eta.v1beta1.Estimate
	3, // 3: api.eta.v1beta1.EstimatedJob.pickup_to_drop_off_estimate:type_name -> api.eta.v1beta1.Estimate
	2, // 4: api.eta.v1beta1.EstimatedJob.worker_location:type_name -> api.eta.v1beta1.Location
	2, // 5: api.eta.v1beta1.EstimatedJob.pickup_location:type_name -> api.eta.v1beta1.Location
	2, // 6: api.eta.v1beta1.EstimatedJob.drop_off_location:type_name -> api.eta.v1beta1.Location
	4, // 7: api.eta.v1beta1.GetEstimatedJobsResponse.jobs:type_name -> api.eta.v1beta1.EstimatedJob
	1, // 8: api.eta.v1beta1.EstimatedJobService.GetEstimatedJobs:input_type -> api.eta.v1beta1.GetEstimatedJobsRequest
	5, // 9: api.eta.v1beta1.EstimatedJobService.GetEstimatedJobs:output_type -> api.eta.v1beta1.GetEstimatedJobsResponse
	9, // [9:10] is the sub-list for method output_type
	8, // [8:9] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_api_eta_v1beta1_eta_proto_init() }
func file_api_eta_v1beta1_eta_proto_init() {
	if File_api_eta_v1beta1_eta_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_eta_v1beta1_eta_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEstimatedJobsRequest); i {
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
		file_api_eta_v1beta1_eta_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
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
		file_api_eta_v1beta1_eta_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Estimate); i {
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
		file_api_eta_v1beta1_eta_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EstimatedJob); i {
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
		file_api_eta_v1beta1_eta_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEstimatedJobsResponse); i {
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
			RawDescriptor: file_api_eta_v1beta1_eta_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_eta_v1beta1_eta_proto_goTypes,
		DependencyIndexes: file_api_eta_v1beta1_eta_proto_depIdxs,
		EnumInfos:         file_api_eta_v1beta1_eta_proto_enumTypes,
		MessageInfos:      file_api_eta_v1beta1_eta_proto_msgTypes,
	}.Build()
	File_api_eta_v1beta1_eta_proto = out.File
	file_api_eta_v1beta1_eta_proto_rawDesc = nil
	file_api_eta_v1beta1_eta_proto_goTypes = nil
	file_api_eta_v1beta1_eta_proto_depIdxs = nil
}
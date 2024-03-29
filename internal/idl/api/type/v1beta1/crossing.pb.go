// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: api/type/v1beta1/crossing.proto

package v1beta1

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

type Crossing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TollgateId string                 `protobuf:"bytes,2,opt,name=tollgate_id,json=tollgateId,proto3" json:"tollgate_id,omitempty"`
	WorkerId   string                 `protobuf:"bytes,3,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
	Direction  string                 `protobuf:"bytes,4,opt,name=direction,proto3" json:"direction,omitempty"`
	Alg        string                 `protobuf:"bytes,5,opt,name=alg,proto3" json:"alg,omitempty"`
	Movement   *Movement              `protobuf:"bytes,6,opt,name=movement,proto3" json:"movement,omitempty"`
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
}

func (x *Crossing) Reset() {
	*x = Crossing{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_type_v1beta1_crossing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Crossing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Crossing) ProtoMessage() {}

func (x *Crossing) ProtoReflect() protoreflect.Message {
	mi := &file_api_type_v1beta1_crossing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Crossing.ProtoReflect.Descriptor instead.
func (*Crossing) Descriptor() ([]byte, []int) {
	return file_api_type_v1beta1_crossing_proto_rawDescGZIP(), []int{0}
}

func (x *Crossing) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Crossing) GetTollgateId() string {
	if x != nil {
		return x.TollgateId
	}
	return ""
}

func (x *Crossing) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

func (x *Crossing) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

func (x *Crossing) GetAlg() string {
	if x != nil {
		return x.Alg
	}
	return ""
}

func (x *Crossing) GetMovement() *Movement {
	if x != nil {
		return x.Movement
	}
	return nil
}

func (x *Crossing) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

var File_api_type_v1beta1_crossing_proto protoreflect.FileDescriptor

var file_api_type_v1beta1_crossing_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x10, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfd, 0x01, 0x0a, 0x08, 0x43, 0x72, 0x6f, 0x73, 0x73, 0x69,
	0x6e, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x6c, 0x6c, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x6f, 0x6c, 0x6c, 0x67, 0x61, 0x74,
	0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10,
	0x0a, 0x03, 0x61, 0x6c, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x6c, 0x67,
	0x12, 0x36, 0x0a, 0x08, 0x6d, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4d, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08,
	0x6d, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x3b, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x56, 0x5a, 0x54, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x6d, 0x61,
	0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69, 0x64, 0x6c, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_type_v1beta1_crossing_proto_rawDescOnce sync.Once
	file_api_type_v1beta1_crossing_proto_rawDescData = file_api_type_v1beta1_crossing_proto_rawDesc
)

func file_api_type_v1beta1_crossing_proto_rawDescGZIP() []byte {
	file_api_type_v1beta1_crossing_proto_rawDescOnce.Do(func() {
		file_api_type_v1beta1_crossing_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_type_v1beta1_crossing_proto_rawDescData)
	})
	return file_api_type_v1beta1_crossing_proto_rawDescData
}

var file_api_type_v1beta1_crossing_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_type_v1beta1_crossing_proto_goTypes = []interface{}{
	(*Crossing)(nil),              // 0: api.type.v1beta1.Crossing
	(*Movement)(nil),              // 1: api.type.v1beta1.Movement
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_api_type_v1beta1_crossing_proto_depIdxs = []int32{
	1, // 0: api.type.v1beta1.Crossing.movement:type_name -> api.type.v1beta1.Movement
	2, // 1: api.type.v1beta1.Crossing.create_time:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_type_v1beta1_crossing_proto_init() }
func file_api_type_v1beta1_crossing_proto_init() {
	if File_api_type_v1beta1_crossing_proto != nil {
		return
	}
	file_api_type_v1beta1_movement_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_type_v1beta1_crossing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Crossing); i {
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
			RawDescriptor: file_api_type_v1beta1_crossing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_type_v1beta1_crossing_proto_goTypes,
		DependencyIndexes: file_api_type_v1beta1_crossing_proto_depIdxs,
		MessageInfos:      file_api_type_v1beta1_crossing_proto_msgTypes,
	}.Build()
	File_api_type_v1beta1_crossing_proto = out.File
	file_api_type_v1beta1_crossing_proto_rawDesc = nil
	file_api_type_v1beta1_crossing_proto_goTypes = nil
	file_api_type_v1beta1_crossing_proto_depIdxs = nil
}

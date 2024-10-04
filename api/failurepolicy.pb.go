//
//Copyright (c) 2024 Diagrid Inc.
//Licensed under the MIT License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.4
// source: proto/api/failurepolicy.proto

package api

import (
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

// FailurePolicy defines the policy to apply when a job fails to trigger.
type FailurePolicy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// policy is the policy to apply when a job fails to trigger.
	//
	// Types that are assignable to Policy:
	//
	//	*FailurePolicy_Drop
	//	*FailurePolicy_Constant
	Policy isFailurePolicy_Policy `protobuf_oneof:"policy"`
}

func (x *FailurePolicy) Reset() {
	*x = FailurePolicy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_failurepolicy_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailurePolicy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailurePolicy) ProtoMessage() {}

func (x *FailurePolicy) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_failurepolicy_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailurePolicy.ProtoReflect.Descriptor instead.
func (*FailurePolicy) Descriptor() ([]byte, []int) {
	return file_proto_api_failurepolicy_proto_rawDescGZIP(), []int{0}
}

func (m *FailurePolicy) GetPolicy() isFailurePolicy_Policy {
	if m != nil {
		return m.Policy
	}
	return nil
}

func (x *FailurePolicy) GetDrop() *FailurePolicyDrop {
	if x, ok := x.GetPolicy().(*FailurePolicy_Drop); ok {
		return x.Drop
	}
	return nil
}

func (x *FailurePolicy) GetConstant() *FailurePolicyConstant {
	if x, ok := x.GetPolicy().(*FailurePolicy_Constant); ok {
		return x.Constant
	}
	return nil
}

type isFailurePolicy_Policy interface {
	isFailurePolicy_Policy()
}

type FailurePolicy_Drop struct {
	Drop *FailurePolicyDrop `protobuf:"bytes,1,opt,name=drop,proto3,oneof"`
}

type FailurePolicy_Constant struct {
	Constant *FailurePolicyConstant `protobuf:"bytes,2,opt,name=constant,proto3,oneof"`
}

func (*FailurePolicy_Drop) isFailurePolicy_Policy() {}

func (*FailurePolicy_Constant) isFailurePolicy_Policy() {}

// FailurePolicyDrop is a policy which drops the job tick when the job fails to
// trigger.
type FailurePolicyDrop struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FailurePolicyDrop) Reset() {
	*x = FailurePolicyDrop{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_failurepolicy_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailurePolicyDrop) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailurePolicyDrop) ProtoMessage() {}

func (x *FailurePolicyDrop) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_failurepolicy_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailurePolicyDrop.ProtoReflect.Descriptor instead.
func (*FailurePolicyDrop) Descriptor() ([]byte, []int) {
	return file_proto_api_failurepolicy_proto_rawDescGZIP(), []int{1}
}

// FailurePolicyRetry is a policy which retries the job at a consistent
// delay when the job fails to trigger.
type FailurePolicyConstant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// delay is the constant delay to wait before retrying the job.
	Delay *durationpb.Duration `protobuf:"bytes,1,opt,name=delay,proto3" json:"delay,omitempty"`
	// max_retries is the optional maximum number of retries to attempt before
	// giving up.
	// If unset, the Job will be retried indefinitely.
	MaxRetries *uint32 `protobuf:"varint,2,opt,name=max_retries,json=maxRetries,proto3,oneof" json:"max_retries,omitempty"`
}

func (x *FailurePolicyConstant) Reset() {
	*x = FailurePolicyConstant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_failurepolicy_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailurePolicyConstant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailurePolicyConstant) ProtoMessage() {}

func (x *FailurePolicyConstant) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_failurepolicy_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailurePolicyConstant.ProtoReflect.Descriptor instead.
func (*FailurePolicyConstant) Descriptor() ([]byte, []int) {
	return file_proto_api_failurepolicy_proto_rawDescGZIP(), []int{2}
}

func (x *FailurePolicyConstant) GetDelay() *durationpb.Duration {
	if x != nil {
		return x.Delay
	}
	return nil
}

func (x *FailurePolicyConstant) GetMaxRetries() uint32 {
	if x != nil && x.MaxRetries != nil {
		return *x.MaxRetries
	}
	return 0
}

var File_proto_api_failurepolicy_proto protoreflect.FileDescriptor

var file_proto_api_failurepolicy_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x61, 0x69, 0x6c,
	0x75, 0x72, 0x65, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x61, 0x70, 0x69, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x01, 0x0a, 0x0d, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x2c, 0x0a, 0x04, 0x64, 0x72, 0x6f, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x46, 0x61, 0x69, 0x6c, 0x75,
	0x72, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x44, 0x72, 0x6f, 0x70, 0x48, 0x00, 0x52, 0x04,
	0x64, 0x72, 0x6f, 0x70, 0x12, 0x38, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x46, 0x61, 0x69,
	0x6c, 0x75, 0x72, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x74, 0x48, 0x00, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x42, 0x08,
	0x0a, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x22, 0x13, 0x0a, 0x11, 0x46, 0x61, 0x69, 0x6c,
	0x75, 0x72, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x44, 0x72, 0x6f, 0x70, 0x22, 0x7e, 0x0a,
	0x15, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x43, 0x6f,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x05, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x12, 0x24, 0x0a, 0x0b, 0x6d, 0x61, 0x78, 0x5f, 0x72,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x0a,
	0x6d, 0x61, 0x78, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x88, 0x01, 0x01, 0x42, 0x0e, 0x0a,
	0x0c, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x72, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x42, 0x27, 0x5a,
	0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x69, 0x61, 0x67,
	0x72, 0x69, 0x64, 0x69, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x65, 0x74, 0x63, 0x64, 0x2d, 0x63, 0x72,
	0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_api_failurepolicy_proto_rawDescOnce sync.Once
	file_proto_api_failurepolicy_proto_rawDescData = file_proto_api_failurepolicy_proto_rawDesc
)

func file_proto_api_failurepolicy_proto_rawDescGZIP() []byte {
	file_proto_api_failurepolicy_proto_rawDescOnce.Do(func() {
		file_proto_api_failurepolicy_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_api_failurepolicy_proto_rawDescData)
	})
	return file_proto_api_failurepolicy_proto_rawDescData
}

var file_proto_api_failurepolicy_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_api_failurepolicy_proto_goTypes = []interface{}{
	(*FailurePolicy)(nil),         // 0: api.FailurePolicy
	(*FailurePolicyDrop)(nil),     // 1: api.FailurePolicyDrop
	(*FailurePolicyConstant)(nil), // 2: api.FailurePolicyConstant
	(*durationpb.Duration)(nil),   // 3: google.protobuf.Duration
}
var file_proto_api_failurepolicy_proto_depIdxs = []int32{
	1, // 0: api.FailurePolicy.drop:type_name -> api.FailurePolicyDrop
	2, // 1: api.FailurePolicy.constant:type_name -> api.FailurePolicyConstant
	3, // 2: api.FailurePolicyConstant.delay:type_name -> google.protobuf.Duration
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_api_failurepolicy_proto_init() }
func file_proto_api_failurepolicy_proto_init() {
	if File_proto_api_failurepolicy_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_api_failurepolicy_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailurePolicy); i {
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
		file_proto_api_failurepolicy_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailurePolicyDrop); i {
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
		file_proto_api_failurepolicy_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailurePolicyConstant); i {
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
	file_proto_api_failurepolicy_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*FailurePolicy_Drop)(nil),
		(*FailurePolicy_Constant)(nil),
	}
	file_proto_api_failurepolicy_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_api_failurepolicy_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_api_failurepolicy_proto_goTypes,
		DependencyIndexes: file_proto_api_failurepolicy_proto_depIdxs,
		MessageInfos:      file_proto_api_failurepolicy_proto_msgTypes,
	}.Build()
	File_proto_api_failurepolicy_proto = out.File
	file_proto_api_failurepolicy_proto_rawDesc = nil
	file_proto_api_failurepolicy_proto_goTypes = nil
	file_proto_api_failurepolicy_proto_depIdxs = nil
}

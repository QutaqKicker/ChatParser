// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.2
// source: audit.proto

package auditv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AuditType int32

const (
	AuditType_INFO    AuditType = 0
	AuditType_WARNING AuditType = 1
	AuditType_ERROR   AuditType = 2
)

// Enum value maps for AuditType.
var (
	AuditType_name = map[int32]string{
		0: "INFO",
		1: "WARNING",
		2: "ERROR",
	}
	AuditType_value = map[string]int32{
		"INFO":    0,
		"WARNING": 1,
		"ERROR":   2,
	}
)

func (x AuditType) Enum() *AuditType {
	p := new(AuditType)
	*p = x
	return p
}

func (x AuditType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AuditType) Descriptor() protoreflect.EnumDescriptor {
	return file_audit_proto_enumTypes[0].Descriptor()
}

func (AuditType) Type() protoreflect.EnumType {
	return &file_audit_proto_enumTypes[0]
}

func (x AuditType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AuditType.Descriptor instead.
func (AuditType) EnumDescriptor() ([]byte, []int) {
	return file_audit_proto_rawDescGZIP(), []int{0}
}

type AuditInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string    `protobuf:"bytes,1,opt,name=ServiceName,proto3" json:"ServiceName,omitempty"`
	Type        AuditType `protobuf:"varint,2,opt,name=Type,proto3,enum=audit.AuditType" json:"Type,omitempty"`
	Message     string    `protobuf:"bytes,3,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *AuditInfoRequest) Reset() {
	*x = AuditInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuditInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuditInfoRequest) ProtoMessage() {}

func (x *AuditInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_audit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuditInfoRequest.ProtoReflect.Descriptor instead.
func (*AuditInfoRequest) Descriptor() ([]byte, []int) {
	return file_audit_proto_rawDescGZIP(), []int{0}
}

func (x *AuditInfoRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *AuditInfoRequest) GetType() AuditType {
	if x != nil {
		return x.Type
	}
	return AuditType_INFO
}

func (x *AuditInfoRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type AuditInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=Ok,proto3" json:"Ok,omitempty"`
}

func (x *AuditInfoResponse) Reset() {
	*x = AuditInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuditInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuditInfoResponse) ProtoMessage() {}

func (x *AuditInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_audit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuditInfoResponse.ProtoReflect.Descriptor instead.
func (*AuditInfoResponse) Descriptor() ([]byte, []int) {
	return file_audit_proto_rawDescGZIP(), []int{1}
}

func (x *AuditInfoResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_audit_proto protoreflect.FileDescriptor

var file_audit_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61,
	0x75, 0x64, 0x69, 0x74, 0x22, 0x74, 0x0a, 0x10, 0x41, 0x75, 0x64, 0x69, 0x74, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74,
	0x2e, 0x41, 0x75, 0x64, 0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x23, 0x0a, 0x11, 0x41, 0x75,
	0x64, 0x69, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x4f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x6b, 0x2a,
	0x2d, 0x0a, 0x09, 0x41, 0x75, 0x64, 0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04,
	0x49, 0x4e, 0x46, 0x4f, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e,
	0x47, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x02, 0x32, 0x47,
	0x0a, 0x05, 0x41, 0x75, 0x64, 0x69, 0x74, 0x12, 0x3e, 0x0a, 0x09, 0x41, 0x75, 0x64, 0x69, 0x74,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x41, 0x75, 0x64,
	0x69, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x51, 0x75, 0x74, 0x61, 0x71, 0x4b, 0x69, 0x63, 0x6b, 0x65,
	0x72, 0x2f, 0x43, 0x68, 0x61, 0x74, 0x50, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2f, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74,
	0x3b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_audit_proto_rawDescOnce sync.Once
	file_audit_proto_rawDescData = file_audit_proto_rawDesc
)

func file_audit_proto_rawDescGZIP() []byte {
	file_audit_proto_rawDescOnce.Do(func() {
		file_audit_proto_rawDescData = protoimpl.X.CompressGZIP(file_audit_proto_rawDescData)
	})
	return file_audit_proto_rawDescData
}

var file_audit_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_audit_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_audit_proto_goTypes = []interface{}{
	(AuditType)(0),            // 0: audit.AuditType
	(*AuditInfoRequest)(nil),  // 1: audit.AuditInfoRequest
	(*AuditInfoResponse)(nil), // 2: audit.AuditInfoResponse
}
var file_audit_proto_depIdxs = []int32{
	0, // 0: audit.AuditInfoRequest.Type:type_name -> audit.AuditType
	1, // 1: audit.Audit.AuditInfo:input_type -> audit.AuditInfoRequest
	2, // 2: audit.Audit.AuditInfo:output_type -> audit.AuditInfoResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_audit_proto_init() }
func file_audit_proto_init() {
	if File_audit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_audit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuditInfoRequest); i {
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
		file_audit_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuditInfoResponse); i {
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
			RawDescriptor: file_audit_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_audit_proto_goTypes,
		DependencyIndexes: file_audit_proto_depIdxs,
		EnumInfos:         file_audit_proto_enumTypes,
		MessageInfos:      file_audit_proto_msgTypes,
	}.Build()
	File_audit_proto = out.File
	file_audit_proto_rawDesc = nil
	file_audit_proto_goTypes = nil
	file_audit_proto_depIdxs = nil
}

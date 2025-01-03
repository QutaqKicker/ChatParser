// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.1
// source: common.proto

package commonv1

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

type MessagesFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             int32                  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	MinCreatedDate *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=MinCreatedDate,proto3" json:"MinCreatedDate,omitempty"`
	MaxCreatedDate *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=MaxCreatedDate,proto3" json:"MaxCreatedDate,omitempty"`
	SubText        string                 `protobuf:"bytes,4,opt,name=SubText,proto3" json:"SubText,omitempty"`
	UserId         string                 `protobuf:"bytes,5,opt,name=UserId,proto3" json:"UserId,omitempty"`
	UserIds        []string               `protobuf:"bytes,6,rep,name=UserIds,proto3" json:"UserIds,omitempty"`
	ChatIds        []int32                `protobuf:"varint,7,rep,packed,name=ChatIds,proto3" json:"ChatIds,omitempty"`
}

func (x *MessagesFilter) Reset() {
	*x = MessagesFilter{}
	mi := &file_common_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessagesFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessagesFilter) ProtoMessage() {}

func (x *MessagesFilter) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessagesFilter.ProtoReflect.Descriptor instead.
func (*MessagesFilter) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

func (x *MessagesFilter) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MessagesFilter) GetMinCreatedDate() *timestamppb.Timestamp {
	if x != nil {
		return x.MinCreatedDate
	}
	return nil
}

func (x *MessagesFilter) GetMaxCreatedDate() *timestamppb.Timestamp {
	if x != nil {
		return x.MaxCreatedDate
	}
	return nil
}

func (x *MessagesFilter) GetSubText() string {
	if x != nil {
		return x.SubText
	}
	return ""
}

func (x *MessagesFilter) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MessagesFilter) GetUserIds() []string {
	if x != nil {
		return x.UserIds
	}
	return nil
}

func (x *MessagesFilter) GetChatIds() []int32 {
	if x != nil {
		return x.ChatIds
	}
	return nil
}

var File_common_proto protoreflect.FileDescriptor

var file_common_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8e, 0x02, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x64, 0x12, 0x42, 0x0a, 0x0e, 0x4d, 0x69,
	0x6e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e,
	0x4d, 0x69, 0x6e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x12, 0x42,
	0x0a, 0x0e, 0x4d, 0x61, 0x78, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0e, 0x4d, 0x61, 0x78, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61,
	0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x62, 0x54, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x75, 0x62, 0x54, 0x65, 0x78, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x05, 0x52,
	0x07, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x73, 0x42, 0x1f, 0x5a, 0x1d, 0x63, 0x68, 0x61, 0x74,
	0x50, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_common_proto_rawDescOnce sync.Once
	file_common_proto_rawDescData = file_common_proto_rawDesc
)

func file_common_proto_rawDescGZIP() []byte {
	file_common_proto_rawDescOnce.Do(func() {
		file_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_rawDescData)
	})
	return file_common_proto_rawDescData
}

var file_common_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_proto_goTypes = []any{
	(*MessagesFilter)(nil),        // 0: common.MessagesFilter
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_common_proto_depIdxs = []int32{
	1, // 0: common.MessagesFilter.MinCreatedDate:type_name -> google.protobuf.Timestamp
	1, // 1: common.MessagesFilter.MaxCreatedDate:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_common_proto_init() }
func file_common_proto_init() {
	if File_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_proto_goTypes,
		DependencyIndexes: file_common_proto_depIdxs,
		MessageInfos:      file_common_proto_msgTypes,
	}.Build()
	File_common_proto = out.File
	file_common_proto_rawDesc = nil
	file_common_proto_goTypes = nil
	file_common_proto_depIdxs = nil
}

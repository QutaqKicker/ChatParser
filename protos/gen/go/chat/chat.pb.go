// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.1
// source: chat.proto

package chatv1

import (
	chatParser_common_v1 "chatParser.common.v1"
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

type SearchAction int32

const (
	SearchAction_None   SearchAction = 0
	SearchAction_Select SearchAction = 1
	SearchAction_Delete SearchAction = 2
	SearchAction_Backup SearchAction = 3
)

// Enum value maps for SearchAction.
var (
	SearchAction_name = map[int32]string{
		0: "None",
		1: "Select",
		2: "Delete",
		3: "Backup",
	}
	SearchAction_value = map[string]int32{
		"None":   0,
		"Select": 1,
		"Delete": 2,
		"Backup": 3,
	}
)

func (x SearchAction) Enum() *SearchAction {
	p := new(SearchAction)
	*p = x
	return p
}

func (x SearchAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SearchAction) Descriptor() protoreflect.EnumDescriptor {
	return file_chat_proto_enumTypes[0].Descriptor()
}

func (SearchAction) Type() protoreflect.EnumType {
	return &file_chat_proto_enumTypes[0]
}

func (x SearchAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SearchAction.Descriptor instead.
func (SearchAction) EnumDescriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

type ParseFromDirRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DirPath string `protobuf:"bytes,1,opt,name=DirPath,proto3" json:"DirPath,omitempty"`
}

func (x *ParseFromDirRequest) Reset() {
	*x = ParseFromDirRequest{}
	mi := &file_chat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ParseFromDirRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseFromDirRequest) ProtoMessage() {}

func (x *ParseFromDirRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseFromDirRequest.ProtoReflect.Descriptor instead.
func (*ParseFromDirRequest) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

func (x *ParseFromDirRequest) GetDirPath() string {
	if x != nil {
		return x.DirPath
	}
	return ""
}

type ParseFromDirResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=Ok,proto3" json:"Ok,omitempty"`
}

func (x *ParseFromDirResponse) Reset() {
	*x = ParseFromDirResponse{}
	mi := &file_chat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ParseFromDirResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseFromDirResponse) ProtoMessage() {}

func (x *ParseFromDirResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseFromDirResponse.ProtoReflect.Descriptor instead.
func (*ParseFromDirResponse) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{1}
}

func (x *ParseFromDirResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type SearchMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action SearchAction                         `protobuf:"varint,1,opt,name=Action,proto3,enum=chat.SearchAction" json:"Action,omitempty"`
	Filter *chatParser_common_v1.MessagesFilter `protobuf:"bytes,2,opt,name=Filter,proto3" json:"Filter,omitempty"`
	Take   int32                                `protobuf:"varint,3,opt,name=Take,proto3" json:"Take,omitempty"`
	Skip   int32                                `protobuf:"varint,4,opt,name=Skip,proto3" json:"Skip,omitempty"`
	Sorts  []string                             `protobuf:"bytes,5,rep,name=Sorts,proto3" json:"Sorts,omitempty"`
}

func (x *SearchMessagesRequest) Reset() {
	*x = SearchMessagesRequest{}
	mi := &file_chat_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchMessagesRequest) ProtoMessage() {}

func (x *SearchMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchMessagesRequest.ProtoReflect.Descriptor instead.
func (*SearchMessagesRequest) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{2}
}

func (x *SearchMessagesRequest) GetAction() SearchAction {
	if x != nil {
		return x.Action
	}
	return SearchAction_None
}

func (x *SearchMessagesRequest) GetFilter() *chatParser_common_v1.MessagesFilter {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *SearchMessagesRequest) GetTake() int32 {
	if x != nil {
		return x.Take
	}
	return 0
}

func (x *SearchMessagesRequest) GetSkip() int32 {
	if x != nil {
		return x.Skip
	}
	return 0
}

func (x *SearchMessagesRequest) GetSorts() []string {
	if x != nil {
		return x.Sorts
	}
	return nil
}

type GetMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*ChatMessage `protobuf:"bytes,1,rep,name=Messages,proto3" json:"Messages,omitempty"`
}

func (x *GetMessagesResponse) Reset() {
	*x = GetMessagesResponse{}
	mi := &file_chat_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessagesResponse) ProtoMessage() {}

func (x *GetMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessagesResponse.ProtoReflect.Descriptor instead.
func (*GetMessagesResponse) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{3}
}

func (x *GetMessagesResponse) GetMessages() []*ChatMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

type ChatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               int32                  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	ChatId           int32                  `protobuf:"varint,2,opt,name=ChatId,proto3" json:"ChatId,omitempty"`
	ChatName         string                 `protobuf:"bytes,3,opt,name=ChatName,proto3" json:"ChatName,omitempty"`
	UserId           string                 `protobuf:"bytes,4,opt,name=UserId,proto3" json:"UserId,omitempty"`
	UserName         string                 `protobuf:"bytes,5,opt,name=UserName,proto3" json:"UserName,omitempty"`
	ReplyToMessageId int32                  `protobuf:"varint,6,opt,name=ReplyToMessageId,proto3" json:"ReplyToMessageId,omitempty"`
	Text             string                 `protobuf:"bytes,7,opt,name=Text,proto3" json:"Text,omitempty"`
	Created          *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=Created,proto3" json:"Created,omitempty"`
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	mi := &file_chat_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{4}
}

func (x *ChatMessage) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ChatMessage) GetChatId() int32 {
	if x != nil {
		return x.ChatId
	}
	return 0
}

func (x *ChatMessage) GetChatName() string {
	if x != nil {
		return x.ChatName
	}
	return ""
}

func (x *ChatMessage) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ChatMessage) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *ChatMessage) GetReplyToMessageId() int32 {
	if x != nil {
		return x.ReplyToMessageId
	}
	return 0
}

func (x *ChatMessage) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *ChatMessage) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

type DeleteMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=Ok,proto3" json:"Ok,omitempty"`
}

func (x *DeleteMessageResponse) Reset() {
	*x = DeleteMessageResponse{}
	mi := &file_chat_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMessageResponse) ProtoMessage() {}

func (x *DeleteMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMessageResponse.ProtoReflect.Descriptor instead.
func (*DeleteMessageResponse) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteMessageResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_chat_proto protoreflect.FileDescriptor

var file_chat_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x68,
	0x61, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x2f, 0x0a, 0x13, 0x50, 0x61, 0x72, 0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x44, 0x69,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x44, 0x69, 0x72, 0x50,
	0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x44, 0x69, 0x72, 0x50, 0x61,
	0x74, 0x68, 0x22, 0x26, 0x0a, 0x14, 0x50, 0x61, 0x72, 0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x44,
	0x69, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x4f, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x6b, 0x22, 0xb1, 0x01, 0x0a, 0x15, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x2e, 0x0a, 0x06, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x06, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x54, 0x61, 0x6b, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x54, 0x61, 0x6b, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x6b, 0x69, 0x70, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x53, 0x6b, 0x69, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x6f, 0x72, 0x74,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x53, 0x6f, 0x72, 0x74, 0x73, 0x22, 0x44,
	0x0a, 0x13, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x43,
	0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x22, 0xfb, 0x01, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x02, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x43, 0x68, 0x61, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x43, 0x68, 0x61, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x10,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x78, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12, 0x34, 0x0a, 0x07,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x22, 0x27, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x4f,
	0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x6b, 0x2a, 0x3c, 0x0a, 0x0c, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x08, 0x0a, 0x04, 0x4e,
	0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x10, 0x02, 0x12, 0x0a, 0x0a,
	0x06, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x10, 0x03, 0x32, 0xe0, 0x01, 0x0a, 0x04, 0x43, 0x68,
	0x61, 0x74, 0x12, 0x45, 0x0a, 0x0c, 0x50, 0x61, 0x72, 0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x44,
	0x69, 0x72, 0x12, 0x19, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x73, 0x65, 0x46,
	0x72, 0x6f, 0x6d, 0x44, 0x69, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e,
	0x63, 0x68, 0x61, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x73, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x44, 0x69,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x47, 0x65, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x4a, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x12, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1b, 0x5a, 0x19,
	0x63, 0x68, 0x61, 0x74, 0x50, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x76, 0x31, 0x3b, 0x63, 0x68, 0x61, 0x74, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_chat_proto_rawDescOnce sync.Once
	file_chat_proto_rawDescData = file_chat_proto_rawDesc
)

func file_chat_proto_rawDescGZIP() []byte {
	file_chat_proto_rawDescOnce.Do(func() {
		file_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_proto_rawDescData)
	})
	return file_chat_proto_rawDescData
}

var file_chat_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_chat_proto_goTypes = []any{
	(SearchAction)(0),                           // 0: chat.SearchAction
	(*ParseFromDirRequest)(nil),                 // 1: chat.ParseFromDirRequest
	(*ParseFromDirResponse)(nil),                // 2: chat.ParseFromDirResponse
	(*SearchMessagesRequest)(nil),               // 3: chat.SearchMessagesRequest
	(*GetMessagesResponse)(nil),                 // 4: chat.GetMessagesResponse
	(*ChatMessage)(nil),                         // 5: chat.ChatMessage
	(*DeleteMessageResponse)(nil),               // 6: chat.DeleteMessageResponse
	(*chatParser_common_v1.MessagesFilter)(nil), // 7: common.MessagesFilter
	(*timestamppb.Timestamp)(nil),               // 8: google.protobuf.Timestamp
}
var file_chat_proto_depIdxs = []int32{
	0, // 0: chat.SearchMessagesRequest.Action:type_name -> chat.SearchAction
	7, // 1: chat.SearchMessagesRequest.Filter:type_name -> common.MessagesFilter
	5, // 2: chat.GetMessagesResponse.Messages:type_name -> chat.ChatMessage
	8, // 3: chat.ChatMessage.Created:type_name -> google.protobuf.Timestamp
	1, // 4: chat.Chat.ParseFromDir:input_type -> chat.ParseFromDirRequest
	3, // 5: chat.Chat.GetMessages:input_type -> chat.SearchMessagesRequest
	3, // 6: chat.Chat.DeleteMessages:input_type -> chat.SearchMessagesRequest
	2, // 7: chat.Chat.ParseFromDir:output_type -> chat.ParseFromDirResponse
	4, // 8: chat.Chat.GetMessages:output_type -> chat.GetMessagesResponse
	6, // 9: chat.Chat.DeleteMessages:output_type -> chat.DeleteMessageResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_chat_proto_init() }
func file_chat_proto_init() {
	if File_chat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chat_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_proto_goTypes,
		DependencyIndexes: file_chat_proto_depIdxs,
		EnumInfos:         file_chat_proto_enumTypes,
		MessageInfos:      file_chat_proto_msgTypes,
	}.Build()
	File_chat_proto = out.File
	file_chat_proto_rawDesc = nil
	file_chat_proto_goTypes = nil
	file_chat_proto_depIdxs = nil
}

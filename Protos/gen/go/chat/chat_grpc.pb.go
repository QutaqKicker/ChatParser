// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.2
// source: chat.proto

package chatv1

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

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	ParseFromDir(ctx context.Context, in *ParseFromDirRequest, opts ...grpc.CallOption) (*ParseFromDirResponse, error)
	GetMessagesCount(ctx context.Context, in *GetMessagesCountRequest, opts ...grpc.CallOption) (*GetMessagesCountResponse, error)
	SearchMessages(ctx context.Context, in *SearchMessagesRequest, opts ...grpc.CallOption) (*SearchMessagesResponse, error)
	DeleteMessages(ctx context.Context, in *SearchMessagesRequest, opts ...grpc.CallOption) (*DeleteMessageResponse, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) ParseFromDir(ctx context.Context, in *ParseFromDirRequest, opts ...grpc.CallOption) (*ParseFromDirResponse, error) {
	out := new(ParseFromDirResponse)
	err := c.cc.Invoke(ctx, "/chat.Chat/ParseFromDir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) GetMessagesCount(ctx context.Context, in *GetMessagesCountRequest, opts ...grpc.CallOption) (*GetMessagesCountResponse, error) {
	out := new(GetMessagesCountResponse)
	err := c.cc.Invoke(ctx, "/chat.Chat/GetMessagesCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) SearchMessages(ctx context.Context, in *SearchMessagesRequest, opts ...grpc.CallOption) (*SearchMessagesResponse, error) {
	out := new(SearchMessagesResponse)
	err := c.cc.Invoke(ctx, "/chat.Chat/SearchMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) DeleteMessages(ctx context.Context, in *SearchMessagesRequest, opts ...grpc.CallOption) (*DeleteMessageResponse, error) {
	out := new(DeleteMessageResponse)
	err := c.cc.Invoke(ctx, "/chat.Chat/DeleteMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility
type ChatServer interface {
	ParseFromDir(context.Context, *ParseFromDirRequest) (*ParseFromDirResponse, error)
	GetMessagesCount(context.Context, *GetMessagesCountRequest) (*GetMessagesCountResponse, error)
	SearchMessages(context.Context, *SearchMessagesRequest) (*SearchMessagesResponse, error)
	DeleteMessages(context.Context, *SearchMessagesRequest) (*DeleteMessageResponse, error)
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (UnimplementedChatServer) ParseFromDir(context.Context, *ParseFromDirRequest) (*ParseFromDirResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseFromDir not implemented")
}
func (UnimplementedChatServer) GetMessagesCount(context.Context, *GetMessagesCountRequest) (*GetMessagesCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessagesCount not implemented")
}
func (UnimplementedChatServer) SearchMessages(context.Context, *SearchMessagesRequest) (*SearchMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchMessages not implemented")
}
func (UnimplementedChatServer) DeleteMessages(context.Context, *SearchMessagesRequest) (*DeleteMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessages not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_ParseFromDir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseFromDirRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).ParseFromDir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Chat/ParseFromDir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).ParseFromDir(ctx, req.(*ParseFromDirRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_GetMessagesCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMessagesCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetMessagesCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Chat/GetMessagesCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetMessagesCount(ctx, req.(*GetMessagesCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_SearchMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).SearchMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Chat/SearchMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).SearchMessages(ctx, req.(*SearchMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_DeleteMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).DeleteMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Chat/DeleteMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).DeleteMessages(ctx, req.(*SearchMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ParseFromDir",
			Handler:    _Chat_ParseFromDir_Handler,
		},
		{
			MethodName: "GetMessagesCount",
			Handler:    _Chat_GetMessagesCount_Handler,
		},
		{
			MethodName: "SearchMessages",
			Handler:    _Chat_SearchMessages_Handler,
		},
		{
			MethodName: "DeleteMessages",
			Handler:    _Chat_DeleteMessages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat.proto",
}

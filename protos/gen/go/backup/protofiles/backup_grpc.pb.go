// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.2
// source: protofiles/backup.proto

package backupv1

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
	ExportToFile(ctx context.Context, in *ExportToFileRequest, opts ...grpc.CallOption) (*ExportToFileResponse, error)
	ImportFromFile(ctx context.Context, in *ImportFromFileRequest, opts ...grpc.CallOption) (*ImportFromFileResponse, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) ExportToFile(ctx context.Context, in *ExportToFileRequest, opts ...grpc.CallOption) (*ExportToFileResponse, error) {
	out := new(ExportToFileResponse)
	err := c.cc.Invoke(ctx, "/backup.Chat/ExportToFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) ImportFromFile(ctx context.Context, in *ImportFromFileRequest, opts ...grpc.CallOption) (*ImportFromFileResponse, error) {
	out := new(ImportFromFileResponse)
	err := c.cc.Invoke(ctx, "/backup.Chat/ImportFromFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility
type ChatServer interface {
	ExportToFile(context.Context, *ExportToFileRequest) (*ExportToFileResponse, error)
	ImportFromFile(context.Context, *ImportFromFileRequest) (*ImportFromFileResponse, error)
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (UnimplementedChatServer) ExportToFile(context.Context, *ExportToFileRequest) (*ExportToFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportToFile not implemented")
}
func (UnimplementedChatServer) ImportFromFile(context.Context, *ImportFromFileRequest) (*ImportFromFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportFromFile not implemented")
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

func _Chat_ExportToFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportToFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).ExportToFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup.Chat/ExportToFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).ExportToFile(ctx, req.(*ExportToFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_ImportFromFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportFromFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).ImportFromFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup.Chat/ImportFromFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).ImportFromFile(ctx, req.(*ImportFromFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "backup.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExportToFile",
			Handler:    _Chat_ExportToFile_Handler,
		},
		{
			MethodName: "ImportFromFile",
			Handler:    _Chat_ImportFromFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protofiles/backup.proto",
}

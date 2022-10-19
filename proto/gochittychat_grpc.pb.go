// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/gochittychat.proto

package goChittyChat

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

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	Connect(ctx context.Context, in *Post, opts ...grpc.CallOption) (ChatService_ConnectClient, error)
	Disconnect(ctx context.Context, in *Post, opts ...grpc.CallOption) (*Empty, error)
	Messages(ctx context.Context, opts ...grpc.CallOption) (ChatService_MessagesClient, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) Connect(ctx context.Context, in *Post, opts ...grpc.CallOption) (ChatService_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], "/goChittyChat.ChatService/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceConnectClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChatService_ConnectClient interface {
	Recv() (*Post, error)
	grpc.ClientStream
}

type chatServiceConnectClient struct {
	grpc.ClientStream
}

func (x *chatServiceConnectClient) Recv() (*Post, error) {
	m := new(Post)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) Disconnect(ctx context.Context, in *Post, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/goChittyChat.ChatService/Disconnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) Messages(ctx context.Context, opts ...grpc.CallOption) (ChatService_MessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[1], "/goChittyChat.ChatService/Messages", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceMessagesClient{stream}
	return x, nil
}

type ChatService_MessagesClient interface {
	Send(*Post) error
	CloseAndRecv() (*Post, error)
	grpc.ClientStream
}

type chatServiceMessagesClient struct {
	grpc.ClientStream
}

func (x *chatServiceMessagesClient) Send(m *Post) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceMessagesClient) CloseAndRecv() (*Post, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Post)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	Connect(*Post, ChatService_ConnectServer) error
	Disconnect(context.Context, *Post) (*Empty, error)
	Messages(ChatService_MessagesServer) error
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (UnimplementedChatServiceServer) Connect(*Post, ChatService_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedChatServiceServer) Disconnect(context.Context, *Post) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disconnect not implemented")
}
func (UnimplementedChatServiceServer) Messages(ChatService_MessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method Messages not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Post)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServiceServer).Connect(m, &chatServiceConnectServer{stream})
}

type ChatService_ConnectServer interface {
	Send(*Post) error
	grpc.ServerStream
}

type chatServiceConnectServer struct {
	grpc.ServerStream
}

func (x *chatServiceConnectServer) Send(m *Post) error {
	return x.ServerStream.SendMsg(m)
}

func _ChatService_Disconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Post)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).Disconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/goChittyChat.ChatService/Disconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).Disconnect(ctx, req.(*Post))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_Messages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).Messages(&chatServiceMessagesServer{stream})
}

type ChatService_MessagesServer interface {
	SendAndClose(*Post) error
	Recv() (*Post, error)
	grpc.ServerStream
}

type chatServiceMessagesServer struct {
	grpc.ServerStream
}

func (x *chatServiceMessagesServer) SendAndClose(m *Post) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceMessagesServer) Recv() (*Post, error) {
	m := new(Post)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "goChittyChat.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Disconnect",
			Handler:    _ChatService_Disconnect_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _ChatService_Connect_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Messages",
			Handler:       _ChatService_Messages_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proto/gochittychat.proto",
}

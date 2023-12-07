// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: keye.proto

package keye

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

const (
	Api_Get_FullMethodName   = "/keye.Api/Get"
	Api_Put_FullMethodName   = "/keye.Api/Put"
	Api_Del_FullMethodName   = "/keye.Api/Del"
	Api_Watch_FullMethodName = "/keye.Api/Watch"
)

// ApiClient is the client API for Api service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApiClient interface {
	Get(ctx context.Context, in *GetParams, opts ...grpc.CallOption) (*GetResponse, error)
	Put(ctx context.Context, in *PutParams, opts ...grpc.CallOption) (*PutResponse, error)
	Del(ctx context.Context, in *DelParams, opts ...grpc.CallOption) (*DelResponse, error)
	Watch(ctx context.Context, in *WatchParams, opts ...grpc.CallOption) (Api_WatchClient, error)
}

type apiClient struct {
	cc grpc.ClientConnInterface
}

func NewApiClient(cc grpc.ClientConnInterface) ApiClient {
	return &apiClient{cc}
}

func (c *apiClient) Get(ctx context.Context, in *GetParams, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, Api_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Put(ctx context.Context, in *PutParams, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := c.cc.Invoke(ctx, Api_Put_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Del(ctx context.Context, in *DelParams, opts ...grpc.CallOption) (*DelResponse, error) {
	out := new(DelResponse)
	err := c.cc.Invoke(ctx, Api_Del_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Watch(ctx context.Context, in *WatchParams, opts ...grpc.CallOption) (Api_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Api_ServiceDesc.Streams[0], Api_Watch_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &apiWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Api_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type apiWatchClient struct {
	grpc.ClientStream
}

func (x *apiWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ApiServer is the server API for Api service.
// All implementations must embed UnimplementedApiServer
// for forward compatibility
type ApiServer interface {
	Get(context.Context, *GetParams) (*GetResponse, error)
	Put(context.Context, *PutParams) (*PutResponse, error)
	Del(context.Context, *DelParams) (*DelResponse, error)
	Watch(*WatchParams, Api_WatchServer) error
	mustEmbedUnimplementedApiServer()
}

// UnimplementedApiServer must be embedded to have forward compatible implementations.
type UnimplementedApiServer struct {
}

func (UnimplementedApiServer) Get(context.Context, *GetParams) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedApiServer) Put(context.Context, *PutParams) (*PutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedApiServer) Del(context.Context, *DelParams) (*DelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedApiServer) Watch(*WatchParams, Api_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedApiServer) mustEmbedUnimplementedApiServer() {}

// UnsafeApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ApiServer will
// result in compilation errors.
type UnsafeApiServer interface {
	mustEmbedUnimplementedApiServer()
}

func RegisterApiServer(s grpc.ServiceRegistrar, srv ApiServer) {
	s.RegisterService(&Api_ServiceDesc, srv)
}

func _Api_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Api_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Get(ctx, req.(*GetParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Api_Put_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Put(ctx, req.(*PutParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Api_Del_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Del(ctx, req.(*DelParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchParams)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ApiServer).Watch(m, &apiWatchServer{stream})
}

type Api_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type apiWatchServer struct {
	grpc.ServerStream
}

func (x *apiWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Api_ServiceDesc is the grpc.ServiceDesc for Api service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Api_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "keye.Api",
	HandlerType: (*ApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Api_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _Api_Put_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _Api_Del_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _Api_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "keye.proto",
}

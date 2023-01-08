// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/shortener.proto

package proto

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

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	CreateLink(ctx context.Context, in *CreateLinkRequest, opts ...grpc.CallOption) (*CreateLinkResponse, error)
	CreateLinkJSON(ctx context.Context, in *CreateLinkJSONRequest, opts ...grpc.CallOption) (*CreateLinkJSONResponse, error)
	GetLink(ctx context.Context, in *GetLinkRequest, opts ...grpc.CallOption) (*GetLinkResponse, error)
	GetManyLinks(ctx context.Context, in *GetManyLinksRequest, opts ...grpc.CallOption) (*GetManyLinksResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	CreateManyLinks(ctx context.Context, in *CreateManyRequest, opts ...grpc.CallOption) (*CreateManyResponse, error)
	DeleteMany(ctx context.Context, in *DeleteManyRequest, opts ...grpc.CallOption) (*DeleteManyResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) CreateLink(ctx context.Context, in *CreateLinkRequest, opts ...grpc.CallOption) (*CreateLinkResponse, error) {
	out := new(CreateLinkResponse)
	err := c.cc.Invoke(ctx, "/shortener.proto.Shortener/CreateLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) CreateLinkJSON(ctx context.Context, in *CreateLinkJSONRequest, opts ...grpc.CallOption) (*CreateLinkJSONResponse, error) {
	out := new(CreateLinkJSONResponse)
	err := c.cc.Invoke(ctx, "/shortener.proto.Shortener/CreateLinkJSON", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetLink(ctx context.Context, in *GetLinkRequest, opts ...grpc.CallOption) (*GetLinkResponse, error) {
	out := new(GetLinkResponse)
	err := c.cc.Invoke(ctx, "/shortener.proto.Shortener/GetLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetManyLinks(ctx context.Context, in *GetManyLinksRequest, opts ...grpc.CallOption) (*GetManyLinksResponse, error) {
	out := new(GetManyLinksResponse)
	err := c.cc.Invoke(ctx, "/shortener.proto.Shortener/GetManyLinks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/shortener.proto.Shortener/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) CreateManyLinks(ctx context.Context, in *CreateManyRequest, opts ...grpc.CallOption) (*CreateManyResponse, error) {
	out := new(CreateManyResponse)
	err := c.cc.Invoke(ctx, "/shortener.proto.Shortener/CreateManyLinks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) DeleteMany(ctx context.Context, in *DeleteManyRequest, opts ...grpc.CallOption) (*DeleteManyResponse, error) {
	out := new(DeleteManyResponse)
	err := c.cc.Invoke(ctx, "/shortener.proto.Shortener/DeleteMany", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	CreateLink(context.Context, *CreateLinkRequest) (*CreateLinkResponse, error)
	CreateLinkJSON(context.Context, *CreateLinkJSONRequest) (*CreateLinkJSONResponse, error)
	GetLink(context.Context, *GetLinkRequest) (*GetLinkResponse, error)
	GetManyLinks(context.Context, *GetManyLinksRequest) (*GetManyLinksResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	CreateManyLinks(context.Context, *CreateManyRequest) (*CreateManyResponse, error)
	DeleteMany(context.Context, *DeleteManyRequest) (*DeleteManyResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) CreateLink(context.Context, *CreateLinkRequest) (*CreateLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLink not implemented")
}
func (UnimplementedShortenerServer) CreateLinkJSON(context.Context, *CreateLinkJSONRequest) (*CreateLinkJSONResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLinkJSON not implemented")
}
func (UnimplementedShortenerServer) GetLink(context.Context, *GetLinkRequest) (*GetLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLink not implemented")
}
func (UnimplementedShortenerServer) GetManyLinks(context.Context, *GetManyLinksRequest) (*GetManyLinksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetManyLinks not implemented")
}
func (UnimplementedShortenerServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedShortenerServer) CreateManyLinks(context.Context, *CreateManyRequest) (*CreateManyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateManyLinks not implemented")
}
func (UnimplementedShortenerServer) DeleteMany(context.Context, *DeleteManyRequest) (*DeleteManyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMany not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_CreateLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).CreateLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.proto.Shortener/CreateLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).CreateLink(ctx, req.(*CreateLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_CreateLinkJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLinkJSONRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).CreateLinkJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.proto.Shortener/CreateLinkJSON",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).CreateLinkJSON(ctx, req.(*CreateLinkJSONRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.proto.Shortener/GetLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetLink(ctx, req.(*GetLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetManyLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetManyLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetManyLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.proto.Shortener/GetManyLinks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetManyLinks(ctx, req.(*GetManyLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.proto.Shortener/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_CreateManyLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateManyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).CreateManyLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.proto.Shortener/CreateManyLinks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).CreateManyLinks(ctx, req.(*CreateManyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_DeleteMany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteManyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).DeleteMany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.proto.Shortener/DeleteMany",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).DeleteMany(ctx, req.(*DeleteManyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.proto.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLink",
			Handler:    _Shortener_CreateLink_Handler,
		},
		{
			MethodName: "CreateLinkJSON",
			Handler:    _Shortener_CreateLinkJSON_Handler,
		},
		{
			MethodName: "GetLink",
			Handler:    _Shortener_GetLink_Handler,
		},
		{
			MethodName: "GetManyLinks",
			Handler:    _Shortener_GetManyLinks_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Shortener_Ping_Handler,
		},
		{
			MethodName: "CreateManyLinks",
			Handler:    _Shortener_CreateManyLinks_Handler,
		},
		{
			MethodName: "DeleteMany",
			Handler:    _Shortener_DeleteMany_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}

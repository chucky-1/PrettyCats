// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

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

// CatsCrudClient is the client API for CatsCrud service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CatsCrudClient interface {
	GetAllCats(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	CreateCats(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	GetCat(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	UpdateCat(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	DeleteCat(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type catsCrudClient struct {
	cc grpc.ClientConnInterface
}

func NewCatsCrudClient(cc grpc.ClientConnInterface) CatsCrudClient {
	return &catsCrudClient{cc}
}

func (c *catsCrudClient) GetAllCats(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.CatsCrud/GetAllCats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catsCrudClient) CreateCats(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.CatsCrud/CreateCats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catsCrudClient) GetCat(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.CatsCrud/GetCat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catsCrudClient) UpdateCat(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.CatsCrud/UpdateCat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catsCrudClient) DeleteCat(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.CatsCrud/DeleteCat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CatsCrudServer is the server API for CatsCrud service.
// All implementations must embed UnimplementedCatsCrudServer
// for forward compatibility
type CatsCrudServer interface {
	GetAllCats(context.Context, *Request) (*Response, error)
	CreateCats(context.Context, *Request) (*Response, error)
	GetCat(context.Context, *Request) (*Response, error)
	UpdateCat(context.Context, *Request) (*Response, error)
	DeleteCat(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedCatsCrudServer()
}

// UnimplementedCatsCrudServer must be embedded to have forward compatible implementations.
type UnimplementedCatsCrudServer struct {
}

func (UnimplementedCatsCrudServer) GetAllCats(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllCats not implemented")
}
func (UnimplementedCatsCrudServer) CreateCats(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCats not implemented")
}
func (UnimplementedCatsCrudServer) GetCat(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCat not implemented")
}
func (UnimplementedCatsCrudServer) UpdateCat(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCat not implemented")
}
func (UnimplementedCatsCrudServer) DeleteCat(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCat not implemented")
}
func (UnimplementedCatsCrudServer) mustEmbedUnimplementedCatsCrudServer() {}

// UnsafeCatsCrudServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CatsCrudServer will
// result in compilation errors.
type UnsafeCatsCrudServer interface {
	mustEmbedUnimplementedCatsCrudServer()
}

func RegisterCatsCrudServer(s grpc.ServiceRegistrar, srv CatsCrudServer) {
	s.RegisterService(&CatsCrud_ServiceDesc, srv)
}

func _CatsCrud_GetAllCats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatsCrudServer).GetAllCats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.CatsCrud/GetAllCats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatsCrudServer).GetAllCats(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatsCrud_CreateCats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatsCrudServer).CreateCats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.CatsCrud/CreateCats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatsCrudServer).CreateCats(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatsCrud_GetCat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatsCrudServer).GetCat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.CatsCrud/GetCat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatsCrudServer).GetCat(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatsCrud_UpdateCat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatsCrudServer).UpdateCat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.CatsCrud/UpdateCat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatsCrudServer).UpdateCat(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatsCrud_DeleteCat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatsCrudServer).DeleteCat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.CatsCrud/DeleteCat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatsCrudServer).DeleteCat(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// CatsCrud_ServiceDesc is the grpc.ServiceDesc for CatsCrud service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CatsCrud_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.CatsCrud",
	HandlerType: (*CatsCrudServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllCats",
			Handler:    _CatsCrud_GetAllCats_Handler,
		},
		{
			MethodName: "CreateCats",
			Handler:    _CatsCrud_CreateCats_Handler,
		},
		{
			MethodName: "GetCat",
			Handler:    _CatsCrud_GetCat_Handler,
		},
		{
			MethodName: "UpdateCat",
			Handler:    _CatsCrud_UpdateCat_Handler,
		},
		{
			MethodName: "DeleteCat",
			Handler:    _CatsCrud_DeleteCat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/catscrud.proto",
}

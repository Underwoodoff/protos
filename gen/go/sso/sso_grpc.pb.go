// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.6.1
// source: sso/sso.proto

package ssov1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Ayth_Register_FullMethodName = "/auth.Ayth/Register"
	Ayth_Login_FullMethodName    = "/auth.Ayth/Login"
	Ayth_IsAdmin_FullMethodName  = "/auth.Ayth/IsAdmin"
)

// AythClient is the client API for Ayth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AythClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	IsAdmin(ctx context.Context, in *IsAdminRequest, opts ...grpc.CallOption) (*IsAdminResponse, error)
}

type aythClient struct {
	cc grpc.ClientConnInterface
}

func NewAythClient(cc grpc.ClientConnInterface) AythClient {
	return &aythClient{cc}
}

func (c *aythClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, Ayth_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aythClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, Ayth_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aythClient) IsAdmin(ctx context.Context, in *IsAdminRequest, opts ...grpc.CallOption) (*IsAdminResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IsAdminResponse)
	err := c.cc.Invoke(ctx, Ayth_IsAdmin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AythServer is the server API for Ayth service.
// All implementations must embed UnimplementedAythServer
// for forward compatibility.
type AythServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	IsAdmin(context.Context, *IsAdminRequest) (*IsAdminResponse, error)
	mustEmbedUnimplementedAythServer()
}

// UnimplementedAythServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAythServer struct{}

func (UnimplementedAythServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAythServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAythServer) IsAdmin(context.Context, *IsAdminRequest) (*IsAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAdmin not implemented")
}
func (UnimplementedAythServer) mustEmbedUnimplementedAythServer() {}
func (UnimplementedAythServer) testEmbeddedByValue()              {}

// UnsafeAythServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AythServer will
// result in compilation errors.
type UnsafeAythServer interface {
	mustEmbedUnimplementedAythServer()
}

func RegisterAythServer(s grpc.ServiceRegistrar, srv AythServer) {
	// If the following call pancis, it indicates UnimplementedAythServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Ayth_ServiceDesc, srv)
}

func _Ayth_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AythServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ayth_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AythServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ayth_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AythServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ayth_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AythServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ayth_IsAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AythServer).IsAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ayth_IsAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AythServer).IsAdmin(ctx, req.(*IsAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Ayth_ServiceDesc is the grpc.ServiceDesc for Ayth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ayth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.Ayth",
	HandlerType: (*AythServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Ayth_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Ayth_Login_Handler,
		},
		{
			MethodName: "IsAdmin",
			Handler:    _Ayth_IsAdmin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sso/sso.proto",
}

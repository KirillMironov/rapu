// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// UsersClient is the client API for Users service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersClient interface {
	SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*Response, error)
	SignIn(ctx context.Context, in *SignInRequest, opts ...grpc.CallOption) (*Response, error)
}

type usersClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersClient(cc grpc.ClientConnInterface) UsersClient {
	return &usersClient{cc}
}

func (c *usersClient) SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Users/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) SignIn(ctx context.Context, in *SignInRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Users/SignIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServer is the server API for Users service.
// All implementations must embed UnimplementedUsersServer
// for forward compatibility
type UsersServer interface {
	SignUp(context.Context, *SignUpRequest) (*Response, error)
	SignIn(context.Context, *SignInRequest) (*Response, error)
	mustEmbedUnimplementedUsersServer()
}

// UnimplementedUsersServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServer struct {
}

func (UnimplementedUsersServer) SignUp(context.Context, *SignUpRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedUsersServer) SignIn(context.Context, *SignInRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignIn not implemented")
}
func (UnimplementedUsersServer) mustEmbedUnimplementedUsersServer() {}

// UnsafeUsersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServer will
// result in compilation errors.
type UnsafeUsersServer interface {
	mustEmbedUnimplementedUsersServer()
}

func RegisterUsersServer(s grpc.ServiceRegistrar, srv UsersServer) {
	s.RegisterService(&Users_ServiceDesc, srv)
}

func _Users_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Users/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).SignUp(ctx, req.(*SignUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_SignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).SignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Users/SignIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).SignIn(ctx, req.(*SignInRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Users_ServiceDesc is the grpc.ServiceDesc for Users service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Users_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Users",
	HandlerType: (*UsersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignUp",
			Handler:    _Users_SignUp_Handler,
		},
		{
			MethodName: "SignIn",
			Handler:    _Users_SignIn_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/delivery/proto/users.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: service_simple_bank.proto

package pb

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
	SimpleBankService_CreateUser_FullMethodName = "/pb.SimpleBankService/CreateUser"
	SimpleBankService_LoginUser_FullMethodName  = "/pb.SimpleBankService/LoginUser"
)

// SimpleBankServiceClient is the client API for SimpleBankService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SimpleBankServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error)
}

type simpleBankServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSimpleBankServiceClient(cc grpc.ClientConnInterface) SimpleBankServiceClient {
	return &simpleBankServiceClient{cc}
}

func (c *simpleBankServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, SimpleBankService_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simpleBankServiceClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	out := new(LoginUserResponse)
	err := c.cc.Invoke(ctx, SimpleBankService_LoginUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimpleBankServiceServer is the server API for SimpleBankService service.
// All implementations must embed UnimplementedSimpleBankServiceServer
// for forward compatibility
type SimpleBankServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error)
	mustEmbedUnimplementedSimpleBankServiceServer()
}

// UnimplementedSimpleBankServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSimpleBankServiceServer struct {
}

func (UnimplementedSimpleBankServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedSimpleBankServiceServer) LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedSimpleBankServiceServer) mustEmbedUnimplementedSimpleBankServiceServer() {}

// UnsafeSimpleBankServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SimpleBankServiceServer will
// result in compilation errors.
type UnsafeSimpleBankServiceServer interface {
	mustEmbedUnimplementedSimpleBankServiceServer()
}

func RegisterSimpleBankServiceServer(s grpc.ServiceRegistrar, srv SimpleBankServiceServer) {
	s.RegisterService(&SimpleBankService_ServiceDesc, srv)
}

func _SimpleBankService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleBankServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SimpleBankService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleBankServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimpleBankService_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleBankServiceServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SimpleBankService_LoginUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleBankServiceServer).LoginUser(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SimpleBankService_ServiceDesc is the grpc.ServiceDesc for SimpleBankService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SimpleBankService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.SimpleBankService",
	HandlerType: (*SimpleBankServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _SimpleBankService_CreateUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _SimpleBankService_LoginUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_simple_bank.proto",
}

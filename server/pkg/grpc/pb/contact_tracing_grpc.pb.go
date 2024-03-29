// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: pkg/grpc/proto/contact_tracing.proto

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

// ContactTracingClient is the client API for ContactTracing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ContactTracingClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResult, error)
	ReportInfection(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResult, error)
}

type contactTracingClient struct {
	cc grpc.ClientConnInterface
}

func NewContactTracingClient(cc grpc.ClientConnInterface) ContactTracingClient {
	return &contactTracingClient{cc}
}

func (c *contactTracingClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResult, error) {
	out := new(RegisterResult)
	err := c.cc.Invoke(ctx, "/pb.ContactTracing/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contactTracingClient) ReportInfection(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResult, error) {
	out := new(ReportResult)
	err := c.cc.Invoke(ctx, "/pb.ContactTracing/ReportInfection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContactTracingServer is the server API for ContactTracing service.
// All implementations must embed UnimplementedContactTracingServer
// for forward compatibility
type ContactTracingServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResult, error)
	ReportInfection(context.Context, *ReportRequest) (*ReportResult, error)
	mustEmbedUnimplementedContactTracingServer()
}

// UnimplementedContactTracingServer must be embedded to have forward compatible implementations.
type UnimplementedContactTracingServer struct {
}

func (UnimplementedContactTracingServer) Register(context.Context, *RegisterRequest) (*RegisterResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedContactTracingServer) ReportInfection(context.Context, *ReportRequest) (*ReportResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportInfection not implemented")
}
func (UnimplementedContactTracingServer) mustEmbedUnimplementedContactTracingServer() {}

// UnsafeContactTracingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContactTracingServer will
// result in compilation errors.
type UnsafeContactTracingServer interface {
	mustEmbedUnimplementedContactTracingServer()
}

func RegisterContactTracingServer(s grpc.ServiceRegistrar, srv ContactTracingServer) {
	s.RegisterService(&ContactTracing_ServiceDesc, srv)
}

func _ContactTracing_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContactTracingServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ContactTracing/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContactTracingServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContactTracing_ReportInfection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContactTracingServer).ReportInfection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ContactTracing/ReportInfection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContactTracingServer).ReportInfection(ctx, req.(*ReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ContactTracing_ServiceDesc is the grpc.ServiceDesc for ContactTracing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ContactTracing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.ContactTracing",
	HandlerType: (*ContactTracingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _ContactTracing_Register_Handler,
		},
		{
			MethodName: "ReportInfection",
			Handler:    _ContactTracing_ReportInfection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/grpc/proto/contact_tracing.proto",
}

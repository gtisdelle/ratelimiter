// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: proto/ratelimit/v1/ratelimit.proto

package ratelimitv1

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
	RateLimitService_ShouldRateLimit_FullMethodName = "/ratelimit.v1.RateLimitService/ShouldRateLimit"
)

// RateLimitServiceClient is the client API for RateLimitService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RateLimitServiceClient interface {
	ShouldRateLimit(ctx context.Context, in *RateLimitRequest, opts ...grpc.CallOption) (*RateLimitResponse, error)
}

type rateLimitServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRateLimitServiceClient(cc grpc.ClientConnInterface) RateLimitServiceClient {
	return &rateLimitServiceClient{cc}
}

func (c *rateLimitServiceClient) ShouldRateLimit(ctx context.Context, in *RateLimitRequest, opts ...grpc.CallOption) (*RateLimitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RateLimitResponse)
	err := c.cc.Invoke(ctx, RateLimitService_ShouldRateLimit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RateLimitServiceServer is the server API for RateLimitService service.
// All implementations must embed UnimplementedRateLimitServiceServer
// for forward compatibility.
type RateLimitServiceServer interface {
	ShouldRateLimit(context.Context, *RateLimitRequest) (*RateLimitResponse, error)
	mustEmbedUnimplementedRateLimitServiceServer()
}

// UnimplementedRateLimitServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRateLimitServiceServer struct{}

func (UnimplementedRateLimitServiceServer) ShouldRateLimit(context.Context, *RateLimitRequest) (*RateLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShouldRateLimit not implemented")
}
func (UnimplementedRateLimitServiceServer) mustEmbedUnimplementedRateLimitServiceServer() {}
func (UnimplementedRateLimitServiceServer) testEmbeddedByValue()                          {}

// UnsafeRateLimitServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RateLimitServiceServer will
// result in compilation errors.
type UnsafeRateLimitServiceServer interface {
	mustEmbedUnimplementedRateLimitServiceServer()
}

func RegisterRateLimitServiceServer(s grpc.ServiceRegistrar, srv RateLimitServiceServer) {
	// If the following call pancis, it indicates UnimplementedRateLimitServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RateLimitService_ServiceDesc, srv)
}

func _RateLimitService_ShouldRateLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimitServiceServer).ShouldRateLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimitService_ShouldRateLimit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimitServiceServer).ShouldRateLimit(ctx, req.(*RateLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RateLimitService_ServiceDesc is the grpc.ServiceDesc for RateLimitService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RateLimitService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ratelimit.v1.RateLimitService",
	HandlerType: (*RateLimitServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShouldRateLimit",
			Handler:    _RateLimitService_ShouldRateLimit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ratelimit/v1/ratelimit.proto",
}

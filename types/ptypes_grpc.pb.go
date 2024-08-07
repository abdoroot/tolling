// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: types/ptypes.proto

package types

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
	Aggreagator_AggregateDistance_FullMethodName = "/Aggreagator/AggregateDistance"
)

// AggreagatorClient is the client API for Aggreagator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AggreagatorClient interface {
	AggregateDistance(ctx context.Context, in *DistanceRequest, opts ...grpc.CallOption) (*None, error)
}

type aggreagatorClient struct {
	cc grpc.ClientConnInterface
}

func NewAggreagatorClient(cc grpc.ClientConnInterface) AggreagatorClient {
	return &aggreagatorClient{cc}
}

func (c *aggreagatorClient) AggregateDistance(ctx context.Context, in *DistanceRequest, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, Aggreagator_AggregateDistance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AggreagatorServer is the server API for Aggreagator service.
// All implementations must embed UnimplementedAggreagatorServer
// for forward compatibility
type AggreagatorServer interface {
	AggregateDistance(context.Context, *DistanceRequest) (*None, error)
	mustEmbedUnimplementedAggreagatorServer()
}

// UnimplementedAggreagatorServer must be embedded to have forward compatible implementations.
type UnimplementedAggreagatorServer struct {
}

func (UnimplementedAggreagatorServer) AggregateDistance(context.Context, *DistanceRequest) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AggregateDistance not implemented")
}
func (UnimplementedAggreagatorServer) mustEmbedUnimplementedAggreagatorServer() {}

// UnsafeAggreagatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AggreagatorServer will
// result in compilation errors.
type UnsafeAggreagatorServer interface {
	mustEmbedUnimplementedAggreagatorServer()
}

func RegisterAggreagatorServer(s grpc.ServiceRegistrar, srv AggreagatorServer) {
	s.RegisterService(&Aggreagator_ServiceDesc, srv)
}

func _Aggreagator_AggregateDistance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DistanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AggreagatorServer).AggregateDistance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Aggreagator_AggregateDistance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AggreagatorServer).AggregateDistance(ctx, req.(*DistanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Aggreagator_ServiceDesc is the grpc.ServiceDesc for Aggreagator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Aggreagator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Aggreagator",
	HandlerType: (*AggreagatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AggregateDistance",
			Handler:    _Aggreagator_AggregateDistance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "types/ptypes.proto",
}

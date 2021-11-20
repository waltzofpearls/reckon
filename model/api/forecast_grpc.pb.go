// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// ForecastClient is the client API for Forecast service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ForecastClient interface {
	Prophet(ctx context.Context, in *ProphetRequest, opts ...grpc.CallOption) (*ProphetReply, error)
}

type forecastClient struct {
	cc grpc.ClientConnInterface
}

func NewForecastClient(cc grpc.ClientConnInterface) ForecastClient {
	return &forecastClient{cc}
}

func (c *forecastClient) Prophet(ctx context.Context, in *ProphetRequest, opts ...grpc.CallOption) (*ProphetReply, error) {
	out := new(ProphetReply)
	err := c.cc.Invoke(ctx, "/api.Forecast/Prophet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ForecastServer is the server API for Forecast service.
// All implementations must embed UnimplementedForecastServer
// for forward compatibility
type ForecastServer interface {
	Prophet(context.Context, *ProphetRequest) (*ProphetReply, error)
	mustEmbedUnimplementedForecastServer()
}

// UnimplementedForecastServer must be embedded to have forward compatible implementations.
type UnimplementedForecastServer struct {
}

func (UnimplementedForecastServer) Prophet(context.Context, *ProphetRequest) (*ProphetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prophet not implemented")
}
func (UnimplementedForecastServer) mustEmbedUnimplementedForecastServer() {}

// UnsafeForecastServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ForecastServer will
// result in compilation errors.
type UnsafeForecastServer interface {
	mustEmbedUnimplementedForecastServer()
}

func RegisterForecastServer(s grpc.ServiceRegistrar, srv ForecastServer) {
	s.RegisterService(&Forecast_ServiceDesc, srv)
}

func _Forecast_Prophet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProphetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForecastServer).Prophet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Forecast/Prophet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForecastServer).Prophet(ctx, req.(*ProphetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Forecast_ServiceDesc is the grpc.ServiceDesc for Forecast service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Forecast_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Forecast",
	HandlerType: (*ForecastServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Prophet",
			Handler:    _Forecast_Prophet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "model/api/forecast.proto",
}

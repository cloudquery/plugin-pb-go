// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.23.4
// source: plugin-pb/destination/v1/destination.proto

package destination

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
	Destination_GetName_FullMethodName     = "/cloudquery.destination.v1.Destination/GetName"
	Destination_GetVersion_FullMethodName  = "/cloudquery.destination.v1.Destination/GetVersion"
	Destination_Configure_FullMethodName   = "/cloudquery.destination.v1.Destination/Configure"
	Destination_Migrate_FullMethodName     = "/cloudquery.destination.v1.Destination/Migrate"
	Destination_Write_FullMethodName       = "/cloudquery.destination.v1.Destination/Write"
	Destination_Close_FullMethodName       = "/cloudquery.destination.v1.Destination/Close"
	Destination_DeleteStale_FullMethodName = "/cloudquery.destination.v1.Destination/DeleteStale"
	Destination_GetMetrics_FullMethodName  = "/cloudquery.destination.v1.Destination/GetMetrics"
)

// DestinationClient is the client API for Destination service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DestinationClient interface {
	// Get the name of the plugin
	GetName(ctx context.Context, in *GetName_Request, opts ...grpc.CallOption) (*GetName_Response, error)
	// Get the current version of the plugin
	GetVersion(ctx context.Context, in *GetVersion_Request, opts ...grpc.CallOption) (*GetVersion_Response, error)
	// Configure the plugin with the given credentials and mode
	Configure(ctx context.Context, in *Configure_Request, opts ...grpc.CallOption) (*Configure_Response, error)
	// Migrate tables to the given plugin version
	Migrate(ctx context.Context, in *Migrate_Request, opts ...grpc.CallOption) (*Migrate_Response, error)
	// Write resources
	Write(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[Write_Request, Write_Response], error)
	// Send signal to flush and close open connections
	Close(ctx context.Context, in *Close_Request, opts ...grpc.CallOption) (*Close_Response, error)
	// DeleteStale deletes stale data that was inserted by a given source
	// and is older than the given timestamp
	DeleteStale(ctx context.Context, in *DeleteStale_Request, opts ...grpc.CallOption) (*DeleteStale_Response, error)
	// Get metrics for the source plugin
	GetMetrics(ctx context.Context, in *GetDestinationMetrics_Request, opts ...grpc.CallOption) (*GetDestinationMetrics_Response, error)
}

type destinationClient struct {
	cc grpc.ClientConnInterface
}

func NewDestinationClient(cc grpc.ClientConnInterface) DestinationClient {
	return &destinationClient{cc}
}

func (c *destinationClient) GetName(ctx context.Context, in *GetName_Request, opts ...grpc.CallOption) (*GetName_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetName_Response)
	err := c.cc.Invoke(ctx, Destination_GetName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) GetVersion(ctx context.Context, in *GetVersion_Request, opts ...grpc.CallOption) (*GetVersion_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetVersion_Response)
	err := c.cc.Invoke(ctx, Destination_GetVersion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) Configure(ctx context.Context, in *Configure_Request, opts ...grpc.CallOption) (*Configure_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Configure_Response)
	err := c.cc.Invoke(ctx, Destination_Configure_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) Migrate(ctx context.Context, in *Migrate_Request, opts ...grpc.CallOption) (*Migrate_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Migrate_Response)
	err := c.cc.Invoke(ctx, Destination_Migrate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) Write(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[Write_Request, Write_Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Destination_ServiceDesc.Streams[0], Destination_Write_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Write_Request, Write_Response]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Destination_WriteClient = grpc.ClientStreamingClient[Write_Request, Write_Response]

func (c *destinationClient) Close(ctx context.Context, in *Close_Request, opts ...grpc.CallOption) (*Close_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Close_Response)
	err := c.cc.Invoke(ctx, Destination_Close_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) DeleteStale(ctx context.Context, in *DeleteStale_Request, opts ...grpc.CallOption) (*DeleteStale_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteStale_Response)
	err := c.cc.Invoke(ctx, Destination_DeleteStale_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) GetMetrics(ctx context.Context, in *GetDestinationMetrics_Request, opts ...grpc.CallOption) (*GetDestinationMetrics_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDestinationMetrics_Response)
	err := c.cc.Invoke(ctx, Destination_GetMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DestinationServer is the server API for Destination service.
// All implementations must embed UnimplementedDestinationServer
// for forward compatibility.
type DestinationServer interface {
	// Get the name of the plugin
	GetName(context.Context, *GetName_Request) (*GetName_Response, error)
	// Get the current version of the plugin
	GetVersion(context.Context, *GetVersion_Request) (*GetVersion_Response, error)
	// Configure the plugin with the given credentials and mode
	Configure(context.Context, *Configure_Request) (*Configure_Response, error)
	// Migrate tables to the given plugin version
	Migrate(context.Context, *Migrate_Request) (*Migrate_Response, error)
	// Write resources
	Write(grpc.ClientStreamingServer[Write_Request, Write_Response]) error
	// Send signal to flush and close open connections
	Close(context.Context, *Close_Request) (*Close_Response, error)
	// DeleteStale deletes stale data that was inserted by a given source
	// and is older than the given timestamp
	DeleteStale(context.Context, *DeleteStale_Request) (*DeleteStale_Response, error)
	// Get metrics for the source plugin
	GetMetrics(context.Context, *GetDestinationMetrics_Request) (*GetDestinationMetrics_Response, error)
	mustEmbedUnimplementedDestinationServer()
}

// UnimplementedDestinationServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDestinationServer struct{}

func (UnimplementedDestinationServer) GetName(context.Context, *GetName_Request) (*GetName_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetName not implemented")
}
func (UnimplementedDestinationServer) GetVersion(context.Context, *GetVersion_Request) (*GetVersion_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedDestinationServer) Configure(context.Context, *Configure_Request) (*Configure_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configure not implemented")
}
func (UnimplementedDestinationServer) Migrate(context.Context, *Migrate_Request) (*Migrate_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Migrate not implemented")
}
func (UnimplementedDestinationServer) Write(grpc.ClientStreamingServer[Write_Request, Write_Response]) error {
	return status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedDestinationServer) Close(context.Context, *Close_Request) (*Close_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedDestinationServer) DeleteStale(context.Context, *DeleteStale_Request) (*DeleteStale_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStale not implemented")
}
func (UnimplementedDestinationServer) GetMetrics(context.Context, *GetDestinationMetrics_Request) (*GetDestinationMetrics_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetrics not implemented")
}
func (UnimplementedDestinationServer) mustEmbedUnimplementedDestinationServer() {}
func (UnimplementedDestinationServer) testEmbeddedByValue()                     {}

// UnsafeDestinationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DestinationServer will
// result in compilation errors.
type UnsafeDestinationServer interface {
	mustEmbedUnimplementedDestinationServer()
}

func RegisterDestinationServer(s grpc.ServiceRegistrar, srv DestinationServer) {
	// If the following call pancis, it indicates UnimplementedDestinationServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Destination_ServiceDesc, srv)
}

func _Destination_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetName_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Destination_GetName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).GetName(ctx, req.(*GetName_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVersion_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Destination_GetVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).GetVersion(ctx, req.(*GetVersion_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Configure_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Destination_Configure_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).Configure(ctx, req.(*Configure_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_Migrate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Migrate_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).Migrate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Destination_Migrate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).Migrate(ctx, req.(*Migrate_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_Write_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DestinationServer).Write(&grpc.GenericServerStream[Write_Request, Write_Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Destination_WriteServer = grpc.ClientStreamingServer[Write_Request, Write_Response]

func _Destination_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Close_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Destination_Close_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).Close(ctx, req.(*Close_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_DeleteStale_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStale_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).DeleteStale(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Destination_DeleteStale_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).DeleteStale(ctx, req.(*DeleteStale_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_GetMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDestinationMetrics_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).GetMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Destination_GetMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).GetMetrics(ctx, req.(*GetDestinationMetrics_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Destination_ServiceDesc is the grpc.ServiceDesc for Destination service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Destination_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cloudquery.destination.v1.Destination",
	HandlerType: (*DestinationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetName",
			Handler:    _Destination_GetName_Handler,
		},
		{
			MethodName: "GetVersion",
			Handler:    _Destination_GetVersion_Handler,
		},
		{
			MethodName: "Configure",
			Handler:    _Destination_Configure_Handler,
		},
		{
			MethodName: "Migrate",
			Handler:    _Destination_Migrate_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Destination_Close_Handler,
		},
		{
			MethodName: "DeleteStale",
			Handler:    _Destination_DeleteStale_Handler,
		},
		{
			MethodName: "GetMetrics",
			Handler:    _Destination_GetMetrics_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Write",
			Handler:       _Destination_Write_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "plugin-pb/destination/v1/destination.proto",
}

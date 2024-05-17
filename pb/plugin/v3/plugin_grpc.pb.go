// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: plugin-pb/plugin/v3/plugin.proto

package plugin

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
	Plugin_GetName_FullMethodName        = "/cloudquery.plugin.v3.Plugin/GetName"
	Plugin_GetVersion_FullMethodName     = "/cloudquery.plugin.v3.Plugin/GetVersion"
	Plugin_GetSpecSchema_FullMethodName  = "/cloudquery.plugin.v3.Plugin/GetSpecSchema"
	Plugin_Init_FullMethodName           = "/cloudquery.plugin.v3.Plugin/Init"
	Plugin_GetTables_FullMethodName      = "/cloudquery.plugin.v3.Plugin/GetTables"
	Plugin_Sync_FullMethodName           = "/cloudquery.plugin.v3.Plugin/Sync"
	Plugin_Read_FullMethodName           = "/cloudquery.plugin.v3.Plugin/Read"
	Plugin_Write_FullMethodName          = "/cloudquery.plugin.v3.Plugin/Write"
	Plugin_Close_FullMethodName          = "/cloudquery.plugin.v3.Plugin/Close"
	Plugin_TestConnection_FullMethodName = "/cloudquery.plugin.v3.Plugin/TestConnection"
)

// PluginClient is the client API for Plugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginClient interface {
	// Get the name of the plugin
	GetName(ctx context.Context, in *GetName_Request, opts ...grpc.CallOption) (*GetName_Response, error)
	// Get the current version of the plugin
	GetVersion(ctx context.Context, in *GetVersion_Request, opts ...grpc.CallOption) (*GetVersion_Response, error)
	// Get plugin spec schema.
	// This will allow validating the input even before calling Init.
	// Should be called before Init.
	GetSpecSchema(ctx context.Context, in *GetSpecSchema_Request, opts ...grpc.CallOption) (*GetSpecSchema_Response, error)
	// Configure the plugin with the given credentials and mode
	Init(ctx context.Context, in *Init_Request, opts ...grpc.CallOption) (*Init_Response, error)
	// Get all tables the source plugin supports. Must be called after Init
	GetTables(ctx context.Context, in *GetTables_Request, opts ...grpc.CallOption) (*GetTables_Response, error)
	// Start a sync on the source plugin. It streams messages as output.
	Sync(ctx context.Context, in *Sync_Request, opts ...grpc.CallOption) (Plugin_SyncClient, error)
	// Start a Read on the source plugin for a given table and schema. It streams messages as output.
	// The plugin assume that that schema was used to also write the data beforehand
	Read(ctx context.Context, in *Read_Request, opts ...grpc.CallOption) (Plugin_ReadClient, error)
	// Write resources. Write is the mirror of Sync, expecting a stream of messages as input.
	Write(ctx context.Context, opts ...grpc.CallOption) (Plugin_WriteClient, error)
	// Send signal to flush and close open connections
	Close(ctx context.Context, in *Close_Request, opts ...grpc.CallOption) (*Close_Response, error)
	// Validate and test the connections used by the plugin
	TestConnection(ctx context.Context, in *TestConnection_Request, opts ...grpc.CallOption) (*TestConnection_Response, error)
}

type pluginClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginClient(cc grpc.ClientConnInterface) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) GetName(ctx context.Context, in *GetName_Request, opts ...grpc.CallOption) (*GetName_Response, error) {
	out := new(GetName_Response)
	err := c.cc.Invoke(ctx, Plugin_GetName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetVersion(ctx context.Context, in *GetVersion_Request, opts ...grpc.CallOption) (*GetVersion_Response, error) {
	out := new(GetVersion_Response)
	err := c.cc.Invoke(ctx, Plugin_GetVersion_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetSpecSchema(ctx context.Context, in *GetSpecSchema_Request, opts ...grpc.CallOption) (*GetSpecSchema_Response, error) {
	out := new(GetSpecSchema_Response)
	err := c.cc.Invoke(ctx, Plugin_GetSpecSchema_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Init(ctx context.Context, in *Init_Request, opts ...grpc.CallOption) (*Init_Response, error) {
	out := new(Init_Response)
	err := c.cc.Invoke(ctx, Plugin_Init_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetTables(ctx context.Context, in *GetTables_Request, opts ...grpc.CallOption) (*GetTables_Response, error) {
	out := new(GetTables_Response)
	err := c.cc.Invoke(ctx, Plugin_GetTables_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Sync(ctx context.Context, in *Sync_Request, opts ...grpc.CallOption) (Plugin_SyncClient, error) {
	stream, err := c.cc.NewStream(ctx, &Plugin_ServiceDesc.Streams[0], Plugin_Sync_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pluginSyncClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Plugin_SyncClient interface {
	Recv() (*Sync_Response, error)
	grpc.ClientStream
}

type pluginSyncClient struct {
	grpc.ClientStream
}

func (x *pluginSyncClient) Recv() (*Sync_Response, error) {
	m := new(Sync_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pluginClient) Read(ctx context.Context, in *Read_Request, opts ...grpc.CallOption) (Plugin_ReadClient, error) {
	stream, err := c.cc.NewStream(ctx, &Plugin_ServiceDesc.Streams[1], Plugin_Read_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pluginReadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Plugin_ReadClient interface {
	Recv() (*Read_Response, error)
	grpc.ClientStream
}

type pluginReadClient struct {
	grpc.ClientStream
}

func (x *pluginReadClient) Recv() (*Read_Response, error) {
	m := new(Read_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pluginClient) Write(ctx context.Context, opts ...grpc.CallOption) (Plugin_WriteClient, error) {
	stream, err := c.cc.NewStream(ctx, &Plugin_ServiceDesc.Streams[2], Plugin_Write_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pluginWriteClient{stream}
	return x, nil
}

type Plugin_WriteClient interface {
	Send(*Write_Request) error
	CloseAndRecv() (*Write_Response, error)
	grpc.ClientStream
}

type pluginWriteClient struct {
	grpc.ClientStream
}

func (x *pluginWriteClient) Send(m *Write_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *pluginWriteClient) CloseAndRecv() (*Write_Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Write_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pluginClient) Close(ctx context.Context, in *Close_Request, opts ...grpc.CallOption) (*Close_Response, error) {
	out := new(Close_Response)
	err := c.cc.Invoke(ctx, Plugin_Close_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) TestConnection(ctx context.Context, in *TestConnection_Request, opts ...grpc.CallOption) (*TestConnection_Response, error) {
	out := new(TestConnection_Response)
	err := c.cc.Invoke(ctx, Plugin_TestConnection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
// All implementations must embed UnimplementedPluginServer
// for forward compatibility
type PluginServer interface {
	// Get the name of the plugin
	GetName(context.Context, *GetName_Request) (*GetName_Response, error)
	// Get the current version of the plugin
	GetVersion(context.Context, *GetVersion_Request) (*GetVersion_Response, error)
	// Get plugin spec schema.
	// This will allow validating the input even before calling Init.
	// Should be called before Init.
	GetSpecSchema(context.Context, *GetSpecSchema_Request) (*GetSpecSchema_Response, error)
	// Configure the plugin with the given credentials and mode
	Init(context.Context, *Init_Request) (*Init_Response, error)
	// Get all tables the source plugin supports. Must be called after Init
	GetTables(context.Context, *GetTables_Request) (*GetTables_Response, error)
	// Start a sync on the source plugin. It streams messages as output.
	Sync(*Sync_Request, Plugin_SyncServer) error
	// Start a Read on the source plugin for a given table and schema. It streams messages as output.
	// The plugin assume that that schema was used to also write the data beforehand
	Read(*Read_Request, Plugin_ReadServer) error
	// Write resources. Write is the mirror of Sync, expecting a stream of messages as input.
	Write(Plugin_WriteServer) error
	// Send signal to flush and close open connections
	Close(context.Context, *Close_Request) (*Close_Response, error)
	// Validate and test the connections used by the plugin
	TestConnection(context.Context, *TestConnection_Request) (*TestConnection_Response, error)
	mustEmbedUnimplementedPluginServer()
}

// UnimplementedPluginServer must be embedded to have forward compatible implementations.
type UnimplementedPluginServer struct {
}

func (UnimplementedPluginServer) GetName(context.Context, *GetName_Request) (*GetName_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetName not implemented")
}
func (UnimplementedPluginServer) GetVersion(context.Context, *GetVersion_Request) (*GetVersion_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedPluginServer) GetSpecSchema(context.Context, *GetSpecSchema_Request) (*GetSpecSchema_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSpecSchema not implemented")
}
func (UnimplementedPluginServer) Init(context.Context, *Init_Request) (*Init_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedPluginServer) GetTables(context.Context, *GetTables_Request) (*GetTables_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTables not implemented")
}
func (UnimplementedPluginServer) Sync(*Sync_Request, Plugin_SyncServer) error {
	return status.Errorf(codes.Unimplemented, "method Sync not implemented")
}
func (UnimplementedPluginServer) Read(*Read_Request, Plugin_ReadServer) error {
	return status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedPluginServer) Write(Plugin_WriteServer) error {
	return status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedPluginServer) Close(context.Context, *Close_Request) (*Close_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedPluginServer) TestConnection(context.Context, *TestConnection_Request) (*TestConnection_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestConnection not implemented")
}
func (UnimplementedPluginServer) mustEmbedUnimplementedPluginServer() {}

// UnsafePluginServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginServer will
// result in compilation errors.
type UnsafePluginServer interface {
	mustEmbedUnimplementedPluginServer()
}

func RegisterPluginServer(s grpc.ServiceRegistrar, srv PluginServer) {
	s.RegisterService(&Plugin_ServiceDesc, srv)
}

func _Plugin_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetName_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_GetName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetName(ctx, req.(*GetName_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVersion_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_GetVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetVersion(ctx, req.(*GetVersion_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetSpecSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSpecSchema_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetSpecSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_GetSpecSchema_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetSpecSchema(ctx, req.(*GetSpecSchema_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Init_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_Init_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Init(ctx, req.(*Init_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTables_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_GetTables_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetTables(ctx, req.(*GetTables_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Sync_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Sync_Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PluginServer).Sync(m, &pluginSyncServer{stream})
}

type Plugin_SyncServer interface {
	Send(*Sync_Response) error
	grpc.ServerStream
}

type pluginSyncServer struct {
	grpc.ServerStream
}

func (x *pluginSyncServer) Send(m *Sync_Response) error {
	return x.ServerStream.SendMsg(m)
}

func _Plugin_Read_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Read_Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PluginServer).Read(m, &pluginReadServer{stream})
}

type Plugin_ReadServer interface {
	Send(*Read_Response) error
	grpc.ServerStream
}

type pluginReadServer struct {
	grpc.ServerStream
}

func (x *pluginReadServer) Send(m *Read_Response) error {
	return x.ServerStream.SendMsg(m)
}

func _Plugin_Write_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PluginServer).Write(&pluginWriteServer{stream})
}

type Plugin_WriteServer interface {
	SendAndClose(*Write_Response) error
	Recv() (*Write_Request, error)
	grpc.ServerStream
}

type pluginWriteServer struct {
	grpc.ServerStream
}

func (x *pluginWriteServer) SendAndClose(m *Write_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *pluginWriteServer) Recv() (*Write_Request, error) {
	m := new(Write_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Plugin_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Close_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_Close_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Close(ctx, req.(*Close_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_TestConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestConnection_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).TestConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_TestConnection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).TestConnection(ctx, req.(*TestConnection_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Plugin_ServiceDesc is the grpc.ServiceDesc for Plugin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Plugin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cloudquery.plugin.v3.Plugin",
	HandlerType: (*PluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetName",
			Handler:    _Plugin_GetName_Handler,
		},
		{
			MethodName: "GetVersion",
			Handler:    _Plugin_GetVersion_Handler,
		},
		{
			MethodName: "GetSpecSchema",
			Handler:    _Plugin_GetSpecSchema_Handler,
		},
		{
			MethodName: "Init",
			Handler:    _Plugin_Init_Handler,
		},
		{
			MethodName: "GetTables",
			Handler:    _Plugin_GetTables_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Plugin_Close_Handler,
		},
		{
			MethodName: "TestConnection",
			Handler:    _Plugin_TestConnection_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Sync",
			Handler:       _Plugin_Sync_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Read",
			Handler:       _Plugin_Read_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Write",
			Handler:       _Plugin_Write_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "plugin-pb/plugin/v3/plugin.proto",
}

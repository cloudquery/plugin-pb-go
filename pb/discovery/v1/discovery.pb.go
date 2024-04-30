// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.0
// 	protoc        v4.23.4
// source: plugin-pb/discovery/v1/discovery.proto

package discovery

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetVersions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVersions) Reset() {
	*x = GetVersions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_pb_discovery_v1_discovery_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersions) ProtoMessage() {}

func (x *GetVersions) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_discovery_v1_discovery_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersions.ProtoReflect.Descriptor instead.
func (*GetVersions) Descriptor() ([]byte, []int) {
	return file_plugin_pb_discovery_v1_discovery_proto_rawDescGZIP(), []int{0}
}

type GetVersions_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVersions_Request) Reset() {
	*x = GetVersions_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_pb_discovery_v1_discovery_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersions_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersions_Request) ProtoMessage() {}

func (x *GetVersions_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_discovery_v1_discovery_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersions_Request.ProtoReflect.Descriptor instead.
func (*GetVersions_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_discovery_v1_discovery_proto_rawDescGZIP(), []int{0, 0}
}

type GetVersions_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Versions []int32 `protobuf:"varint,1,rep,packed,name=versions,proto3" json:"versions,omitempty"`
}

func (x *GetVersions_Response) Reset() {
	*x = GetVersions_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_pb_discovery_v1_discovery_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersions_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersions_Response) ProtoMessage() {}

func (x *GetVersions_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_discovery_v1_discovery_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersions_Response.ProtoReflect.Descriptor instead.
func (*GetVersions_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_discovery_v1_discovery_proto_rawDescGZIP(), []int{0, 1}
}

func (x *GetVersions_Response) GetVersions() []int32 {
	if x != nil {
		return x.Versions
	}
	return nil
}

var File_plugin_pb_discovery_v1_discovery_proto protoreflect.FileDescriptor

var file_plugin_pb_discovery_v1_discovery_proto_rawDesc = []byte{
	0x0a, 0x26, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2d, 0x70, 0x62, 0x2f, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76,
	0x31, 0x22, 0x40, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x32, 0x77, 0x0a, 0x09, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x12, 0x6a, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x2c, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x5c, 0x0a, 0x1a,
	0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x50, 0x01, 0x5a, 0x3c, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2d, 0x70, 0x62, 0x2d, 0x67, 0x6f,
	0x2f, 0x70, 0x62, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x76, 0x31,
	0x3b, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_plugin_pb_discovery_v1_discovery_proto_rawDescOnce sync.Once
	file_plugin_pb_discovery_v1_discovery_proto_rawDescData = file_plugin_pb_discovery_v1_discovery_proto_rawDesc
)

func file_plugin_pb_discovery_v1_discovery_proto_rawDescGZIP() []byte {
	file_plugin_pb_discovery_v1_discovery_proto_rawDescOnce.Do(func() {
		file_plugin_pb_discovery_v1_discovery_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_pb_discovery_v1_discovery_proto_rawDescData)
	})
	return file_plugin_pb_discovery_v1_discovery_proto_rawDescData
}

var file_plugin_pb_discovery_v1_discovery_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_plugin_pb_discovery_v1_discovery_proto_goTypes = []interface{}{
	(*GetVersions)(nil),          // 0: cloudquery.discovery.v1.GetVersions
	(*GetVersions_Request)(nil),  // 1: cloudquery.discovery.v1.GetVersions.Request
	(*GetVersions_Response)(nil), // 2: cloudquery.discovery.v1.GetVersions.Response
}
var file_plugin_pb_discovery_v1_discovery_proto_depIdxs = []int32{
	1, // 0: cloudquery.discovery.v1.Discovery.GetVersions:input_type -> cloudquery.discovery.v1.GetVersions.Request
	2, // 1: cloudquery.discovery.v1.Discovery.GetVersions:output_type -> cloudquery.discovery.v1.GetVersions.Response
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_plugin_pb_discovery_v1_discovery_proto_init() }
func file_plugin_pb_discovery_v1_discovery_proto_init() {
	if File_plugin_pb_discovery_v1_discovery_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plugin_pb_discovery_v1_discovery_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_pb_discovery_v1_discovery_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersions_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_pb_discovery_v1_discovery_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersions_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_plugin_pb_discovery_v1_discovery_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_pb_discovery_v1_discovery_proto_goTypes,
		DependencyIndexes: file_plugin_pb_discovery_v1_discovery_proto_depIdxs,
		MessageInfos:      file_plugin_pb_discovery_v1_discovery_proto_msgTypes,
	}.Build()
	File_plugin_pb_discovery_v1_discovery_proto = out.File
	file_plugin_pb_discovery_v1_discovery_proto_rawDesc = nil
	file_plugin_pb_discovery_v1_discovery_proto_goTypes = nil
	file_plugin_pb_discovery_v1_discovery_proto_depIdxs = nil
}

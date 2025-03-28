// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.23.4
// source: plugin-pb/destination/v1/destination.proto

package destination

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetName struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetName) Reset() {
	*x = GetName{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetName) ProtoMessage() {}

func (x *GetName) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetName.ProtoReflect.Descriptor instead.
func (*GetName) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{0}
}

type GetVersion struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetVersion) Reset() {
	*x = GetVersion{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersion) ProtoMessage() {}

func (x *GetVersion) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersion.ProtoReflect.Descriptor instead.
func (*GetVersion) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{1}
}

type Configure struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Configure) Reset() {
	*x = Configure{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Configure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configure) ProtoMessage() {}

func (x *Configure) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configure.ProtoReflect.Descriptor instead.
func (*Configure) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{2}
}

type Migrate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Migrate) Reset() {
	*x = Migrate{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Migrate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate) ProtoMessage() {}

func (x *Migrate) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate.ProtoReflect.Descriptor instead.
func (*Migrate) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{3}
}

type Write struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Write) Reset() {
	*x = Write{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Write) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write) ProtoMessage() {}

func (x *Write) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write.ProtoReflect.Descriptor instead.
func (*Write) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{4}
}

type Close struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Close) Reset() {
	*x = Close{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Close) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Close) ProtoMessage() {}

func (x *Close) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Close.ProtoReflect.Descriptor instead.
func (*Close) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{5}
}

type DeleteStale struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteStale) Reset() {
	*x = DeleteStale{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteStale) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStale) ProtoMessage() {}

func (x *DeleteStale) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStale.ProtoReflect.Descriptor instead.
func (*DeleteStale) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{6}
}

type GetDestinationMetrics struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDestinationMetrics) Reset() {
	*x = GetDestinationMetrics{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDestinationMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDestinationMetrics) ProtoMessage() {}

func (x *GetDestinationMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDestinationMetrics.ProtoReflect.Descriptor instead.
func (*GetDestinationMetrics) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{7}
}

type GetName_Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetName_Request) Reset() {
	*x = GetName_Request{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetName_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetName_Request) ProtoMessage() {}

func (x *GetName_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetName_Request.ProtoReflect.Descriptor instead.
func (*GetName_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{0, 0}
}

type GetName_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetName_Response) Reset() {
	*x = GetName_Response{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetName_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetName_Response) ProtoMessage() {}

func (x *GetName_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetName_Response.ProtoReflect.Descriptor instead.
func (*GetName_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{0, 1}
}

func (x *GetName_Response) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetVersion_Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetVersion_Request) Reset() {
	*x = GetVersion_Request{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetVersion_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersion_Request) ProtoMessage() {}

func (x *GetVersion_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersion_Request.ProtoReflect.Descriptor instead.
func (*GetVersion_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{1, 0}
}

type GetVersion_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Version       string                 `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetVersion_Response) Reset() {
	*x = GetVersion_Response{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetVersion_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersion_Response) ProtoMessage() {}

func (x *GetVersion_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersion_Response.ProtoReflect.Descriptor instead.
func (*GetVersion_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{1, 1}
}

func (x *GetVersion_Response) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type Configure_Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Holds information such as credentials, regions, accounts, etc'
	// Marshalled spec.SourceSpec or spec.DestinationSpec
	Config        []byte `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Configure_Request) Reset() {
	*x = Configure_Request{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Configure_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configure_Request) ProtoMessage() {}

func (x *Configure_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configure_Request.ProtoReflect.Descriptor instead.
func (*Configure_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Configure_Request) GetConfig() []byte {
	if x != nil {
		return x.Config
	}
	return nil
}

type Configure_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Configure_Response) Reset() {
	*x = Configure_Response{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Configure_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configure_Response) ProtoMessage() {}

func (x *Configure_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configure_Response.ProtoReflect.Descriptor instead.
func (*Configure_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{2, 1}
}

type Migrate_Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version       string                 `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Tables        [][]byte               `protobuf:"bytes,3,rep,name=tables,proto3" json:"tables,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Migrate_Request) Reset() {
	*x = Migrate_Request{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Migrate_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate_Request) ProtoMessage() {}

func (x *Migrate_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate_Request.ProtoReflect.Descriptor instead.
func (*Migrate_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{3, 0}
}

func (x *Migrate_Request) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Migrate_Request) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Migrate_Request) GetTables() [][]byte {
	if x != nil {
		return x.Tables
	}
	return nil
}

type Migrate_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Migrate_Response) Reset() {
	*x = Migrate_Response{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Migrate_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate_Response) ProtoMessage() {}

func (x *Migrate_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate_Response.ProtoReflect.Descriptor instead.
func (*Migrate_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{3, 1}
}

type Write_Request struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Source    string                 `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// marshaled arrow.schema
	Tables [][]byte `protobuf:"bytes,3,rep,name=tables,proto3" json:"tables,omitempty"`
	// marshalled *schema.Resources
	Resource []byte `protobuf:"bytes,4,opt,name=resource,proto3" json:"resource,omitempty"`
	// marshalled specs.Source
	SourceSpec    []byte `protobuf:"bytes,5,opt,name=source_spec,json=sourceSpec,proto3" json:"source_spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Write_Request) Reset() {
	*x = Write_Request{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Write_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write_Request) ProtoMessage() {}

func (x *Write_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write_Request.ProtoReflect.Descriptor instead.
func (*Write_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{4, 0}
}

func (x *Write_Request) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Write_Request) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Write_Request) GetTables() [][]byte {
	if x != nil {
		return x.Tables
	}
	return nil
}

func (x *Write_Request) GetResource() []byte {
	if x != nil {
		return x.Resource
	}
	return nil
}

func (x *Write_Request) GetSourceSpec() []byte {
	if x != nil {
		return x.SourceSpec
	}
	return nil
}

type Write_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Write_Response) Reset() {
	*x = Write_Response{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Write_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write_Response) ProtoMessage() {}

func (x *Write_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write_Response.ProtoReflect.Descriptor instead.
func (*Write_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{4, 1}
}

type Close_Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Close_Request) Reset() {
	*x = Close_Request{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Close_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Close_Request) ProtoMessage() {}

func (x *Close_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Close_Request.ProtoReflect.Descriptor instead.
func (*Close_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{5, 0}
}

type Close_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Close_Response) Reset() {
	*x = Close_Response{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[19]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Close_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Close_Response) ProtoMessage() {}

func (x *Close_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[19]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Close_Response.ProtoReflect.Descriptor instead.
func (*Close_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{5, 1}
}

type DeleteStale_Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Source        string                 `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Tables        [][]byte               `protobuf:"bytes,3,rep,name=tables,proto3" json:"tables,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteStale_Request) Reset() {
	*x = DeleteStale_Request{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[20]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteStale_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStale_Request) ProtoMessage() {}

func (x *DeleteStale_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[20]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStale_Request.ProtoReflect.Descriptor instead.
func (*DeleteStale_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{6, 0}
}

func (x *DeleteStale_Request) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *DeleteStale_Request) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *DeleteStale_Request) GetTables() [][]byte {
	if x != nil {
		return x.Tables
	}
	return nil
}

type DeleteStale_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FailedDeletes uint64                 `protobuf:"varint,1,opt,name=failed_deletes,json=failedDeletes,proto3" json:"failed_deletes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteStale_Response) Reset() {
	*x = DeleteStale_Response{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[21]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteStale_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStale_Response) ProtoMessage() {}

func (x *DeleteStale_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[21]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStale_Response.ProtoReflect.Descriptor instead.
func (*DeleteStale_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{6, 1}
}

func (x *DeleteStale_Response) GetFailedDeletes() uint64 {
	if x != nil {
		return x.FailedDeletes
	}
	return 0
}

type GetDestinationMetrics_Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDestinationMetrics_Request) Reset() {
	*x = GetDestinationMetrics_Request{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[22]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDestinationMetrics_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDestinationMetrics_Request) ProtoMessage() {}

func (x *GetDestinationMetrics_Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[22]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDestinationMetrics_Request.ProtoReflect.Descriptor instead.
func (*GetDestinationMetrics_Request) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{7, 0}
}

type GetDestinationMetrics_Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// marshalled json of plugins.DestinationMetrics
	Metrics       []byte `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDestinationMetrics_Response) Reset() {
	*x = GetDestinationMetrics_Response{}
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[23]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDestinationMetrics_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDestinationMetrics_Response) ProtoMessage() {}

func (x *GetDestinationMetrics_Response) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_pb_destination_v1_destination_proto_msgTypes[23]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDestinationMetrics_Response.ProtoReflect.Descriptor instead.
func (*GetDestinationMetrics_Response) Descriptor() ([]byte, []int) {
	return file_plugin_pb_destination_v1_destination_proto_rawDescGZIP(), []int{7, 1}
}

func (x *GetDestinationMetrics_Response) GetMetrics() []byte {
	if x != nil {
		return x.Metrics
	}
	return nil
}

var File_plugin_pb_destination_v1_destination_proto protoreflect.FileDescriptor

const file_plugin_pb_destination_v1_destination_proto_rawDesc = "" +
	"\n" +
	"*plugin-pb/destination/v1/destination.proto\x12\x19cloudquery.destination.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"4\n" +
	"\aGetName\x1a\t\n" +
	"\aRequest\x1a\x1e\n" +
	"\bResponse\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"=\n" +
	"\n" +
	"GetVersion\x1a\t\n" +
	"\aRequest\x1a$\n" +
	"\bResponse\x12\x18\n" +
	"\aversion\x18\x01 \x01(\tR\aversion\":\n" +
	"\tConfigure\x1a!\n" +
	"\aRequest\x12\x16\n" +
	"\x06config\x18\x01 \x01(\fR\x06config\x1a\n" +
	"\n" +
	"\bResponse\"f\n" +
	"\aMigrate\x1aO\n" +
	"\aRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x18\n" +
	"\aversion\x18\x02 \x01(\tR\aversion\x12\x16\n" +
	"\x06tables\x18\x03 \x03(\fR\x06tables\x1a\n" +
	"\n" +
	"\bResponse\"\xc6\x01\n" +
	"\x05Write\x1a\xb0\x01\n" +
	"\aRequest\x12\x16\n" +
	"\x06source\x18\x01 \x01(\tR\x06source\x128\n" +
	"\ttimestamp\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\ttimestamp\x12\x16\n" +
	"\x06tables\x18\x03 \x03(\fR\x06tables\x12\x1a\n" +
	"\bresource\x18\x04 \x01(\fR\bresource\x12\x1f\n" +
	"\vsource_spec\x18\x05 \x01(\fR\n" +
	"sourceSpec\x1a\n" +
	"\n" +
	"\bResponse\"\x1e\n" +
	"\x05Close\x1a\t\n" +
	"\aRequest\x1a\n" +
	"\n" +
	"\bResponse\"\xb5\x01\n" +
	"\vDeleteStale\x1as\n" +
	"\aRequest\x12\x16\n" +
	"\x06source\x18\x01 \x01(\tR\x06source\x128\n" +
	"\ttimestamp\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\ttimestamp\x12\x16\n" +
	"\x06tables\x18\x03 \x03(\fR\x06tables\x1a1\n" +
	"\bResponse\x12%\n" +
	"\x0efailed_deletes\x18\x01 \x01(\x04R\rfailedDeletes\"H\n" +
	"\x15GetDestinationMetrics\x1a\t\n" +
	"\aRequest\x1a$\n" +
	"\bResponse\x12\x18\n" +
	"\ametrics\x18\x01 \x01(\fR\ametrics2\xde\x06\n" +
	"\vDestination\x12b\n" +
	"\aGetName\x12*.cloudquery.destination.v1.GetName.Request\x1a+.cloudquery.destination.v1.GetName.Response\x12k\n" +
	"\n" +
	"GetVersion\x12-.cloudquery.destination.v1.GetVersion.Request\x1a..cloudquery.destination.v1.GetVersion.Response\x12h\n" +
	"\tConfigure\x12,.cloudquery.destination.v1.Configure.Request\x1a-.cloudquery.destination.v1.Configure.Response\x12b\n" +
	"\aMigrate\x12*.cloudquery.destination.v1.Migrate.Request\x1a+.cloudquery.destination.v1.Migrate.Response\x12^\n" +
	"\x05Write\x12(.cloudquery.destination.v1.Write.Request\x1a).cloudquery.destination.v1.Write.Response(\x01\x12\\\n" +
	"\x05Close\x12(.cloudquery.destination.v1.Close.Request\x1a).cloudquery.destination.v1.Close.Response\x12n\n" +
	"\vDeleteStale\x12..cloudquery.destination.v1.DeleteStale.Request\x1a/.cloudquery.destination.v1.DeleteStale.Response\x12\x81\x01\n" +
	"\n" +
	"GetMetrics\x128.cloudquery.destination.v1.GetDestinationMetrics.Request\x1a9.cloudquery.destination.v1.GetDestinationMetrics.ResponseBBZ@github.com/cloudquery/plugin-pb-go/pb/destination/v1;destinationb\x06proto3"

var (
	file_plugin_pb_destination_v1_destination_proto_rawDescOnce sync.Once
	file_plugin_pb_destination_v1_destination_proto_rawDescData []byte
)

func file_plugin_pb_destination_v1_destination_proto_rawDescGZIP() []byte {
	file_plugin_pb_destination_v1_destination_proto_rawDescOnce.Do(func() {
		file_plugin_pb_destination_v1_destination_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_plugin_pb_destination_v1_destination_proto_rawDesc), len(file_plugin_pb_destination_v1_destination_proto_rawDesc)))
	})
	return file_plugin_pb_destination_v1_destination_proto_rawDescData
}

var file_plugin_pb_destination_v1_destination_proto_msgTypes = make([]protoimpl.MessageInfo, 24)
var file_plugin_pb_destination_v1_destination_proto_goTypes = []any{
	(*GetName)(nil),                        // 0: cloudquery.destination.v1.GetName
	(*GetVersion)(nil),                     // 1: cloudquery.destination.v1.GetVersion
	(*Configure)(nil),                      // 2: cloudquery.destination.v1.Configure
	(*Migrate)(nil),                        // 3: cloudquery.destination.v1.Migrate
	(*Write)(nil),                          // 4: cloudquery.destination.v1.Write
	(*Close)(nil),                          // 5: cloudquery.destination.v1.Close
	(*DeleteStale)(nil),                    // 6: cloudquery.destination.v1.DeleteStale
	(*GetDestinationMetrics)(nil),          // 7: cloudquery.destination.v1.GetDestinationMetrics
	(*GetName_Request)(nil),                // 8: cloudquery.destination.v1.GetName.Request
	(*GetName_Response)(nil),               // 9: cloudquery.destination.v1.GetName.Response
	(*GetVersion_Request)(nil),             // 10: cloudquery.destination.v1.GetVersion.Request
	(*GetVersion_Response)(nil),            // 11: cloudquery.destination.v1.GetVersion.Response
	(*Configure_Request)(nil),              // 12: cloudquery.destination.v1.Configure.Request
	(*Configure_Response)(nil),             // 13: cloudquery.destination.v1.Configure.Response
	(*Migrate_Request)(nil),                // 14: cloudquery.destination.v1.Migrate.Request
	(*Migrate_Response)(nil),               // 15: cloudquery.destination.v1.Migrate.Response
	(*Write_Request)(nil),                  // 16: cloudquery.destination.v1.Write.Request
	(*Write_Response)(nil),                 // 17: cloudquery.destination.v1.Write.Response
	(*Close_Request)(nil),                  // 18: cloudquery.destination.v1.Close.Request
	(*Close_Response)(nil),                 // 19: cloudquery.destination.v1.Close.Response
	(*DeleteStale_Request)(nil),            // 20: cloudquery.destination.v1.DeleteStale.Request
	(*DeleteStale_Response)(nil),           // 21: cloudquery.destination.v1.DeleteStale.Response
	(*GetDestinationMetrics_Request)(nil),  // 22: cloudquery.destination.v1.GetDestinationMetrics.Request
	(*GetDestinationMetrics_Response)(nil), // 23: cloudquery.destination.v1.GetDestinationMetrics.Response
	(*timestamppb.Timestamp)(nil),          // 24: google.protobuf.Timestamp
}
var file_plugin_pb_destination_v1_destination_proto_depIdxs = []int32{
	24, // 0: cloudquery.destination.v1.Write.Request.timestamp:type_name -> google.protobuf.Timestamp
	24, // 1: cloudquery.destination.v1.DeleteStale.Request.timestamp:type_name -> google.protobuf.Timestamp
	8,  // 2: cloudquery.destination.v1.Destination.GetName:input_type -> cloudquery.destination.v1.GetName.Request
	10, // 3: cloudquery.destination.v1.Destination.GetVersion:input_type -> cloudquery.destination.v1.GetVersion.Request
	12, // 4: cloudquery.destination.v1.Destination.Configure:input_type -> cloudquery.destination.v1.Configure.Request
	14, // 5: cloudquery.destination.v1.Destination.Migrate:input_type -> cloudquery.destination.v1.Migrate.Request
	16, // 6: cloudquery.destination.v1.Destination.Write:input_type -> cloudquery.destination.v1.Write.Request
	18, // 7: cloudquery.destination.v1.Destination.Close:input_type -> cloudquery.destination.v1.Close.Request
	20, // 8: cloudquery.destination.v1.Destination.DeleteStale:input_type -> cloudquery.destination.v1.DeleteStale.Request
	22, // 9: cloudquery.destination.v1.Destination.GetMetrics:input_type -> cloudquery.destination.v1.GetDestinationMetrics.Request
	9,  // 10: cloudquery.destination.v1.Destination.GetName:output_type -> cloudquery.destination.v1.GetName.Response
	11, // 11: cloudquery.destination.v1.Destination.GetVersion:output_type -> cloudquery.destination.v1.GetVersion.Response
	13, // 12: cloudquery.destination.v1.Destination.Configure:output_type -> cloudquery.destination.v1.Configure.Response
	15, // 13: cloudquery.destination.v1.Destination.Migrate:output_type -> cloudquery.destination.v1.Migrate.Response
	17, // 14: cloudquery.destination.v1.Destination.Write:output_type -> cloudquery.destination.v1.Write.Response
	19, // 15: cloudquery.destination.v1.Destination.Close:output_type -> cloudquery.destination.v1.Close.Response
	21, // 16: cloudquery.destination.v1.Destination.DeleteStale:output_type -> cloudquery.destination.v1.DeleteStale.Response
	23, // 17: cloudquery.destination.v1.Destination.GetMetrics:output_type -> cloudquery.destination.v1.GetDestinationMetrics.Response
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_plugin_pb_destination_v1_destination_proto_init() }
func file_plugin_pb_destination_v1_destination_proto_init() {
	if File_plugin_pb_destination_v1_destination_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_plugin_pb_destination_v1_destination_proto_rawDesc), len(file_plugin_pb_destination_v1_destination_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   24,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_pb_destination_v1_destination_proto_goTypes,
		DependencyIndexes: file_plugin_pb_destination_v1_destination_proto_depIdxs,
		MessageInfos:      file_plugin_pb_destination_v1_destination_proto_msgTypes,
	}.Build()
	File_plugin_pb_destination_v1_destination_proto = out.File
	file_plugin_pb_destination_v1_destination_proto_goTypes = nil
	file_plugin_pb_destination_v1_destination_proto_depIdxs = nil
}

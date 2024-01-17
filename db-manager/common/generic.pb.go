// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proto/generic.proto

package common

import (
	any1 "github.com/golang/protobuf/ptypes/any"
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

// Insert data
type ProtobufInsertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keyspace  string      `protobuf:"bytes,1,opt,name=keyspace,proto3" json:"keyspace,omitempty"`
	Protobufs []*any1.Any `protobuf:"bytes,2,rep,name=protobufs,proto3" json:"protobufs,omitempty"`
}

func (x *ProtobufInsertRequest) Reset() {
	*x = ProtobufInsertRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_generic_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtobufInsertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtobufInsertRequest) ProtoMessage() {}

func (x *ProtobufInsertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_generic_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtobufInsertRequest.ProtoReflect.Descriptor instead.
func (*ProtobufInsertRequest) Descriptor() ([]byte, []int) {
	return file_proto_generic_proto_rawDescGZIP(), []int{0}
}

func (x *ProtobufInsertRequest) GetKeyspace() string {
	if x != nil {
		return x.Keyspace
	}
	return ""
}

func (x *ProtobufInsertRequest) GetProtobufs() []*any1.Any {
	if x != nil {
		return x.Protobufs
	}
	return nil
}

type ProtobufInsertResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errs []string `protobuf:"bytes,1,rep,name=errs,proto3" json:"errs,omitempty"`
}

func (x *ProtobufInsertResponse) Reset() {
	*x = ProtobufInsertResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_generic_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtobufInsertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtobufInsertResponse) ProtoMessage() {}

func (x *ProtobufInsertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_generic_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtobufInsertResponse.ProtoReflect.Descriptor instead.
func (*ProtobufInsertResponse) Descriptor() ([]byte, []int) {
	return file_proto_generic_proto_rawDescGZIP(), []int{1}
}

func (x *ProtobufInsertResponse) GetErrs() []string {
	if x != nil {
		return x.Errs
	}
	return nil
}

// Delete data
type ProtobufDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keyspace   string `protobuf:"bytes,1,opt,name=keyspace,proto3" json:"keyspace,omitempty"`
	Table      string `protobuf:"bytes,2,opt,name=table,proto3" json:"table,omitempty"`
	Column     string `protobuf:"bytes,3,opt,name=column,proto3" json:"column,omitempty"`
	Constraint string `protobuf:"bytes,4,opt,name=constraint,proto3" json:"constraint,omitempty"`
}

func (x *ProtobufDeleteRequest) Reset() {
	*x = ProtobufDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_generic_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtobufDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtobufDeleteRequest) ProtoMessage() {}

func (x *ProtobufDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_generic_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtobufDeleteRequest.ProtoReflect.Descriptor instead.
func (*ProtobufDeleteRequest) Descriptor() ([]byte, []int) {
	return file_proto_generic_proto_rawDescGZIP(), []int{2}
}

func (x *ProtobufDeleteRequest) GetKeyspace() string {
	if x != nil {
		return x.Keyspace
	}
	return ""
}

func (x *ProtobufDeleteRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *ProtobufDeleteRequest) GetColumn() string {
	if x != nil {
		return x.Column
	}
	return ""
}

func (x *ProtobufDeleteRequest) GetConstraint() string {
	if x != nil {
		return x.Constraint
	}
	return ""
}

type ProtobufDeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errs []string `protobuf:"bytes,1,rep,name=errs,proto3" json:"errs,omitempty"`
}

func (x *ProtobufDeleteResponse) Reset() {
	*x = ProtobufDeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_generic_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtobufDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtobufDeleteResponse) ProtoMessage() {}

func (x *ProtobufDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_generic_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtobufDeleteResponse.ProtoReflect.Descriptor instead.
func (*ProtobufDeleteResponse) Descriptor() ([]byte, []int) {
	return file_proto_generic_proto_rawDescGZIP(), []int{3}
}

func (x *ProtobufDeleteResponse) GetErrs() []string {
	if x != nil {
		return x.Errs
	}
	return nil
}

// Drop Table (clear data)
type ProtobufDroptableRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keyspace string `protobuf:"bytes,1,opt,name=keyspace,proto3" json:"keyspace,omitempty"`
	Table    string `protobuf:"bytes,2,opt,name=table,proto3" json:"table,omitempty"`
}

func (x *ProtobufDroptableRequest) Reset() {
	*x = ProtobufDroptableRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_generic_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtobufDroptableRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtobufDroptableRequest) ProtoMessage() {}

func (x *ProtobufDroptableRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_generic_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtobufDroptableRequest.ProtoReflect.Descriptor instead.
func (*ProtobufDroptableRequest) Descriptor() ([]byte, []int) {
	return file_proto_generic_proto_rawDescGZIP(), []int{4}
}

func (x *ProtobufDroptableRequest) GetKeyspace() string {
	if x != nil {
		return x.Keyspace
	}
	return ""
}

func (x *ProtobufDroptableRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

type ProtobufDroptableResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errs []string `protobuf:"bytes,1,rep,name=errs,proto3" json:"errs,omitempty"`
}

func (x *ProtobufDroptableResponse) Reset() {
	*x = ProtobufDroptableResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_generic_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtobufDroptableResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtobufDroptableResponse) ProtoMessage() {}

func (x *ProtobufDroptableResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_generic_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtobufDroptableResponse.ProtoReflect.Descriptor instead.
func (*ProtobufDroptableResponse) Descriptor() ([]byte, []int) {
	return file_proto_generic_proto_rawDescGZIP(), []int{5}
}

func (x *ProtobufDroptableResponse) GetErrs() []string {
	if x != nil {
		return x.Errs
	}
	return nil
}

var File_proto_generic_proto protoreflect.FileDescriptor

var file_proto_generic_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x69, 0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x69, 0x6e, 0x73,
	0x65, 0x72, 0x74, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6b,
	0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b,
	0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x52, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x22, 0x2e, 0x0a, 0x18, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x5f, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x72, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x65, 0x72, 0x72, 0x73, 0x22, 0x83, 0x01, 0x0a, 0x17,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6c,
	0x75, 0x6d, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d,
	0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e,
	0x74, 0x22, 0x2e, 0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x65, 0x72, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x65, 0x72, 0x72,
	0x73, 0x22, 0x4e, 0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x64, 0x72,
	0x6f, 0x70, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x22, 0x31, 0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x64, 0x72,
	0x6f, 0x70, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x65, 0x72, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x65, 0x72, 0x72, 0x73, 0x32, 0xd1, 0x01, 0x0a, 0x09, 0x44, 0x42, 0x47, 0x65, 0x6e, 0x65, 0x72,
	0x69, 0x63, 0x12, 0x3d, 0x0a, 0x06, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x18, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x5f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x5f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3d, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x46, 0x0a, 0x09, 0x44, 0x72, 0x6f, 0x70, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x1b, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x64, 0x72, 0x6f, 0x70, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x5f, 0x64, 0x72, 0x6f, 0x70, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2e, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_generic_proto_rawDescOnce sync.Once
	file_proto_generic_proto_rawDescData = file_proto_generic_proto_rawDesc
)

func file_proto_generic_proto_rawDescGZIP() []byte {
	file_proto_generic_proto_rawDescOnce.Do(func() {
		file_proto_generic_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_generic_proto_rawDescData)
	})
	return file_proto_generic_proto_rawDescData
}

var file_proto_generic_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_generic_proto_goTypes = []interface{}{
	(*ProtobufInsertRequest)(nil),     // 0: protobuf_insert_request
	(*ProtobufInsertResponse)(nil),    // 1: protobuf_insert_response
	(*ProtobufDeleteRequest)(nil),     // 2: protobuf_delete_request
	(*ProtobufDeleteResponse)(nil),    // 3: protobuf_delete_response
	(*ProtobufDroptableRequest)(nil),  // 4: protobuf_droptable_request
	(*ProtobufDroptableResponse)(nil), // 5: protobuf_droptable_response
	(*any1.Any)(nil),                  // 6: google.protobuf.Any
}
var file_proto_generic_proto_depIdxs = []int32{
	6, // 0: protobuf_insert_request.protobufs:type_name -> google.protobuf.Any
	0, // 1: DBGeneric.Insert:input_type -> protobuf_insert_request
	2, // 2: DBGeneric.Delete:input_type -> protobuf_delete_request
	4, // 3: DBGeneric.DropTable:input_type -> protobuf_droptable_request
	1, // 4: DBGeneric.Insert:output_type -> protobuf_insert_response
	3, // 5: DBGeneric.Delete:output_type -> protobuf_delete_response
	5, // 6: DBGeneric.DropTable:output_type -> protobuf_droptable_response
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_generic_proto_init() }
func file_proto_generic_proto_init() {
	if File_proto_generic_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_generic_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtobufInsertRequest); i {
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
		file_proto_generic_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtobufInsertResponse); i {
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
		file_proto_generic_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtobufDeleteRequest); i {
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
		file_proto_generic_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtobufDeleteResponse); i {
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
		file_proto_generic_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtobufDroptableRequest); i {
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
		file_proto_generic_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtobufDroptableResponse); i {
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
			RawDescriptor: file_proto_generic_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_generic_proto_goTypes,
		DependencyIndexes: file_proto_generic_proto_depIdxs,
		MessageInfos:      file_proto_generic_proto_msgTypes,
	}.Build()
	File_proto_generic_proto = out.File
	file_proto_generic_proto_rawDesc = nil
	file_proto_generic_proto_goTypes = nil
	file_proto_generic_proto_depIdxs = nil
}

//
//Trigger is the service for creating triggers to start deployments.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: trigger.proto

package trigger

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type TriggerStatus int32

const (
	TriggerStatus_NOT_STARTED       TriggerStatus = 0
	TriggerStatus_IN_PROGRESS       TriggerStatus = 1
	TriggerStatus_AWAITING_APPROVAL TriggerStatus = 2
	TriggerStatus_FAILED            TriggerStatus = 3
	TriggerStatus_SUCCEEDED         TriggerStatus = 4
	TriggerStatus_SKIPPED           TriggerStatus = 5
)

// Enum value maps for TriggerStatus.
var (
	TriggerStatus_name = map[int32]string{
		0: "NOT_STARTED",
		1: "IN_PROGRESS",
		2: "AWAITING_APPROVAL",
		3: "FAILED",
		4: "SUCCEEDED",
		5: "SKIPPED",
	}
	TriggerStatus_value = map[string]int32{
		"NOT_STARTED":       0,
		"IN_PROGRESS":       1,
		"AWAITING_APPROVAL": 2,
		"FAILED":            3,
		"SUCCEEDED":         4,
		"SKIPPED":           5,
	}
)

func (x TriggerStatus) Enum() *TriggerStatus {
	p := new(TriggerStatus)
	*p = x
	return p
}

func (x TriggerStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TriggerStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_trigger_proto_enumTypes[0].Descriptor()
}

func (TriggerStatus) Type() protoreflect.EnumType {
	return &file_trigger_proto_enumTypes[0]
}

func (x TriggerStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TriggerStatus.Descriptor instead.
func (TriggerStatus) EnumDescriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{0}
}

type TriggerManual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Args string `protobuf:"bytes,3,opt,name=args,proto3" json:"args,omitempty"`
}

func (x *TriggerManual) Reset() {
	*x = TriggerManual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trigger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerManual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerManual) ProtoMessage() {}

func (x *TriggerManual) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerManual.ProtoReflect.Descriptor instead.
func (*TriggerManual) Descriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{0}
}

func (x *TriggerManual) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *TriggerManual) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TriggerManual) GetArgs() string {
	if x != nil {
		return x.Args
	}
	return ""
}

type ManualTriggered struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RunUuid string `protobuf:"bytes,1,opt,name=run_uuid,json=runUuid,proto3" json:"run_uuid,omitempty"`
}

func (x *ManualTriggered) Reset() {
	*x = ManualTriggered{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trigger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManualTriggered) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManualTriggered) ProtoMessage() {}

func (x *ManualTriggered) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManualTriggered.ProtoReflect.Descriptor instead.
func (*ManualTriggered) Descriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{1}
}

func (x *ManualTriggered) GetRunUuid() string {
	if x != nil {
		return x.RunUuid
	}
	return ""
}

type TriggerWeb struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid  string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Token string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	Args  string `protobuf:"bytes,4,opt,name=args,proto3" json:"args,omitempty"`
}

func (x *TriggerWeb) Reset() {
	*x = TriggerWeb{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trigger_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerWeb) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerWeb) ProtoMessage() {}

func (x *TriggerWeb) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerWeb.ProtoReflect.Descriptor instead.
func (*TriggerWeb) Descriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{2}
}

func (x *TriggerWeb) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *TriggerWeb) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TriggerWeb) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *TriggerWeb) GetArgs() string {
	if x != nil {
		return x.Args
	}
	return ""
}

type WebTriggered struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RunUuid string `protobuf:"bytes,1,opt,name=run_uuid,json=runUuid,proto3" json:"run_uuid,omitempty"`
}

func (x *WebTriggered) Reset() {
	*x = WebTriggered{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trigger_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebTriggered) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebTriggered) ProtoMessage() {}

func (x *WebTriggered) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebTriggered.ProtoReflect.Descriptor instead.
func (*WebTriggered) Descriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{3}
}

func (x *WebTriggered) GetRunUuid() string {
	if x != nil {
		return x.RunUuid
	}
	return ""
}

type ReadStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeploymentUuid  string `protobuf:"bytes,1,opt,name=deployment_uuid,json=deploymentUuid,proto3" json:"deployment_uuid,omitempty"`
	DeploymentToken string `protobuf:"bytes,2,opt,name=deployment_token,json=deploymentToken,proto3" json:"deployment_token,omitempty"`
	RunUuid         string `protobuf:"bytes,3,opt,name=run_uuid,json=runUuid,proto3" json:"run_uuid,omitempty"`
}

func (x *ReadStatus) Reset() {
	*x = ReadStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trigger_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadStatus) ProtoMessage() {}

func (x *ReadStatus) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadStatus.ProtoReflect.Descriptor instead.
func (*ReadStatus) Descriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{4}
}

func (x *ReadStatus) GetDeploymentUuid() string {
	if x != nil {
		return x.DeploymentUuid
	}
	return ""
}

func (x *ReadStatus) GetDeploymentToken() string {
	if x != nil {
		return x.DeploymentToken
	}
	return ""
}

func (x *ReadStatus) GetRunUuid() string {
	if x != nil {
		return x.RunUuid
	}
	return ""
}

type StatusRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status TriggerStatus `protobuf:"varint,1,opt,name=status,proto3,enum=redsail.bosn.TriggerStatus" json:"status,omitempty"`
}

func (x *StatusRead) Reset() {
	*x = StatusRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trigger_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusRead) ProtoMessage() {}

func (x *StatusRead) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusRead.ProtoReflect.Descriptor instead.
func (*StatusRead) Descriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{5}
}

func (x *StatusRead) GetStatus() TriggerStatus {
	if x != nil {
		return x.Status
	}
	return TriggerStatus_NOT_STARTED
}

var File_trigger_proto protoreflect.FileDescriptor

var file_trigger_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0c, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x22, 0x4b, 0x0a,
	0x0d, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x4d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x12, 0x12,
	0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x22, 0x2c, 0x0a, 0x0f, 0x4d, 0x61,
	0x6e, 0x75, 0x61, 0x6c, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x65, 0x64, 0x12, 0x19, 0x0a,
	0x08, 0x72, 0x75, 0x6e, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x72, 0x75, 0x6e, 0x55, 0x75, 0x69, 0x64, 0x22, 0x5e, 0x0a, 0x0a, 0x54, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x57, 0x65, 0x62, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x22, 0x29, 0x0a, 0x0c, 0x57, 0x65, 0x62, 0x54,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x65, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x75, 0x6e, 0x5f,
	0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x75, 0x6e, 0x55,
	0x75, 0x69, 0x64, 0x22, 0x7b, 0x0a, 0x0a, 0x52, 0x65, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x27, 0x0a, 0x0f, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x75, 0x69, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x64, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x75, 0x6e, 0x5f, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x75, 0x6e, 0x55, 0x75, 0x69, 0x64,
	0x22, 0x41, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x61, 0x64, 0x12, 0x33,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b,
	0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x54, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x2a, 0x70, 0x0a, 0x0d, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x0f, 0x0a, 0x0b, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x52,
	0x54, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47,
	0x52, 0x45, 0x53, 0x53, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x57, 0x41, 0x49, 0x54, 0x49,
	0x4e, 0x47, 0x5f, 0x41, 0x50, 0x50, 0x52, 0x4f, 0x56, 0x41, 0x4c, 0x10, 0x02, 0x12, 0x0a, 0x0a,
	0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x55, 0x43,
	0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x4b, 0x49, 0x50,
	0x50, 0x45, 0x44, 0x10, 0x05, 0x32, 0xca, 0x01, 0x0a, 0x07, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x12, 0x44, 0x0a, 0x06, 0x4d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x12, 0x1b, 0x2e, 0x72, 0x65,
	0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x4d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x1a, 0x1d, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61,
	0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x4d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x54, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x65, 0x64, 0x12, 0x3b, 0x0a, 0x03, 0x57, 0x65, 0x62, 0x12, 0x18,
	0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x54, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x57, 0x65, 0x62, 0x1a, 0x1a, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61,
	0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x57, 0x65, 0x62, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x65, 0x64, 0x12, 0x3c, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18,
	0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65,
	0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x18, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61,
	0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65,
	0x61, 0x64, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x74, 0x65, 0x63, 0x68, 0x6e, 0x6f, 0x6c, 0x6f,
	0x67, 0x69, 0x65, 0x73, 0x2f, 0x62, 0x6f, 0x61, 0x74, 0x73, 0x77, 0x61, 0x69, 0x6e, 0x2f, 0x72,
	0x70, 0x63, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x3b, 0x74, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trigger_proto_rawDescOnce sync.Once
	file_trigger_proto_rawDescData = file_trigger_proto_rawDesc
)

func file_trigger_proto_rawDescGZIP() []byte {
	file_trigger_proto_rawDescOnce.Do(func() {
		file_trigger_proto_rawDescData = protoimpl.X.CompressGZIP(file_trigger_proto_rawDescData)
	})
	return file_trigger_proto_rawDescData
}

var file_trigger_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_trigger_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_trigger_proto_goTypes = []interface{}{
	(TriggerStatus)(0),      // 0: redsail.bosn.TriggerStatus
	(*TriggerManual)(nil),   // 1: redsail.bosn.TriggerManual
	(*ManualTriggered)(nil), // 2: redsail.bosn.ManualTriggered
	(*TriggerWeb)(nil),      // 3: redsail.bosn.TriggerWeb
	(*WebTriggered)(nil),    // 4: redsail.bosn.WebTriggered
	(*ReadStatus)(nil),      // 5: redsail.bosn.ReadStatus
	(*StatusRead)(nil),      // 6: redsail.bosn.StatusRead
}
var file_trigger_proto_depIdxs = []int32{
	0, // 0: redsail.bosn.StatusRead.status:type_name -> redsail.bosn.TriggerStatus
	1, // 1: redsail.bosn.Trigger.Manual:input_type -> redsail.bosn.TriggerManual
	3, // 2: redsail.bosn.Trigger.Web:input_type -> redsail.bosn.TriggerWeb
	5, // 3: redsail.bosn.Trigger.Status:input_type -> redsail.bosn.ReadStatus
	2, // 4: redsail.bosn.Trigger.Manual:output_type -> redsail.bosn.ManualTriggered
	4, // 5: redsail.bosn.Trigger.Web:output_type -> redsail.bosn.WebTriggered
	6, // 6: redsail.bosn.Trigger.Status:output_type -> redsail.bosn.StatusRead
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_trigger_proto_init() }
func file_trigger_proto_init() {
	if File_trigger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trigger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TriggerManual); i {
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
		file_trigger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManualTriggered); i {
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
		file_trigger_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TriggerWeb); i {
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
		file_trigger_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebTriggered); i {
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
		file_trigger_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadStatus); i {
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
		file_trigger_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusRead); i {
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
			RawDescriptor: file_trigger_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trigger_proto_goTypes,
		DependencyIndexes: file_trigger_proto_depIdxs,
		EnumInfos:         file_trigger_proto_enumTypes,
		MessageInfos:      file_trigger_proto_msgTypes,
	}.Build()
	File_trigger_proto = out.File
	file_trigger_proto_rawDesc = nil
	file_trigger_proto_goTypes = nil
	file_trigger_proto_depIdxs = nil
}

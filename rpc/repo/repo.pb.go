//
//Repo is the service managing external repositories, such as helm.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: repo.proto

package repo

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

type CreateRepo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`         // name of the repo
	Endpoint string `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"` // repo endpoint
}

func (x *CreateRepo) Reset() {
	*x = CreateRepo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRepo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRepo) ProtoMessage() {}

func (x *CreateRepo) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRepo.ProtoReflect.Descriptor instead.
func (*CreateRepo) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{0}
}

func (x *CreateRepo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateRepo) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

type RepoCreated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RepoCreated) Reset() {
	*x = RepoCreated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepoCreated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoCreated) ProtoMessage() {}

func (x *RepoCreated) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepoCreated.ProtoReflect.Descriptor instead.
func (*RepoCreated) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{1}
}

type UpdateRepo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`         // unique id of the repo
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`         // name of the repo
	Endpoint string `protobuf:"bytes,3,opt,name=endpoint,proto3" json:"endpoint,omitempty"` // repo endpoint
}

func (x *UpdateRepo) Reset() {
	*x = UpdateRepo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRepo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRepo) ProtoMessage() {}

func (x *UpdateRepo) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRepo.ProtoReflect.Descriptor instead.
func (*UpdateRepo) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateRepo) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *UpdateRepo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateRepo) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

type RepoUpdated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RepoUpdated) Reset() {
	*x = RepoUpdated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepoUpdated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoUpdated) ProtoMessage() {}

func (x *RepoUpdated) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepoUpdated.ProtoReflect.Descriptor instead.
func (*RepoUpdated) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{3}
}

type DestroyRepo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"` // unique id of the repo
}

func (x *DestroyRepo) Reset() {
	*x = DestroyRepo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DestroyRepo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DestroyRepo) ProtoMessage() {}

func (x *DestroyRepo) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DestroyRepo.ProtoReflect.Descriptor instead.
func (*DestroyRepo) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{4}
}

func (x *DestroyRepo) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type RepoDestroyed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RepoDestroyed) Reset() {
	*x = RepoDestroyed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepoDestroyed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoDestroyed) ProtoMessage() {}

func (x *RepoDestroyed) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepoDestroyed.ProtoReflect.Descriptor instead.
func (*RepoDestroyed) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{5}
}

type ReadRepo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"` // unique id of the repo
}

func (x *ReadRepo) Reset() {
	*x = ReadRepo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRepo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRepo) ProtoMessage() {}

func (x *ReadRepo) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRepo.ProtoReflect.Descriptor instead.
func (*ReadRepo) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{6}
}

func (x *ReadRepo) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type RepoRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`         // unique id of the repo
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`         // name of the repo
	Endpoint string `protobuf:"bytes,3,opt,name=endpoint,proto3" json:"endpoint,omitempty"` // repo endpoint
	Ready    bool   `protobuf:"varint,6,opt,name=ready,proto3" json:"ready,omitempty"`      // repo ready status, based on whether index.yaml can be fetched
}

func (x *RepoRead) Reset() {
	*x = RepoRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepoRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoRead) ProtoMessage() {}

func (x *RepoRead) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepoRead.ProtoReflect.Descriptor instead.
func (*RepoRead) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{7}
}

func (x *RepoRead) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *RepoRead) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RepoRead) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *RepoRead) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

type ReadRepos struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReadRepos) Reset() {
	*x = ReadRepos{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRepos) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRepos) ProtoMessage() {}

func (x *ReadRepos) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRepos.ProtoReflect.Descriptor instead.
func (*ReadRepos) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{8}
}

type ReposRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repos []*RepoRead `protobuf:"bytes,1,rep,name=repos,proto3" json:"repos,omitempty"` // repos read
}

func (x *ReposRead) Reset() {
	*x = ReposRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReposRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReposRead) ProtoMessage() {}

func (x *ReposRead) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReposRead.ProtoReflect.Descriptor instead.
func (*ReposRead) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{9}
}

func (x *ReposRead) GetRepos() []*RepoRead {
	if x != nil {
		return x.Repos
	}
	return nil
}

type ChartRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string         `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`         // the chart name
	Versions []*VersionRead `protobuf:"bytes,2,rep,name=versions,proto3" json:"versions,omitempty"` // the versions available for this chart
}

func (x *ChartRead) Reset() {
	*x = ChartRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChartRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChartRead) ProtoMessage() {}

func (x *ChartRead) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChartRead.ProtoReflect.Descriptor instead.
func (*ChartRead) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{10}
}

func (x *ChartRead) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChartRead) GetVersions() []*VersionRead {
	if x != nil {
		return x.Versions
	}
	return nil
}

type VersionRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`                                     // the name of the chart
	ChartVersion string `protobuf:"bytes,2,opt,name=chart_version,json=chartVersion,proto3" json:"chart_version,omitempty"` // the chart version
	AppVersion   string `protobuf:"bytes,3,opt,name=app_version,json=appVersion,proto3" json:"app_version,omitempty"`       // the chart's default app version
	Description  string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`                       // description of the chart
	Url          string `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`                                       // the url for this specific version of the chart
}

func (x *VersionRead) Reset() {
	*x = VersionRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionRead) ProtoMessage() {}

func (x *VersionRead) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionRead.ProtoReflect.Descriptor instead.
func (*VersionRead) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{11}
}

func (x *VersionRead) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *VersionRead) GetChartVersion() string {
	if x != nil {
		return x.ChartVersion
	}
	return ""
}

func (x *VersionRead) GetAppVersion() string {
	if x != nil {
		return x.AppVersion
	}
	return ""
}

func (x *VersionRead) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *VersionRead) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type ChartsRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Charts []*ChartRead `protobuf:"bytes,1,rep,name=charts,proto3" json:"charts,omitempty"` // the list of charts
}

func (x *ChartsRead) Reset() {
	*x = ChartsRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChartsRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChartsRead) ProtoMessage() {}

func (x *ChartsRead) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChartsRead.ProtoReflect.Descriptor instead.
func (*ChartsRead) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{12}
}

func (x *ChartsRead) GetCharts() []*ChartRead {
	if x != nil {
		return x.Charts
	}
	return nil
}

var File_repo_proto protoreflect.FileDescriptor

var file_repo_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x72, 0x65,
	0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x22, 0x3c, 0x0a, 0x0a, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0x0d, 0x0a, 0x0b, 0x52, 0x65, 0x70, 0x6f,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x22, 0x50, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0x0d, 0x0a, 0x0b, 0x52, 0x65, 0x70,
	0x6f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x22, 0x21, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x74,
	0x72, 0x6f, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x22, 0x0f, 0x0a, 0x0d, 0x52,
	0x65, 0x70, 0x6f, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x22, 0x1e, 0x0a, 0x08,
	0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x22, 0x64, 0x0a, 0x08,
	0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61,
	0x64, 0x79, 0x22, 0x0b, 0x0a, 0x09, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x22,
	0x39, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x52, 0x65, 0x61, 0x64, 0x12, 0x2c, 0x0a, 0x05,
	0x72, 0x65, 0x70, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x72, 0x65,
	0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x52,
	0x65, 0x61, 0x64, 0x52, 0x05, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x22, 0x56, 0x0a, 0x09, 0x43, 0x68,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x22, 0x9b, 0x01, 0x0a, 0x0b, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x68, 0x61, 0x72, 0x74, 0x5f,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63,
	0x68, 0x61, 0x72, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x61,
	0x70, 0x70, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x61, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c,
	0x22, 0x3d, 0x0a, 0x0a, 0x43, 0x68, 0x61, 0x72, 0x74, 0x73, 0x52, 0x65, 0x61, 0x64, 0x12, 0x2f,
	0x0a, 0x06, 0x63, 0x68, 0x61, 0x72, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x43, 0x68,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x61, 0x64, 0x52, 0x06, 0x63, 0x68, 0x61, 0x72, 0x74, 0x73, 0x32,
	0xf4, 0x02, 0x0a, 0x04, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x3d, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x12, 0x18, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73,
	0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x19, 0x2e, 0x72,
	0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x3d, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x18, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x19, 0x2e, 0x72, 0x65,
	0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x41, 0x0a, 0x07, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f,
	0x79, 0x12, 0x19, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e,
	0x2e, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x1b, 0x2e, 0x72,
	0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x12, 0x36, 0x0a, 0x04, 0x52, 0x65, 0x61,
	0x64, 0x12, 0x16, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e,
	0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x16, 0x2e, 0x72, 0x65, 0x64, 0x73,
	0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x61,
	0x64, 0x12, 0x37, 0x0a, 0x03, 0x41, 0x6c, 0x6c, 0x12, 0x17, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61,
	0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x1a, 0x17, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e,
	0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x52, 0x65, 0x61, 0x64, 0x12, 0x3a, 0x0a, 0x06, 0x43, 0x68,
	0x61, 0x72, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62,
	0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x18, 0x2e, 0x72,
	0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x43, 0x68, 0x61, 0x72,
	0x74, 0x73, 0x52, 0x65, 0x61, 0x64, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x74, 0x65, 0x63, 0x68,
	0x6e, 0x6f, 0x6c, 0x6f, 0x67, 0x69, 0x65, 0x73, 0x2f, 0x62, 0x6f, 0x61, 0x74, 0x73, 0x77, 0x61,
	0x69, 0x6e, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x65, 0x70, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_repo_proto_rawDescOnce sync.Once
	file_repo_proto_rawDescData = file_repo_proto_rawDesc
)

func file_repo_proto_rawDescGZIP() []byte {
	file_repo_proto_rawDescOnce.Do(func() {
		file_repo_proto_rawDescData = protoimpl.X.CompressGZIP(file_repo_proto_rawDescData)
	})
	return file_repo_proto_rawDescData
}

var file_repo_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_repo_proto_goTypes = []interface{}{
	(*CreateRepo)(nil),    // 0: redsail.bosn.CreateRepo
	(*RepoCreated)(nil),   // 1: redsail.bosn.RepoCreated
	(*UpdateRepo)(nil),    // 2: redsail.bosn.UpdateRepo
	(*RepoUpdated)(nil),   // 3: redsail.bosn.RepoUpdated
	(*DestroyRepo)(nil),   // 4: redsail.bosn.DestroyRepo
	(*RepoDestroyed)(nil), // 5: redsail.bosn.RepoDestroyed
	(*ReadRepo)(nil),      // 6: redsail.bosn.ReadRepo
	(*RepoRead)(nil),      // 7: redsail.bosn.RepoRead
	(*ReadRepos)(nil),     // 8: redsail.bosn.ReadRepos
	(*ReposRead)(nil),     // 9: redsail.bosn.ReposRead
	(*ChartRead)(nil),     // 10: redsail.bosn.ChartRead
	(*VersionRead)(nil),   // 11: redsail.bosn.VersionRead
	(*ChartsRead)(nil),    // 12: redsail.bosn.ChartsRead
}
var file_repo_proto_depIdxs = []int32{
	7,  // 0: redsail.bosn.ReposRead.repos:type_name -> redsail.bosn.RepoRead
	11, // 1: redsail.bosn.ChartRead.versions:type_name -> redsail.bosn.VersionRead
	10, // 2: redsail.bosn.ChartsRead.charts:type_name -> redsail.bosn.ChartRead
	0,  // 3: redsail.bosn.Repo.Create:input_type -> redsail.bosn.CreateRepo
	2,  // 4: redsail.bosn.Repo.Update:input_type -> redsail.bosn.UpdateRepo
	4,  // 5: redsail.bosn.Repo.Destroy:input_type -> redsail.bosn.DestroyRepo
	6,  // 6: redsail.bosn.Repo.Read:input_type -> redsail.bosn.ReadRepo
	8,  // 7: redsail.bosn.Repo.All:input_type -> redsail.bosn.ReadRepos
	6,  // 8: redsail.bosn.Repo.Charts:input_type -> redsail.bosn.ReadRepo
	1,  // 9: redsail.bosn.Repo.Create:output_type -> redsail.bosn.RepoCreated
	3,  // 10: redsail.bosn.Repo.Update:output_type -> redsail.bosn.RepoUpdated
	5,  // 11: redsail.bosn.Repo.Destroy:output_type -> redsail.bosn.RepoDestroyed
	7,  // 12: redsail.bosn.Repo.Read:output_type -> redsail.bosn.RepoRead
	9,  // 13: redsail.bosn.Repo.All:output_type -> redsail.bosn.ReposRead
	12, // 14: redsail.bosn.Repo.Charts:output_type -> redsail.bosn.ChartsRead
	9,  // [9:15] is the sub-list for method output_type
	3,  // [3:9] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_repo_proto_init() }
func file_repo_proto_init() {
	if File_repo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_repo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRepo); i {
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
		file_repo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepoCreated); i {
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
		file_repo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRepo); i {
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
		file_repo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepoUpdated); i {
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
		file_repo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DestroyRepo); i {
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
		file_repo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepoDestroyed); i {
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
		file_repo_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadRepo); i {
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
		file_repo_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepoRead); i {
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
		file_repo_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadRepos); i {
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
		file_repo_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReposRead); i {
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
		file_repo_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChartRead); i {
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
		file_repo_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionRead); i {
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
		file_repo_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChartsRead); i {
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
			RawDescriptor: file_repo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_repo_proto_goTypes,
		DependencyIndexes: file_repo_proto_depIdxs,
		MessageInfos:      file_repo_proto_msgTypes,
	}.Build()
	File_repo_proto = out.File
	file_repo_proto_rawDesc = nil
	file_repo_proto_goTypes = nil
	file_repo_proto_depIdxs = nil
}

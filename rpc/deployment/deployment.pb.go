//
//Deployment is the service for creation and management of application installs/upgrades.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: deployment.proto

package delivery

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

type CreateDeployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`                         // the name of this deployment
	RepoId   string `protobuf:"bytes,2,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"`       // the unique id of the repo to get the deployment yaml from
	Branch   string `protobuf:"bytes,3,opt,name=branch,proto3" json:"branch,omitempty"`                     // the branch from the repo to get the file from
	FilePath string `protobuf:"bytes,4,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"` // the path to the deployment file
}

func (x *CreateDeployment) Reset() {
	*x = CreateDeployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateDeployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDeployment) ProtoMessage() {}

func (x *CreateDeployment) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDeployment.ProtoReflect.Descriptor instead.
func (*CreateDeployment) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{0}
}

func (x *CreateDeployment) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateDeployment) GetRepoId() string {
	if x != nil {
		return x.RepoId
	}
	return ""
}

func (x *CreateDeployment) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *CreateDeployment) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

type DeploymentCreated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeploymentCreated) Reset() {
	*x = DeploymentCreated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentCreated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentCreated) ProtoMessage() {}

func (x *DeploymentCreated) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentCreated.ProtoReflect.Descriptor instead.
func (*DeploymentCreated) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{1}
}

type UpdateDeployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`                         // unique id of the deployment
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`                         // the name of this deployment
	RepoId   string `protobuf:"bytes,3,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"`       // the unique id of the repo to get the deployment yaml from
	Branch   string `protobuf:"bytes,4,opt,name=branch,proto3" json:"branch,omitempty"`                     // the branch from the repo to get the file from
	FilePath string `protobuf:"bytes,5,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"` // the path to the deployment file
}

func (x *UpdateDeployment) Reset() {
	*x = UpdateDeployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateDeployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDeployment) ProtoMessage() {}

func (x *UpdateDeployment) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDeployment.ProtoReflect.Descriptor instead.
func (*UpdateDeployment) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateDeployment) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *UpdateDeployment) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateDeployment) GetRepoId() string {
	if x != nil {
		return x.RepoId
	}
	return ""
}

func (x *UpdateDeployment) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *UpdateDeployment) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

type DeploymentUpdated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeploymentUpdated) Reset() {
	*x = DeploymentUpdated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentUpdated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentUpdated) ProtoMessage() {}

func (x *DeploymentUpdated) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentUpdated.ProtoReflect.Descriptor instead.
func (*DeploymentUpdated) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{3}
}

type DestroyDeployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"` // unique id of the deployment
}

func (x *DestroyDeployment) Reset() {
	*x = DestroyDeployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DestroyDeployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DestroyDeployment) ProtoMessage() {}

func (x *DestroyDeployment) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DestroyDeployment.ProtoReflect.Descriptor instead.
func (*DestroyDeployment) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{4}
}

func (x *DestroyDeployment) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type DeploymentDestroyed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeploymentDestroyed) Reset() {
	*x = DeploymentDestroyed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentDestroyed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentDestroyed) ProtoMessage() {}

func (x *DeploymentDestroyed) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentDestroyed.ProtoReflect.Descriptor instead.
func (*DeploymentDestroyed) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{5}
}

type ReadDeployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"` // unique id of the deployment
}

func (x *ReadDeployment) Reset() {
	*x = ReadDeployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadDeployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadDeployment) ProtoMessage() {}

func (x *ReadDeployment) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadDeployment.ProtoReflect.Descriptor instead.
func (*ReadDeployment) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{6}
}

func (x *ReadDeployment) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type DeploymentRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`                         // unique id of the deployment
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`                         // the name of this deployment
	RepoId   string `protobuf:"bytes,3,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"`       // the unique id of the repo to get the deployment yaml from
	RepoName string `protobuf:"bytes,4,opt,name=repo_name,json=repoName,proto3" json:"repo_name,omitempty"` // the name of the repo
	Branch   string `protobuf:"bytes,5,opt,name=branch,proto3" json:"branch,omitempty"`                     // the branch from the repo to get the file from
	FilePath string `protobuf:"bytes,6,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"` // the path to the deployment file
	Yaml     []byte `protobuf:"bytes,7,opt,name=yaml,proto3" json:"yaml,omitempty"`                         // the templated yaml of this deployment
}

func (x *DeploymentRead) Reset() {
	*x = DeploymentRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentRead) ProtoMessage() {}

func (x *DeploymentRead) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentRead.ProtoReflect.Descriptor instead.
func (*DeploymentRead) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{7}
}

func (x *DeploymentRead) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *DeploymentRead) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DeploymentRead) GetRepoId() string {
	if x != nil {
		return x.RepoId
	}
	return ""
}

func (x *DeploymentRead) GetRepoName() string {
	if x != nil {
		return x.RepoName
	}
	return ""
}

func (x *DeploymentRead) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *DeploymentRead) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *DeploymentRead) GetYaml() []byte {
	if x != nil {
		return x.Yaml
	}
	return nil
}

type DeploymentReadSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`                         // unique id of the deployment
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`                         // name of the deployment
	RepoId   string `protobuf:"bytes,3,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"`       // the name of the repo
	RepoName string `protobuf:"bytes,4,opt,name=repo_name,json=repoName,proto3" json:"repo_name,omitempty"` // the name of the repo
	Branch   string `protobuf:"bytes,5,opt,name=branch,proto3" json:"branch,omitempty"`                     // the branch from the repo to get the file from
	FilePath string `protobuf:"bytes,6,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"` // the path to the deployment file
}

func (x *DeploymentReadSummary) Reset() {
	*x = DeploymentReadSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentReadSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentReadSummary) ProtoMessage() {}

func (x *DeploymentReadSummary) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentReadSummary.ProtoReflect.Descriptor instead.
func (*DeploymentReadSummary) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{8}
}

func (x *DeploymentReadSummary) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *DeploymentReadSummary) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DeploymentReadSummary) GetRepoId() string {
	if x != nil {
		return x.RepoId
	}
	return ""
}

func (x *DeploymentReadSummary) GetRepoName() string {
	if x != nil {
		return x.RepoName
	}
	return ""
}

func (x *DeploymentReadSummary) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *DeploymentReadSummary) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

type ReadDeployments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReadDeployments) Reset() {
	*x = ReadDeployments{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadDeployments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadDeployments) ProtoMessage() {}

func (x *ReadDeployments) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadDeployments.ProtoReflect.Descriptor instead.
func (*ReadDeployments) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{9}
}

type DeploymentsRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deployments []*DeploymentReadSummary `protobuf:"bytes,1,rep,name=deployments,proto3" json:"deployments,omitempty"` // the list of deployments
}

func (x *DeploymentsRead) Reset() {
	*x = DeploymentsRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentsRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentsRead) ProtoMessage() {}

func (x *DeploymentsRead) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentsRead.ProtoReflect.Descriptor instead.
func (*DeploymentsRead) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{10}
}

func (x *DeploymentsRead) GetDeployments() []*DeploymentReadSummary {
	if x != nil {
		return x.Deployments
	}
	return nil
}

var File_deployment_proto protoreflect.FileDescriptor

var file_deployment_proto_rawDesc = []byte{
	0x0a, 0x10, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0c, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e,
	0x22, 0x74, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x65, 0x70, 0x6f,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0x13, 0x0a, 0x11, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x22, 0x88, 0x01, 0x0a, 0x10,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x65, 0x70, 0x6f,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0x13, 0x0a, 0x11, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x22, 0x27, 0x0a, 0x11, 0x44,
	0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x75, 0x75, 0x69, 0x64, 0x22, 0x15, 0x0a, 0x13, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x22, 0x24, 0x0a, 0x0e, 0x52,
	0x65, 0x61, 0x64, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69,
	0x64, 0x22, 0xb7, 0x01, 0x0a, 0x0e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07,
	0x72, 0x65, 0x70, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72,
	0x65, 0x70, 0x6f, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x70, 0x6f, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x70, 0x6f, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x79, 0x61, 0x6d, 0x6c, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x79, 0x61, 0x6d, 0x6c, 0x22, 0xaa, 0x01, 0x0a, 0x15,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x61, 0x64, 0x53, 0x75,
	0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a,
	0x07, 0x72, 0x65, 0x70, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x70, 0x6f, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x70, 0x6f, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x70, 0x6f, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0x11, 0x0a, 0x0f, 0x52, 0x65, 0x61, 0x64,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x58, 0x0a, 0x0f, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x61, 0x64, 0x12, 0x45,
	0x0a, 0x0b, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f,
	0x73, 0x6e, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x61,
	0x64, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x0b, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x32, 0xfa, 0x02, 0x0a, 0x0a, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x49, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1e,
	0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x1f,
	0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x44, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12,
	0x49, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x72, 0x65, 0x64, 0x73,
	0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x1f, 0x2e, 0x72, 0x65, 0x64, 0x73,
	0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x4d, 0x0a, 0x07, 0x44, 0x65,
	0x73, 0x74, 0x72, 0x6f, 0x79, 0x12, 0x1f, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e,
	0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x44, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x21, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c,
	0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x12, 0x42, 0x0a, 0x04, 0x52, 0x65, 0x61,
	0x64, 0x12, 0x1c, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e,
	0x2e, 0x52, 0x65, 0x61, 0x64, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x1a,
	0x1c, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x61, 0x64, 0x12, 0x43, 0x0a,
	0x03, 0x41, 0x6c, 0x6c, 0x12, 0x1d, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62,
	0x6f, 0x73, 0x6e, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x1a, 0x1d, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f,
	0x73, 0x6e, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x61, 0x64, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x74, 0x65, 0x63, 0x68, 0x6e, 0x6f, 0x6c, 0x6f,
	0x67, 0x69, 0x65, 0x73, 0x2f, 0x62, 0x6f, 0x61, 0x74, 0x73, 0x77, 0x61, 0x69, 0x6e, 0x2f, 0x72,
	0x70, 0x63, 0x2f, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x3b, 0x64, 0x65, 0x6c, 0x69,
	0x76, 0x65, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_deployment_proto_rawDescOnce sync.Once
	file_deployment_proto_rawDescData = file_deployment_proto_rawDesc
)

func file_deployment_proto_rawDescGZIP() []byte {
	file_deployment_proto_rawDescOnce.Do(func() {
		file_deployment_proto_rawDescData = protoimpl.X.CompressGZIP(file_deployment_proto_rawDescData)
	})
	return file_deployment_proto_rawDescData
}

var file_deployment_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_deployment_proto_goTypes = []interface{}{
	(*CreateDeployment)(nil),      // 0: redsail.bosn.CreateDeployment
	(*DeploymentCreated)(nil),     // 1: redsail.bosn.DeploymentCreated
	(*UpdateDeployment)(nil),      // 2: redsail.bosn.UpdateDeployment
	(*DeploymentUpdated)(nil),     // 3: redsail.bosn.DeploymentUpdated
	(*DestroyDeployment)(nil),     // 4: redsail.bosn.DestroyDeployment
	(*DeploymentDestroyed)(nil),   // 5: redsail.bosn.DeploymentDestroyed
	(*ReadDeployment)(nil),        // 6: redsail.bosn.ReadDeployment
	(*DeploymentRead)(nil),        // 7: redsail.bosn.DeploymentRead
	(*DeploymentReadSummary)(nil), // 8: redsail.bosn.DeploymentReadSummary
	(*ReadDeployments)(nil),       // 9: redsail.bosn.ReadDeployments
	(*DeploymentsRead)(nil),       // 10: redsail.bosn.DeploymentsRead
}
var file_deployment_proto_depIdxs = []int32{
	8,  // 0: redsail.bosn.DeploymentsRead.deployments:type_name -> redsail.bosn.DeploymentReadSummary
	0,  // 1: redsail.bosn.Deployment.Create:input_type -> redsail.bosn.CreateDeployment
	2,  // 2: redsail.bosn.Deployment.Update:input_type -> redsail.bosn.UpdateDeployment
	4,  // 3: redsail.bosn.Deployment.Destroy:input_type -> redsail.bosn.DestroyDeployment
	6,  // 4: redsail.bosn.Deployment.Read:input_type -> redsail.bosn.ReadDeployment
	9,  // 5: redsail.bosn.Deployment.All:input_type -> redsail.bosn.ReadDeployments
	1,  // 6: redsail.bosn.Deployment.Create:output_type -> redsail.bosn.DeploymentCreated
	3,  // 7: redsail.bosn.Deployment.Update:output_type -> redsail.bosn.DeploymentUpdated
	5,  // 8: redsail.bosn.Deployment.Destroy:output_type -> redsail.bosn.DeploymentDestroyed
	7,  // 9: redsail.bosn.Deployment.Read:output_type -> redsail.bosn.DeploymentRead
	10, // 10: redsail.bosn.Deployment.All:output_type -> redsail.bosn.DeploymentsRead
	6,  // [6:11] is the sub-list for method output_type
	1,  // [1:6] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_deployment_proto_init() }
func file_deployment_proto_init() {
	if File_deployment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_deployment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateDeployment); i {
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
		file_deployment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentCreated); i {
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
		file_deployment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateDeployment); i {
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
		file_deployment_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentUpdated); i {
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
		file_deployment_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DestroyDeployment); i {
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
		file_deployment_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentDestroyed); i {
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
		file_deployment_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadDeployment); i {
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
		file_deployment_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentRead); i {
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
		file_deployment_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentReadSummary); i {
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
		file_deployment_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadDeployments); i {
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
		file_deployment_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentsRead); i {
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
			RawDescriptor: file_deployment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_deployment_proto_goTypes,
		DependencyIndexes: file_deployment_proto_depIdxs,
		MessageInfos:      file_deployment_proto_msgTypes,
	}.Build()
	File_deployment_proto = out.File
	file_deployment_proto_rawDesc = nil
	file_deployment_proto_goTypes = nil
	file_deployment_proto_depIdxs = nil
}

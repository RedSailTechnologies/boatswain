// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: kraken/service.proto

package kraken

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

type Cluster struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Endpoint string `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Ready    bool   `protobuf:"varint,3,opt,name=ready,proto3" json:"ready,omitempty"`
}

func (x *Cluster) Reset() {
	*x = Cluster{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kraken_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cluster) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cluster) ProtoMessage() {}

func (x *Cluster) ProtoReflect() protoreflect.Message {
	mi := &file_kraken_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cluster.ProtoReflect.Descriptor instead.
func (*Cluster) Descriptor() ([]byte, []int) {
	return file_kraken_service_proto_rawDescGZIP(), []int{0}
}

func (x *Cluster) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Cluster) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Cluster) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

type ClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClusterRequest) Reset() {
	*x = ClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kraken_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterRequest) ProtoMessage() {}

func (x *ClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kraken_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterRequest.ProtoReflect.Descriptor instead.
func (*ClusterRequest) Descriptor() ([]byte, []int) {
	return file_kraken_service_proto_rawDescGZIP(), []int{1}
}

type ClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clusters []*Cluster `protobuf:"bytes,1,rep,name=clusters,proto3" json:"clusters,omitempty"`
}

func (x *ClusterResponse) Reset() {
	*x = ClusterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kraken_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterResponse) ProtoMessage() {}

func (x *ClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kraken_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterResponse.ProtoReflect.Descriptor instead.
func (*ClusterResponse) Descriptor() ([]byte, []int) {
	return file_kraken_service_proto_rawDescGZIP(), []int{2}
}

func (x *ClusterResponse) GetClusters() []*Cluster {
	if x != nil {
		return x.Clusters
	}
	return nil
}

type Deployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Version   string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *Deployment) Reset() {
	*x = Deployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kraken_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deployment) ProtoMessage() {}

func (x *Deployment) ProtoReflect() protoreflect.Message {
	mi := &file_kraken_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deployment.ProtoReflect.Descriptor instead.
func (*Deployment) Descriptor() ([]byte, []int) {
	return file_kraken_service_proto_rawDescGZIP(), []int{3}
}

func (x *Deployment) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Deployment) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Deployment) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type DeploymentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cluster string `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
}

func (x *DeploymentRequest) Reset() {
	*x = DeploymentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kraken_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentRequest) ProtoMessage() {}

func (x *DeploymentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kraken_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentRequest.ProtoReflect.Descriptor instead.
func (*DeploymentRequest) Descriptor() ([]byte, []int) {
	return file_kraken_service_proto_rawDescGZIP(), []int{4}
}

func (x *DeploymentRequest) GetCluster() string {
	if x != nil {
		return x.Cluster
	}
	return ""
}

type DeploymentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deployments []*Deployment `protobuf:"bytes,1,rep,name=deployments,proto3" json:"deployments,omitempty"`
}

func (x *DeploymentResponse) Reset() {
	*x = DeploymentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kraken_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentResponse) ProtoMessage() {}

func (x *DeploymentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kraken_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentResponse.ProtoReflect.Descriptor instead.
func (*DeploymentResponse) Descriptor() ([]byte, []int) {
	return file_kraken_service_proto_rawDescGZIP(), []int{5}
}

func (x *DeploymentResponse) GetDeployments() []*Deployment {
	if x != nil {
		return x.Deployments
	}
	return nil
}

var File_kraken_service_proto protoreflect.FileDescriptor

var file_kraken_service_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6b, 0x72, 0x61, 0x6b, 0x65, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e,
	0x62, 0x6f, 0x73, 0x6e, 0x22, 0x4f, 0x0a, 0x07, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x72, 0x65, 0x61, 0x64, 0x79, 0x22, 0x10, 0x0a, 0x0e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x44, 0x0a, 0x0f, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x08, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72,
	0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x52, 0x08, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x22, 0x58, 0x0a,
	0x0a, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x2d, 0x0a, 0x11, 0x44, 0x65, 0x70, 0x6c, 0x6f,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x22, 0x50, 0x0a, 0x12, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x0b,
	0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e,
	0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0b, 0x64, 0x65, 0x70,
	0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x32, 0xa3, 0x01, 0x0a, 0x06, 0x4b, 0x72, 0x61,
	0x6b, 0x65, 0x6e, 0x12, 0x47, 0x0a, 0x08, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x12,
	0x1c, 0x2e, 0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e,
	0x72, 0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x0b,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1f, 0x2e, 0x72, 0x65,
	0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x72,
	0x65, 0x64, 0x73, 0x61, 0x69, 0x6c, 0x2e, 0x62, 0x6f, 0x73, 0x6e, 0x2e, 0x44, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x35,
	0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x64,
	0x73, 0x61, 0x69, 0x6c, 0x74, 0x65, 0x63, 0x68, 0x6e, 0x6f, 0x6c, 0x6f, 0x67, 0x69, 0x65, 0x73,
	0x2f, 0x62, 0x6f, 0x61, 0x74, 0x73, 0x77, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6b,
	0x72, 0x61, 0x6b, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kraken_service_proto_rawDescOnce sync.Once
	file_kraken_service_proto_rawDescData = file_kraken_service_proto_rawDesc
)

func file_kraken_service_proto_rawDescGZIP() []byte {
	file_kraken_service_proto_rawDescOnce.Do(func() {
		file_kraken_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_kraken_service_proto_rawDescData)
	})
	return file_kraken_service_proto_rawDescData
}

var file_kraken_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_kraken_service_proto_goTypes = []interface{}{
	(*Cluster)(nil),            // 0: redsail.bosn.Cluster
	(*ClusterRequest)(nil),     // 1: redsail.bosn.ClusterRequest
	(*ClusterResponse)(nil),    // 2: redsail.bosn.ClusterResponse
	(*Deployment)(nil),         // 3: redsail.bosn.Deployment
	(*DeploymentRequest)(nil),  // 4: redsail.bosn.DeploymentRequest
	(*DeploymentResponse)(nil), // 5: redsail.bosn.DeploymentResponse
}
var file_kraken_service_proto_depIdxs = []int32{
	0, // 0: redsail.bosn.ClusterResponse.clusters:type_name -> redsail.bosn.Cluster
	3, // 1: redsail.bosn.DeploymentResponse.deployments:type_name -> redsail.bosn.Deployment
	1, // 2: redsail.bosn.Kraken.Clusters:input_type -> redsail.bosn.ClusterRequest
	4, // 3: redsail.bosn.Kraken.Deployments:input_type -> redsail.bosn.DeploymentRequest
	2, // 4: redsail.bosn.Kraken.Clusters:output_type -> redsail.bosn.ClusterResponse
	5, // 5: redsail.bosn.Kraken.Deployments:output_type -> redsail.bosn.DeploymentResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_kraken_service_proto_init() }
func file_kraken_service_proto_init() {
	if File_kraken_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kraken_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cluster); i {
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
		file_kraken_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterRequest); i {
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
		file_kraken_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterResponse); i {
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
		file_kraken_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Deployment); i {
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
		file_kraken_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentRequest); i {
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
		file_kraken_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentResponse); i {
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
			RawDescriptor: file_kraken_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kraken_service_proto_goTypes,
		DependencyIndexes: file_kraken_service_proto_depIdxs,
		MessageInfos:      file_kraken_service_proto_msgTypes,
	}.Build()
	File_kraken_service_proto = out.File
	file_kraken_service_proto_rawDesc = nil
	file_kraken_service_proto_goTypes = nil
	file_kraken_service_proto_depIdxs = nil
}

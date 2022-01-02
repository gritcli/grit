// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: github.com/gritcli/grit/internal/api/api.proto

package api

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

type SourcesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SourcesRequest) Reset() {
	*x = SourcesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SourcesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourcesRequest) ProtoMessage() {}

func (x *SourcesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourcesRequest.ProtoReflect.Descriptor instead.
func (*SourcesRequest) Descriptor() ([]byte, []int) {
	return file_github_com_gritcli_grit_internal_api_api_proto_rawDescGZIP(), []int{0}
}

type SourcesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sources []*Source `protobuf:"bytes,1,rep,name=sources,proto3" json:"sources,omitempty"`
}

func (x *SourcesResponse) Reset() {
	*x = SourcesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SourcesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourcesResponse) ProtoMessage() {}

func (x *SourcesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourcesResponse.ProtoReflect.Descriptor instead.
func (*SourcesResponse) Descriptor() ([]byte, []int) {
	return file_github_com_gritcli_grit_internal_api_api_proto_rawDescGZIP(), []int{1}
}

func (x *SourcesResponse) GetSources() []*Source {
	if x != nil {
		return x.Sources
	}
	return nil
}

type Source struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Source) Reset() {
	*x = Source{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_github_com_gritcli_grit_internal_api_api_proto_rawDescGZIP(), []int{2}
}

func (x *Source) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Source) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type Repo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourceName  string `protobuf:"bytes,1,opt,name=source_name,json=sourceName,proto3" json:"source_name,omitempty"`
	RepoId      string `protobuf:"bytes,2,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"`
	RepoName    string `protobuf:"bytes,3,opt,name=repo_name,json=repoName,proto3" json:"repo_name,omitempty"`
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	WebUrl      string `protobuf:"bytes,5,opt,name=web_url,json=webUrl,proto3" json:"web_url,omitempty"`
}

func (x *Repo) Reset() {
	*x = Repo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Repo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Repo) ProtoMessage() {}

func (x *Repo) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Repo.ProtoReflect.Descriptor instead.
func (*Repo) Descriptor() ([]byte, []int) {
	return file_github_com_gritcli_grit_internal_api_api_proto_rawDescGZIP(), []int{3}
}

func (x *Repo) GetSourceName() string {
	if x != nil {
		return x.SourceName
	}
	return ""
}

func (x *Repo) GetRepoId() string {
	if x != nil {
		return x.RepoId
	}
	return ""
}

func (x *Repo) GetRepoName() string {
	if x != nil {
		return x.RepoName
	}
	return ""
}

func (x *Repo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Repo) GetWebUrl() string {
	if x != nil {
		return x.WebUrl
	}
	return ""
}

type ResolveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *ResolveRequest) Reset() {
	*x = ResolveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResolveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveRequest) ProtoMessage() {}

func (x *ResolveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveRequest.ProtoReflect.Descriptor instead.
func (*ResolveRequest) Descriptor() ([]byte, []int) {
	return file_github_com_gritcli_grit_internal_api_api_proto_rawDescGZIP(), []int{4}
}

func (x *ResolveRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type ResolveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repo `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
}

func (x *ResolveResponse) Reset() {
	*x = ResolveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResolveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveResponse) ProtoMessage() {}

func (x *ResolveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveResponse.ProtoReflect.Descriptor instead.
func (*ResolveResponse) Descriptor() ([]byte, []int) {
	return file_github_com_gritcli_grit_internal_api_api_proto_rawDescGZIP(), []int{5}
}

func (x *ResolveResponse) GetRepo() *Repo {
	if x != nil {
		return x.Repo
	}
	return nil
}

var File_github_com_gritcli_grit_internal_api_api_proto protoreflect.FileDescriptor

var file_github_com_gritcli_grit_internal_api_api_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x69,
	0x74, 0x63, 0x6c, 0x69, 0x2f, 0x67, 0x72, 0x69, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x67, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x32, 0x2e, 0x61, 0x70, 0x69, 0x22, 0x10, 0x0a,
	0x0e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x40, 0x0a, 0x0f, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2d, 0x0a, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x67, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x32, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x22, 0x3e, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x98, 0x01, 0x0a, 0x04, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x72,
	0x65, 0x70, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65,
	0x70, 0x6f, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x70, 0x6f, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x70, 0x6f, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x77, 0x65, 0x62, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x77, 0x65, 0x62, 0x55, 0x72, 0x6c, 0x22, 0x26, 0x0a, 0x0e,
	0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x22, 0x38, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x32, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x32, 0x93,
	0x01, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x44, 0x0a, 0x07, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x12, 0x1b, 0x2e, 0x67, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x32, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x67, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x32, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x07,
	0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x12, 0x1b, 0x2e, 0x67, 0x72, 0x69, 0x74, 0x2e, 0x76,
	0x32, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x67, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x32, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x30, 0x01, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x69, 0x74, 0x63, 0x6c, 0x69, 0x2f, 0x67, 0x72, 0x69, 0x74, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_gritcli_grit_internal_api_api_proto_rawDescOnce sync.Once
	file_github_com_gritcli_grit_internal_api_api_proto_rawDescData = file_github_com_gritcli_grit_internal_api_api_proto_rawDesc
)

func file_github_com_gritcli_grit_internal_api_api_proto_rawDescGZIP() []byte {
	file_github_com_gritcli_grit_internal_api_api_proto_rawDescOnce.Do(func() {
		file_github_com_gritcli_grit_internal_api_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_gritcli_grit_internal_api_api_proto_rawDescData)
	})
	return file_github_com_gritcli_grit_internal_api_api_proto_rawDescData
}

var file_github_com_gritcli_grit_internal_api_api_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_github_com_gritcli_grit_internal_api_api_proto_goTypes = []interface{}{
	(*SourcesRequest)(nil),  // 0: grit.v2.api.SourcesRequest
	(*SourcesResponse)(nil), // 1: grit.v2.api.SourcesResponse
	(*Source)(nil),          // 2: grit.v2.api.Source
	(*Repo)(nil),            // 3: grit.v2.api.Repo
	(*ResolveRequest)(nil),  // 4: grit.v2.api.ResolveRequest
	(*ResolveResponse)(nil), // 5: grit.v2.api.ResolveResponse
}
var file_github_com_gritcli_grit_internal_api_api_proto_depIdxs = []int32{
	2, // 0: grit.v2.api.SourcesResponse.sources:type_name -> grit.v2.api.Source
	3, // 1: grit.v2.api.ResolveResponse.repo:type_name -> grit.v2.api.Repo
	0, // 2: grit.v2.api.API.Sources:input_type -> grit.v2.api.SourcesRequest
	4, // 3: grit.v2.api.API.Resolve:input_type -> grit.v2.api.ResolveRequest
	1, // 4: grit.v2.api.API.Sources:output_type -> grit.v2.api.SourcesResponse
	5, // 5: grit.v2.api.API.Resolve:output_type -> grit.v2.api.ResolveResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_github_com_gritcli_grit_internal_api_api_proto_init() }
func file_github_com_gritcli_grit_internal_api_api_proto_init() {
	if File_github_com_gritcli_grit_internal_api_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SourcesRequest); i {
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
		file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SourcesResponse); i {
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
		file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Source); i {
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
		file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Repo); i {
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
		file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResolveRequest); i {
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
		file_github_com_gritcli_grit_internal_api_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResolveResponse); i {
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
			RawDescriptor: file_github_com_gritcli_grit_internal_api_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_gritcli_grit_internal_api_api_proto_goTypes,
		DependencyIndexes: file_github_com_gritcli_grit_internal_api_api_proto_depIdxs,
		MessageInfos:      file_github_com_gritcli_grit_internal_api_api_proto_msgTypes,
	}.Build()
	File_github_com_gritcli_grit_internal_api_api_proto = out.File
	file_github_com_gritcli_grit_internal_api_api_proto_rawDesc = nil
	file_github_com_gritcli_grit_internal_api_api_proto_goTypes = nil
	file_github_com_gritcli_grit_internal_api_api_proto_depIdxs = nil
}

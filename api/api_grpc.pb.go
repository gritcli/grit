// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.2
// source: github.com/gritcli/grit/api/api.proto

package api

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

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIClient interface {
	// ListSources lists the configured repository sources.
	ListSources(ctx context.Context, in *ListSourcesRequest, opts ...grpc.CallOption) (*ListSourcesResponse, error)
	// ResolveLocalRepo resolves repository name, URL or other identifier to a
	// list if local repository clones.
	ResolveLocalRepo(ctx context.Context, in *ResolveLocalRepoRequest, opts ...grpc.CallOption) (API_ResolveLocalRepoClient, error)
	// ResolveRemoteRepo resolves repository name, URL or other identifier to a
	// list of remote repositories.
	ResolveRemoteRepo(ctx context.Context, in *ResolveRemoteRepoRequest, opts ...grpc.CallOption) (API_ResolveRemoteRepoClient, error)
	// CloneRemoteRepo makes a local clone of a repository from a source.
	CloneRemoteRepo(ctx context.Context, in *CloneRemoteRepoRequest, opts ...grpc.CallOption) (API_CloneRemoteRepoClient, error)
	// SuggestRepo returns a list of repository names to be used as
	// suggestions for completing a partial repository name.
	SuggestRepo(ctx context.Context, in *SuggestRepoRequest, opts ...grpc.CallOption) (*SuggestResponse, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) ListSources(ctx context.Context, in *ListSourcesRequest, opts ...grpc.CallOption) (*ListSourcesResponse, error) {
	out := new(ListSourcesResponse)
	err := c.cc.Invoke(ctx, "/grit.v2.api.API/ListSources", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) ResolveLocalRepo(ctx context.Context, in *ResolveLocalRepoRequest, opts ...grpc.CallOption) (API_ResolveLocalRepoClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[0], "/grit.v2.api.API/ResolveLocalRepo", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIResolveLocalRepoClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_ResolveLocalRepoClient interface {
	Recv() (*ResolveLocalRepoResponse, error)
	grpc.ClientStream
}

type aPIResolveLocalRepoClient struct {
	grpc.ClientStream
}

func (x *aPIResolveLocalRepoClient) Recv() (*ResolveLocalRepoResponse, error) {
	m := new(ResolveLocalRepoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) ResolveRemoteRepo(ctx context.Context, in *ResolveRemoteRepoRequest, opts ...grpc.CallOption) (API_ResolveRemoteRepoClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[1], "/grit.v2.api.API/ResolveRemoteRepo", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIResolveRemoteRepoClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_ResolveRemoteRepoClient interface {
	Recv() (*ResolveRemoteRepoResponse, error)
	grpc.ClientStream
}

type aPIResolveRemoteRepoClient struct {
	grpc.ClientStream
}

func (x *aPIResolveRemoteRepoClient) Recv() (*ResolveRemoteRepoResponse, error) {
	m := new(ResolveRemoteRepoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) CloneRemoteRepo(ctx context.Context, in *CloneRemoteRepoRequest, opts ...grpc.CallOption) (API_CloneRemoteRepoClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[2], "/grit.v2.api.API/CloneRemoteRepo", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPICloneRemoteRepoClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_CloneRemoteRepoClient interface {
	Recv() (*CloneRemoteRepoResponse, error)
	grpc.ClientStream
}

type aPICloneRemoteRepoClient struct {
	grpc.ClientStream
}

func (x *aPICloneRemoteRepoClient) Recv() (*CloneRemoteRepoResponse, error) {
	m := new(CloneRemoteRepoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SuggestRepo(ctx context.Context, in *SuggestRepoRequest, opts ...grpc.CallOption) (*SuggestResponse, error) {
	out := new(SuggestResponse)
	err := c.cc.Invoke(ctx, "/grit.v2.api.API/SuggestRepo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// APIServer is the server API for API service.
// All implementations should embed UnimplementedAPIServer
// for forward compatibility
type APIServer interface {
	// ListSources lists the configured repository sources.
	ListSources(context.Context, *ListSourcesRequest) (*ListSourcesResponse, error)
	// ResolveLocalRepo resolves repository name, URL or other identifier to a
	// list if local repository clones.
	ResolveLocalRepo(*ResolveLocalRepoRequest, API_ResolveLocalRepoServer) error
	// ResolveRemoteRepo resolves repository name, URL or other identifier to a
	// list of remote repositories.
	ResolveRemoteRepo(*ResolveRemoteRepoRequest, API_ResolveRemoteRepoServer) error
	// CloneRemoteRepo makes a local clone of a repository from a source.
	CloneRemoteRepo(*CloneRemoteRepoRequest, API_CloneRemoteRepoServer) error
	// SuggestRepo returns a list of repository names to be used as
	// suggestions for completing a partial repository name.
	SuggestRepo(context.Context, *SuggestRepoRequest) (*SuggestResponse, error)
}

// UnimplementedAPIServer should be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (UnimplementedAPIServer) ListSources(context.Context, *ListSourcesRequest) (*ListSourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSources not implemented")
}
func (UnimplementedAPIServer) ResolveLocalRepo(*ResolveLocalRepoRequest, API_ResolveLocalRepoServer) error {
	return status.Errorf(codes.Unimplemented, "method ResolveLocalRepo not implemented")
}
func (UnimplementedAPIServer) ResolveRemoteRepo(*ResolveRemoteRepoRequest, API_ResolveRemoteRepoServer) error {
	return status.Errorf(codes.Unimplemented, "method ResolveRemoteRepo not implemented")
}
func (UnimplementedAPIServer) CloneRemoteRepo(*CloneRemoteRepoRequest, API_CloneRemoteRepoServer) error {
	return status.Errorf(codes.Unimplemented, "method CloneRemoteRepo not implemented")
}
func (UnimplementedAPIServer) SuggestRepo(context.Context, *SuggestRepoRequest) (*SuggestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuggestRepo not implemented")
}

// UnsafeAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIServer will
// result in compilation errors.
type UnsafeAPIServer interface {
	mustEmbedUnimplementedAPIServer()
}

func RegisterAPIServer(s grpc.ServiceRegistrar, srv APIServer) {
	s.RegisterService(&API_ServiceDesc, srv)
}

func _API_ListSources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).ListSources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grit.v2.api.API/ListSources",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).ListSources(ctx, req.(*ListSourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_ResolveLocalRepo_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ResolveLocalRepoRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).ResolveLocalRepo(m, &aPIResolveLocalRepoServer{stream})
}

type API_ResolveLocalRepoServer interface {
	Send(*ResolveLocalRepoResponse) error
	grpc.ServerStream
}

type aPIResolveLocalRepoServer struct {
	grpc.ServerStream
}

func (x *aPIResolveLocalRepoServer) Send(m *ResolveLocalRepoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _API_ResolveRemoteRepo_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ResolveRemoteRepoRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).ResolveRemoteRepo(m, &aPIResolveRemoteRepoServer{stream})
}

type API_ResolveRemoteRepoServer interface {
	Send(*ResolveRemoteRepoResponse) error
	grpc.ServerStream
}

type aPIResolveRemoteRepoServer struct {
	grpc.ServerStream
}

func (x *aPIResolveRemoteRepoServer) Send(m *ResolveRemoteRepoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _API_CloneRemoteRepo_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CloneRemoteRepoRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).CloneRemoteRepo(m, &aPICloneRemoteRepoServer{stream})
}

type API_CloneRemoteRepoServer interface {
	Send(*CloneRemoteRepoResponse) error
	grpc.ServerStream
}

type aPICloneRemoteRepoServer struct {
	grpc.ServerStream
}

func (x *aPICloneRemoteRepoServer) Send(m *CloneRemoteRepoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SuggestRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SuggestRepoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).SuggestRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grit.v2.api.API/SuggestRepo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).SuggestRepo(ctx, req.(*SuggestRepoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// API_ServiceDesc is the grpc.ServiceDesc for API service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var API_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grit.v2.api.API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSources",
			Handler:    _API_ListSources_Handler,
		},
		{
			MethodName: "SuggestRepo",
			Handler:    _API_SuggestRepo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ResolveLocalRepo",
			Handler:       _API_ResolveLocalRepo_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ResolveRemoteRepo",
			Handler:       _API_ResolveRemoteRepo_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "CloneRemoteRepo",
			Handler:       _API_CloneRemoteRepo_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/gritcli/grit/api/api.proto",
}

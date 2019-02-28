// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: project.proto

package proto

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/googleapis/google/api"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"
import _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
import _ "github.com/mwitkow/go-proto-validators"
import types "minibox.ai/minibox/pkg/api/v1/types"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProjectServiceClient is the client API for ProjectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProjectServiceClient interface {
	ListProjects(ctx context.Context, in *types.ListProjectsRequest, opts ...grpc.CallOption) (*types.ProjectsReply, error)
	CreateProject(ctx context.Context, in *types.CreateProjectRequest, opts ...grpc.CallOption) (*types.Project, error)
	GetProject(ctx context.Context, in *types.GetProjectRequest, opts ...grpc.CallOption) (*types.Project, error)
	DeleteProject(ctx context.Context, in *types.DeleteProjectRequest, opts ...grpc.CallOption) (*types.Project, error)
	UpdateProject(ctx context.Context, in *types.UpdateProjectRequest, opts ...grpc.CallOption) (*types.Project, error)
}

type projectServiceClient struct {
	cc *grpc.ClientConn
}

func NewProjectServiceClient(cc *grpc.ClientConn) ProjectServiceClient {
	return &projectServiceClient{cc}
}

func (c *projectServiceClient) ListProjects(ctx context.Context, in *types.ListProjectsRequest, opts ...grpc.CallOption) (*types.ProjectsReply, error) {
	out := new(types.ProjectsReply)
	err := c.cc.Invoke(ctx, "/proto.ProjectService/ListProjects", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) CreateProject(ctx context.Context, in *types.CreateProjectRequest, opts ...grpc.CallOption) (*types.Project, error) {
	out := new(types.Project)
	err := c.cc.Invoke(ctx, "/proto.ProjectService/CreateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetProject(ctx context.Context, in *types.GetProjectRequest, opts ...grpc.CallOption) (*types.Project, error) {
	out := new(types.Project)
	err := c.cc.Invoke(ctx, "/proto.ProjectService/GetProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) DeleteProject(ctx context.Context, in *types.DeleteProjectRequest, opts ...grpc.CallOption) (*types.Project, error) {
	out := new(types.Project)
	err := c.cc.Invoke(ctx, "/proto.ProjectService/DeleteProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) UpdateProject(ctx context.Context, in *types.UpdateProjectRequest, opts ...grpc.CallOption) (*types.Project, error) {
	out := new(types.Project)
	err := c.cc.Invoke(ctx, "/proto.ProjectService/UpdateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectServiceServer is the server API for ProjectService service.
type ProjectServiceServer interface {
	ListProjects(context.Context, *types.ListProjectsRequest) (*types.ProjectsReply, error)
	CreateProject(context.Context, *types.CreateProjectRequest) (*types.Project, error)
	GetProject(context.Context, *types.GetProjectRequest) (*types.Project, error)
	DeleteProject(context.Context, *types.DeleteProjectRequest) (*types.Project, error)
	UpdateProject(context.Context, *types.UpdateProjectRequest) (*types.Project, error)
}

func RegisterProjectServiceServer(s *grpc.Server, srv ProjectServiceServer) {
	s.RegisterService(&_ProjectService_serviceDesc, srv)
}

func _ProjectService_ListProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.ListProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).ListProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ProjectService/ListProjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).ListProjects(ctx, req.(*types.ListProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.CreateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ProjectService/CreateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).CreateProject(ctx, req.(*types.CreateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.GetProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ProjectService/GetProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetProject(ctx, req.(*types.GetProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.DeleteProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ProjectService/DeleteProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).DeleteProject(ctx, req.(*types.DeleteProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_UpdateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.UpdateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).UpdateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ProjectService/UpdateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).UpdateProject(ctx, req.(*types.UpdateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProjectService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ProjectService",
	HandlerType: (*ProjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListProjects",
			Handler:    _ProjectService_ListProjects_Handler,
		},
		{
			MethodName: "CreateProject",
			Handler:    _ProjectService_CreateProject_Handler,
		},
		{
			MethodName: "GetProject",
			Handler:    _ProjectService_GetProject_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _ProjectService_DeleteProject_Handler,
		},
		{
			MethodName: "UpdateProject",
			Handler:    _ProjectService_UpdateProject_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "project.proto",
}

func init() { proto.RegisterFile("project.proto", fileDescriptor_project_d38e1cf70f92a2d8) }
func init() { golang_proto.RegisterFile("project.proto", fileDescriptor_project_d38e1cf70f92a2d8) }

var fileDescriptor_project_d38e1cf70f92a2d8 = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x4f, 0x6b, 0x13, 0x41,
	0x14, 0x67, 0x5b, 0xea, 0x61, 0x34, 0x45, 0x06, 0x95, 0x90, 0xe8, 0x36, 0x04, 0x04, 0x09, 0xee,
	0x8e, 0x8d, 0xe0, 0x21, 0x07, 0xa1, 0x2a, 0x88, 0xa0, 0x50, 0x94, 0x1e, 0x14, 0xa4, 0x4c, 0x76,
	0x5f, 0xa6, 0x63, 0x77, 0xf7, 0x8d, 0x33, 0x2f, 0x89, 0x41, 0xbc, 0xf8, 0x11, 0xf4, 0x0b, 0x79,
	0xec, 0x51, 0xf0, 0x0b, 0x48, 0xea, 0x37, 0xf0, 0x0b, 0x88, 0x93, 0x89, 0x69, 0xd2, 0xd6, 0x9e,
	0xf6, 0xbd, 0xdf, 0xdf, 0xe1, 0xb1, 0xac, 0x66, 0x2c, 0xbe, 0x83, 0x8c, 0x52, 0x63, 0x91, 0x90,
	0x6f, 0xf8, 0x4f, 0xa3, 0xa9, 0x10, 0x55, 0x01, 0xc2, 0x6f, 0xfd, 0xe1, 0x40, 0x40, 0x69, 0x68,
	0x32, 0xd3, 0x34, 0xb6, 0x56, 0x49, 0xd2, 0x25, 0x38, 0x92, 0xa5, 0x09, 0x82, 0x78, 0x55, 0x90,
	0x0f, 0xad, 0x24, 0x8d, 0x55, 0xe0, 0x5b, 0xab, 0xfc, 0x40, 0x43, 0x91, 0xef, 0x97, 0xd2, 0x1d,
	0x06, 0xc5, 0xcd, 0xa0, 0x90, 0x46, 0x0b, 0x59, 0x55, 0x48, 0xde, 0xee, 0x02, 0x7b, 0xd7, 0x7f,
	0xb2, 0x44, 0x41, 0x95, 0xb8, 0xb1, 0x54, 0x0a, 0xac, 0x40, 0xe3, 0x15, 0x67, 0xa8, 0x13, 0xa5,
	0xe9, 0x60, 0xd8, 0x4f, 0x33, 0x2c, 0x85, 0x42, 0x85, 0x8b, 0xda, 0xbf, 0x9b, 0x5f, 0xfc, 0x14,
	0xe4, 0x0f, 0x4e, 0xc8, 0xcb, 0xb1, 0xa6, 0x43, 0x1c, 0x0b, 0x85, 0x89, 0x27, 0x93, 0x91, 0x2c,
	0x74, 0x2e, 0x09, 0xad, 0x13, 0xff, 0xc6, 0xe0, 0xbb, 0x4c, 0x13, 0x03, 0xa1, 0xb3, 0xfb, 0x7b,
	0x9d, 0x6d, 0xee, 0xce, 0x0e, 0xfb, 0x0a, 0xec, 0x48, 0x67, 0xc0, 0xdf, 0xb0, 0x2b, 0xcf, 0xb5,
	0xa3, 0x80, 0x3a, 0xde, 0x48, 0x67, 0x86, 0x93, 0xe0, 0x4b, 0x78, 0x3f, 0x04, 0x47, 0x8d, 0x6b,
	0x81, 0x5b, 0xe0, 0xa6, 0x98, 0xb4, 0xeb, 0x9f, 0x7f, 0xfc, 0xfa, 0xba, 0xc6, 0xf9, 0x55, 0x7f,
	0x96, 0xd1, 0xb6, 0x30, 0xf3, 0xac, 0xd7, 0xac, 0xf6, 0xd8, 0x82, 0x24, 0x08, 0x06, 0xde, 0x0c,
	0x01, 0x4b, 0xe8, 0x3c, 0x7d, 0x73, 0x39, 0xbd, 0xdd, 0xf4, 0xb9, 0xd7, 0xdb, 0xa7, 0x72, 0x7b,
	0x51, 0x87, 0xbf, 0x65, 0xec, 0x29, 0xcc, 0x1f, 0xc8, 0xeb, 0xc1, 0xba, 0x80, 0xce, 0x0b, 0xbd,
	0xed, 0x43, 0xb7, 0xf8, 0xad, 0xd5, 0x50, 0xf1, 0x31, 0x4c, 0xfb, 0x3a, 0xff, 0xc4, 0x33, 0x56,
	0x7b, 0x02, 0x05, 0x9c, 0x7e, 0xf9, 0x12, 0x7a, 0x41, 0x49, 0xe7, 0x82, 0x92, 0x01, 0xab, 0xed,
	0x99, 0xfc, 0x8c, 0xf3, 0x2c, 0xa1, 0xe7, 0x95, 0xdc, 0xf1, 0x25, 0xed, 0xc6, 0xff, 0x4b, 0x7a,
	0x51, 0xe7, 0xd1, 0xde, 0x97, 0x9d, 0x87, 0x7c, 0xa3, 0xbb, 0xbe, 0x9d, 0xde, 0xeb, 0x44, 0x6b,
	0xb6, 0xcb, 0x6e, 0xbc, 0xd0, 0x95, 0xee, 0xe3, 0x87, 0xd6, 0xce, 0xee, 0xb3, 0x96, 0x05, 0x83,
	0x4e, 0x13, 0xda, 0x09, 0xaf, 0x1f, 0x10, 0x19, 0xd7, 0x13, 0x22, 0xc7, 0xcc, 0xa5, 0xe5, 0x4c,
	0x94, 0x4a, 0x2d, 0x8e, 0xa6, 0x71, 0xf4, 0x7d, 0x1a, 0x47, 0x3f, 0xa7, 0x71, 0xf4, 0xed, 0x38,
	0x8e, 0x8e, 0x8e, 0xe3, 0xa8, 0x7f, 0xc9, 0xff, 0x53, 0xf7, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff,
	0x34, 0xcb, 0x0b, 0x68, 0xab, 0x03, 0x00, 0x00,
}

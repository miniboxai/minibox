// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dataset.proto

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
import types "minibox.ai/pkg/api/v1/types"

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

// DatasetServiceClient is the client API for DatasetService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DatasetServiceClient interface {
	CreateDataset(ctx context.Context, in *types.CreateDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error)
	ListDatasets(ctx context.Context, in *types.ListDatasetsRequest, opts ...grpc.CallOption) (*types.DatasetsReply, error)
	DeleteDataset(ctx context.Context, in *types.DeleteDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error)
	RestoreDataset(ctx context.Context, in *types.RestoreDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error)
	ListDatasetObjects(ctx context.Context, in *types.ListObjectsRequest, opts ...grpc.CallOption) (*types.ObjectsReply, error)
	GetDatasetObject(ctx context.Context, in *types.GetObjectRequest, opts ...grpc.CallOption) (DatasetService_GetDatasetObjectClient, error)
	PutDatasetObject(ctx context.Context, opts ...grpc.CallOption) (DatasetService_PutDatasetObjectClient, error)
	DeleteDatasetObject(ctx context.Context, in *types.DeleteObjectRequest, opts ...grpc.CallOption) (*types.ObjectReply, error)
}

type datasetServiceClient struct {
	cc *grpc.ClientConn
}

func NewDatasetServiceClient(cc *grpc.ClientConn) DatasetServiceClient {
	return &datasetServiceClient{cc}
}

func (c *datasetServiceClient) CreateDataset(ctx context.Context, in *types.CreateDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error) {
	out := new(types.Dataset)
	err := c.cc.Invoke(ctx, "/proto.DatasetService/CreateDataset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceClient) ListDatasets(ctx context.Context, in *types.ListDatasetsRequest, opts ...grpc.CallOption) (*types.DatasetsReply, error) {
	out := new(types.DatasetsReply)
	err := c.cc.Invoke(ctx, "/proto.DatasetService/ListDatasets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceClient) DeleteDataset(ctx context.Context, in *types.DeleteDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error) {
	out := new(types.Dataset)
	err := c.cc.Invoke(ctx, "/proto.DatasetService/DeleteDataset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceClient) RestoreDataset(ctx context.Context, in *types.RestoreDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error) {
	out := new(types.Dataset)
	err := c.cc.Invoke(ctx, "/proto.DatasetService/RestoreDataset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceClient) ListDatasetObjects(ctx context.Context, in *types.ListObjectsRequest, opts ...grpc.CallOption) (*types.ObjectsReply, error) {
	out := new(types.ObjectsReply)
	err := c.cc.Invoke(ctx, "/proto.DatasetService/ListDatasetObjects", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datasetServiceClient) GetDatasetObject(ctx context.Context, in *types.GetObjectRequest, opts ...grpc.CallOption) (DatasetService_GetDatasetObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DatasetService_serviceDesc.Streams[0], "/proto.DatasetService/GetDatasetObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &datasetServiceGetDatasetObjectClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DatasetService_GetDatasetObjectClient interface {
	Recv() (*types.GetObjectReply, error)
	grpc.ClientStream
}

type datasetServiceGetDatasetObjectClient struct {
	grpc.ClientStream
}

func (x *datasetServiceGetDatasetObjectClient) Recv() (*types.GetObjectReply, error) {
	m := new(types.GetObjectReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *datasetServiceClient) PutDatasetObject(ctx context.Context, opts ...grpc.CallOption) (DatasetService_PutDatasetObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DatasetService_serviceDesc.Streams[1], "/proto.DatasetService/PutDatasetObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &datasetServicePutDatasetObjectClient{stream}
	return x, nil
}

type DatasetService_PutDatasetObjectClient interface {
	Send(*types.PutObjectRequest) error
	CloseAndRecv() (*types.PutObjectReply, error)
	grpc.ClientStream
}

type datasetServicePutDatasetObjectClient struct {
	grpc.ClientStream
}

func (x *datasetServicePutDatasetObjectClient) Send(m *types.PutObjectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *datasetServicePutDatasetObjectClient) CloseAndRecv() (*types.PutObjectReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(types.PutObjectReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *datasetServiceClient) DeleteDatasetObject(ctx context.Context, in *types.DeleteObjectRequest, opts ...grpc.CallOption) (*types.ObjectReply, error) {
	out := new(types.ObjectReply)
	err := c.cc.Invoke(ctx, "/proto.DatasetService/DeleteDatasetObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatasetServiceServer is the server API for DatasetService service.
type DatasetServiceServer interface {
	CreateDataset(context.Context, *types.CreateDatasetRequest) (*types.Dataset, error)
	ListDatasets(context.Context, *types.ListDatasetsRequest) (*types.DatasetsReply, error)
	DeleteDataset(context.Context, *types.DeleteDatasetRequest) (*types.Dataset, error)
	RestoreDataset(context.Context, *types.RestoreDatasetRequest) (*types.Dataset, error)
	ListDatasetObjects(context.Context, *types.ListObjectsRequest) (*types.ObjectsReply, error)
	GetDatasetObject(*types.GetObjectRequest, DatasetService_GetDatasetObjectServer) error
	PutDatasetObject(DatasetService_PutDatasetObjectServer) error
	DeleteDatasetObject(context.Context, *types.DeleteObjectRequest) (*types.ObjectReply, error)
}

func RegisterDatasetServiceServer(s *grpc.Server, srv DatasetServiceServer) {
	s.RegisterService(&_DatasetService_serviceDesc, srv)
}

func _DatasetService_CreateDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.CreateDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceServer).CreateDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DatasetService/CreateDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceServer).CreateDataset(ctx, req.(*types.CreateDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetService_ListDatasets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.ListDatasetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceServer).ListDatasets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DatasetService/ListDatasets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceServer).ListDatasets(ctx, req.(*types.ListDatasetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetService_DeleteDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.DeleteDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceServer).DeleteDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DatasetService/DeleteDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceServer).DeleteDataset(ctx, req.(*types.DeleteDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetService_RestoreDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.RestoreDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceServer).RestoreDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DatasetService/RestoreDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceServer).RestoreDataset(ctx, req.(*types.RestoreDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetService_ListDatasetObjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.ListObjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceServer).ListDatasetObjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DatasetService/ListDatasetObjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceServer).ListDatasetObjects(ctx, req.(*types.ListObjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatasetService_GetDatasetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(types.GetObjectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DatasetServiceServer).GetDatasetObject(m, &datasetServiceGetDatasetObjectServer{stream})
}

type DatasetService_GetDatasetObjectServer interface {
	Send(*types.GetObjectReply) error
	grpc.ServerStream
}

type datasetServiceGetDatasetObjectServer struct {
	grpc.ServerStream
}

func (x *datasetServiceGetDatasetObjectServer) Send(m *types.GetObjectReply) error {
	return x.ServerStream.SendMsg(m)
}

func _DatasetService_PutDatasetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DatasetServiceServer).PutDatasetObject(&datasetServicePutDatasetObjectServer{stream})
}

type DatasetService_PutDatasetObjectServer interface {
	SendAndClose(*types.PutObjectReply) error
	Recv() (*types.PutObjectRequest, error)
	grpc.ServerStream
}

type datasetServicePutDatasetObjectServer struct {
	grpc.ServerStream
}

func (x *datasetServicePutDatasetObjectServer) SendAndClose(m *types.PutObjectReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *datasetServicePutDatasetObjectServer) Recv() (*types.PutObjectRequest, error) {
	m := new(types.PutObjectRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DatasetService_DeleteDatasetObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.DeleteObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatasetServiceServer).DeleteDatasetObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DatasetService/DeleteDatasetObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatasetServiceServer).DeleteDatasetObject(ctx, req.(*types.DeleteObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DatasetService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DatasetService",
	HandlerType: (*DatasetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateDataset",
			Handler:    _DatasetService_CreateDataset_Handler,
		},
		{
			MethodName: "ListDatasets",
			Handler:    _DatasetService_ListDatasets_Handler,
		},
		{
			MethodName: "DeleteDataset",
			Handler:    _DatasetService_DeleteDataset_Handler,
		},
		{
			MethodName: "RestoreDataset",
			Handler:    _DatasetService_RestoreDataset_Handler,
		},
		{
			MethodName: "ListDatasetObjects",
			Handler:    _DatasetService_ListDatasetObjects_Handler,
		},
		{
			MethodName: "DeleteDatasetObject",
			Handler:    _DatasetService_DeleteDatasetObject_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetDatasetObject",
			Handler:       _DatasetService_GetDatasetObject_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PutDatasetObject",
			Handler:       _DatasetService_PutDatasetObject_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "dataset.proto",
}

func init() { proto.RegisterFile("dataset.proto", fileDescriptor_dataset_19dad26abfe56544) }
func init() { golang_proto.RegisterFile("dataset.proto", fileDescriptor_dataset_19dad26abfe56544) }

var fileDescriptor_dataset_19dad26abfe56544 = []byte{
	// 570 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x93, 0x4f, 0x6f, 0xd3, 0x3e,
	0x1c, 0xc6, 0x95, 0xfd, 0xb4, 0x1d, 0xf2, 0xa3, 0x55, 0xe5, 0x31, 0xfe, 0xa4, 0x23, 0xab, 0x2a,
	0x26, 0x50, 0x21, 0x71, 0xb7, 0x49, 0x3b, 0x4c, 0x02, 0x69, 0x30, 0x69, 0x02, 0x81, 0xa8, 0xca,
	0x09, 0x0e, 0x4c, 0x4e, 0xf2, 0x5d, 0x66, 0x96, 0xd4, 0xc6, 0x76, 0xda, 0x95, 0xb1, 0x0b, 0x2f,
	0x01, 0xde, 0x10, 0xc7, 0x1d, 0x91, 0xb8, 0x70, 0x44, 0x1d, 0x2f, 0x04, 0xd5, 0x75, 0x4b, 0xd2,
	0xd2, 0x5e, 0x5a, 0xfb, 0x79, 0x1e, 0x3f, 0x9f, 0xe8, 0xeb, 0xc4, 0x2e, 0x45, 0x44, 0x11, 0x09,
	0xca, 0xe7, 0x82, 0x29, 0x86, 0x96, 0xf5, 0x9f, 0x53, 0x8d, 0x19, 0x8b, 0x13, 0xc0, 0x7a, 0x17,
	0x64, 0xc7, 0x18, 0x52, 0xae, 0xfa, 0xa3, 0x8c, 0xb3, 0x31, 0x6d, 0x2a, 0x9a, 0x82, 0x54, 0x24,
	0xe5, 0x26, 0xe0, 0x4e, 0x07, 0xa2, 0x4c, 0x10, 0x45, 0x59, 0xc7, 0xf8, 0xb5, 0x69, 0xff, 0x98,
	0x42, 0x12, 0x1d, 0xa5, 0x44, 0x9e, 0x9a, 0xc4, 0xba, 0x49, 0x10, 0x4e, 0x31, 0xe9, 0x74, 0x98,
	0xd2, 0xc7, 0xa5, 0x71, 0x1f, 0xea, 0xbf, 0xd0, 0x8b, 0xa1, 0xe3, 0xc9, 0x1e, 0x89, 0x63, 0x10,
	0x98, 0x71, 0x9d, 0xf8, 0x47, 0xda, 0x8b, 0xa9, 0x3a, 0xc9, 0x02, 0x3f, 0x64, 0x29, 0x8e, 0x59,
	0xcc, 0xfe, 0x62, 0x87, 0x3b, 0xbd, 0xd1, 0x2b, 0x13, 0xdf, 0xcd, 0xc5, 0xd3, 0x1e, 0x55, 0xa7,
	0xac, 0x87, 0x63, 0xe6, 0x69, 0xd3, 0xeb, 0x92, 0x84, 0x46, 0x44, 0x31, 0x21, 0xf1, 0x64, 0x69,
	0xce, 0xfd, 0xaf, 0xfa, 0x1c, 0x0c, 0x73, 0xfb, 0xe7, 0x8a, 0x5d, 0x3e, 0x18, 0x0d, 0xf6, 0x35,
	0x88, 0x2e, 0x0d, 0x01, 0xbd, 0xb1, 0x4b, 0x4f, 0x05, 0x10, 0x05, 0x46, 0x47, 0x55, 0x7f, 0x74,
	0xa2, 0xa0, 0xb6, 0xe1, 0x43, 0x06, 0x52, 0x39, 0x65, 0x63, 0x1a, 0xb9, 0x5e, 0xfd, 0xfc, 0xe3,
	0xf7, 0xd7, 0xa5, 0xb5, 0x7a, 0x45, 0x8f, 0xa4, 0xbb, 0x85, 0xcd, 0xb5, 0xc9, 0x3d, 0xab, 0x81,
	0xde, 0xda, 0xd7, 0x5e, 0x50, 0xa9, 0x4c, 0x56, 0x22, 0xc7, 0x1c, 0xce, 0x8b, 0xe3, 0xe2, 0xeb,
	0xc5, 0x62, 0xd9, 0x06, 0x9e, 0xf4, 0xeb, 0xb7, 0x74, 0x3d, 0x42, 0x33, 0xf5, 0x28, 0xb4, 0x4b,
	0x07, 0x90, 0xc0, 0xec, 0x63, 0x17, 0xd4, 0x79, 0x8f, 0xbd, 0xa9, 0x7b, 0x37, 0x1a, 0x77, 0xa6,
	0x7b, 0xf1, 0xb9, 0x59, 0x1d, 0xd1, 0xe8, 0x02, 0xa5, 0x76, 0xb9, 0x0d, 0x52, 0x31, 0x31, 0xa1,
	0xac, 0x9b, 0xa2, 0xa2, 0x3c, 0x0f, 0xe3, 0x69, 0xcc, 0x3d, 0x67, 0x73, 0x21, 0x06, 0x8b, 0x51,
	0x19, 0xfa, 0x64, 0xa3, 0xdc, 0x68, 0x5e, 0x05, 0xef, 0x21, 0x54, 0x12, 0xdd, 0xce, 0x4d, 0xcd,
	0x68, 0x63, 0xde, 0xaa, 0xb1, 0x26, 0xf2, 0x70, 0x66, 0xbb, 0x1a, 0xda, 0x44, 0xfe, 0x2c, 0x94,
	0x46, 0x7e, 0x9e, 0x1b, 0x24, 0x2c, 0xc0, 0xe7, 0xc3, 0xdf, 0x47, 0x8d, 0x0b, 0xf4, 0xce, 0xae,
	0x1c, 0x42, 0x11, 0x8e, 0x6e, 0x1a, 0xc0, 0xe1, 0x58, 0x19, 0x93, 0xd7, 0x66, 0x8d, 0x85, 0xf7,
	0xd5, 0xb4, 0x10, 0xb1, 0x2b, 0xad, 0x6c, 0x4e, 0x7f, 0x2b, 0x9b, 0xd3, 0x9f, 0x33, 0x86, 0xfd,
	0x8b, 0x5e, 0xb7, 0xfb, 0x16, 0xfa, 0x68, 0xaf, 0x16, 0xae, 0xdf, 0x50, 0x9c, 0xc2, 0xab, 0x51,
	0x04, 0xa1, 0xc2, 0x08, 0x47, 0x94, 0x1d, 0x4d, 0xf1, 0x1a, 0x0f, 0x16, 0x5f, 0x5b, 0x61, 0x7c,
	0x4f, 0xc8, 0x97, 0xfd, 0xe7, 0x68, 0x79, 0xfb, 0xbf, 0x2d, 0xbf, 0xd9, 0xb0, 0x96, 0xc4, 0x63,
	0xfb, 0xc6, 0x4b, 0xda, 0xa1, 0x01, 0x3b, 0xab, 0xed, 0xb7, 0x9e, 0xd5, 0x04, 0x70, 0x26, 0xa9,
	0x62, 0xa2, 0x8f, 0xee, 0x9e, 0x28, 0xc5, 0xe5, 0x1e, 0xc6, 0xd3, 0xdf, 0x7f, 0x2c, 0x78, 0xe8,
	0xc1, 0x19, 0x49, 0x79, 0x02, 0x97, 0x03, 0xd7, 0xfa, 0x3e, 0x70, 0xad, 0x5f, 0x03, 0xd7, 0xfa,
	0x76, 0xe5, 0x5a, 0x97, 0x57, 0xae, 0x15, 0xac, 0xe8, 0x8f, 0x78, 0xe7, 0x4f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x97, 0x35, 0x96, 0x89, 0x1c, 0x05, 0x00, 0x00,
}
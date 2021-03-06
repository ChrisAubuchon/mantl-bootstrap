// Code generated by protoc-gen-go.
// source: bootstrap.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	bootstrap.proto

It has these top-level messages:
	ConsulConfig
	FileData
	ShutdownMsg
	Response
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto1.ProtoPackageIsVersion1

type ShutdownMsg_Code int32

const (
	ShutdownMsg_Success ShutdownMsg_Code = 0
	ShutdownMsg_Failure ShutdownMsg_Code = 1
)

var ShutdownMsg_Code_name = map[int32]string{
	0: "Success",
	1: "Failure",
}
var ShutdownMsg_Code_value = map[string]int32{
	"Success": 0,
	"Failure": 1,
}

func (x ShutdownMsg_Code) String() string {
	return proto1.EnumName(ShutdownMsg_Code_name, int32(x))
}
func (ShutdownMsg_Code) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

type Response_Code int32

const (
	Response_Success Response_Code = 0
	Response_Failure Response_Code = 1
)

var Response_Code_name = map[int32]string{
	0: "Success",
	1: "Failure",
}
var Response_Code_value = map[string]int32{
	"Success": 0,
	"Failure": 1,
}

func (x Response_Code) String() string {
	return proto1.EnumName(Response_Code_name, int32(x))
}
func (Response_Code) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

type ConsulConfig struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ConsulConfig) Reset()                    { *m = ConsulConfig{} }
func (m *ConsulConfig) String() string            { return proto1.CompactTextString(m) }
func (*ConsulConfig) ProtoMessage()               {}
func (*ConsulConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type FileData struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	Mode uint32 `protobuf:"varint,3,opt,name=mode" json:"mode,omitempty"`
}

func (m *FileData) Reset()                    { *m = FileData{} }
func (m *FileData) String() string            { return proto1.CompactTextString(m) }
func (*FileData) ProtoMessage()               {}
func (*FileData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ShutdownMsg struct {
	Code ShutdownMsg_Code `protobuf:"varint,1,opt,name=code,enum=proto.ShutdownMsg_Code" json:"code,omitempty"`
}

func (m *ShutdownMsg) Reset()                    { *m = ShutdownMsg{} }
func (m *ShutdownMsg) String() string            { return proto1.CompactTextString(m) }
func (*ShutdownMsg) ProtoMessage()               {}
func (*ShutdownMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Response struct {
	Code Response_Code `protobuf:"varint,1,opt,name=code,enum=proto.Response_Code" json:"code,omitempty"`
	Mesg string        `protobuf:"bytes,2,opt,name=mesg" json:"mesg,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto1.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto1.RegisterType((*ConsulConfig)(nil), "proto.ConsulConfig")
	proto1.RegisterType((*FileData)(nil), "proto.FileData")
	proto1.RegisterType((*ShutdownMsg)(nil), "proto.ShutdownMsg")
	proto1.RegisterType((*Response)(nil), "proto.Response")
	proto1.RegisterEnum("proto.ShutdownMsg_Code", ShutdownMsg_Code_name, ShutdownMsg_Code_value)
	proto1.RegisterEnum("proto.Response_Code", Response_Code_name, Response_Code_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for BootstrapRPC service

type BootstrapRPCClient interface {
	ConfigureConsul(ctx context.Context, in *ConsulConfig, opts ...grpc.CallOption) (*Response, error)
	WriteFile(ctx context.Context, in *FileData, opts ...grpc.CallOption) (*Response, error)
	Shutdown(ctx context.Context, in *ShutdownMsg, opts ...grpc.CallOption) (*Response, error)
}

type bootstrapRPCClient struct {
	cc *grpc.ClientConn
}

func NewBootstrapRPCClient(cc *grpc.ClientConn) BootstrapRPCClient {
	return &bootstrapRPCClient{cc}
}

func (c *bootstrapRPCClient) ConfigureConsul(ctx context.Context, in *ConsulConfig, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/proto.BootstrapRPC/ConfigureConsul", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bootstrapRPCClient) WriteFile(ctx context.Context, in *FileData, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/proto.BootstrapRPC/WriteFile", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bootstrapRPCClient) Shutdown(ctx context.Context, in *ShutdownMsg, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/proto.BootstrapRPC/Shutdown", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for BootstrapRPC service

type BootstrapRPCServer interface {
	ConfigureConsul(context.Context, *ConsulConfig) (*Response, error)
	WriteFile(context.Context, *FileData) (*Response, error)
	Shutdown(context.Context, *ShutdownMsg) (*Response, error)
}

func RegisterBootstrapRPCServer(s *grpc.Server, srv BootstrapRPCServer) {
	s.RegisterService(&_BootstrapRPC_serviceDesc, srv)
}

func _BootstrapRPC_ConfigureConsul_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConsulConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BootstrapRPCServer).ConfigureConsul(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BootstrapRPC/ConfigureConsul",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BootstrapRPCServer).ConfigureConsul(ctx, req.(*ConsulConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _BootstrapRPC_WriteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BootstrapRPCServer).WriteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BootstrapRPC/WriteFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BootstrapRPCServer).WriteFile(ctx, req.(*FileData))
	}
	return interceptor(ctx, in, info, handler)
}

func _BootstrapRPC_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BootstrapRPCServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BootstrapRPC/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BootstrapRPCServer).Shutdown(ctx, req.(*ShutdownMsg))
	}
	return interceptor(ctx, in, info, handler)
}

var _BootstrapRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.BootstrapRPC",
	HandlerType: (*BootstrapRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConfigureConsul",
			Handler:    _BootstrapRPC_ConfigureConsul_Handler,
		},
		{
			MethodName: "WriteFile",
			Handler:    _BootstrapRPC_WriteFile_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _BootstrapRPC_Shutdown_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x90, 0x41, 0x4f, 0x83, 0x40,
	0x10, 0x85, 0x45, 0x51, 0xe9, 0x14, 0xa5, 0x19, 0x4d, 0x24, 0x9e, 0x9a, 0x3d, 0x35, 0x31, 0x62,
	0x52, 0x0f, 0xde, 0xc5, 0x70, 0x33, 0x31, 0xdb, 0x83, 0x17, 0x2f, 0x5b, 0x58, 0x29, 0x06, 0xbb,
	0x84, 0x5d, 0xe2, 0x4f, 0xf2, 0x6f, 0xba, 0xbb, 0x74, 0x13, 0x14, 0x93, 0x5e, 0xe0, 0xcd, 0xbc,
	0x61, 0xe6, 0xf1, 0x41, 0xb4, 0x16, 0x42, 0x49, 0xd5, 0xb2, 0x26, 0x69, 0x5a, 0xa1, 0x04, 0x1e,
	0xdb, 0x17, 0x21, 0x10, 0xa6, 0x62, 0x2b, 0xbb, 0x5a, 0x3f, 0xdf, 0xab, 0x12, 0x11, 0xfc, 0x82,
	0x29, 0x16, 0x7b, 0x73, 0x6f, 0x11, 0x52, 0xab, 0x49, 0x06, 0x41, 0x56, 0xd5, 0xfc, 0x49, 0xeb,
	0xff, 0x7c, 0xd3, 0x6b, 0x98, 0xda, 0xc4, 0x87, 0xba, 0x37, 0xa1, 0x56, 0x9b, 0xde, 0xa7, 0x28,
	0x78, 0x7c, 0xa4, 0x7b, 0x67, 0xd4, 0x6a, 0xf2, 0x06, 0xd3, 0xd5, 0xa6, 0x53, 0x85, 0xf8, 0xda,
	0x3e, 0xcb, 0x12, 0x6f, 0xc0, 0xcf, 0xcd, 0x88, 0x59, 0x75, 0xbe, 0xbc, 0xea, 0x73, 0x25, 0x83,
	0x89, 0x24, 0xd5, 0x36, 0xb5, 0x43, 0x64, 0x0e, 0xbe, 0xa9, 0x70, 0x0a, 0xa7, 0xab, 0x2e, 0xcf,
	0xb9, 0x94, 0xb3, 0x03, 0x53, 0x64, 0xac, 0xaa, 0xbb, 0x96, 0xcf, 0x3c, 0xf2, 0x01, 0x01, 0xe5,
	0xb2, 0xd1, 0x3f, 0xc3, 0x71, 0xf1, 0x6b, 0xf5, 0xe5, 0x6e, 0xb5, 0xb3, 0x07, 0x7b, 0x6d, 0x4e,
	0x2e, 0x4b, 0x97, 0xdd, 0xe8, 0xfd, 0xb7, 0x96, 0xdf, 0x1e, 0x84, 0x8f, 0x0e, 0x28, 0x7d, 0x49,
	0xf1, 0x01, 0xa2, 0x1e, 0xa0, 0xf6, 0x7b, 0x9e, 0x78, 0xb1, 0xbb, 0x3a, 0xc4, 0x7b, 0x1d, 0xfd,
	0x89, 0x82, 0xb7, 0x30, 0x79, 0x6d, 0x2b, 0xc5, 0x0d, 0x60, 0x74, 0xae, 0xa3, 0x3d, 0x1e, 0xbf,
	0x83, 0xc0, 0x01, 0x42, 0x1c, 0x13, 0x1b, 0x7d, 0xb0, 0x3e, 0xb1, 0xf5, 0xfd, 0x4f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x0e, 0xe5, 0x66, 0xc3, 0x00, 0x02, 0x00, 0x00,
}

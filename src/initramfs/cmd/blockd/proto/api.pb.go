// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the process name.
type ResizePartitionRequest struct {
	Number               int32    `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Size                 int64    `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResizePartitionRequest) Reset()         { *m = ResizePartitionRequest{} }
func (m *ResizePartitionRequest) String() string { return proto.CompactTextString(m) }
func (*ResizePartitionRequest) ProtoMessage()    {}
func (*ResizePartitionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *ResizePartitionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResizePartitionRequest.Unmarshal(m, b)
}
func (m *ResizePartitionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResizePartitionRequest.Marshal(b, m, deterministic)
}
func (m *ResizePartitionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResizePartitionRequest.Merge(m, src)
}
func (m *ResizePartitionRequest) XXX_Size() int {
	return xxx_messageInfo_ResizePartitionRequest.Size(m)
}
func (m *ResizePartitionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ResizePartitionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ResizePartitionRequest proto.InternalMessageInfo

func (m *ResizePartitionRequest) GetNumber() int32 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *ResizePartitionRequest) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func init() {
	proto.RegisterType((*ResizePartitionRequest)(nil), "proto.ResizePartitionRequest")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 163 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xd2, 0xe9, 0xf9, 0xf9, 0xe9, 0x39,
	0xa9, 0xfa, 0x60, 0x5e, 0x52, 0x69, 0x9a, 0x7e, 0x6a, 0x6e, 0x41, 0x49, 0x25, 0x44, 0x8d, 0x92,
	0x0b, 0x97, 0x58, 0x50, 0x6a, 0x71, 0x66, 0x55, 0x6a, 0x40, 0x62, 0x51, 0x49, 0x66, 0x49, 0x66,
	0x7e, 0x5e, 0x50, 0x6a, 0x61, 0x69, 0x6a, 0x71, 0x89, 0x90, 0x18, 0x17, 0x5b, 0x5e, 0x69, 0x6e,
	0x52, 0x6a, 0x91, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x6b, 0x10, 0x94, 0x27, 0x24, 0xc4, 0xc5, 0x02,
	0x52, 0x2f, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x1c, 0x04, 0x66, 0x1b, 0x79, 0x73, 0xb1, 0x39, 0xe5,
	0xe4, 0x27, 0x67, 0xa7, 0x08, 0x39, 0x72, 0xb1, 0x41, 0xcc, 0x13, 0x92, 0x85, 0xd8, 0xa0, 0x87,
	0xdd, 0x78, 0x29, 0x31, 0x3d, 0x88, 0xb3, 0xf4, 0x60, 0xce, 0xd2, 0x73, 0x05, 0x39, 0x4b, 0x89,
	0x21, 0x89, 0x0d, 0x2c, 0x62, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xa6, 0xfc, 0x26, 0x27, 0xca,
	0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BlockdClient is the client API for Blockd service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BlockdClient interface {
	Resize(ctx context.Context, in *ResizePartitionRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type blockdClient struct {
	cc *grpc.ClientConn
}

func NewBlockdClient(cc *grpc.ClientConn) BlockdClient {
	return &blockdClient{cc}
}

func (c *blockdClient) Resize(ctx context.Context, in *ResizePartitionRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/proto.Blockd/Resize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockdServer is the server API for Blockd service.
type BlockdServer interface {
	Resize(context.Context, *ResizePartitionRequest) (*empty.Empty, error)
}

func RegisterBlockdServer(s *grpc.Server, srv BlockdServer) {
	s.RegisterService(&_Blockd_serviceDesc, srv)
}

func _Blockd_Resize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResizePartitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockdServer).Resize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Blockd/Resize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockdServer).Resize(ctx, req.(*ResizePartitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Blockd_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Blockd",
	HandlerType: (*BlockdServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Resize",
			Handler:    _Blockd_Resize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

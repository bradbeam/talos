// Code generated by protoc-gen-go. DO NOT EDIT.
// source: network/network.proto

package network

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"

	common "github.com/talos-systems/talos/api/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AddressFamily int32

const (
	AddressFamily_AF_UNSPEC AddressFamily = 0
	AddressFamily_AF_INET   AddressFamily = 2
	AddressFamily_IPV4      AddressFamily = 2
	AddressFamily_AF_INET6  AddressFamily = 10
	AddressFamily_IPV6      AddressFamily = 10
)

var AddressFamily_name = map[int32]string{
	0: "AF_UNSPEC",
	2: "AF_INET",
	// Duplicate value: 2: "IPV4",
	10: "AF_INET6",
	// Duplicate value: 10: "IPV6",
}

var AddressFamily_value = map[string]int32{
	"AF_UNSPEC": 0,
	"AF_INET":   2,
	"IPV4":      2,
	"AF_INET6":  10,
	"IPV6":      10,
}

func (x AddressFamily) String() string {
	return proto.EnumName(AddressFamily_name, int32(x))
}

func (AddressFamily) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{0}
}

type RouteProtocol int32

const (
	RouteProtocol_RTPROT_UNSPEC   RouteProtocol = 0
	RouteProtocol_RTPROT_REDIRECT RouteProtocol = 1
	RouteProtocol_RTPROT_KERNEL   RouteProtocol = 2
	RouteProtocol_RTPROT_BOOT     RouteProtocol = 3
	RouteProtocol_RTPROT_STATIC   RouteProtocol = 4
	RouteProtocol_RTPROT_GATED    RouteProtocol = 8
	RouteProtocol_RTPROT_RA       RouteProtocol = 9
	RouteProtocol_RTPROT_MRT      RouteProtocol = 10
	RouteProtocol_RTPROT_ZEBRA    RouteProtocol = 11
	RouteProtocol_RTPROT_BIRD     RouteProtocol = 12
	RouteProtocol_RTPROT_DNROUTED RouteProtocol = 13
	RouteProtocol_RTPROT_XORP     RouteProtocol = 14
	RouteProtocol_RTPROT_NTK      RouteProtocol = 15
	RouteProtocol_RTPROT_DHCP     RouteProtocol = 16
	RouteProtocol_RTPROT_MROUTED  RouteProtocol = 17
	RouteProtocol_RTPROT_BABEL    RouteProtocol = 42
)

var RouteProtocol_name = map[int32]string{
	0:  "RTPROT_UNSPEC",
	1:  "RTPROT_REDIRECT",
	2:  "RTPROT_KERNEL",
	3:  "RTPROT_BOOT",
	4:  "RTPROT_STATIC",
	8:  "RTPROT_GATED",
	9:  "RTPROT_RA",
	10: "RTPROT_MRT",
	11: "RTPROT_ZEBRA",
	12: "RTPROT_BIRD",
	13: "RTPROT_DNROUTED",
	14: "RTPROT_XORP",
	15: "RTPROT_NTK",
	16: "RTPROT_DHCP",
	17: "RTPROT_MROUTED",
	42: "RTPROT_BABEL",
}

var RouteProtocol_value = map[string]int32{
	"RTPROT_UNSPEC":   0,
	"RTPROT_REDIRECT": 1,
	"RTPROT_KERNEL":   2,
	"RTPROT_BOOT":     3,
	"RTPROT_STATIC":   4,
	"RTPROT_GATED":    8,
	"RTPROT_RA":       9,
	"RTPROT_MRT":      10,
	"RTPROT_ZEBRA":    11,
	"RTPROT_BIRD":     12,
	"RTPROT_DNROUTED": 13,
	"RTPROT_XORP":     14,
	"RTPROT_NTK":      15,
	"RTPROT_DHCP":     16,
	"RTPROT_MROUTED":  17,
	"RTPROT_BABEL":    42,
}

func (x RouteProtocol) String() string {
	return proto.EnumName(RouteProtocol_name, int32(x))
}

func (RouteProtocol) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{1}
}

type InterfaceFlags int32

const (
	InterfaceFlags_FLAG_UNKNOWN        InterfaceFlags = 0
	InterfaceFlags_FLAG_UP             InterfaceFlags = 1
	InterfaceFlags_FLAG_BROADCAST      InterfaceFlags = 2
	InterfaceFlags_FLAG_LOOPBACK       InterfaceFlags = 3
	InterfaceFlags_FLAG_POINT_TO_POINT InterfaceFlags = 4
	InterfaceFlags_FLAG_MULTICAST      InterfaceFlags = 5
)

var InterfaceFlags_name = map[int32]string{
	0: "FLAG_UNKNOWN",
	1: "FLAG_UP",
	2: "FLAG_BROADCAST",
	3: "FLAG_LOOPBACK",
	4: "FLAG_POINT_TO_POINT",
	5: "FLAG_MULTICAST",
}

var InterfaceFlags_value = map[string]int32{
	"FLAG_UNKNOWN":        0,
	"FLAG_UP":             1,
	"FLAG_BROADCAST":      2,
	"FLAG_LOOPBACK":       3,
	"FLAG_POINT_TO_POINT": 4,
	"FLAG_MULTICAST":      5,
}

func (x InterfaceFlags) String() string {
	return proto.EnumName(InterfaceFlags_name, int32(x))
}

func (InterfaceFlags) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{2}
}

// The response message containing the routes.
type RoutesReply struct {
	Response             []*RoutesResponse `protobuf:"bytes,1,rep,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RoutesReply) Reset()         { *m = RoutesReply{} }
func (m *RoutesReply) String() string { return proto.CompactTextString(m) }
func (*RoutesReply) ProtoMessage()    {}
func (*RoutesReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{0}
}

func (m *RoutesReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoutesReply.Unmarshal(m, b)
}

func (m *RoutesReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoutesReply.Marshal(b, m, deterministic)
}

func (m *RoutesReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoutesReply.Merge(m, src)
}

func (m *RoutesReply) XXX_Size() int {
	return xxx_messageInfo_RoutesReply.Size(m)
}

func (m *RoutesReply) XXX_DiscardUnknown() {
	xxx_messageInfo_RoutesReply.DiscardUnknown(m)
}

var xxx_messageInfo_RoutesReply proto.InternalMessageInfo

func (m *RoutesReply) GetResponse() []*RoutesResponse {
	if m != nil {
		return m.Response
	}
	return nil
}

type RoutesResponse struct {
	Metadata             *common.NodeMetadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Routes               []*Route             `protobuf:"bytes,2,rep,name=routes,proto3" json:"routes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *RoutesResponse) Reset()         { *m = RoutesResponse{} }
func (m *RoutesResponse) String() string { return proto.CompactTextString(m) }
func (*RoutesResponse) ProtoMessage()    {}
func (*RoutesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{1}
}

func (m *RoutesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoutesResponse.Unmarshal(m, b)
}

func (m *RoutesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoutesResponse.Marshal(b, m, deterministic)
}

func (m *RoutesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoutesResponse.Merge(m, src)
}

func (m *RoutesResponse) XXX_Size() int {
	return xxx_messageInfo_RoutesResponse.Size(m)
}

func (m *RoutesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RoutesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RoutesResponse proto.InternalMessageInfo

func (m *RoutesResponse) GetMetadata() *common.NodeMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *RoutesResponse) GetRoutes() []*Route {
	if m != nil {
		return m.Routes
	}
	return nil
}

// The response message containing a route.
type Route struct {
	// Interface is the interface over which traffic to this destination should be sent
	Interface string `protobuf:"bytes,1,opt,name=interface,proto3" json:"interface,omitempty"`
	// Destination is the network prefix CIDR which this route provides
	Destination string `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	// Gateway is the gateway address to which traffic to this destination should be sent
	Gateway string `protobuf:"bytes,3,opt,name=gateway,proto3" json:"gateway,omitempty"`
	// Metric is the priority of the route, where lower metrics have higher priorities
	Metric uint32 `protobuf:"varint,4,opt,name=metric,proto3" json:"metric,omitempty"`
	// Scope desribes the scope of this route
	Scope uint32 `protobuf:"varint,5,opt,name=scope,proto3" json:"scope,omitempty"`
	// Source is the source prefix CIDR for the route, if one is defined
	Source string `protobuf:"bytes,6,opt,name=source,proto3" json:"source,omitempty"`
	// Family is the address family of the route.  Currently, the only options are AF_INET (IPV4) and AF_INET6 (IPV6).
	Family AddressFamily `protobuf:"varint,7,opt,name=family,proto3,enum=network.AddressFamily" json:"family,omitempty"`
	// Protocol is the protocol by which this route came to be in place
	Protocol RouteProtocol `protobuf:"varint,8,opt,name=protocol,proto3,enum=network.RouteProtocol" json:"protocol,omitempty"`
	// Flags indicate any special flags on the route
	Flags                uint32   `protobuf:"varint,9,opt,name=flags,proto3" json:"flags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Route) Reset()         { *m = Route{} }
func (m *Route) String() string { return proto.CompactTextString(m) }
func (*Route) ProtoMessage()    {}
func (*Route) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{2}
}

func (m *Route) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Route.Unmarshal(m, b)
}

func (m *Route) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Route.Marshal(b, m, deterministic)
}

func (m *Route) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Route.Merge(m, src)
}

func (m *Route) XXX_Size() int {
	return xxx_messageInfo_Route.Size(m)
}

func (m *Route) XXX_DiscardUnknown() {
	xxx_messageInfo_Route.DiscardUnknown(m)
}

var xxx_messageInfo_Route proto.InternalMessageInfo

func (m *Route) GetInterface() string {
	if m != nil {
		return m.Interface
	}
	return ""
}

func (m *Route) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *Route) GetGateway() string {
	if m != nil {
		return m.Gateway
	}
	return ""
}

func (m *Route) GetMetric() uint32 {
	if m != nil {
		return m.Metric
	}
	return 0
}

func (m *Route) GetScope() uint32 {
	if m != nil {
		return m.Scope
	}
	return 0
}

func (m *Route) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Route) GetFamily() AddressFamily {
	if m != nil {
		return m.Family
	}
	return AddressFamily_AF_UNSPEC
}

func (m *Route) GetProtocol() RouteProtocol {
	if m != nil {
		return m.Protocol
	}
	return RouteProtocol_RTPROT_UNSPEC
}

func (m *Route) GetFlags() uint32 {
	if m != nil {
		return m.Flags
	}
	return 0
}

type InterfacesReply struct {
	Response             []*InterfacesResponse `protobuf:"bytes,1,rep,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *InterfacesReply) Reset()         { *m = InterfacesReply{} }
func (m *InterfacesReply) String() string { return proto.CompactTextString(m) }
func (*InterfacesReply) ProtoMessage()    {}
func (*InterfacesReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{3}
}

func (m *InterfacesReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InterfacesReply.Unmarshal(m, b)
}

func (m *InterfacesReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InterfacesReply.Marshal(b, m, deterministic)
}

func (m *InterfacesReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InterfacesReply.Merge(m, src)
}

func (m *InterfacesReply) XXX_Size() int {
	return xxx_messageInfo_InterfacesReply.Size(m)
}

func (m *InterfacesReply) XXX_DiscardUnknown() {
	xxx_messageInfo_InterfacesReply.DiscardUnknown(m)
}

var xxx_messageInfo_InterfacesReply proto.InternalMessageInfo

func (m *InterfacesReply) GetResponse() []*InterfacesResponse {
	if m != nil {
		return m.Response
	}
	return nil
}

type InterfacesResponse struct {
	Metadata             *common.NodeMetadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Interfaces           []*Interface         `protobuf:"bytes,2,rep,name=interfaces,proto3" json:"interfaces,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *InterfacesResponse) Reset()         { *m = InterfacesResponse{} }
func (m *InterfacesResponse) String() string { return proto.CompactTextString(m) }
func (*InterfacesResponse) ProtoMessage()    {}
func (*InterfacesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{4}
}

func (m *InterfacesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InterfacesResponse.Unmarshal(m, b)
}

func (m *InterfacesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InterfacesResponse.Marshal(b, m, deterministic)
}

func (m *InterfacesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InterfacesResponse.Merge(m, src)
}

func (m *InterfacesResponse) XXX_Size() int {
	return xxx_messageInfo_InterfacesResponse.Size(m)
}

func (m *InterfacesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InterfacesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InterfacesResponse proto.InternalMessageInfo

func (m *InterfacesResponse) GetMetadata() *common.NodeMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *InterfacesResponse) GetInterfaces() []*Interface {
	if m != nil {
		return m.Interfaces
	}
	return nil
}

// Interface represents a net.Interface
type Interface struct {
	Index                uint32         `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Mtu                  uint32         `protobuf:"varint,2,opt,name=mtu,proto3" json:"mtu,omitempty"`
	Name                 string         `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Hardwareaddr         string         `protobuf:"bytes,4,opt,name=hardwareaddr,proto3" json:"hardwareaddr,omitempty"`
	Flags                InterfaceFlags `protobuf:"varint,5,opt,name=flags,proto3,enum=network.InterfaceFlags" json:"flags,omitempty"`
	Ipaddress            []string       `protobuf:"bytes,6,rep,name=ipaddress,proto3" json:"ipaddress,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Interface) Reset()         { *m = Interface{} }
func (m *Interface) String() string { return proto.CompactTextString(m) }
func (*Interface) ProtoMessage()    {}
func (*Interface) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ad937ae012c472, []int{5}
}

func (m *Interface) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Interface.Unmarshal(m, b)
}

func (m *Interface) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Interface.Marshal(b, m, deterministic)
}

func (m *Interface) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Interface.Merge(m, src)
}

func (m *Interface) XXX_Size() int {
	return xxx_messageInfo_Interface.Size(m)
}

func (m *Interface) XXX_DiscardUnknown() {
	xxx_messageInfo_Interface.DiscardUnknown(m)
}

var xxx_messageInfo_Interface proto.InternalMessageInfo

func (m *Interface) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Interface) GetMtu() uint32 {
	if m != nil {
		return m.Mtu
	}
	return 0
}

func (m *Interface) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Interface) GetHardwareaddr() string {
	if m != nil {
		return m.Hardwareaddr
	}
	return ""
}

func (m *Interface) GetFlags() InterfaceFlags {
	if m != nil {
		return m.Flags
	}
	return InterfaceFlags_FLAG_UNKNOWN
}

func (m *Interface) GetIpaddress() []string {
	if m != nil {
		return m.Ipaddress
	}
	return nil
}

func init() {
	proto.RegisterEnum("network.AddressFamily", AddressFamily_name, AddressFamily_value)
	proto.RegisterEnum("network.RouteProtocol", RouteProtocol_name, RouteProtocol_value)
	proto.RegisterEnum("network.InterfaceFlags", InterfaceFlags_name, InterfaceFlags_value)
	proto.RegisterType((*RoutesReply)(nil), "network.RoutesReply")
	proto.RegisterType((*RoutesResponse)(nil), "network.RoutesResponse")
	proto.RegisterType((*Route)(nil), "network.Route")
	proto.RegisterType((*InterfacesReply)(nil), "network.InterfacesReply")
	proto.RegisterType((*InterfacesResponse)(nil), "network.InterfacesResponse")
	proto.RegisterType((*Interface)(nil), "network.Interface")
}

func init() { proto.RegisterFile("network/network.proto", fileDescriptor_96ad937ae012c472) }

var fileDescriptor_96ad937ae012c472 = []byte{
	// 830 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x4d, 0x6f, 0xeb, 0x44,
	0x14, 0xc5, 0xf9, 0xce, 0x4d, 0x9d, 0x4c, 0xa7, 0xa5, 0xb5, 0xfa, 0x58, 0x44, 0x59, 0xa0, 0x2a,
	0xd2, 0x4b, 0x50, 0x1e, 0x2a, 0x3b, 0x24, 0x3b, 0x71, 0x1e, 0x21, 0xa9, 0x6d, 0xa6, 0x2e, 0xa0,
	0xb7, 0xa9, 0xa6, 0xf1, 0x24, 0xcf, 0x10, 0x7f, 0xc8, 0x76, 0x54, 0xc2, 0x0e, 0xf1, 0x4b, 0xf8,
	0x05, 0x6c, 0xf9, 0x79, 0xc8, 0xe3, 0xb1, 0xeb, 0x10, 0xd8, 0xb0, 0xf2, 0xdc, 0x73, 0xce, 0xdc,
	0x7b, 0x67, 0xce, 0xcc, 0x18, 0x3e, 0xf5, 0x59, 0xf2, 0x12, 0x44, 0x3f, 0x8f, 0xc5, 0x77, 0x14,
	0x46, 0x41, 0x12, 0xe0, 0xa6, 0x08, 0x6f, 0xde, 0x6c, 0x83, 0x60, 0xbb, 0x63, 0x63, 0x0e, 0x3f,
	0xef, 0x37, 0x63, 0xe6, 0x85, 0xc9, 0x21, 0x53, 0xdd, 0x5c, 0xac, 0x03, 0xcf, 0x0b, 0xfc, 0x71,
	0xf6, 0xc9, 0xc0, 0x81, 0x06, 0x1d, 0x12, 0xec, 0x13, 0x16, 0x13, 0x16, 0xee, 0x0e, 0xf8, 0x1d,
	0xb4, 0x22, 0x16, 0x87, 0x81, 0x1f, 0x33, 0x45, 0xea, 0x57, 0x6f, 0x3b, 0x93, 0xeb, 0x51, 0x5e,
	0x2b, 0xd7, 0x65, 0x34, 0x29, 0x84, 0x83, 0x9f, 0xa0, 0x7b, 0xcc, 0xe1, 0x2f, 0xa0, 0xe5, 0xb1,
	0x84, 0x3a, 0x34, 0xa1, 0x8a, 0xd4, 0x97, 0x6e, 0x3b, 0x93, 0xcb, 0x91, 0x28, 0x6b, 0x04, 0x0e,
	0xbb, 0x17, 0x1c, 0x29, 0x54, 0xf8, 0x73, 0x68, 0x44, 0x3c, 0x87, 0x52, 0xe1, 0x65, 0xbb, 0xc7,
	0x65, 0x89, 0x60, 0x07, 0x7f, 0x54, 0xa0, 0xce, 0x11, 0xfc, 0x19, 0xb4, 0x5d, 0x3f, 0x61, 0xd1,
	0x86, 0xae, 0x19, 0x2f, 0xd2, 0x26, 0xaf, 0x00, 0xee, 0x43, 0xc7, 0x61, 0x71, 0xe2, 0xfa, 0x34,
	0x71, 0x03, 0x5f, 0xa9, 0x70, 0xbe, 0x0c, 0x61, 0x05, 0x9a, 0x5b, 0x9a, 0xb0, 0x17, 0x7a, 0x50,
	0xaa, 0x9c, 0xcd, 0x43, 0x7c, 0x05, 0x0d, 0x8f, 0x25, 0x91, 0xbb, 0x56, 0x6a, 0x7d, 0xe9, 0x56,
	0x26, 0x22, 0xc2, 0x97, 0x50, 0x8f, 0xd7, 0x41, 0xc8, 0x94, 0x3a, 0x87, 0xb3, 0x20, 0x55, 0xc7,
	0xc1, 0x3e, 0x5a, 0x33, 0xa5, 0xc1, 0xd3, 0x88, 0x08, 0x8f, 0xa0, 0xb1, 0xa1, 0x9e, 0xbb, 0x3b,
	0x28, 0xcd, 0xbe, 0x74, 0xdb, 0x9d, 0x5c, 0x15, 0x2b, 0x52, 0x1d, 0x27, 0x62, 0x71, 0x3c, 0xe7,
	0x2c, 0x11, 0x2a, 0x3c, 0x81, 0x16, 0xb7, 0x64, 0x1d, 0xec, 0x94, 0xd6, 0x3f, 0x66, 0xf0, 0x15,
	0x5b, 0x82, 0x25, 0x85, 0x2e, 0xed, 0x68, 0xb3, 0xa3, 0xdb, 0x58, 0x69, 0x67, 0x1d, 0xf1, 0x60,
	0xf0, 0x2d, 0xf4, 0x16, 0xf9, 0x46, 0x08, 0x5f, 0xbf, 0x3a, 0xf1, 0xf5, 0x4d, 0x91, 0xbc, 0xac,
	0x3d, 0xf1, 0xf6, 0x57, 0xc0, 0xa7, 0xfc, 0xff, 0xf0, 0x77, 0x02, 0x50, 0x98, 0x93, 0x7b, 0x8c,
	0x4f, 0x5b, 0x20, 0x25, 0xd5, 0xe0, 0x2f, 0x09, 0xda, 0x05, 0x93, 0xae, 0xd5, 0xf5, 0x1d, 0xf6,
	0x0b, 0x2f, 0x28, 0x93, 0x2c, 0xc0, 0x08, 0xaa, 0x5e, 0xb2, 0xe7, 0xfe, 0xca, 0x24, 0x1d, 0x62,
	0x0c, 0x35, 0x9f, 0x7a, 0x4c, 0x98, 0xca, 0xc7, 0x78, 0x00, 0x67, 0x1f, 0x69, 0xe4, 0xbc, 0xd0,
	0x88, 0x51, 0xc7, 0x89, 0xb8, 0xaf, 0x6d, 0x72, 0x84, 0xe1, 0xb7, 0xf9, 0x5e, 0xd6, 0xf9, 0xe6,
	0x5f, 0x9f, 0x36, 0x37, 0x4f, 0x69, 0xb1, 0xc9, 0xfc, 0xf8, 0x85, 0x34, 0x73, 0x52, 0x69, 0xf4,
	0xab, 0xfc, 0xf8, 0xe5, 0xc0, 0xf0, 0x3b, 0x90, 0x8f, 0x5c, 0xc6, 0x32, 0xb4, 0xd5, 0xf9, 0xd3,
	0xa3, 0xf1, 0x60, 0xe9, 0x53, 0xf4, 0x09, 0xee, 0x40, 0x53, 0x9d, 0x3f, 0x2d, 0x0c, 0xdd, 0x46,
	0x15, 0xdc, 0x82, 0xda, 0xc2, 0xfa, 0xfe, 0x4b, 0x54, 0xc1, 0x67, 0xd0, 0x12, 0xf0, 0x1d, 0x02,
	0x81, 0xdf, 0x21, 0xb8, 0xa9, 0x20, 0x69, 0xf8, 0x67, 0x05, 0xe4, 0xa3, 0x73, 0x80, 0xcf, 0x41,
	0x26, 0xb6, 0x45, 0x4c, 0xfb, 0x35, 0xef, 0x05, 0xf4, 0x04, 0x44, 0xf4, 0xd9, 0x82, 0xe8, 0x53,
	0x1b, 0x49, 0x25, 0xdd, 0x52, 0x27, 0x86, 0xbe, 0x42, 0x15, 0xdc, 0x83, 0x8e, 0x80, 0x34, 0xd3,
	0xb4, 0x51, 0xb5, 0xa4, 0x79, 0xb0, 0x55, 0x7b, 0x31, 0x45, 0x35, 0x8c, 0xe0, 0x4c, 0x40, 0xef,
	0x55, 0x5b, 0x9f, 0xa1, 0x56, 0xba, 0x88, 0x3c, 0xbb, 0x8a, 0xda, 0xb8, 0x0b, 0x20, 0xc2, 0x7b,
	0x62, 0x23, 0x28, 0x4d, 0xf8, 0xa0, 0x6b, 0x44, 0x45, 0x9d, 0x72, 0x99, 0x05, 0x99, 0xa1, 0xb3,
	0x52, 0x7f, 0x33, 0x83, 0x98, 0x8f, 0x69, 0x5a, 0xb9, 0xa4, 0xfa, 0xd1, 0x24, 0x16, 0xea, 0x96,
	0x12, 0x1b, 0xf6, 0x12, 0xf5, 0x4a, 0x82, 0xd9, 0x37, 0x53, 0x0b, 0x21, 0x8c, 0xa1, 0x5b, 0x54,
	0xce, 0xb2, 0x9c, 0x97, 0xaa, 0x6b, 0xaa, 0xa6, 0xaf, 0xd0, 0x70, 0xf8, 0xbb, 0x04, 0xdd, 0x63,
	0xf3, 0x52, 0xd1, 0x7c, 0xa5, 0xbe, 0x7f, 0x7a, 0x34, 0x96, 0x86, 0xf9, 0x83, 0x91, 0x39, 0x91,
	0x21, 0x16, 0x92, 0xd2, 0xbc, 0x3c, 0xd0, 0x88, 0xa9, 0xce, 0xa6, 0xea, 0x43, 0xea, 0xce, 0x39,
	0xc8, 0x1c, 0x5b, 0x99, 0xa6, 0xa5, 0xa9, 0xd3, 0x25, 0xaa, 0xe2, 0x6b, 0xb8, 0xe0, 0x90, 0x65,
	0x2e, 0x0c, 0xfb, 0xc9, 0x36, 0xb3, 0x01, 0xaa, 0x15, 0xf3, 0xef, 0x1f, 0x57, 0xf6, 0x82, 0xcf,
	0xaf, 0x4f, 0x7e, 0x93, 0xa0, 0x69, 0x64, 0x47, 0x09, 0xdf, 0x41, 0x23, 0x7b, 0x29, 0xf1, 0xd5,
	0x28, 0x7b, 0xaa, 0x47, 0xf9, 0x53, 0x3d, 0xd2, 0xd3, 0xa7, 0xfa, 0xe6, 0xf2, 0xe4, 0xb9, 0x4d,
	0xaf, 0xef, 0xd7, 0x00, 0xaf, 0xb7, 0xf0, 0x3f, 0xe7, 0x2a, 0xff, 0x7a, 0xa5, 0xc3, 0xdd, 0x41,
	0x5b, 0x42, 0x6f, 0x1d, 0x78, 0x05, 0x4d, 0x43, 0x57, 0x03, 0xd1, 0x93, 0x1a, 0xba, 0x96, 0xf4,
	0x61, 0xb8, 0x75, 0x93, 0x8f, 0xfb, 0xe7, 0xf4, 0x12, 0x8f, 0x13, 0xba, 0x0b, 0xe2, 0xb7, 0xf1,
	0x21, 0x4e, 0x98, 0x17, 0x67, 0xd1, 0x98, 0x86, 0x6e, 0xfe, 0xcf, 0x79, 0x6e, 0xf0, 0xb2, 0xef,
	0xfe, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xf3, 0xb1, 0x46, 0x21, 0x8d, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NetworkClient is the client API for Network service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NetworkClient interface {
	Routes(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*RoutesReply, error)
	Interfaces(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*InterfacesReply, error)
}

type networkClient struct {
	cc *grpc.ClientConn
}

func NewNetworkClient(cc *grpc.ClientConn) NetworkClient {
	return &networkClient{cc}
}

func (c *networkClient) Routes(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*RoutesReply, error) {
	out := new(RoutesReply)
	err := c.cc.Invoke(ctx, "/network.Network/Routes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) Interfaces(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*InterfacesReply, error) {
	out := new(InterfacesReply)
	err := c.cc.Invoke(ctx, "/network.Network/Interfaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkServer is the server API for Network service.
type NetworkServer interface {
	Routes(context.Context, *empty.Empty) (*RoutesReply, error)
	Interfaces(context.Context, *empty.Empty) (*InterfacesReply, error)
}

func RegisterNetworkServer(s *grpc.Server, srv NetworkServer) {
	s.RegisterService(&_Network_serviceDesc, srv)
}

func _Network_Routes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).Routes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/network.Network/Routes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).Routes(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_Interfaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).Interfaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/network.Network/Interfaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).Interfaces(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Network_serviceDesc = grpc.ServiceDesc{
	ServiceName: "network.Network",
	HandlerType: (*NetworkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Routes",
			Handler:    _Network_Routes_Handler,
		},
		{
			MethodName: "Interfaces",
			Handler:    _Network_Interfaces_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "network/network.proto",
}

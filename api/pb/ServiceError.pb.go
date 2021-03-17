// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ServiceError.proto

package casper_v1

import (
	_ "dawn1806/rchain/api/pb/scalapb"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ServiceError struct {
	Messages             []string `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServiceError) Reset()         { *m = ServiceError{} }
func (m *ServiceError) String() string { return proto.CompactTextString(m) }
func (*ServiceError) ProtoMessage()    {}
func (*ServiceError) Descriptor() ([]byte, []int) {
	return fileDescriptor_4208cfdfcb156969, []int{0}
}

func (m *ServiceError) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceError.Unmarshal(m, b)
}
func (m *ServiceError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceError.Marshal(b, m, deterministic)
}
func (m *ServiceError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceError.Merge(m, src)
}
func (m *ServiceError) XXX_Size() int {
	return xxx_messageInfo_ServiceError.Size(m)
}
func (m *ServiceError) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceError.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceError proto.InternalMessageInfo

func (m *ServiceError) GetMessages() []string {
	if m != nil {
		return m.Messages
	}
	return nil
}

func init() {
	proto.RegisterType((*ServiceError)(nil), "ServiceError")
}

func init() { proto.RegisterFile("ServiceError.proto", fileDescriptor_4208cfdfcb156969) }

var fileDescriptor_4208cfdfcb156969 = []byte{
	// 120 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x0a, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0x75, 0x2d, 0x2a, 0xca, 0x2f, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0x12,
	0x2d, 0x4e, 0x4e, 0xcc, 0x49, 0x2c, 0x48, 0xd2, 0x87, 0xd2, 0x10, 0x61, 0x25, 0x2d, 0x2e, 0x1e,
	0x64, 0xc5, 0x42, 0x52, 0x5c, 0x1c, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0xc5, 0x12, 0x8c,
	0x0a, 0xcc, 0x1a, 0x9c, 0x41, 0x70, 0xbe, 0x93, 0xd2, 0x23, 0x7b, 0x79, 0x2e, 0xe9, 0xe4, 0xfc,
	0xfc, 0x02, 0xbd, 0xa2, 0xe4, 0x8c, 0xc4, 0xcc, 0x3c, 0xbd, 0xe4, 0xc4, 0xe2, 0x82, 0x54, 0xa8,
	0x0d, 0xc9, 0xf9, 0x39, 0x1e, 0x0c, 0x49, 0x6c, 0x60, 0xb6, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff,
	0x5b, 0xe4, 0x31, 0xc1, 0x83, 0x00, 0x00, 0x00,
}
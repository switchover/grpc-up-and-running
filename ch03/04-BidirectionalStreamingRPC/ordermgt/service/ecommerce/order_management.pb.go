// Code generated by protoc-gen-go. DO NOT EDIT.
// source: order_management.proto

package ecommerce

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Order struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Items                []string `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Price                float32  `protobuf:"fixed32,4,opt,name=price,proto3" json:"price,omitempty"`
	Destination          string   `protobuf:"bytes,5,opt,name=destination,proto3" json:"destination,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_6653354279552460, []int{0}
}

func (m *Order) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Order.Unmarshal(m, b)
}
func (m *Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Order.Marshal(b, m, deterministic)
}
func (m *Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Order.Merge(m, src)
}
func (m *Order) XXX_Size() int {
	return xxx_messageInfo_Order.Size(m)
}
func (m *Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Order) GetItems() []string {
	if m != nil {
		return m.Items
	}
	return nil
}

func (m *Order) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Order) GetPrice() float32 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *Order) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

type CombinedShipment struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	OrdersList           []*Order `protobuf:"bytes,3,rep,name=ordersList,proto3" json:"ordersList,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CombinedShipment) Reset()         { *m = CombinedShipment{} }
func (m *CombinedShipment) String() string { return proto.CompactTextString(m) }
func (*CombinedShipment) ProtoMessage()    {}
func (*CombinedShipment) Descriptor() ([]byte, []int) {
	return fileDescriptor_6653354279552460, []int{1}
}

func (m *CombinedShipment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CombinedShipment.Unmarshal(m, b)
}
func (m *CombinedShipment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CombinedShipment.Marshal(b, m, deterministic)
}
func (m *CombinedShipment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CombinedShipment.Merge(m, src)
}
func (m *CombinedShipment) XXX_Size() int {
	return xxx_messageInfo_CombinedShipment.Size(m)
}
func (m *CombinedShipment) XXX_DiscardUnknown() {
	xxx_messageInfo_CombinedShipment.DiscardUnknown(m)
}

var xxx_messageInfo_CombinedShipment proto.InternalMessageInfo

func (m *CombinedShipment) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CombinedShipment) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *CombinedShipment) GetOrdersList() []*Order {
	if m != nil {
		return m.OrdersList
	}
	return nil
}

func init() {
	proto.RegisterType((*Order)(nil), "ecommerce.Order")
	proto.RegisterType((*CombinedShipment)(nil), "ecommerce.CombinedShipment")
}

func init() {
	proto.RegisterFile("order_management.proto", fileDescriptor_6653354279552460)
}

var fileDescriptor_6653354279552460 = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xdd, 0x4a, 0xc3, 0x40,
	0x10, 0x85, 0xd9, 0xc4, 0x14, 0xb2, 0x45, 0x2d, 0x8b, 0x94, 0x50, 0x45, 0x42, 0xaf, 0x72, 0xb5,
	0x2d, 0xf5, 0x11, 0xbc, 0xf5, 0x07, 0x52, 0xf0, 0x56, 0x36, 0xc9, 0x18, 0x07, 0xb2, 0x3f, 0xec,
	0x6e, 0xf0, 0x11, 0x7c, 0x6d, 0xc9, 0xa6, 0x0d, 0xc1, 0x5e, 0xce, 0x99, 0x73, 0x86, 0x33, 0x1f,
	0x5d, 0x6b, 0xdb, 0x80, 0xfd, 0x94, 0x42, 0x89, 0x16, 0x24, 0x28, 0xcf, 0x8d, 0xd5, 0x5e, 0xb3,
	0x14, 0x6a, 0x2d, 0x25, 0xd8, 0x1a, 0x36, 0x8f, 0xad, 0xd6, 0x6d, 0x07, 0xbb, 0xb0, 0xa8, 0xfa,
	0xaf, 0xdd, 0x8f, 0x15, 0xc6, 0x80, 0x75, 0xa3, 0x75, 0xfb, 0x4b, 0x68, 0xf2, 0x3e, 0x5c, 0x61,
	0x37, 0x34, 0xc2, 0x26, 0x23, 0x39, 0x29, 0xd2, 0x32, 0xc2, 0x86, 0xdd, 0xd1, 0x04, 0x3d, 0x48,
	0x97, 0x45, 0x79, 0x5c, 0xa4, 0xe5, 0x38, 0xb0, 0x9c, 0x2e, 0x1b, 0x70, 0xb5, 0x45, 0xe3, 0x51,
	0xab, 0x2c, 0x0e, 0xf6, 0xb9, 0x34, 0xe4, 0x8c, 0xc5, 0x1a, 0xb2, 0xab, 0x9c, 0x14, 0x51, 0x39,
	0x0e, 0xa7, 0x9c, 0x47, 0x25, 0x42, 0x2e, 0x99, 0x72, 0x67, 0x69, 0xdb, 0xd1, 0xd5, 0xb3, 0x96,
	0x15, 0x2a, 0x68, 0x8e, 0xdf, 0x68, 0x86, 0x77, 0x2e, 0x3a, 0xad, 0xe9, 0xc2, 0x79, 0xe1, 0xfb,
	0xa1, 0xd4, 0xa0, 0x9d, 0x26, 0xb6, 0xa7, 0x34, 0xa0, 0x70, 0x2f, 0xe8, 0x7c, 0x16, 0xe7, 0x71,
	0xb1, 0x3c, 0xac, 0xf8, 0x44, 0x81, 0x87, 0x0f, 0xcb, 0x99, 0xe7, 0x20, 0xe8, 0x6d, 0x10, 0x5f,
	0x27, 0x76, 0xec, 0x8d, 0x5e, 0x1b, 0xab, 0x6b, 0x70, 0x2e, 0x6c, 0x1c, 0x7b, 0xe0, 0x23, 0x3c,
	0x7e, 0x86, 0xc7, 0x8f, 0xde, 0xa2, 0x6a, 0x3f, 0x44, 0xd7, 0xc3, 0xe6, 0x7e, 0x76, 0xff, 0x7f,
	0xf1, 0x82, 0xec, 0x49, 0xb5, 0x08, 0xb1, 0xa7, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xed, 0xe8,
	0xc2, 0xa6, 0xa6, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// OrderManagementClient is the client API for OrderManagement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrderManagementClient interface {
	// ...
	ProcessOrders(ctx context.Context, opts ...grpc.CallOption) (OrderManagement_ProcessOrdersClient, error)
}

type orderManagementClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderManagementClient(cc grpc.ClientConnInterface) OrderManagementClient {
	return &orderManagementClient{cc}
}

func (c *orderManagementClient) ProcessOrders(ctx context.Context, opts ...grpc.CallOption) (OrderManagement_ProcessOrdersClient, error) {
	stream, err := c.cc.NewStream(ctx, &_OrderManagement_serviceDesc.Streams[0], "/ecommerce.OrderManagement/processOrders", opts...)
	if err != nil {
		return nil, err
	}
	x := &orderManagementProcessOrdersClient{stream}
	return x, nil
}

type OrderManagement_ProcessOrdersClient interface {
	Send(*wrappers.StringValue) error
	Recv() (*CombinedShipment, error)
	grpc.ClientStream
}

type orderManagementProcessOrdersClient struct {
	grpc.ClientStream
}

func (x *orderManagementProcessOrdersClient) Send(m *wrappers.StringValue) error {
	return x.ClientStream.SendMsg(m)
}

func (x *orderManagementProcessOrdersClient) Recv() (*CombinedShipment, error) {
	m := new(CombinedShipment)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OrderManagementServer is the server API for OrderManagement service.
type OrderManagementServer interface {
	// ...
	ProcessOrders(OrderManagement_ProcessOrdersServer) error
}

// UnimplementedOrderManagementServer can be embedded to have forward compatible implementations.
type UnimplementedOrderManagementServer struct {
}

func (*UnimplementedOrderManagementServer) ProcessOrders(srv OrderManagement_ProcessOrdersServer) error {
	return status.Errorf(codes.Unimplemented, "method ProcessOrders not implemented")
}

func RegisterOrderManagementServer(s *grpc.Server, srv OrderManagementServer) {
	s.RegisterService(&_OrderManagement_serviceDesc, srv)
}

func _OrderManagement_ProcessOrders_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OrderManagementServer).ProcessOrders(&orderManagementProcessOrdersServer{stream})
}

type OrderManagement_ProcessOrdersServer interface {
	Send(*CombinedShipment) error
	Recv() (*wrappers.StringValue, error)
	grpc.ServerStream
}

type orderManagementProcessOrdersServer struct {
	grpc.ServerStream
}

func (x *orderManagementProcessOrdersServer) Send(m *CombinedShipment) error {
	return x.ServerStream.SendMsg(m)
}

func (x *orderManagementProcessOrdersServer) Recv() (*wrappers.StringValue, error) {
	m := new(wrappers.StringValue)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _OrderManagement_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ecommerce.OrderManagement",
	HandlerType: (*OrderManagementServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "processOrders",
			Handler:       _OrderManagement_ProcessOrders_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "order_management.proto",
}

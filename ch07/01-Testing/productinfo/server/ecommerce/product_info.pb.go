// Code generated by protoc-gen-go. DO NOT EDIT.
// source: product_info.proto

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

type Product struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Price                float32  `protobuf:"fixed32,4,opt,name=price,proto3" json:"price,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Product) Reset()         { *m = Product{} }
func (m *Product) String() string { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()    {}
func (*Product) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a4d768ec9cb4951, []int{0}
}

func (m *Product) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Product.Unmarshal(m, b)
}
func (m *Product) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Product.Marshal(b, m, deterministic)
}
func (m *Product) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Product.Merge(m, src)
}
func (m *Product) XXX_Size() int {
	return xxx_messageInfo_Product.Size(m)
}
func (m *Product) XXX_DiscardUnknown() {
	xxx_messageInfo_Product.DiscardUnknown(m)
}

var xxx_messageInfo_Product proto.InternalMessageInfo

func (m *Product) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Product) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Product) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Product) GetPrice() float32 {
	if m != nil {
		return m.Price
	}
	return 0
}

func init() {
	proto.RegisterType((*Product)(nil), "ecommerce.Product")
}

func init() {
	proto.RegisterFile("product_info.proto", fileDescriptor_9a4d768ec9cb4951)
}

var fileDescriptor_9a4d768ec9cb4951 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x4e, 0xbd, 0x4a, 0x04, 0x31,
	0x10, 0x66, 0xd7, 0x53, 0xb9, 0x39, 0xb0, 0x18, 0x2c, 0xc2, 0x22, 0xb2, 0x58, 0x5d, 0x95, 0x03,
	0xed, 0xed, 0xed, 0xe4, 0x04, 0x5b, 0xc9, 0x25, 0xb3, 0x21, 0x70, 0x9b, 0x09, 0x73, 0x59, 0x7c,
	0x0a, 0xdf, 0x59, 0x4c, 0x76, 0xc5, 0x42, 0xae, 0x9b, 0xf9, 0xfe, 0x01, 0x93, 0xb0, 0x9b, 0x6c,
	0xfe, 0x08, 0x71, 0x60, 0x9d, 0x84, 0x33, 0xe3, 0x9a, 0x2c, 0x8f, 0x23, 0x89, 0xa5, 0xee, 0xde,
	0x33, 0xfb, 0x23, 0xed, 0x0a, 0x71, 0x98, 0x86, 0xdd, 0xa7, 0x98, 0x94, 0x48, 0x4e, 0x55, 0xfa,
	0x40, 0x70, 0xfd, 0x5a, 0x03, 0xf0, 0x06, 0xda, 0xe0, 0x54, 0xd3, 0x37, 0xdb, 0xf5, 0xbe, 0x0d,
	0x0e, 0x11, 0x56, 0xd1, 0x8c, 0xa4, 0xda, 0x82, 0x94, 0x1b, 0x7b, 0xd8, 0x38, 0x3a, 0x59, 0x09,
	0x29, 0x07, 0x8e, 0xea, 0xa2, 0x50, 0x7f, 0x21, 0xbc, 0x85, 0xcb, 0x24, 0xc1, 0x92, 0x5a, 0xf5,
	0xcd, 0xb6, 0xdd, 0xd7, 0xe7, 0xf1, 0xab, 0x81, 0xcd, 0xdc, 0xf3, 0x12, 0x07, 0xc6, 0x67, 0x00,
	0xe3, 0xdc, 0xd2, 0x8c, 0xfa, 0x77, 0xb0, 0x9e, 0xb1, 0xee, 0x4e, 0xd7, 0xe5, 0x7a, 0x59, 0xae,
	0xdf, 0xb2, 0x84, 0xe8, 0xdf, 0xcd, 0x71, 0xa2, 0x1f, 0xbf, 0xa7, 0xbc, 0xf8, 0xcf, 0x6a, 0xbb,
	0x7f, 0xd2, 0x0f, 0x57, 0x45, 0xf9, 0xf4, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x2e, 0xd2, 0x99, 0xac,
	0x3e, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProductInfoClient is the client API for ProductInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProductInfoClient interface {
	AddProduct(ctx context.Context, in *Product, opts ...grpc.CallOption) (*wrappers.StringValue, error)
	GetProduct(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*Product, error)
}

type productInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewProductInfoClient(cc grpc.ClientConnInterface) ProductInfoClient {
	return &productInfoClient{cc}
}

func (c *productInfoClient) AddProduct(ctx context.Context, in *Product, opts ...grpc.CallOption) (*wrappers.StringValue, error) {
	out := new(wrappers.StringValue)
	err := c.cc.Invoke(ctx, "/ecommerce.ProductInfo/addProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productInfoClient) GetProduct(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*Product, error) {
	out := new(Product)
	err := c.cc.Invoke(ctx, "/ecommerce.ProductInfo/getProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductInfoServer is the server API for ProductInfo service.
type ProductInfoServer interface {
	AddProduct(context.Context, *Product) (*wrappers.StringValue, error)
	GetProduct(context.Context, *wrappers.StringValue) (*Product, error)
}

// UnimplementedProductInfoServer can be embedded to have forward compatible implementations.
type UnimplementedProductInfoServer struct {
}

func (*UnimplementedProductInfoServer) AddProduct(ctx context.Context, req *Product) (*wrappers.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProduct not implemented")
}
func (*UnimplementedProductInfoServer) GetProduct(ctx context.Context, req *wrappers.StringValue) (*Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProduct not implemented")
}

func RegisterProductInfoServer(s *grpc.Server, srv ProductInfoServer) {
	s.RegisterService(&_ProductInfo_serviceDesc, srv)
}

func _ProductInfo_AddProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Product)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductInfoServer).AddProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ecommerce.ProductInfo/AddProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductInfoServer).AddProduct(ctx, req.(*Product))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductInfo_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrappers.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductInfoServer).GetProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ecommerce.ProductInfo/GetProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductInfoServer).GetProduct(ctx, req.(*wrappers.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProductInfo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ecommerce.ProductInfo",
	HandlerType: (*ProductInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "addProduct",
			Handler:    _ProductInfo_AddProduct_Handler,
		},
		{
			MethodName: "getProduct",
			Handler:    _ProductInfo_GetProduct_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "product_info.proto",
}

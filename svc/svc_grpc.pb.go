// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: svc.proto

package svc

import (
	context "context"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ConvGetImageClient is the client API for ConvGetImage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConvGetImageClient interface {
	Convert(ctx context.Context, in *ConvertRequest, opts ...grpc.CallOption) (*ConvertResponse, error)
	Image(ctx context.Context, in *ImageRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
}

type convGetImageClient struct {
	cc grpc.ClientConnInterface
}

func NewConvGetImageClient(cc grpc.ClientConnInterface) ConvGetImageClient {
	return &convGetImageClient{cc}
}

func (c *convGetImageClient) Convert(ctx context.Context, in *ConvertRequest, opts ...grpc.CallOption) (*ConvertResponse, error) {
	out := new(ConvertResponse)
	err := c.cc.Invoke(ctx, "/conv_get_img.ConvGetImage/Convert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *convGetImageClient) Image(ctx context.Context, in *ImageRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/conv_get_img.ConvGetImage/Image", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConvGetImageServer is the server API for ConvGetImage service.
// All implementations must embed UnimplementedConvGetImageServer
// for forward compatibility
type ConvGetImageServer interface {
	Convert(context.Context, *ConvertRequest) (*ConvertResponse, error)
	Image(context.Context, *ImageRequest) (*httpbody.HttpBody, error)
	mustEmbedUnimplementedConvGetImageServer()
}

// UnimplementedConvGetImageServer must be embedded to have forward compatible implementations.
type UnimplementedConvGetImageServer struct {
}

func (UnimplementedConvGetImageServer) Convert(context.Context, *ConvertRequest) (*ConvertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Convert not implemented")
}
func (UnimplementedConvGetImageServer) Image(context.Context, *ImageRequest) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Image not implemented")
}
func (UnimplementedConvGetImageServer) mustEmbedUnimplementedConvGetImageServer() {}

// UnsafeConvGetImageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConvGetImageServer will
// result in compilation errors.
type UnsafeConvGetImageServer interface {
	mustEmbedUnimplementedConvGetImageServer()
}

func RegisterConvGetImageServer(s grpc.ServiceRegistrar, srv ConvGetImageServer) {
	s.RegisterService(&ConvGetImage_ServiceDesc, srv)
}

func _ConvGetImage_Convert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConvertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConvGetImageServer).Convert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/conv_get_img.ConvGetImage/Convert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConvGetImageServer).Convert(ctx, req.(*ConvertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConvGetImage_Image_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConvGetImageServer).Image(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/conv_get_img.ConvGetImage/Image",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConvGetImageServer).Image(ctx, req.(*ImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConvGetImage_ServiceDesc is the grpc.ServiceDesc for ConvGetImage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConvGetImage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "conv_get_img.ConvGetImage",
	HandlerType: (*ConvGetImageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Convert",
			Handler:    _ConvGetImage_Convert_Handler,
		},
		{
			MethodName: "Image",
			Handler:    _ConvGetImage_Image_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "svc.proto",
}

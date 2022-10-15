// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: protos/multiple_issuer.proto

package protos

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

// MultipleIssuerClient is the client API for MultipleIssuer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MultipleIssuerClient interface {
	IssueMultipleVC(ctx context.Context, in *MsgRequestMultipleVC, opts ...grpc.CallOption) (*MsgResponseMultipleVC, error)
}

type multipleIssuerClient struct {
	cc grpc.ClientConnInterface
}

func NewMultipleIssuerClient(cc grpc.ClientConnInterface) MultipleIssuerClient {
	return &multipleIssuerClient{cc}
}

func (c *multipleIssuerClient) IssueMultipleVC(ctx context.Context, in *MsgRequestMultipleVC, opts ...grpc.CallOption) (*MsgResponseMultipleVC, error) {
	out := new(MsgResponseMultipleVC)
	err := c.cc.Invoke(ctx, "/MultipleIssuer.MultipleIssuer/IssueMultipleVC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MultipleIssuerServer is the server API for MultipleIssuer service.
// All implementations must embed UnimplementedMultipleIssuerServer
// for forward compatibility
type MultipleIssuerServer interface {
	IssueMultipleVC(context.Context, *MsgRequestMultipleVC) (*MsgResponseMultipleVC, error)
	mustEmbedUnimplementedMultipleIssuerServer()
}

// UnimplementedMultipleIssuerServer must be embedded to have forward compatible implementations.
type UnimplementedMultipleIssuerServer struct {
}

func (UnimplementedMultipleIssuerServer) IssueMultipleVC(context.Context, *MsgRequestMultipleVC) (*MsgResponseMultipleVC, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueMultipleVC not implemented")
}
func (UnimplementedMultipleIssuerServer) mustEmbedUnimplementedMultipleIssuerServer() {}

// UnsafeMultipleIssuerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MultipleIssuerServer will
// result in compilation errors.
type UnsafeMultipleIssuerServer interface {
	mustEmbedUnimplementedMultipleIssuerServer()
}

func RegisterMultipleIssuerServer(s grpc.ServiceRegistrar, srv MultipleIssuerServer) {
	s.RegisterService(&MultipleIssuer_ServiceDesc, srv)
}

func _MultipleIssuer_IssueMultipleVC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRequestMultipleVC)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultipleIssuerServer).IssueMultipleVC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MultipleIssuer.MultipleIssuer/IssueMultipleVC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultipleIssuerServer).IssueMultipleVC(ctx, req.(*MsgRequestMultipleVC))
	}
	return interceptor(ctx, in, info, handler)
}

// MultipleIssuer_ServiceDesc is the grpc.ServiceDesc for MultipleIssuer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MultipleIssuer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MultipleIssuer.MultipleIssuer",
	HandlerType: (*MultipleIssuerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IssueMultipleVC",
			Handler:    _MultipleIssuer_IssueMultipleVC_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/multiple_issuer.proto",
}
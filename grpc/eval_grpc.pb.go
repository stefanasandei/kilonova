// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// EvalClient is the client API for Eval service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EvalClient interface {
	// Compile compiles a program, to be used for later execution
	Compile(ctx context.Context, in *CompileRequest, opts ...grpc.CallOption) (*CompileResponse, error)
	// Execute runs a test, returning their output
	Execute(ctx context.Context, in *Test, opts ...grpc.CallOption) (*TestResponse, error)
	Clean(ctx context.Context, in *CleanArgs, opts ...grpc.CallOption) (*Empty, error)
}

type evalClient struct {
	cc grpc.ClientConnInterface
}

func NewEvalClient(cc grpc.ClientConnInterface) EvalClient {
	return &evalClient{cc}
}

func (c *evalClient) Compile(ctx context.Context, in *CompileRequest, opts ...grpc.CallOption) (*CompileResponse, error) {
	out := new(CompileResponse)
	err := c.cc.Invoke(ctx, "/eval.Eval/Compile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evalClient) Execute(ctx context.Context, in *Test, opts ...grpc.CallOption) (*TestResponse, error) {
	out := new(TestResponse)
	err := c.cc.Invoke(ctx, "/eval.Eval/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evalClient) Clean(ctx context.Context, in *CleanArgs, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/eval.Eval/Clean", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EvalServer is the server API for Eval service.
// All implementations must embed UnimplementedEvalServer
// for forward compatibility
type EvalServer interface {
	// Compile compiles a program, to be used for later execution
	Compile(context.Context, *CompileRequest) (*CompileResponse, error)
	// Execute runs a test, returning their output
	Execute(context.Context, *Test) (*TestResponse, error)
	Clean(context.Context, *CleanArgs) (*Empty, error)
	mustEmbedUnimplementedEvalServer()
}

// UnimplementedEvalServer must be embedded to have forward compatible implementations.
type UnimplementedEvalServer struct {
}

func (UnimplementedEvalServer) Compile(context.Context, *CompileRequest) (*CompileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Compile not implemented")
}
func (UnimplementedEvalServer) Execute(context.Context, *Test) (*TestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedEvalServer) Clean(context.Context, *CleanArgs) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Clean not implemented")
}
func (UnimplementedEvalServer) mustEmbedUnimplementedEvalServer() {}

// UnsafeEvalServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EvalServer will
// result in compilation errors.
type UnsafeEvalServer interface {
	mustEmbedUnimplementedEvalServer()
}

func RegisterEvalServer(s grpc.ServiceRegistrar, srv EvalServer) {
	s.RegisterService(&_Eval_serviceDesc, srv)
}

func _Eval_Compile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvalServer).Compile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eval.Eval/Compile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvalServer).Compile(ctx, req.(*CompileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Eval_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Test)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvalServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eval.Eval/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvalServer).Execute(ctx, req.(*Test))
	}
	return interceptor(ctx, in, info, handler)
}

func _Eval_Clean_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CleanArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvalServer).Clean(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eval.Eval/Clean",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvalServer).Clean(ctx, req.(*CleanArgs))
	}
	return interceptor(ctx, in, info, handler)
}

var _Eval_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eval.Eval",
	HandlerType: (*EvalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Compile",
			Handler:    _Eval_Compile_Handler,
		},
		{
			MethodName: "Execute",
			Handler:    _Eval_Execute_Handler,
		},
		{
			MethodName: "Clean",
			Handler:    _Eval_Clean_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "eval.proto",
}

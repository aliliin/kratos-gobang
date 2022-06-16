// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.3.1

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationGobangRegister = "/gobang.v1.Gobang/Register"
const OperationGobangLogin = "/gobang.v1.Gobang/Login"
const OperationGobangMemberStatus = "/gobang.v1.Gobang/MemberStatus"
const OperationGobangRoomCreate = "/gobang.v1.Gobang/RoomCreate"
const OperationGobangSayHello = "/gobang.v1.Gobang/SayHello"

type GobangHTTPServer interface {
	Login(context.Context, *LoginReq) (*LoginReply, error)
	MemberStatus(context.Context, *emptypb.Empty) (*StatusReply, error)
	Register(context.Context, *RegisterReq) (*RegisterReply, error)
	RoomCreate(context.Context, *RoomRequest) (*RoomReply, error)
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGobangHTTPServer(s *http.Server, srv GobangHTTPServer) {
	r := s.Route("/")
	r.POST("/member/register", _Gobang_Register0_HTTP_Handler(srv))
	r.POST("/member/login", _Gobang_Login0_HTTP_Handler(srv))
	r.GET("/member/status", _Gobang_MemberStatus0_HTTP_Handler(srv))
	r.POST("/room/create", _Gobang_RoomCreate0_HTTP_Handler(srv))
	r.GET("/helloworld/{name}", _Gobang_SayHello0_HTTP_Handler(srv))
}

func _Gobang_Register0_HTTP_Handler(srv GobangHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RegisterReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGobangRegister)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Register(ctx, req.(*RegisterReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RegisterReply)
		return ctx.Result(200, reply)
	}
}

func _Gobang_Login0_HTTP_Handler(srv GobangHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LoginReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGobangLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*LoginReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginReply)
		return ctx.Result(200, reply)
	}
}

func _Gobang_MemberStatus0_HTTP_Handler(srv GobangHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGobangMemberStatus)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.MemberStatus(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*StatusReply)
		return ctx.Result(200, reply)
	}
}

func _Gobang_RoomCreate0_HTTP_Handler(srv GobangHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RoomRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGobangRoomCreate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RoomCreate(ctx, req.(*RoomRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RoomReply)
		return ctx.Result(200, reply)
	}
}

func _Gobang_SayHello0_HTTP_Handler(srv GobangHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in HelloRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGobangSayHello)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayHello(ctx, req.(*HelloRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*HelloReply)
		return ctx.Result(200, reply)
	}
}

type GobangHTTPClient interface {
	Login(ctx context.Context, req *LoginReq, opts ...http.CallOption) (rsp *LoginReply, err error)
	MemberStatus(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *StatusReply, err error)
	Register(ctx context.Context, req *RegisterReq, opts ...http.CallOption) (rsp *RegisterReply, err error)
	RoomCreate(ctx context.Context, req *RoomRequest, opts ...http.CallOption) (rsp *RoomReply, err error)
	SayHello(ctx context.Context, req *HelloRequest, opts ...http.CallOption) (rsp *HelloReply, err error)
}

type GobangHTTPClientImpl struct {
	cc *http.Client
}

func NewGobangHTTPClient(client *http.Client) GobangHTTPClient {
	return &GobangHTTPClientImpl{client}
}

func (c *GobangHTTPClientImpl) Login(ctx context.Context, in *LoginReq, opts ...http.CallOption) (*LoginReply, error) {
	var out LoginReply
	pattern := "/member/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGobangLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GobangHTTPClientImpl) MemberStatus(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*StatusReply, error) {
	var out StatusReply
	pattern := "/member/status"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGobangMemberStatus))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GobangHTTPClientImpl) Register(ctx context.Context, in *RegisterReq, opts ...http.CallOption) (*RegisterReply, error) {
	var out RegisterReply
	pattern := "/member/register"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGobangRegister))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GobangHTTPClientImpl) RoomCreate(ctx context.Context, in *RoomRequest, opts ...http.CallOption) (*RoomReply, error) {
	var out RoomReply
	pattern := "/room/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGobangRoomCreate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GobangHTTPClientImpl) SayHello(ctx context.Context, in *HelloRequest, opts ...http.CallOption) (*HelloReply, error) {
	var out HelloReply
	pattern := "/helloworld/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGobangSayHello))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

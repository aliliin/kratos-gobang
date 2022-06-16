package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	v1 "service/api/gobang/v1"
	"service/internal/biz"
	"service/internal/pkg/auth"
)

// GobangService is a Gobang service.
type GobangService struct {
	v1.UnimplementedGobangServer

	uc *biz.MemberUsecase

	log *log.Helper
}

// NewGobangService new a Gobang service.
func NewGobangService(uc *biz.MemberUsecase, logger log.Logger) *GobangService {
	return &GobangService{uc: uc, log: log.NewHelper(logger)}
}

// SayHello implements helloworld.GreeterServer.
func (s *GobangService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	return &v1.HelloReply{Message: "Hello "}, nil
}

func (s *GobangService) Register(ctx context.Context, in *v1.RegisterReq) (*v1.RegisterReply, error) {
	member, err := s.uc.CreateMember(ctx, &biz.Member{
		UserName: in.Username,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RegisterReply{
		Username: member.UserName,
	}, nil
}

func (s *GobangService) Login(ctx context.Context, in *v1.LoginReq) (*v1.LoginReply, error) {
	member, err := s.uc.CheckMember(ctx, &biz.Member{
		UserName: in.Username,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{Token: member.Token}, nil
}

func (s *GobangService) MemberStatus(ctx context.Context, in *v1.StatusRequest) (*v1.StatusReply, error) {
	cu := auth.FromContext(ctx)
	fmt.Println("cu.", cu.Username)

	//user, err := uc.uRepo.UserById(ctx, uId)

	//transport.Header.Get(ctx, "X-Session-Id")
	//if header, ok := transport.FromServerContext(ctx); ok {
	//	path := header.RequestHeader().Get("X-Session-Id")
	//	fmt.Println(path)
	//	//fmt.Println(header)
	//} else {
	//	return nil, errors.New(500, "jwt claim missing", "dfasd")
	//}
	//return nil, errors.New(999, "not implemented", "not implemented")
	//return &v1.StatusReply{}, nil
	return nil, nil
}

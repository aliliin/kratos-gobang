package service

import (
	"context"
	v1 "service/api/gobang/v1"
	"service/internal/biz"
)

// GobangService is a Gobang service.
type GobangService struct {
	v1.UnimplementedGobangServer

	uc *biz.MemberUsecase
}

// NewGobangService new a Gobang service.
func NewGobangService(uc *biz.MemberUsecase) *GobangService {
	return &GobangService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GobangService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	return &v1.HelloReply{Message: "Hello "}, nil
}

func (s *GobangService) MemberStatus(ctx context.Context, in *v1.StatusRequest) (*v1.StatusReply, error) {
	return &v1.StatusReply{
		Id:       1,
		Username: "ddddddd",
	}, nil
}

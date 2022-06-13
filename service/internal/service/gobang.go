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

// SayHello implements gobang.GobangServer.
func (s *GobangService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGobang(ctx, &biz.Gobang{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

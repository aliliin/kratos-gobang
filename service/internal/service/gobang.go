package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "service/api/gobang/v1"
	"service/internal/biz"
	"strings"
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

func (s *GobangService) MemberStatus(ctx context.Context, r *emptypb.Empty) (*v1.StatusReply, error) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		tokenString := tr.RequestHeader().Get("Authorization")
		if tokenString != "" {
			auths := strings.SplitN(tokenString, " ", 2)
			if len(auths) != 2 || !strings.EqualFold(auths[0], "Token") {
				return nil, errors.New("jwt token missing")
			}
			token, err := jwt.Parse(auths[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				secret := "gobang"
				return []byte(secret), nil
			})
			if err != nil {
				return nil, err
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if u, ok := claims["username"]; ok {
					fmt.Println(u.(string))
					member, err := s.uc.FindByMember(ctx, u.(string))
					if err != nil {
						return nil, err
					}
					if member != nil {
						return &v1.StatusReply{
							Id:       int32(member.ID),
							Username: member.UserName,
						}, nil
					}
					return nil, nil

				}

			} else {
				return nil, errors.New("Token Invalid")
			}
		}

	}
	return nil, nil

}

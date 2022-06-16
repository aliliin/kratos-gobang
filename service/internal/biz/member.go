package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"service/internal/conf"
	"service/internal/pkg/auth"
	"time"
)

var (
	ErrUserPassword = errors.Forbidden("USER_PASSWORD_ERROR", "用户密码错误")
)

type Member struct {
	ID          int64     `json:"id"`
	UserName    string    `json:"username"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt time.Time `json:"last_login_at"`
	Token       string    `json:"token"`
}

// MemberRepo is a Greater repo.
type MemberRepo interface {
	Save(context.Context, *Member) (*Member, error)
	FindByUsername(context.Context, *Member) (*Member, error)
	FindByMemberName(context.Context, string) (*Member, error)
	ListAll(context.Context) ([]*Member, error)
}

// MemberUsecase is a Member usecase.
type MemberUsecase struct {
	repo MemberRepo
	log  *log.Helper

	jwtc *conf.JWT
}

// NewMemberUsecase new a Member usecase.
func NewMemberUsecase(repo MemberRepo, logger log.Logger, jwt *conf.JWT) *MemberUsecase {
	return &MemberUsecase{repo: repo, log: log.NewHelper(logger), jwtc: jwt}
}

// CreateMember creates a Member, and returns the new Member.
func (uc *MemberUsecase) CreateMember(ctx context.Context, r *Member) (*Member, error) {
	return uc.repo.Save(ctx, r)
}

func (uc *MemberUsecase) CheckMember(ctx context.Context, r *Member) (*Member, error) {
	user, err := uc.repo.FindByUsername(ctx, r)
	if err != nil {
		return nil, err
	}
	return &Member{
		ID:          user.ID,
		UserName:    user.UserName,
		CreatedAt:   user.CreatedAt,
		LastLoginAt: user.LastLoginAt,
		Token:       uc.generateToken(r.UserName),
	}, nil
}

func (uc *MemberUsecase) FindByMember(ctx context.Context, username string) (*Member, error) {

	user, err := uc.repo.FindByMemberName(ctx, username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return &Member{
			ID:          user.ID,
			UserName:    user.UserName,
			CreatedAt:   user.CreatedAt,
			LastLoginAt: user.LastLoginAt,
		}, nil
	}
	return nil, nil
}

func (uc *MemberUsecase) generateToken(username string) string {
	return auth.GenerateToken(uc.jwtc.Secret, username)
}

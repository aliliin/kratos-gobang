package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "service/api/gobang/v1"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Member is a Member model.
type Member struct {
	Hello string
}

// MemberRepo is a Greater repo.
type MemberRepo interface {
	Save(context.Context, *Member) (*Member, error)
	Update(context.Context, *Member) (*Member, error)
	FindByID(context.Context, int64) (*Member, error)
	ListByHello(context.Context, string) ([]*Member, error)
	ListAll(context.Context) ([]*Member, error)
}

// MemberUsecase is a Member usecase.
type MemberUsecase struct {
	repo MemberRepo
	log  *log.Helper
}

// NewMemberUsecase new a Member usecase.
func NewMemberUsecase(repo MemberRepo, logger log.Logger) *MemberUsecase {
	return &MemberUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateMember creates a Member, and returns the new Member.
func (uc *MemberUsecase) CreateMember(ctx context.Context, g *Member) (*Member, error) {
	uc.log.WithContext(ctx).Infof("CreateMember: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

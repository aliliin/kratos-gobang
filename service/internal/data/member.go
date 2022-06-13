package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"service/internal/biz"
)

type memberRepo struct {
	data *Data
	log  *log.Helper
}

// NewMemberRepo .
func NewMemberRepo(data *Data, logger log.Logger) biz.MemberRepo {
	return &memberRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *memberRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *memberRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *memberRepo) FindByID(context.Context, int64) (*biz.Greeter, error) {
	return nil, nil
}

func (r *memberRepo) ListByHello(context.Context, string) ([]*biz.Greeter, error) {
	return nil, nil
}

func (r *memberRepo) ListAll(context.Context) ([]*biz.Greeter, error) {
	return nil, nil
}

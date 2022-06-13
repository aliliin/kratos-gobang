package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"service/internal/biz"
)

type WatchMemberId []int

// Room 房间模型
type Room struct {
	RoomId           int32
	CreatorId        int32
	PlayerIdOne      int32 // 玩家1
	PlayerOneReady   bool
	PlayerIdTwo      int32
	PlayerTwoReady   bool // 玩家2是否准备
	Status           int
	Title            string
	Creator          string
	StatusText       string      // 状态文本
	playerOne        interface{} // 玩家信息1
	playerTwo        interface{}
	PersonNumber     int           // 人数
	WatchMemberIds   WatchMemberId // 观战用户列表
	WatchMemberInfos interface{}   // 观战用户信息列表

}

type RoomRepo struct {
	data *Data
	log  *log.Helper
}

// NewRoomRepo .
func NewRoomRepo(data *Data, logger log.Logger) biz.RoomRepo {
	return &RoomRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *RoomRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *RoomRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *RoomRepo) FindByID(context.Context, int64) (*biz.Greeter, error) {
	return nil, nil
}

func (r *RoomRepo) ListByHello(context.Context, string) ([]*biz.Greeter, error) {
	return nil, nil
}

func (r *RoomRepo) ListAll(context.Context) ([]*biz.Greeter, error) {
	return nil, nil
}

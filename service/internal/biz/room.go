package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type WatchMemberId []int

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

// MemberRepo is a Greater repo.
type RoomRepo interface {
	Save(context.Context, *Room) (*Room, error)
	Update(context.Context, *Room) (*Room, error)
	FindByID(context.Context, int64) (*Room, error)
	ListByHello(context.Context, string) ([]*Room, error)
	ListAll(context.Context) ([]*Room, error)
}

// RoomUsecase is a Member usecase.
type RoomUsecase struct {
	repo RoomRepo
	log  *log.Helper
}

// NewRoomUsecase new a Member usecase.
func NewRoomUsecase(repo RoomRepo, logger log.Logger) *RoomUsecase {
	return &RoomUsecase{repo: repo, log: log.NewHelper(logger)}
}

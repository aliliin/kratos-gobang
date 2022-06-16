package data

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"service/internal/biz"
	"strings"
	"time"
)

type Member struct {
	ID          int64     `gorm:"primarykey;auto_increment"`
	Username    string    `gorm:"column:username;type:varchar(50) not null;default:'';index;comment '用户名称'"`
	Password    string    `gorm:"type:varchar(100) not null; default:'';comment '用户密码'"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;default:current_timestamp;comment '创建时间'"`
	LastLoginAt time.Time `gorm:"column:last_login_at;type:timestamp;comment '最后登录时间'"`
}

func (Member) TableName() string {
	return "member"
}

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

func (r *memberRepo) Save(ctx context.Context, u *biz.Member) (*biz.Member, error) {
	var member Member
	result := r.data.db.Where(&Member{Username: u.UserName}).First(&member)
	if result.RowsAffected == 1 {
		return nil, errors.New(500, "USER_EXIST", "用户名已存在"+u.UserName)
	}
	member.Username = u.UserName
	member.Password = encrypt(u.Password) // 密码加密
	res := r.data.db.Create(&member)
	if res.Error != nil {
		return nil, errors.New(500, "CREAT_USER_ERROR", "用户创建失败")
	}

	return &biz.Member{
		UserName: member.Username,
	}, nil
}

func encrypt(psd string) string {
	options := &password.Options{SaltLen: 16, Iterations: 10000, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(psd, options)
	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}

func (r *memberRepo) FindByUsername(ctx context.Context, u *biz.Member) (*biz.Member, error) {
	var member Member
	result := r.data.db.Where(&Member{Username: u.UserName}).First(&member)
	if result.RowsAffected == 0 {
		return nil, errors.New(500, "USER_EXIST", "用户不存在"+u.UserName)
	}
	if checkPassword(u.Password, member.Password) {
		member.LastLoginAt = time.Now()
		r.data.db.Save(&member)
		return &biz.Member{
			ID:          member.ID,
			UserName:    member.Username,
			Password:    member.Password,
			CreatedAt:   member.CreatedAt,
			LastLoginAt: member.LastLoginAt,
		}, nil
	}
	return nil, errors.New(500, "PASSWORD_ERROR", "密码错误")
}

func checkPassword(psd, encryptedPassword string) bool {
	options := &password.Options{SaltLen: 16, Iterations: 10000, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(encryptedPassword, "$")
	check := password.Verify(psd, passwordInfo[2], passwordInfo[3], options)
	return check
}

func (r *memberRepo) ListAll(context.Context) ([]*biz.Member, error) {
	return nil, nil
}

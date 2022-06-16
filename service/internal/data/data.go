package data

import (
	"service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v9"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewMemberRepo)

// Data .
type Data struct {
	db  *gorm.DB
	rdb *redis.Client
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, rdb *redis.Client) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "service/data"))

	return &Data{db: db, rdb: rdb, log: log}, func() {}, nil
}

// NewDB .
func NewDB(c *conf.Data, logger log.Logger) *gorm.DB {
	// 终端打印输入 sql 执行记录
	log := log.NewHelper(log.With(logger, "module", "gorm"))

	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&Member{})
	if err != nil {
		return nil
	}
	return db
}

func NewRedis(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	//rdb.AddHook(redisotel.TracingHook{})
	if err := rdb.Close(); err != nil {
		log.Error(err)
	}
	return rdb
}

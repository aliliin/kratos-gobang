// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"service/internal/biz"
	"service/internal/conf"
	"service/internal/data"
	"service/internal/server"
	"service/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, jwt *conf.JWT, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData, logger)
	client := data.NewRedis(confData)
	dataData, cleanup, err := data.NewData(confData, logger, db, client)
	if err != nil {
		return nil, nil, err
	}
	memberRepo := data.NewMemberRepo(dataData, logger)
	memberUsecase := biz.NewMemberUsecase(memberRepo, logger, jwt)
	gobangService := service.NewGobangService(memberUsecase, logger)
	httpServer := server.NewHTTPServer(confServer, jwt, gobangService, logger)
	app := newApp(logger, httpServer)
	return app, func() {
		cleanup()
	}, nil
}

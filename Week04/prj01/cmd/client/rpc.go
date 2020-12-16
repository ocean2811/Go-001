//+build wireinject

package main

import (
	"go001/internal/biz"
	"go001/internal/data"
	"go001/internal/service"

	"github.com/google/wire"
)

func InitRPCServer(redisAddr string) service.Server {
	wire.Build(data.NewRedis, biz.NewBiz, service.NewRPCServer)
	return nil
}

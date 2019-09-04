package main

import (
	"mo-gateway/src/config"
	"mo-gateway/src/control/mysql"
	"mo-gateway/src/server"
	"mo-gateway/src/server/gin_server"
	"mo-gateway/src/server/grpc"
	"vilgo/vlog"
)

func main() {
	cfg := config.Init()
	vlog.DefaultLogger()
	if err := server.Init(&mysql.MySqlServe{}); err != nil {
		vlog.LogE("service init fail %v", err)
		return
	}
	s := server.NewServe(gin_server.NewGinServer(cfg), grpc.NewGrpcServer(cfg))
	s.Start()
}

package main

import (
	"mo-gateway/src/control/mysql"
	"mo-gateway/src/control/server"
	"mo-gateway/src/control/server/gin_server"
	"vilgo/vlog"
)

func main() {
	vlog.DefaultLogger()
	if err := server.Init(&mysql.MySqlServe{}); err != nil {
		vlog.LogE("service init fail %v", err)
		return
	}
	s := server.NewServe(gin_server.NewGinServer())
	s.Start()
}

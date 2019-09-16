package main

import (
	"fmt"
	"mo-gateway/src/config"
	"mo-gateway/src/control/mysql"
	"mo-gateway/src/server"
	"mo-gateway/src/server/gin_server"
	"mo-gateway/src/server/grpc"
	"os"
	"os/signal"
	"syscall"
	"vilgo/vlog"
)

func main() {
	sgc := make(chan os.Signal, 1)
	signal.Notify(sgc, os.Interrupt, os.Kill, syscall.SIGQUIT)
	cfg := config.Init()
	vlog.DefaultLogger()
	if err := server.Init(&mysql.MySqlServe{}); err != nil {
		vlog.LogE("service init fail %v", err)
		sg := <-sgc
		fmt.Println("mo-gateway exit ", sg)
		os.Exit(-1)
	}
	s := server.NewServe(gin_server.NewGinServer(cfg), grpc.NewGrpcServer(cfg))
	s.Start()
	fmt.Println("mo-gateway start ok ")
	sg := <-sgc
	fmt.Println("mo-gateway exit ", sg)
}

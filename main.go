package main

import (
	"fmt"
	"mo-gateway/src/config"
	"mo-gateway/src/server"
	"mo-gateway/src/server/gin_server"
	"mo-gateway/src/server/grpc"
	"os"
	"os/signal"
	"syscall"
	"vilgo/vlog"
)

func WaitSignal() {
	sgc := make(chan os.Signal, 1)
	signal.Notify(sgc, os.Interrupt, os.Kill, syscall.SIGQUIT)
	sg := <-sgc
	fmt.Println("exit ", sg)
}

func main() {
	cfg := config.Init()
	vlog.DefaultLogger("info.log")
	s := server.NewServe(gin_server.NewGinServer(cfg), grpc.NewGrpcServer(cfg))
	s.Start()
	fmt.Println("mo-gateway start ok ")
	WaitSignal()
}

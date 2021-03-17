package main

import (
	"fmt"
	"github.com/ville-vv/mo-gateway/src/config"
	"github.com/ville-vv/mo-gateway/src/server"
	"github.com/ville-vv/mo-gateway/src/server/gin_server"
	"github.com/ville-vv/mo-gateway/src/server/grpc"
	"github.com/ville-vv/vilgo/vlog"
	"os"
	"os/signal"
	"syscall"
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

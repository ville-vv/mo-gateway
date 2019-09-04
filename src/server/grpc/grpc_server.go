package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mo-gateway/src/config"
	"mo-gateway/src/server/grpc/health"
	"mo-gateway/src/server/grpc/pb"
	"net"
	"reflect"
)

type GrpcServer struct {
	name string
	addr string
}

func NewGrpcServer(cfg *config.Config) *GrpcServer {
	g := &GrpcServer{}
	g.name = reflect.TypeOf(g).String()
	g.addr = fmt.Sprintf("%s:%s", cfg.GrpcServer.Host, cfg.GrpcServer.Port)
	return g
}

func (sel *GrpcServer) Name() string {
	return fmt.Sprintf("%s with %s", sel.name, sel.addr)
}

func (sel *GrpcServer) Start(...interface{}) error {
	// 注入
	s := grpc.NewServer()
	pb.RegisterMoGatewayServer(s, sel)

	health.New(s, "mo-gateway")

	lis, err := net.Listen("tcp", sel.addr)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

func (sel *GrpcServer) Stop() {
}

func (sel *GrpcServer) DoPing(ctx context.Context, req *pb.Ping) (resp *pb.Pong, err error) {
	return &pb.Pong{}, nil
}

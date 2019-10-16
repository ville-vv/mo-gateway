package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"mo-gateway/src/config"
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
	g.name = reflect.TypeOf(g).Elem().String()
	g.addr = fmt.Sprintf("%s:%s", cfg.GrpcServer.Host, cfg.GrpcServer.Port)
	return g
}

func (sel *GrpcServer) Name() string {
	return fmt.Sprintf("%s with %s", sel.name, sel.addr)
}

func (sel *GrpcServer) healthCheckServer(s *grpc.Server) {
	hsrv := health.NewServer()
	hsrv.SetServingStatus("mo-gateway", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, hsrv)
}

func (sel *GrpcServer) Start(...interface{}) error {
	// 注入
	s := grpc.NewServer()
	pb.RegisterMoGatewayServer(s, sel)
	sel.healthCheckServer(s)
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

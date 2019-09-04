package health

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
}

func New(s *grpc.Server, name string) {
	hlth := health.NewServer()
	hlth.SetServingStatus(name, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, hlth)
}

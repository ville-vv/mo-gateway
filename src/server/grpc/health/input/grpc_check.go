package input

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"mo-gateway/src/server/grpc/health"
	"time"
)

const (
	CHECK_STATUS_NOT_SERVING   int32 = 0
	CHECK_STATUS_SERVING       int32 = 1
	CHECK_STATUS_UNKNOWN       int32 = 2
	CHECK_STATUS_NO_HEALTH_PRO int32 = 3
	CHECK_STATUS_NO_TIMEOUT    int32 = 4
)

type Client struct {
	addr             string
	name             string
	healthCheckParam *healthpb.HealthCheckRequest
	ctx              context.Context
	timeout          time.Duration
}

func NewClient(addr string, name string, timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 30
	}
	return &Client{
		addr:             addr,
		name:             name,
		healthCheckParam: &healthpb.HealthCheckRequest{},
		ctx:              context.Background(),
		timeout:          timeout * time.Second,
	}
}

func (sel *Client) Init() error {
	return nil
}

// check the grpc server
func (sel *Client) check() (sts int32) {
	dialCtx, cencel := context.WithTimeout(context.Background(), sel.timeout)
	defer cencel()
	// 建立grpc链接
	conn, err := grpc.DialContext(dialCtx, sel.addr, grpc.WithInsecure())
	if err != nil {
		if err == context.DeadlineExceeded {
			sts = CHECK_STATUS_NO_TIMEOUT
		} else {
			sts = CHECK_STATUS_NOT_SERVING
		}
		return
	}
	cli := healthpb.NewHealthClient(conn)
	defer conn.Close()
	rpcCtx, cencel2 := context.WithTimeout(sel.ctx, sel.timeout)
	defer cencel2()

	resp, err := cli.Check(rpcCtx, sel.healthCheckParam)
	if err != nil {
		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
			// "error: this server does not implement the grpc health protocol (grpc.health.v1.Health)"
			sts = CHECK_STATUS_NO_HEALTH_PRO
		} else if stat, ok := status.FromError(err); ok && stat.Code() == codes.DeadlineExceeded {
			//"timeout: health rpc did not complete within timeout
			sts = CHECK_STATUS_NO_TIMEOUT
		} else {
			// error: health rpc failed:
			sts = CHECK_STATUS_NOT_SERVING
		}
		return
	}

	switch resp.Status {
	case healthpb.HealthCheckResponse_SERVING:
		sts = CHECK_STATUS_SERVING
	case healthpb.HealthCheckResponse_NOT_SERVING:
		sts = CHECK_STATUS_NOT_SERVING
	default:
		sts = CHECK_STATUS_UNKNOWN
	}
	return
}

// Gather some data
func (sel *Client) Gather(acc health.Accumulator) error {
	sts := sel.check()
	acc.AddStatus(sel.name, map[string]interface{}{"active": sts}, time.Now())
	return nil
}

func (sel *Client) Label() string {
	return ""
}

package health

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"time"
)

const (
	CHECK_STATUS_NOT_SERVING       int32 = 0
	CHECK_STATUS_SERVING           int32 = 1
	CHECK_STATUS_UNKNOWN           int32 = 2
	CHECK_STATUS_NO_HEALTH_PRO     int32 = 3
	CHECK_STATUS_NO_TIMEOUT        int32 = 4
	CHECK_STATUS_HEALTH_RPC_FAILED int32 = 5
)

type Client struct {
	addr             string
	name             string
	cli              healthpb.HealthClient
	healthCheckParam *healthpb.HealthCheckRequest
}

func NewClient(addr string, name string) *Client {
	return &Client{
		addr:             addr,
		name:             name,
		healthCheckParam: &healthpb.HealthCheckRequest{},
	}
}

func (sel *Client) Init() error {
	// 建立grpc链接
	conn, err := grpc.Dial(sel.addr)
	if err != nil {
		return err
	}
	sel.cli = healthpb.NewHealthClient(conn)
	return nil
}

// check the grpc server
func (sel *Client) check() (sts int32, err error) {
	resp, err := sel.cli.Check(context.Background(), sel.healthCheckParam)
	if err != nil {
		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
			// "error: this server does not implement the grpc health protocol (grpc.health.v1.Health)"
			sts = CHECK_STATUS_NO_HEALTH_PRO
		} else if stat, ok := status.FromError(err); ok && stat.Code() == codes.DeadlineExceeded {
			//"timeout: health rpc did not complete within timeout
			sts = CHECK_STATUS_NO_TIMEOUT
		} else {
			// error: health rpc failed:
			sts = CHECK_STATUS_HEALTH_RPC_FAILED
		}
		return
	}

	switch resp.Status {
	case healthpb.HealthCheckResponse_SERVING:
		sts = CHECK_STATUS_SERVING
	case healthpb.HealthCheckResponse_NOT_SERVING:
		sts = CHECK_STATUS_NOT_SERVING
	case healthpb.HealthCheckResponse_UNKNOWN:
		sts = CHECK_STATUS_UNKNOWN
	default:
		err = errors.New("wrong check status")
	}
	return
}

// Gather some data
func (sel *Client) Gather(acc Accumulator) error {
	sts, err := sel.check()
	if err != nil {
		acc.AddError(sel.name, err)
		return nil
	}
	acc.AddStatus(sel.name, map[string]interface{}{"active": sts}, time.Now())
	return nil
}

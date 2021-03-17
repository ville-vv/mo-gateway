package gin_server

import (
	"fmt"
	"github.com/ville-vv/mo-gateway/src/config"
	"github.com/ville-vv/vilgo/vlog"
	"reflect"
)

type GinServer struct {
	name string
	addr string
}

func NewGinServer(cfg *config.Config) *GinServer {
	g := new(GinServer)
	g.name = reflect.TypeOf(g).Elem().Name()
	g.addr = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	return g
}

func (g *GinServer) Name() string {
	return fmt.Sprintf("%s with %s", g.name, g.addr)
}

func (g *GinServer) Start(args ...interface{}) error {
	r := NewRoute()
	r.Use(&UserRouterModel{})
	r.Use(&WebHookRouter{})
	r.RequestLog()
	r.ResponseLog()
	return r.Start(g.addr)
}

func (g *GinServer) Stop() {
	vlog.LogI("%s stop", g.name)
}

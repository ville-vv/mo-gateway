package gin_server

import (
	"fmt"
	"mo-gateway/src/config"
	"reflect"
	"vilgo/vlog"
)

type GinServer struct {
	name string
	addr string
}

func NewGinServer(cfg *config.Config) *GinServer {
	g := new(GinServer)
	g.name = reflect.TypeOf(g).String()
	g.addr = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	return g
}

func (g *GinServer) Name() string {
	return fmt.Sprintf("%s with %s", g.name, g.addr)
}

func (g *GinServer) Start(args ...interface{}) error {
	r := NewRoute()
	r.Use(&UserRouterModel{})
	r.ResponseLog()
	return r.Start(g.addr)
}

func (g *GinServer) Stop() {
	vlog.LogI("%s stop", g.name)
}

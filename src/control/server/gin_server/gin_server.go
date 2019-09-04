package gin_server

import (
	"mo-gateway/src/config"
	"reflect"
	"vilgo/vlog"
)


type GinServer struct {
	name string
	addr string
}

func NewGinServer() *GinServer {
	g := new(GinServer)
	g.name = reflect.TypeOf(g).Name()
	g.addr = config.ServerAddress()
	return g
}

func (g *GinServer) Name() string {
	return g.name
}

func (g *GinServer) Start(args ... interface{}) error {
	r := NewRoute()
	r.Use(&UserRouterModel{})
	r.RequestLog().ResponseLog()
	return r.Start(g.addr)
}

func (g *GinServer) Stop() {
	vlog.LogI("%s stop", g.name)
}


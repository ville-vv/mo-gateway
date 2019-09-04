// Package gin_server 
package gin_server

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"reflect"
	"vilgo/vlog"
)

type IRouter interface {
	Load(*gin.Engine) error
}


type Route struct {
	gEng      *gin.Engine
	loadMp    map[string]IRouter
	accessLog vlog.ILogger
}

func NewRoute() *Route {
	r := new(Route)
	r.gEng = gin.New()
	r.gEng.Use(gin.Recovery())
	// 服务访问日志，只用来记录 请求和响应的日志
	r.accessLog = vlog.NewGoLogger(&vlog.LogCnf{
		ProgramName:   "router",
		OutPutFile:    []string{"./log/access.log"},
		OutPutErrFile: []string{},
		Level:         vlog.LogLevelInfo,
	})
	return r
}

// 设置路由
func (r *Route) Use(m IRouter) *Route {
	if r.loadMp == nil {
		r.loadMp = make(map[string]IRouter)
	}
	key := reflect.TypeOf(m)
	r.loadMp[key.String()] = m
	return r
}

// 启动路由
func (r *Route) Start(addr string) error {
	for _, v := range r.loadMp {
		// 循环注册路由服务
		if err := v.Load(r.gEng); err != nil {
			return err
		}
	}
	return r.gEng.Run(addr)
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	// 这里截取 response 的数据
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// 接口请求日志输出中间层
func (r *Route) reqLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyBytes, _ := ioutil.ReadAll(ctx.Request.Body)
		defer ctx.Request.Body.Close()
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		r.accessLog.LogI("[REQ] [REMOTE %s ] [URI %s ] [BODY %s ]", ctx.Request.RemoteAddr, ctx.Request.URL, string(bodyBytes))
		ctx.Next()
	}
}

// 接口返回日志输出中间层
func (r *Route) respLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bWriter := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = bWriter
		ctx.Next()
		r.accessLog.LogI("[RESP] [REMOTE %s ] [URI %s ] [BODY %s ]", ctx.Request.RemoteAddr, ctx.Request.URL, bWriter.body.String())
	}
}

func (r *Route) RequestLog() *Route {
	r.gEng.Use(r.reqLog())
	return r
}

func (r *Route) ResponseLog() *Route {
	r.gEng.Use(r.respLog())
	return r
}

// @File     : route_webhook
// @Author   : Ville
// @Time     : 19-10-16 上午9:35
// gin_server
package gin_server

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mo-gateway/src/control/errmsg"
	"mo-gateway/src/handler/webhook"
	"net/http"
)

type WebHookRouter struct {
	root string
	Cmd
}

func (sel *WebHookRouter) Load(r *gin.Engine) error {
	sel.root = "/api/webhook"
	sel.Cmd = Cmd{
		Handler: &webhook.Handler{},
	}
	sel.route(r)
	return nil
}

func (sel *WebHookRouter) route(r *gin.Engine) {
	gp := r.Group(sel.root)
	gp.POST("/travis", sel.Travis)
	//gp.GET("/travis", sel.cmd.Travis)
	return
}

type PostHandleFun func([]byte, ...interface{}) (interface{}, error)
type GetHandleFun func(map[string][]string, ...interface{}) (interface{}, error)

type WebHookHandler interface {
	TravisPost([]byte, ...interface{}) (interface{}, error)
	TravisGet(map[string][]string, ...interface{}) (interface{}, error)
}

type Cmd struct {
	Handler WebHookHandler
}

func (sel *Cmd) Travis(ctx *gin.Context) {
	switch ctx.Request.Method {
	case "POST":
		Post(ctx, sel.Handler.TravisPost)
	case "GET":
		Get(ctx, sel.Handler.TravisGet)
	default:
		sayBackErr(ctx, 400, int(errmsg.ServerNotFound), errmsg.ServerNotFound)
	}
}

func Post(ctx *gin.Context, hf PostHandleFun) {
	defer ctx.Request.Body.Close()
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		sayBackErr(ctx, 200, int(errmsg.ParamFormatErr), errmsg.ParamFormatErr)
		return
	}
	resp, err := hf(body)
	if err != nil {
		sayBackErr(ctx, 200, 500, err)
		return
	}
	sayBackOk(ctx, resp)
}

func Get(ctx *gin.Context, hf GetHandleFun) {
	resp, err := hf(ctx.Request.URL.Query())
	if err != nil {
		sayBackErr(ctx, 200, 500, err)
		return
	}
	sayBackOk(ctx, resp)
}

type SayBackData struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"err_msg,omitempty"`
}

func sayBackErr(ctx *gin.Context, httpCode int, errCode int, err error) {
	ctx.JSON(httpCode, &BackDataOk{Status: int(errCode), ErrMsg: err.Error()})
}

func sayBackOk(ctx *gin.Context, body interface{}) {
	ctx.JSON(http.StatusOK, &BackDataOk{Status: 200, Data: body})
}

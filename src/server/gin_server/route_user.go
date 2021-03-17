package gin_server

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/ville-vv/mo-gateway/src/control/errmsg"
	"github.com/ville-vv/mo-gateway/src/handler/user"
	"io/ioutil"
	"net/http"
)

type HandlerFun func([]byte) (interface{}, errmsg.ErrCode)

type UserRouterModel struct {
	UserRouteCmd
	root string
}

func (sel *UserRouterModel) init() {
	sel.root = "/api/user"
	sel.UserRouteCmd = UserRouteCmd{
		handler: &user.User{},
	}
}

func (sel *UserRouterModel) Load(r *gin.Engine) error {
	sel.init()
	gp := r.Group(sel.root)
	gp.POST("/login", sel.Login)
	gp.POST("/register", sel.UserRegister)
	gp.POST("/callback", sel.CallBack)
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
type UserRouteCmd struct {
	handler user.Handler
}

func (sel *UserRouteCmd) Login(ctx *gin.Context) {
	PostSayBack(ctx, sel.handler.Login)
}

func (sel *UserRouteCmd) UserRegister(ctx *gin.Context) {
	PostSayBack(ctx, sel.handler.Register)
}

func (sel *UserRouteCmd) CallBack(ctx *gin.Context) {
	PostSayBack(ctx, sel.handler.CallBack)
}

//----------------------------------------        返回值         -------------------------------------------------

type BackDataOk struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"err_msg,omitempty"`
}

func Back200(ctx *gin.Context, body interface{}) {
	ctx.JSON(http.StatusOK, &BackDataOk{Status: 200, Data: body})
}

func BackError(ctx *gin.Context, code errmsg.ErrCode) {
	ctx.JSON(http.StatusOK, &BackDataOk{Status: int(code), ErrMsg: code.Error()})
}

func GetSayBack(ctx *gin.Context, hf HandlerFun) {
	gd := ctx.Request.URL.Query()
	data, err := jsoniter.Marshal(gd)
	if err != nil {
		BackError(ctx, errmsg.ParamParseErr)
		return
	}
	sayBack(ctx, data, hf)
}

func PostSayBack(ctx *gin.Context, hf HandlerFun) {
	defer ctx.Request.Body.Close()
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		BackError(ctx, errmsg.ParamParseErr)
		return
	}
	sayBack(ctx, body, hf)
	return
}

func sayBack(ctx *gin.Context, req []byte, hf HandlerFun) {
	resp, code := hf(req)
	if code != 0 {
		BackError(ctx, code)
		return
	}
	Back200(ctx, resp)
}

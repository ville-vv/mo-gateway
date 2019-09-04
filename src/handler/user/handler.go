package user

import (
	"github.com/json-iterator/go"
	"mo-gateway/src/control/errmsg"
	"vilgo/vlog"
	"vilgo/vuid"
)

// Handler 业务处理模块接口
type Handler interface {
	Login([]byte) (interface{}, errmsg.ErrCode)
	Register([]byte) (interface{}, errmsg.ErrCode)
	CallBack([]byte) (interface{}, errmsg.ErrCode)
}

type User struct{}

// Login 登录接口
func (*User) Login(body []byte) (interface{}, errmsg.ErrCode) {
	var req LoginReq
	if err := jsoniter.Unmarshal(body, &req); err != nil {
		vlog.LogE("params Unmarshal fail :%v", err)
		return nil, errmsg.ParamFormatErr
	}
	return &LoginResp{Success: true}, 0
}

// Register 注册接口
func (*User) Register(body []byte) (interface{}, errmsg.ErrCode) {
	var req RegisterReq
	if err := jsoniter.Unmarshal(body, &req); err != nil {
		vlog.LogE("params Unmarshal fail :%v", err)
		return nil, errmsg.ParamFormatErr
	}
	if err := req.IsValid(); err != nil {
		vlog.LogE("check params value fail :%v", err)
		return nil, errmsg.InvalidParam
	}
	return &RegisterResp{Account: vuid.GenUUid()}, 0
}

func (*User) CallBack(body []byte) (interface{}, errmsg.ErrCode) {
	vlog.LogI("收到参数： %s", string(body))
	return map[string]interface{}{"status": 200}, 0
}

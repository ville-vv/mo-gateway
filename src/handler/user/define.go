package user

import "errors"

// 登录请求
type LoginReq struct {
	ID  string `json:"id"`
	Pwd string `json:"pwd"`
	Tk  string `json:"tk"`
}
type LoginResp struct {
	Success bool `json:"success"`
}

// 注册请求
type RegisterReq struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	RandCode string `json:"rand_code"`
	SmsCode  string `json:"sms_code"`
}

func (sel *RegisterReq) IsValid() error {
	switch {
	case sel.UserName == "":
		return errors.New("user name is empty")
	case sel.Phone == "":
		return errors.New("phone is empty")
	case sel.Email == "":
		return errors.New("email is empty")
	case sel.SmsCode == "":
		return errors.New("sms code is empty")
	}
	return nil
}

type RegisterResp struct {
	Account int64 `json:"account"`
}

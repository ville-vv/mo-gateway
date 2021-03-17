package user

import (
	"github.com/ville-vv/vilgo/valid"
)

type Login struct {
	tokenVerify valid.Validator
}

func (sel *Login) GenToken() string {
	return sel.tokenVerify.Generate(nil)
}

func (sel *Login) VerifyToken(tkStr string) bool {
	sel.tokenVerify.Verify(tkStr)
	return false
}

package user

import (
	"testing"
	"vilgo/valid"
	"vilgo/vlog"
)

func TestLogin_GenToken(t *testing.T) {
	vlog.DefaultLogger()
	lg := &Login{valid.NewJwtToken()}
	tkstr := lg.GenToken()
	vlog.LogI(tkstr)
	lg.VerifyToken(tkstr)
}

func BenchmarkLogin_GenToken(b *testing.B) {
	lg := &Login{valid.NewJwtToken()}
	for i := 0; i < b.N; i++ {
		lg.GenToken()
	}
	
	//
}

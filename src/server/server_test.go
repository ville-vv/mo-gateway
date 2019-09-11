package server

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 测试一下 context 的使用
func TestServe_context(t *testing.T) {

	ctx1 := context.Background()
	if ctx1.Done() == nil {
		fmt.Println("yes")
	}

	ctx, cancelTimeout := context.WithTimeout(ctx1, time.Second*2)
	defer cancelTimeout()
	ctx2 := context.WithValue(ctx, "aaa", "bbbb")
	if ctx2.Done() == nil {
		fmt.Println("yes")
	}
	tm := time.Now().Add(time.Second * 10)
	ctx3, cancelDeadline := context.WithDeadline(ctx2, tm)
	defer cancelDeadline()
	if ctx3.Done() == nil {
		return
	}
	select {
	case <-ctx3.Done():
		fmt.Println("WithDeadline")
	case <-ctx.Done():
		fmt.Println("WithTimeout")
	}
}

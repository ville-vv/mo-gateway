package errmsg

import (
	"fmt"
	"testing"
)

func TestErrCode_Error(t *testing.T) {
	fmt.Println(StatusOK)
	fmt.Println(SystemInterErr)
	fmt.Println(ServerNotFound)
}

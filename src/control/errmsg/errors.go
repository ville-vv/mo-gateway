package errmsg

type ErrCode int

func (s ErrCode) Error() string {
	return errMap[s]
}

const (
	// 服务 status 3  位 int
	StatusOK       ErrCode = 200
	ParamParseErr  ErrCode = 201
	ServerNotFound ErrCode = 400
	SystemInterErr ErrCode = 500

	// 业务相关的错误 4 位 int
	ParamFormatErr ErrCode = 1201 // 非法参数
	InvalidParam   ErrCode = 1202 // 非法参数
)

var (
	errMap = map[ErrCode]string{
		StatusOK:       "",
		ParamParseErr:  "parameter parse error",
		ServerNotFound: "server not found",
		SystemInterErr: "system internal error",
		ParamFormatErr: "",
		InvalidParam:   "invalid parameter",
	}
)

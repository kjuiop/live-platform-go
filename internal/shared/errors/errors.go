package serrors

import "errors"

const (
	NoError            = ""
	CodeInvalidRequest = "S4401"
	CodeEmptyParam     = "S4402"
	CodeInternalError  = "S5001"
	CodeRedisSave      = "S5002"
	CodeRedisDelete    = "S5003"
)

var codeToMessage = map[string]string{
	NoError:            "ok",
	CodeInvalidRequest: "invalid request body",
	CodeEmptyParam:     "invalid params",
	CodeInternalError:  "internal server error",
	CodeRedisSave:      "internal redis error occurred",
	CodeRedisDelete:    "internal redis error occurred",
}

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrRedisSave      = errors.New("redis save error")
	ErrRedisDelete    = errors.New("redis delete error")
)

func GetCustomMessage(code string) string {
	if msg, ok := codeToMessage[code]; ok {
		return msg
	}
	return "unknown error"
}

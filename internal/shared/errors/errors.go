package serrors

import "errors"

const (
	NoError            = ""
	CodeInvalidRequest = "S4401"
	CodeEmptyParam     = "S4402"
	CodeInternalError  = "S5001"
	CodeRedisSave      = "S5002"
	CodeRedisDelete    = "S5003"
	CodeRedisGet       = "S5004"
)

var codeToMessage = map[string]string{
	NoError:            "ok",
	CodeInvalidRequest: "invalid request body",
	CodeEmptyParam:     "invalid params",
	CodeInternalError:  "internal server error",
	CodeRedisSave:      "internal storage error occurred",
	CodeRedisDelete:    "internal storage error occurred",
	CodeRedisGet:       "internal storage error occurred",
}

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrRepoSave       = errors.New("repository save error")
	ErrRepoDelete     = errors.New("repository delete error")
	ErrRepoGet        = errors.New("repository get error")
)

func GetCustomMessage(code string) string {
	if msg, ok := codeToMessage[code]; ok {
		return msg
	}
	return "unknown error"
}

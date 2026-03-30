package serrors

import "errors"

const (
	NoError            = ""
	CodeInvalidRequest = "S4401"
	CodeEmptyParam     = "S4402"
	CodeInternalError  = "S5001"
	CodeRepoSave       = "S5002"
	CodeRepoDelete     = "S5003"
	CodeRepoGet        = "S5004"
)

var codeToMessage = map[string]string{
	NoError:            "ok",
	CodeInvalidRequest: "invalid request body",
	CodeEmptyParam:     "invalid params",
	CodeInternalError:  "internal server error",
	CodeRepoSave:       "internal storage error occurred",
	CodeRepoDelete:     "internal storage error occurred",
	CodeRepoGet:        "internal storage error occurred",
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

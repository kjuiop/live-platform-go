package serrors

import "errors"

const (
	NoError            = ""
	CodeInvalidRequest = "S4401"
	CodeEmptyParam     = "S4402"
	CodeInternalError  = "S5001"
)

var codeToMessage = map[string]string{
	NoError:            "ok",
	CodeInvalidRequest: "invalid request body",
	CodeEmptyParam:     "invalid params",
	CodeInternalError:  "internal server error",
}

var (
	ErrInvalidRequest = errors.New("invalid request")
)

func GetCustomMessage(code string) string {
	if msg, ok := codeToMessage[code]; ok {
		return msg
	}
	return "unknown error"
}

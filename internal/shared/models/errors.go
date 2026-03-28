package models

import (
	"errors"
	"fmt"
)

type CustomErr struct {
	Code int
	Err  error
}

func (ce *CustomErr) Error() string {
	if ce.Err != nil {
		return ce.Err.Error()
	}
	return fmt.Sprintf("error code %d: Unknown error", ce.Code)
}

var errorList = map[int]CustomErr{
	ErrNotFoundServerInfo: {
		Code: ErrNotFoundServerInfo,
		Err:  errors.New("not found server info"),
	},
}

const (
	InternalError      = "internal server error"
	InternalRedisError = "internal redis error occur"
)

const (
	NoError                int = 0
	ErrParsing             int = 4001
	ErrNotFoundChatRoom    int = 4002
	ErrNotConnectSocket    int = 4003
	ErrEmptyParam          int = 4004
	ErrNotFoundServerInfo  int = 4005
	ErrRedisHMSETError     int = 5001
	ErrRedisExistError     int = 5002
	ErrRedisHMDELError     int = 5003
	ErrInternalServerError int = 5004
)

var codeToMessage = map[int]string{
	NoError:                "ok",
	ErrParsing:             "invalid request body",
	ErrNotFoundChatRoom:    "not found chat room",
	ErrNotConnectSocket:    "not connect socket",
	ErrEmptyParam:          "invalid params",
	ErrRedisHMSETError:     InternalRedisError,
	ErrRedisExistError:     InternalRedisError,
	ErrRedisHMDELError:     InternalRedisError,
	ErrInternalServerError: InternalError,
}

func GetCustomErrMessage(code int, error string) string {
	message, exists := codeToMessage[code]
	if !exists {
		return "Unknown error"
	}

	return fmt.Sprintf("%s, err : %s", message, error)
}

func GetCustomErr(code int) error {
	customErr, exists := errorList[code]
	if !exists || customErr.Err == nil {
		return errors.New("unknown error")
	}
	return customErr.Err
}

func GetCustomMessage(code int) string {
	message, exists := codeToMessage[code]
	if !exists {
		return "Unknown error"
	}

	return message
}

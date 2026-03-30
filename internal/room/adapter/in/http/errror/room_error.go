package errror

import (
	"errors"
	"net/http"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
)

const (
	ErrCodeUnknown        int = 5000
	ErrCodeInvalidRequest int = 4001
	ErrCodeRoomNotFound   int = 4401
	ErrCodeRedisSave      int = 5001
	ErrCodeRedisDelete    int = 5003
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrRedisSave      = errors.New("redis save error")
	ErrRedisDelete    = errors.New("redis delete error")
)

type ErrMapping struct {
	HttpStatus int
	ErrorCode  int
	Msg        string
}

var DomainErrMap = map[error]ErrMapping{
	domain.ErrRoomNotFound: {http.StatusNotFound, ErrCodeRoomNotFound, "not found chat room"},
	ErrInvalidRequest:      {http.StatusBadRequest, ErrCodeInvalidRequest, "invalid request body"},
	ErrRedisSave:           {http.StatusInternalServerError, ErrCodeRedisSave, "internal redis error occurred"},
	ErrRedisDelete:         {http.StatusInternalServerError, ErrCodeRedisDelete, "internal redis error occurred"},
}

func GetMapping(err error) ErrMapping {
	if m, ok := DomainErrMap[err]; ok {
		return m
	}
	return ErrMapping{
		HttpStatus: http.StatusInternalServerError,
		ErrorCode:  ErrCodeUnknown,
		Msg:        "unknown error",
	}
}

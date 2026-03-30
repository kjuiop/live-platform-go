package error

import (
	"net/http"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
	serrors "github.com/kjuiop/live-platform-go/internal/shared/errors"
)

const (
	ErrCodeRoomNotFound = "R4401"
	ErrCodeUnknown      = "UNKNOWN"
)

type ErrMapping struct {
	HttpStatus int
	ErrorCode  string
	Msg        string
}

var DomainErrMap = map[error]ErrMapping{
	domain.ErrRoomNotFound:    {http.StatusNotFound, ErrCodeRoomNotFound, "not found chat room"},
	serrors.ErrInvalidRequest: {http.StatusBadRequest, serrors.CodeInvalidRequest, "invalid request body"},
	serrors.ErrRedisSave:      {http.StatusInternalServerError, serrors.CodeRedisSave, "internal redis error occurred"},
	serrors.ErrRedisDelete:    {http.StatusInternalServerError, serrors.CodeRedisDelete, "internal redis error occurred"},
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

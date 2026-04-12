package error

import (
	"errors"
	"net/http"

	roomerr "github.com/kjuiop/live-platform-go/internal/room/domain/errors"
	serrors "github.com/kjuiop/live-platform-go/internal/shared/errors"
)

const (
	ErrCodeRoomNotFound  = "R4401"
	ErrCodeInvalidParams = "R4402"
	ErrCodeCreateFailed  = "R5001"
	ErrCodeDeleteFailed  = "R5002"
	ErrCodeGetFailed     = "R5003"
	ErrCodeUnknown       = "UNKNOWN"
)

type ErrMapping struct {
	HttpStatus int
	ErrorCode  string
	Msg        string
}

func GetMapping(err error) ErrMapping {
	switch {
	case errors.Is(err, serrors.ErrInvalidRequest):
		return ErrMapping{http.StatusBadRequest, ErrCodeInvalidParams, "invalid request body"}
	case errors.Is(err, roomerr.ErrRoomNotFound):
		return ErrMapping{http.StatusNotFound, ErrCodeRoomNotFound, "not found chat room"}
	case errors.Is(err, roomerr.ErrCreateChatRoomFailed):
		return ErrMapping{http.StatusInternalServerError, ErrCodeCreateFailed, "failed to create chat room"}
	case errors.Is(err, roomerr.ErrDeleteChatRoomFailed):
		return ErrMapping{http.StatusInternalServerError, ErrCodeDeleteFailed, "failed to delete chat room"}
	case errors.Is(err, roomerr.ErrGetChatRoomsFailed):
		return ErrMapping{http.StatusInternalServerError, ErrCodeGetFailed, "failed to get chat rooms"}
	default:
		return ErrMapping{http.StatusInternalServerError, ErrCodeUnknown, "unknown error"}
	}
}

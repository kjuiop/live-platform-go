package errors

import "errors"

var (
	ErrRoomNotFound         = errors.New("room not found")
	ErrCreateChatRoomFailed = errors.New("failed to create chat room")
	ErrDeleteChatRoomFailed = errors.New("failed to delete chat room")
	ErrGetChatRoomsFailed   = errors.New("failed to get chat rooms")
)

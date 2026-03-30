package in

import "errors"

var (
	ErrCreateChatRoomFailed = errors.New("failed to create chat room")
	ErrDeleteChatRoomFailed = errors.New("failed to delete chat room")
	ErrGetChatRoomsFailed   = errors.New("failed to get chat rooms")
)

package in

import (
	"context"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
)

type RoomService interface {
	CreateChatRoom(ctx context.Context, room domain.RoomInfo) error
}

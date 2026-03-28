package out

import (
	"context"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
)

type RoomRepository interface {
	Save(ctx context.Context, room domain.RoomInfo) error
	RegisterRoomMap(ctx context.Context, room domain.RoomInfo) error
}

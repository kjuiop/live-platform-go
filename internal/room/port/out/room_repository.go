package out

import (
	"context"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
)

type RoomRepository interface {
	SaveRoom(ctx context.Context, room domain.RoomInfo) error
	DeleteRoom(ctx context.Context, roomId string) error
	GetRooms(ctx context.Context) ([]domain.RoomInfo, error)
}

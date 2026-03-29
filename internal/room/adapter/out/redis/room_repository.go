package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/kjuiop/live-platform-go/internal/room/adapter/out/redis/lua"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
	roomoutport "github.com/kjuiop/live-platform-go/internal/room/port/out"
	"github.com/kjuiop/live-platform-go/platform/redis"
)

const (
	roomKeyPrefix = "live-platform-room"
	roomMapKey    = "live-platform-room-map"
	RoomExpire    = time.Duration(7) * 24 * time.Hour
)

var _ roomoutport.RoomRepository = (*RoomRedisRepository)(nil)

type RoomRedisRepository struct {
	redis *redis.Client
}

func NewRoomRedisRepository(redis *redis.Client) *RoomRedisRepository {
	return &RoomRedisRepository{
		redis: redis,
	}
}

func roomKey(roomId string) string {
	return fmt.Sprintf("%s:%s", roomKeyPrefix, roomId)
}

func (r *RoomRedisRepository) SaveRoom(ctx context.Context, room domain.RoomInfo) error {
	mapField := fmt.Sprintf("%s:%s", room.ChannelKey, room.BroadcastKey)
	keys := []string{
		roomKey(room.RoomId),
		roomMapKey,
	}
	args := []interface{}{
		room.RoomId,
		room.CustomerId,
		room.ChannelKey,
		room.BroadcastKey,
		room.CreatedAt,
		int64(RoomExpire.Seconds()),
		mapField,
		room.RoomId,
	}
	if err := r.redis.RunScript(ctx, lua.SaveRoomScript, keys, args...); err != nil {
		return fmt.Errorf("failed to save room to redis : %w", err)
	}
	return nil
}

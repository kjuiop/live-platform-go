package redis

import (
	"context"
	"fmt"
	"time"

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

func (r *RoomRedisRepository) Save(ctx context.Context, data domain.RoomInfo) error {
	key := roomKey(data.RoomId)

	if err := r.redis.HSet(ctx, key, data.ConvertRedisData(), RoomExpire); err != nil {
		return fmt.Errorf("failed to save room to redis : %w", err)
	}

	return nil
}

func (r *RoomRedisRepository) RegisterRoomMap(ctx context.Context, room domain.RoomInfo) error {
	field := fmt.Sprintf("%s:%s", room.ChannelKey, room.BroadcastKey)
	if err := r.redis.HSetField(ctx, roomMapKey, field, room.RoomId); err != nil {
		return fmt.Errorf("failed to register room map : %w", err)
	}

	return nil
}

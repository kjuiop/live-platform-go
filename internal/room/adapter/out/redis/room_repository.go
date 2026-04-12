package redis

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/kjuiop/live-platform-go/internal/room/adapter/out/redis/lua"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
	roomerr "github.com/kjuiop/live-platform-go/internal/room/domain/errors"
	roomoutport "github.com/kjuiop/live-platform-go/internal/room/port/out"
	"github.com/kjuiop/live-platform-go/platform/redis"
)

const (
	roomKeyPrefix = "live:rooms"
	roomMapKey    = "live:rooms-map"
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
		return fmt.Errorf("roomRepo.SaveRoom: %w", err)
	}
	return nil
}

func (r *RoomRedisRepository) DeleteRoom(ctx context.Context, roomId string) error {
	keys := []string{
		roomKey(roomId),
		roomMapKey,
	}
	cmd, err := r.redis.RunScriptResult(ctx, lua.DeleteRoomScript, keys, roomId)
	if err != nil {
		return fmt.Errorf("roomRepo:DeleteRoom: %w", err)
	}
	result, err := cmd.Int()
	if err != nil {
		return err
	}
	if result == 0 {
		return fmt.Errorf("roomRepo:DeleteRoom: %w", roomerr.ErrRoomNotFound)
	}
	return nil
}

func (r *RoomRedisRepository) GetRooms(ctx context.Context) ([]domain.RoomInfo, error) {
	// 1. room-map 에서 모든 roomId 조회
	roomIds, err := r.redis.HVals(ctx, roomMapKey)
	if err != nil {
		return nil, fmt.Errorf("roomRepo.GetRooms: %w", err)
	}
	if len(roomIds) == 0 {
		return []domain.RoomInfo{}, nil
	}

	// 2. 각 roomId의 키 생성
	keys := make([]string, len(roomIds))
	for i, id := range roomIds {
		keys[i] = roomKey(id)
	}

	// 3. Pipeline 으로 HGETALL 일괄 조회
	results, err := r.redis.HGetAllPipeline(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("roomRepo.GetRooms:: %w", err)
	}

	// 4. 결과 파싱 (TTL 만료된 방 skip)
	rooms := make([]domain.RoomInfo, 0, len(results))
	for i, m := range results {
		if len(m) == 0 {
			slog.Warn("GetRooms: room key has expired, skipping", "roomKey", keys[i])
			continue
		}
		createdAt, err := strconv.ParseInt(m["created_at"], 10, 64)
		if err != nil {
			slog.Warn("GetRooms: invalid created_at, skipping", "roomId", m["room_id"], "value", m["created_at"])
			continue
		}
		rooms = append(rooms, domain.RoomInfo{
			RoomId:       m["room_id"],
			CustomerId:   m["customer_id"],
			ChannelKey:   m["channel_key"],
			BroadcastKey: m["broadcast_key"],
			CreatedAt:    createdAt,
		})
	}
	return rooms, nil
}

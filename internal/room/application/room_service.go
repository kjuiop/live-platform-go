package application

import (
	"context"
	"fmt"
	"time"

	roomRepo "github.com/kjuiop/live-platform-go/internal/room/adapter/out/redis"
	"github.com/kjuiop/live-platform-go/internal/room/domain"
	roomin "github.com/kjuiop/live-platform-go/internal/room/port/in"
)

var _ roomin.RoomService = (*RoomServiceImpl)(nil)

type RoomServiceImpl struct {
	contextTimeout time.Duration
	roomRepo       *roomRepo.RoomRedisRepository
}

func NewRoomService(timeout time.Duration, roomRepo *roomRepo.RoomRedisRepository) *RoomServiceImpl {
	return &RoomServiceImpl{
		contextTimeout: timeout,
		roomRepo:       roomRepo,
	}
}

func (r *RoomServiceImpl) CreateChatRoom(ctx context.Context, room domain.RoomInfo) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.Save(ctx, room); err != nil {
		return fmt.Errorf("failed to save room to redis : %w", err)
	}

	return nil
}

func (r *RoomServiceImpl) RegisterRoomId(ctx context.Context, room domain.RoomInfo) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.RegisterRoomMap(ctx, room); err != nil {
		return fmt.Errorf("failed to register room id to redis : %w", err)
	}

	return nil
}

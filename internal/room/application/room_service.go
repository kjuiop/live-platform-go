package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
	roomin "github.com/kjuiop/live-platform-go/internal/room/port/in"
	roomoutport "github.com/kjuiop/live-platform-go/internal/room/port/out"
)

var _ roomin.RoomService = (*RoomServiceImpl)(nil)

type RoomServiceImpl struct {
	contextTimeout time.Duration
	roomRepo       roomoutport.RoomRepository
}

func NewRoomService(timeout time.Duration, roomRepo roomoutport.RoomRepository) *RoomServiceImpl {
	return &RoomServiceImpl{
		contextTimeout: timeout,
		roomRepo:       roomRepo,
	}
}

func (r *RoomServiceImpl) CreateChatRoom(ctx context.Context, room domain.RoomInfo) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.SaveRoom(ctx, room); err != nil {
		return fmt.Errorf("failed to save room to redis: %w", err)
	}

	return nil
}

func (r *RoomServiceImpl) DeleteChatRoom(ctx context.Context, roomId string) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.DeleteRoom(ctx, roomId); err != nil {
		if errors.Is(err, domain.ErrRoomNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}

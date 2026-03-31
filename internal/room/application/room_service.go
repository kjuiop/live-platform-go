package application

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	roomerr "github.com/kjuiop/live-platform-go/internal/room/domain/errors"

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
		return fmt.Errorf("roomSrv:CreateChatRoom: %w:%w", roomerr.ErrCreateChatRoomFailed, err)
	}

	return nil
}

func (r *RoomServiceImpl) DeleteChatRoom(ctx context.Context, roomId string) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if err := r.roomRepo.DeleteRoom(ctx, roomId); err != nil {
		if errors.Is(err, roomerr.ErrRoomNotFound) {
			return roomerr.ErrRoomNotFound
		}
		return fmt.Errorf("roomSrv.DeleteChatRoom: %w: %w", roomerr.ErrDeleteChatRoomFailed, err)
	}
	return nil
}

func (r *RoomServiceImpl) GetChatRooms(ctx context.Context) ([]domain.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	rooms, err := r.roomRepo.GetRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("roomSrv.GetChatRooms: %w:%w", roomerr.ErrGetChatRoomsFailed, err)
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].CreatedAt > rooms[j].CreatedAt
	})
	return rooms, nil
}

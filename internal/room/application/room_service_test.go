package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
)

type mockRoomRepository struct {
	saveErr            error
	registerRoomMapErr error
}

func (m *mockRoomRepository) Save(_ context.Context, _ domain.RoomInfo) error {
	return m.saveErr
}

func (m *mockRoomRepository) RegisterRoomMap(_ context.Context, _ domain.RoomInfo) error {
	return m.registerRoomMapErr
}

func TestRoomServiceImpl_CreateChatRoom(t *testing.T) {
	tests := []struct {
		name    string
		saveErr error
		wantErr bool
	}{
		{
			name:    "정상 저장",
			saveErr: nil,
			wantErr: false,
		},
		{
			name:    "Redis 저장 실패",
			saveErr: errors.New("redis error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRoomRepository{saveErr: tt.saveErr}
			svc := NewRoomService(5*time.Second, repo)

			err := svc.CreateChatRoom(context.Background(), domain.RoomInfo{
				RoomId:     "room-1",
				CustomerId: "customer-1",
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateChatRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoomServiceImpl_RegisterRoomId(t *testing.T) {
	tests := []struct {
		name               string
		registerRoomMapErr error
		wantErr            bool
	}{
		{
			name:               "정상 등록",
			registerRoomMapErr: nil,
			wantErr:            false,
		},
		{
			name:               "Redis 등록 실패",
			registerRoomMapErr: errors.New("redis error"),
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRoomRepository{registerRoomMapErr: tt.registerRoomMapErr}
			svc := NewRoomService(5*time.Second, repo)

			err := svc.RegisterRoomId(context.Background(), domain.RoomInfo{
				RoomId:       "room-1",
				ChannelKey:   "ch-abc",
				BroadcastKey: "bc-xyz",
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterRoomId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

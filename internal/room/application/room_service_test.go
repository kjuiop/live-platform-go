package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
)

type mockRoomRepository struct {
	saveErr error
}

func (m *mockRoomRepository) SaveRoom(_ context.Context, _ domain.RoomInfo) error {
	return m.saveErr
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

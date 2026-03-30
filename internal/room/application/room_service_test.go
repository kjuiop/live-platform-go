package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
	serrors "github.com/kjuiop/live-platform-go/internal/shared/errors"
)

type mockRoomRepository struct {
	saveErr   error
	deleteErr error
}

func (m *mockRoomRepository) DeleteRoom(_ context.Context, _ string) error {
	return m.deleteErr
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

func TestRoomServiceImpl_DeleteChatRoom(t *testing.T) {
	repoErr := errors.New("redis error")

	tests := []struct {
		name         string
		deleteErr    error
		wantErr      bool
		wantNotFound bool
		wantWrapped  error
	}{
		{
			name:      "정상 삭제",
			deleteErr: nil,
			wantErr:   false,
		},
		{
			name:         "ErrRoomNotFound 그대로 전달",
			deleteErr:    domain.ErrRoomNotFound,
			wantErr:      true,
			wantNotFound: true,
		},
		{
			name:        "기타 레포지토리 에러 → ErrRedisDelete 로 변환",
			deleteErr:   repoErr,
			wantErr:     true,
			wantWrapped: serrors.ErrRepoDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRoomRepository{deleteErr: tt.deleteErr}
			svc := NewRoomService(5*time.Second, repo)

			err := svc.DeleteChatRoom(context.Background(), "room-1")

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteChatRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNotFound && !errors.Is(err, domain.ErrRoomNotFound) {
				t.Errorf("DeleteChatRoom() expected ErrRoomNotFound, got %v", err)
			}
			if tt.wantWrapped != nil && !errors.Is(err, tt.wantWrapped) {
				t.Errorf("DeleteChatRoom() expected wrapped err %v, got %v", tt.wantWrapped, err)
			}
		})
	}
}

func (m *mockRoomRepository) GetRooms(ctx context.Context) ([]domain.RoomInfo, error) {
	//TODO implement me
	panic("implement me")
}

package application

import (
	"context"
	"errors"
	"testing"
	"time"

	roomerr "github.com/kjuiop/live-platform-go/internal/room/domain/errors"

	"github.com/kjuiop/live-platform-go/internal/room/domain"
)

type mockRoomRepository struct {
	saveErr        error
	deleteErr      error
	getRoomsResult []domain.RoomInfo
	getRoomsErr    error
}

func (m *mockRoomRepository) DeleteRoom(_ context.Context, _ string) error {
	return m.deleteErr
}

func (m *mockRoomRepository) SaveRoom(_ context.Context, _ domain.RoomInfo) error {
	return m.saveErr
}

func (m *mockRoomRepository) GetRooms(_ context.Context) ([]domain.RoomInfo, error) {
	return m.getRoomsResult, m.getRoomsErr
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
			deleteErr:    roomerr.ErrRoomNotFound,
			wantErr:      true,
			wantNotFound: true,
		},
		{
			name:        "기타 레포지토리 에러 → ErrDeleteChatRoomFailed 로 변환",
			deleteErr:   repoErr,
			wantErr:     true,
			wantWrapped: roomerr.ErrDeleteChatRoomFailed,
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
			if tt.wantNotFound && !errors.Is(err, roomerr.ErrRoomNotFound) {
				t.Errorf("DeleteChatRoom() expected ErrRoomNotFound, got %v", err)
			}
			if tt.wantWrapped != nil && !errors.Is(err, tt.wantWrapped) {
				t.Errorf("DeleteChatRoom() expected wrapped err %v, got %v", tt.wantWrapped, err)
			}
		})
	}
}

func TestRoomServiceImpl_GetChatRooms(t *testing.T) {
	tests := []struct {
		name           string
		getRoomsResult []domain.RoomInfo
		getRoomsErr    error
		wantErr        bool
		wantErrIs      error
		wantLen        int
		checkOrder     bool
	}{
		{
			name:           "빈 목록 반환",
			getRoomsResult: []domain.RoomInfo{},
			wantLen:        0,
		},
		{
			name: "최신순 정렬",
			getRoomsResult: []domain.RoomInfo{
				{RoomId: "room-1", CreatedAt: 1000},
				{RoomId: "room-2", CreatedAt: 3000},
				{RoomId: "room-3", CreatedAt: 2000},
			},
			wantLen:    3,
			checkOrder: true,
		},
		{
			name:        "Repository 에러 → ErrGetChatRoomsFailed 래핑",
			getRoomsErr: errors.New("redis error"),
			wantErr:     true,
			wantErrIs:   roomerr.ErrGetChatRoomsFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRoomRepository{
				getRoomsResult: tt.getRoomsResult,
				getRoomsErr:    tt.getRoomsErr,
			}
			svc := NewRoomService(5*time.Second, repo)

			rooms, err := svc.GetChatRooms(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetChatRooms() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("GetChatRooms() expected errors.Is(%v), got %v", tt.wantErrIs, err)
			}
			if !tt.wantErr && len(rooms) != tt.wantLen {
				t.Errorf("GetChatRooms() len = %d, want %d", len(rooms), tt.wantLen)
			}
			if tt.checkOrder {
				for i := 1; i < len(rooms); i++ {
					if rooms[i-1].CreatedAt < rooms[i].CreatedAt {
						t.Errorf("GetChatRooms() not sorted: rooms[%d].CreatedAt=%d < rooms[%d].CreatedAt=%d",
							i-1, rooms[i-1].CreatedAt, i, rooms[i].CreatedAt)
					}
				}
			}
		})
	}
}

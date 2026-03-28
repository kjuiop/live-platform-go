package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/kjuiop/live-platform-go/config"
	"github.com/kjuiop/live-platform-go/internal/room/domain"
)

var testClient *TestClient

type TestClient struct {
	srv         *httptest.Server
	roomHandler *RoomHandler
}

type mockRoomService struct {
	createChatRoomErr error
	registerRoomIdErr error
}

func (m *mockRoomService) CreateChatRoom(_ context.Context, _ domain.RoomInfo) error {
	return m.createChatRoomErr
}

func (m *mockRoomService) RegisterRoomId(_ context.Context, _ domain.RoomInfo) error {
	return m.registerRoomIdErr
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	svc := &mockRoomService{}
	cfg := config.Policy{Prefix: "room,chat,live"}

	testClient = &TestClient{
		roomHandler: NewRoomHandler(cfg, svc),
	}
	testClient.srv = setupTestServer()

	exitCode := m.Run()
	testClient.srv.Close()
	os.Exit(exitCode)
}

func setupTestServer() *httptest.Server {
	r := gin.New()
	testClient.roomHandler.RegisterRoutes(r.Group("/api/v1"))
	return httptest.NewServer(r)
}

// ---- 테스트 케이스 ----

func TestCreateRoom(t *testing.T) {
	tests := []struct {
		name              string
		body              map[string]string
		createChatRoomErr error
		registerRoomIdErr error
		wantCode          int
		wantRoomId        bool
	}{
		{
			name: "정상 생성 - 201",
			body: map[string]string{
				"customer_id":   "customer-1",
				"channel_key":   "ch-abc",
				"broadcast_key": "bc-xyz",
			},
			wantCode:   http.StatusCreated,
			wantRoomId: true,
		},
		{
			name: "필수 필드 누락 - 400",
			body: map[string]string{
				"channel_key":   "ch-abc",
				"broadcast_key": "bc-xyz",
				// customer_id 누락
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "CreateChatRoom 실패 - 500",
			body: map[string]string{
				"customer_id":   "customer-1",
				"channel_key":   "ch-abc",
				"broadcast_key": "bc-xyz",
			},
			createChatRoomErr: errors.New("redis error"),
			wantCode:          http.StatusInternalServerError,
		},
		{
			name: "RegisterRoomId 실패 - 500",
			body: map[string]string{
				"customer_id":   "customer-1",
				"channel_key":   "ch-abc",
				"broadcast_key": "bc-xyz",
			},
			registerRoomIdErr: errors.New("redis error"),
			wantCode:          http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock 서비스 에러 주입
			svc := testClient.roomHandler.service.(*mockRoomService)
			svc.createChatRoomErr = tt.createChatRoomErr
			svc.registerRoomIdErr = tt.registerRoomIdErr

			bodyBytes, _ := json.Marshal(tt.body)
			req, err := http.NewRequest(
				http.MethodPost,
				testClient.srv.URL+"/api/v1/rooms/",
				bytes.NewReader(bodyBytes),
			)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status code: got %d, want %d", resp.StatusCode, tt.wantCode)
			}

			if tt.wantRoomId {
				var body map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
					t.Fatalf("decode body: %v", err)
				}
				result, ok := body["result"].(map[string]interface{})
				if !ok || result["room_id"] == "" {
					t.Error("room_id: should not be empty in response")
				}
			}
		})
	}
}

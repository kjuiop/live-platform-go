package domain

import (
	"strings"
	"testing"

	"github.com/kjuiop/live-platform-go/internal/room/adapter/in/http/form"
)

func TestNewRoomInfo(t *testing.T) {
	req := form.RoomRequest{
		CustomerId:   "customer-1",
		ChannelKey:   "channel-abc",
		BroadcastKey: "broadcast-xyz",
	}
	prefix := "N1,N2"

	info := NewRoomInfo(req, prefix)

	if info.RoomId == "" {
		t.Error("RoomId: should not be empty")
	}

	// prefix 중 하나로 시작하는지 확인
	prefixes := strings.Split(prefix, ",")
	found := false
	for _, p := range prefixes {
		if strings.HasPrefix(info.RoomId, p+"-") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("RoomId: should start with one of prefixes, got %s", info.RoomId)
	}

	if info.CustomerId != req.CustomerId {
		t.Errorf("CustomerId: got %s, want %s", info.CustomerId, req.CustomerId)
	}
	if info.ChannelKey != req.ChannelKey {
		t.Errorf("ChannelKey: got %s, want %s", info.ChannelKey, req.ChannelKey)
	}
	if info.BroadcastKey != req.BroadcastKey {
		t.Errorf("BroadcastKey: got %s, want %s", info.BroadcastKey, req.BroadcastKey)
	}
	if info.CreatedAt == 0 {
		t.Error("CreatedAt: should not be zero")
	}
}

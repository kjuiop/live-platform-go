package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kjuiop/live-platform-go/internal/room/adapter/in/http/form"
	"github.com/kjuiop/live-platform-go/internal/shared/utils"
)

var (
	ErrRoomNotFound = errors.New("room not found")
)

type RoomInfo struct {
	RoomId       string `json:"room_id"`
	CustomerId   string `json:"customer_id"`
	ChannelKey   string `json:"channel_key"`
	BroadcastKey string `json:"broadcast_key"`
	CreatedAt    int64  `json:"created_at"`
}

func NewRoomInfo(req form.RoomRequest, prefix string) *RoomInfo {
	return &RoomInfo{
		RoomId:       fmt.Sprintf("%s-%s", getChatPrefix(prefix), utils.GenUUID()),
		CustomerId:   req.CustomerId,
		ChannelKey:   req.ChannelKey,
		BroadcastKey: req.BroadcastKey,
		CreatedAt:    time.Now().Unix(),
	}
}

func getChatPrefix(prefix string) string {
	array := strings.Split(prefix, ",")
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[rng.Intn(len(array))]
}

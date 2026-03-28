package domain

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kjuiop/live-platform-go/internal/room/adapter/in/http/form"
	"github.com/kjuiop/live-platform-go/internal/shared/utils"
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

type RoomRedisData struct {
	RoomId       string `redis:"room_id"`
	CustomerId   string `redis:"customer_id"`
	ChannelKey   string `redis:"channel_key"`
	BroadcastKey string `redis:"broadcast_key"`
	CreatedAt    int64  `redis:"created_at"`
}

func (r *RoomInfo) ConvertRedisData() RoomRedisData {
	return RoomRedisData{
		RoomId:       r.RoomId,
		CustomerId:   r.CustomerId,
		ChannelKey:   r.ChannelKey,
		BroadcastKey: r.BroadcastKey,
		CreatedAt:    r.CreatedAt,
	}
}

type RoomMapRedisData struct {
	RoomId    string `redis:"room_id"`
	CreatedAt int64  `redis:"created_at"`
}

func (r *RoomInfo) ConvertRedisRoomMapData() RoomMapRedisData {
	return RoomMapRedisData{
		RoomId:    r.RoomId,
		CreatedAt: r.CreatedAt,
	}
}

func getChatPrefix(prefix string) string {
	array := strings.Split(prefix, ",")
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[rng.Intn(len(array))]
}

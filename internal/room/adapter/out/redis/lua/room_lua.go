package lua

import "github.com/redis/go-redis/v9"

var SaveRoomScript = redis.NewScript(`
redis.call('HSET', KEYS[1],
    'room_id',        ARGV[1],
    'customer_id',    ARGV[2],
    'channel_key',    ARGV[3],
    'broadcast_key',  ARGV[4],
    'created_at',     ARGV[5])
redis.call('EXPIRE', KEYS[1], ARGV[6])
redis.call('HSET', KEYS[2], ARGV[7], ARGV[8])
return 1
`)

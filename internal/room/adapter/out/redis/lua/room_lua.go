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

// DeleteRoomScript
/**
 * 1. hgetall 로 room_key 를 조회해 해당 키가 있는지 확인
 * 2. room_key 가 존재한다면 channel_key 와 broadcast_key 를 추출하여 map_key 생성
 * 2-1. 이때 i=1 부터 2칸씩 조회하는 이유는 hgetall 의 결과가 [field1, value1, field2, value2, ...] 형태이기 때문
 * 3. map_field 를 추출
 * 4. room_key 삭제
 * 5. map_key 삭제
 * 6. 성공적으로 삭제되었다면 1 반환, room_key 가 존재하지 않는다면 0 반환
 */
var DeleteRoomScript = redis.NewScript(`
local room = redis.call('HGETALL', KEYS[1])
if #room == 0 then
	return 0
end
local channel_key = ''
local broadcast_key = ''
for i = 1, #room, 2 do
	if room[i] == 'channel_key' then
		channel_key = room[i + 1]
	elseif room[i] == 'broadcast_key' then
		broadcast_key = room[i + 1]
	end
end
local map_field = channel_key .. ':' .. broadcast_key
redis.call('DEL', KEYS[1])
redis.call('HDEL', KEYS[2], map_field)
return 1
`)

 -- 你都验证码在 Redis 上的 Key
local key = KEYS[1]
 -- 验证次数, 我们一个验证码, 最多重复三次, 这个记录验证了几次
local cntKey = key..":cnt"
 -- 你的验证码是 123456
local val = ARGV[1]
 -- 验证码有效期是 10 分钟, 600 秒
local ttl = tonumber(redis.call("ttl", key))
 if ttl == -1 then
     --    key 存在但是没有过期时间(系统错误)
     return -2
     --    -2 是 key 不存在, ttl < 540 是发了一个验证码, 已经超过一分钟了
 elseif ttl == -2 or ttl < 540 then
     --    后续如果验证码 有不同过期时间在这里优化
     redis.call("set", key, val)
     redis.call("expire", key, 600)
     redis.call("set",cntKey, 3)
     redis.call("expire", cntKey, 600)
     return 0
 else
     -- 已经发送一个验证码,但是不到一分钟
     return -1
 end

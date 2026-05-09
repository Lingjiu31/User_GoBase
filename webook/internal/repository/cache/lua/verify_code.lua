local key = KEYS[1]
 -- 用户输入的code
local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
 -- 转换成一个数字
local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 用户一直输错
    return -1
elseif expectedCode == code then
    -- 输入正确
    redis.call("set",cntKey, -1)
    return 0
else
    -- 输入错误
    redis.call("decr",cntKey, -1)
    return -2
end
package redis

import (
	kredis "dgo/framework/tools/db/redis"
	kinit "dgo/work/base/initialize"
	"errors"
	"github.com/garyburd/redigo/redis"
)

var (
	ErrSingleLockOperationFailed = errors.New("SingleLock : Operation failed")
	ErrSingleLockNotLocked       = errors.New("SingleLock : Not locked")
	ErrSingleLockLockIsUnlocked  = errors.New("SingleLock : Lock is unlocked")
)

var redisPool *kredis.RedisPool

func InitGetRedis(name string, dbIds ...int64) *kredis.RedisPool {
	redisPool, _ = kinit.GetRedisConnect(name, dbIds...)
	return redisPool
}

func InitGetReadRedis(name string, dbIds ...int64) *kredis.RedisPool {
	redisPool, _ = kinit.GetRedisReadConnect(name, dbIds...)
	return redisPool
}

func InitGetDivideRedis(name string, dbIds ...int64) *kredis.RedisPool {
	redisPool, _ = kinit.GetRedisDivideConnect(name, dbIds...)
	return redisPool
}

func InitGetReadDivideRedis(name string, dbIds ...int64) *kredis.RedisPool {
	redisPool, _ = kinit.GetRedisReadDivideConnect(name, dbIds...)
	return redisPool
}

func Set(key string, v interface{}, expire int64) error {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	if err := rc.Send("SET", key, v); err != nil {
		kinit.LogInfo.Println(err)
		return err
	}

	if expire > 0 {
		if err := rc.Send("EXPIRE", key, expire); err != nil {
			kinit.LogInfo.Println(err)
			return err
		}
	}
	if err := rc.Flush(); err != nil {
		kinit.LogInfo.Println(err)
		return err
	}
	return nil
}
func Incr(key string) error {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	if _, err := rc.Do("INCR", key); err != nil {
		kinit.LogInfo.Println(err)
		return err
	}
	return nil
}
func Decr(key string) error {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	if _, err := rc.Do("DECR", key); err != nil {
		kinit.LogInfo.Println(err)
		return err
	}
	return nil
}
func GetFloat64(key string) (float64, error) {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	v, err := redis.Float64(rc.Do("get", key))
	return v, err
}

func GetString(key string) (string, error) {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	v, err := redis.String(rc.Do("get", key))
	return v, err
}
func GetInt64(key string) (int64, error) {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	v, err := redis.Int64(rc.Do("get", key))
	return v, err
}

func GetIncr(key string) (int, error) {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	return redis.Int(rc.Do("incr", key))
}

func HgetAll(key string) (map[string]int64, error) {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	return redis.Int64Map(rc.Do("HGETALL", key))
}
func Hset(key string, subKey string, v interface{}) error {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	if err := rc.Send("HSET", key, subKey, v); err != nil {
		kinit.LogInfo.Println(err)
		return err
	}
	if err := rc.Flush(); err != nil {
		kinit.LogInfo.Println(err)
		return err
	}
	return nil
}
func Hdel(key string, subKey string) error {
	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)
	if err := rc.Send("HDEL", key, subKey); err != nil {
		kinit.LogInfo.Println(err)
		return err
	}
	if err := rc.Flush(); err != nil {
		kinit.LogInfo.Println(err)
		return err
	}
	return nil
}

func Lock(keyName string, value string, milliseconds int64) error {
	if len(keyName) == 0 ||
		len(value) == 0 {
		return ErrSingleLockNotLocked
	}

	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)

	//	try to lock
	rpl, err := redis.String(rc.Do("SET", keyName, value, "NX", "PX", milliseconds))
	if nil != err {
		return err
	}
	if rpl != "OK" {
		return ErrSingleLockOperationFailed
	}

	return nil
}

func Unlock(keyName string, value string) error {
	if len(keyName) == 0 ||
		len(value) == 0 {
		return ErrSingleLockNotLocked
	}

	rc := redisPool.GetRedis()
	defer redisPool.CloseRedis(rc)

	luaRelease := redis.NewScript(1, `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)

	rpl, err := redis.Int(luaRelease.Do(rc, keyName, value))
	if nil != err {
		return err
	}

	if rpl != 1 {
		return ErrSingleLockLockIsUnlocked
	}

	return nil
}

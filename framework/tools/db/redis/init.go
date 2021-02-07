package redis

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisPool struct {
	rp *redis.Pool
}

/**
 * option
      db          第几个库
	  MaxIdle     最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
	  MaxActive   最大的激活连接数，表示同时最多有N个连接
	  IdleTimeout 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
*/

func NewRedisPool(host, port, auth string, option ...string) (*RedisPool, error) {

	addr := fmt.Sprintf("%s:%s", host, port)

	countOption := len(option)
	initIdleTimeout, initMaxIdle := 60, 30

	RedisClient := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				log.Println("redis open fail:", err)
				return nil, err
			}
			if auth != "" {
				if _, err := c.Do("AUTH", auth); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			// 选择db
			if countOption > 0 && option[0] != "0" {
				_, _ = c.Do("SELECT", option[0])
			}
			return c, nil
		},
	}

	//设置最大的空闲连接数
	if countOption == 2 || countOption == 3 || countOption == 4 {
		MaxIdle, err := strconv.Atoi(option[1])
		if err != nil {
			log.Println("MaxIdle convert fail:", err)
			return nil, err
		}
		if MaxIdle > 0 {
			initMaxIdle = MaxIdle
		}
	}
	RedisClient.MaxIdle = initMaxIdle

	//设置最大的激活连接数
	if countOption == 3 || countOption == 4 {
		maxActive, err := strconv.Atoi(option[2])
		if err != nil {
			log.Println("MaxActive convert fail:", err)
			return nil, err
		}
		if maxActive > 0 {
			RedisClient.MaxActive = maxActive
		}
	}

	//设置最大的空闲连接等待时间
	if countOption == 4 {
		idleTimeout, err := strconv.Atoi(option[3])
		if err != nil {
			log.Println("IdleTimeout convert fail:", err)
			return nil, err
		}
		if idleTimeout > 0 {
			initIdleTimeout = idleTimeout
		}
	}
	RedisClient.IdleTimeout = time.Duration(initIdleTimeout) * time.Second

	//懒加载
	return &RedisPool{rp: RedisClient}, nil
}

func (ts *RedisPool) GetRedis() redis.Conn {
	rc := ts.rp.Get()
	return rc
}
func (ts *RedisPool) CloseRedis(rc redis.Conn) {
	_ = rc.Close()
}

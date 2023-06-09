package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = ""
)

// 创建 redis 连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		// 最大连接数
		MaxIdle: 50,
		// 在给定时间分配的最大连接数
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 打开连接
			conn, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}

			// 访问认证
			if _, err = conn.Do("AUTH", redisPass); err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}

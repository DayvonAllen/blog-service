package database

import (
	"github.com/gomodule/redigo/redis"
)

func ConnectToRedis() *redis.Pool{
	const maxConnections = 10
	redisPool := &redis.Pool{
		MaxIdle: maxConnections,
		Dial:    func() (redis.Conn, error) { return redis.Dial("tcp", "backend-redis-srv:6379") },
	}
	return redisPool
}

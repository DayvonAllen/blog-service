package database

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

func ConnectToRedis() *cache.Cache {
	red := cache.New(&cache.Options{
		Redis: redis.NewRing(&redis.RingOptions{
			Addrs: map[string]string{
				"server1": ":6379",
			},
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})


	return red
}

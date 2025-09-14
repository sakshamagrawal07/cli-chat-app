package redis

import (
	"time"
)

func MarkUserOnline(username string) {
	RedisClient.Set(ctx, "user:"+username+":online", "1", 5*time.Minute) // expires automatically
}

func MarkUserOffline(username string) {
	RedisClient.Del(ctx, "user:"+username+":online")
}

func IsUserOnline(username string) bool {
	val, err := RedisClient.Get(ctx, "user:"+username+":online").Result()
	return err == nil && val == "1"
}

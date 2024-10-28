package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <file-path>")
		return
	}
	var key string = os.Args[1]
	var redis_client *redis.Client
	limit, _ := strconv.Atoi(os.Args[2])
	connectToRedisServer(&redis_client)
	defer (redis_client).Close()
	for i := 0; i < limit; i++ {
		text, _ := redis_client.LPop(key).Result()
		fmt.Println(text)
	}
}

func connectToRedisServer(redis_client **redis.Client) {
	*redis_client = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := (*redis_client).Ping().Err(); err != nil {
		panic("unable to connect to Redis: " + err.Error())
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file-path>")
		return
	}
	var key string = os.Args[1]
	var redis_client *redis.Client
	var count int
	connectToRedisServer(&redis_client)
	defer (redis_client).Close()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		redis_client.RPush(key, scanner.Text())
		count++
		if count > 20 {
			redis_client.LPop(key)
			count--
		}
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

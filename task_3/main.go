package main

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

func main() {
	var lenOfArgs int
	var ns string
	var index int
	var array []string
	var redis_client *redis.Client
	readFromStdin(&array, &lenOfArgs, &ns)
	connectToRedisServer(&redis_client)
	redis_client.Incr(ns)
	id, _ := redis_client.Get(ns).Result()
	for i := 0; i < lenOfArgs-2; i++ {
		makeNewKey(&redis_client, ns, id, &index, array[i])
	}
}

func readFromStdin(array *[]string, lenOfArgs *int, ns *string) {
	*lenOfArgs = len(os.Args)
	if *lenOfArgs > 1 {
		*ns = os.Args[1]
		for i := 2; i < *lenOfArgs; i++ {
			*array = append(*array, os.Args[i])
		}
	}
}

func connectToRedisServer(redis_client **redis.Client) {
	*redis_client = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := (*redis_client).Ping().Err(); err != nil {
		panic("Unable to connect to Redis: " + err.Error())
	}
}

func makeNewKey(redis_client **redis.Client, ns string, id string, index *int, keyValue string) {
	keyName := ns + "-" + id + "-" + strconv.Itoa(*index)
	err := (*redis_client).Set(keyName, keyValue, 0).Err()
	if err != nil {
		panic(err)
	}
	*index += 1
}

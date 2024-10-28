package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file-path>")
		return
	}
	var filePath string = os.Args[1]
	var redis_client *redis.Client
	connectToRedisServer(&redis_client)
	defer (redis_client).Close()
	if err := printFile(redis_client, filePath); err != nil {
		fmt.Println("error:", err)
		return
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

func printFile(redis_client *redis.Client, filePath string) error {
	fileAbsolutePath := getAbsolutePath(filePath)
	if val, err := redis_client.Get(fileAbsolutePath).Result(); err != redis.Nil {
		fmt.Print(val)
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed getting information about file: %w", err)
		}
		fileSize := fileInfo.Size()
		if fileSize > 100 {
			_, err = file.Seek(-100, io.SeekEnd)
			if err != nil {
				return fmt.Errorf("failed in seeking the file: %w", err)
			}
		}
		data := make([]byte, 100)
		_, err = file.Read(data)
		if err != nil {
			return fmt.Errorf("failed reading data: %w", err)
		}

		fmt.Print(string(data))
		redis_client.Set(fileAbsolutePath, data, time.Second*60)
	}
	return nil
}

func getAbsolutePath(filePath string) string {
	if string(filePath[0]) != "/" {
		firstPath, _ := os.Getwd()
		filePath = firstPath + "/" + filePath
	}
	filePath = normalizePath(filePath)
	return filePath
}

func normalizePath(finalPath string) string {
	var stack []string
	slices := strings.Split(finalPath, "/")
	for i := 0; i < len(slices); i++ {
		if slices[i] == ".." {
			stack = stack[:len(stack)-1]
		} else if slices[i] != "." {
			stack = append(stack, slices[i])
		}
	}
	return strings.Join(stack, "/")
}

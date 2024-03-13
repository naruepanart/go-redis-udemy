package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	CommandsCode("This is CommandsCode Function")
}

func CommandsCode(n string) {
	fmt.Println(n)
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0, // Use the default DB
	})
	defer rdb.Close()

	ctx := context.Background()
	err := rdb.FlushAll(ctx).Err()
	if err != nil {
		panic(err)
	}

	err = rdb.Set(ctx, "key-string", "go-redis-no-expire", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := rdb.Get(ctx, "key-string").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("get key-string:", val)

	fmt.Println("---")

	err = rdb.Set(ctx, "key-string-60-sec", "go-redis-60-sec", 60*time.Second).Err()
	if err != nil {
		panic(err)
	}
	val, err = rdb.Get(ctx, "key-string-60-sec").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("get key-string-60-sec:", val)
}

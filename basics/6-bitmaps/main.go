package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "redis: ", log.Lshortfile)

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379", // Redis server address
		DB:   0,                // Use the default database
	})
	defer rdb.Close() // Close the Redis client connection when main function exits

	// Handle interrupt signal to allow for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Flush all keys from the database (optional)
	err := rdb.FlushAll(ctx).Err()
	if err != nil {
		logger.Fatal(err)
	}

	// Set bit at index
	err = rdb.SetBit(ctx, "bitmap-key", 0, 1).Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Bit set at index 0")

	// Get bit at index
	bit, err := rdb.GetBit(ctx, "bitmap-key", 0).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Bit at index 0:", bit)

	// Count bits in range
	count, err := rdb.BitCount(ctx, "bitmap-key", &redis.BitCount{Start: 0, End: -1}).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Bit count:", count)

	// Bitwise operations
	err = rdb.SetBit(ctx, "bitmap-key-2", 0, 1).Err()
	if err != nil {
		logger.Fatal(err)
	}
	err = rdb.SetBit(ctx, "bitmap-key-2", 1, 0).Err()
	if err != nil {
		logger.Fatal(err)
	}
	result, err := rdb.BitOpAnd(ctx, "bitmap-AND", "bitmap-key", "bitmap-key-2").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Result of BITOP AND:", result)

	// Get bit at index from resulting bitmap
	bitResult, err := rdb.GetBit(ctx, "bitmap-AND", 0).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Bit at index 0 in result of BITOP AND:", bitResult)
}

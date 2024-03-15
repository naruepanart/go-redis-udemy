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

	// Add elements to a HyperLogLog
	err = rdb.PFAdd(ctx, "hll-key", "element1", "element2", "element3").Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("HyperLogLog Add: element1, element2, element3 to hll-key")

	// Count unique elements in HyperLogLog
	count, err := rdb.PFCount(ctx, "hll-key").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("HyperLogLog Count:", count)

	// Add more elements to the HyperLogLog
	err = rdb.PFAdd(ctx, "hll-key", "element4", "element5").Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("HyperLogLog Add: element4, element5 to hll-key")

	// Count unique elements again
	count, err = rdb.PFCount(ctx, "hll-key").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("HyperLogLog Count After Adding More Elements:", count)
}
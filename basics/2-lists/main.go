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

	// Push elements to a list
	err = rdb.LPush(ctx, "key-list", "element1", "element2", "element3").Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("List Push: element1, element2, element3 to key-list")

	// Get range of elements from a list
	listRange, err := rdb.LRange(ctx, "key-list", 0, -1).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("List Range:", listRange)

	// Get length of the list
	listLen, err := rdb.LLen(ctx, "key-list").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("List Length:", listLen)

	// Remove elements from the list
	err = rdb.LRem(ctx, "key-list", 0, "element2").Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("List Remove: element2 from key-list")

	// Trim the list to a specific range
	err = rdb.LTrim(ctx, "key-list", 0, 1).Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("List Trim: key-list trimmed to range 0-1")

	// Insert element before a pivot element
	err = rdb.LInsert(ctx, "key-list", "BEFORE", "element1", "newElement").Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("List Insert Before: newElement inserted before element1")

	// Set the value of an element at a specific index
	err = rdb.LSet(ctx, "key-list", 1, "updatedElement").Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("List Set: element at index 1 updated")

	// Get the element at a specific index
	elementAtIndex, err := rdb.LIndex(ctx, "key-list", 1).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("List Get at Index 1:", elementAtIndex)
}

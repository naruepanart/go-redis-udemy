package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
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

	// Set a string key-value pair without expiration
	val, err := rdb.Set(ctx, "key-string", "go-redis-no-expire", 0).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("String Set: key-string", val)

	// Get value of a string key
	val, err = rdb.Get(ctx, "key-string").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("String Get:", val)

	// Set a string key-value pair with expiration time
	expiration := 60 * time.Second // Expiration time in seconds
	val, err = rdb.Set(ctx, "key-string-60-sec", "go-redis-60-sec", expiration).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("String Set: key-string-60-sec with 60 sec expiration", val)

	// Get value of a string key with expiration
	val, err = rdb.Get(ctx, "key-string-60-sec").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("String Get with expiration:", val)

	// Get the length of a string value
	strLen, err := rdb.StrLen(ctx, "key-string").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Length of key-string:", strLen)

	// Increment a key with integer value by 1
	incrVal, err := rdb.Incr(ctx, "key-incr").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Incremented value of key-incr:", incrVal)

	// Decrement a key with integer value by 1
	decrVal, err := rdb.Decr(ctx, "key-decr").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Decremented value of key-decr:", decrVal)

	// Increment a key with integer value by a specified amount
	incrByVal, err := rdb.IncrBy(ctx, "key-incrby", 5).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Incremented value of key-incrby by 5:", incrByVal)

	// Decrement a key with integer value by a specified amount
	decrByVal, err := rdb.DecrBy(ctx, "key-decrby", 3).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Decremented value of key-decrby by 3:", decrByVal)
}

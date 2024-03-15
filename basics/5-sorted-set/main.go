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

	// Add elements to a sorted set
	err = rdb.ZAdd(ctx, "sorted-set", &redis.Z{Score: 1, Member: "element1"}, &redis.Z{Score: 2, Member: "element2"}, &redis.Z{Score: 3, Member: "element3"}).Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Sorted Set Add: element1, element2, element3 to sorted-set")

	// Get range of elements from a sorted set
	sortedSetRange, err := rdb.ZRange(ctx, "sorted-set", 0, -1).Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Sorted Set Range:", sortedSetRange)

	// Get length of the sorted set
	sortedSetLen, err := rdb.ZCard(ctx, "sorted-set").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Sorted Set Length:", sortedSetLen)

	// Remove elements from the sorted set
	err = rdb.ZRem(ctx, "sorted-set", "element2").Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Sorted Set Remove: element2 from sorted-set")

	// Get the score of an element in the sorted set
	elementScore, err := rdb.ZScore(ctx, "sorted-set", "element1").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Sorted Set Score for element1:", elementScore)

	// Increment the score of an element in the sorted set
	newScore, err := rdb.ZIncrBy(ctx, "sorted-set", 5, "element1").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Sorted Set Increment Score for element1:", newScore)

	// Get the rank of an element in the sorted set
	elementRank, err := rdb.ZRank(ctx, "sorted-set", "element3").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Sorted Set Rank for element3:", elementRank)
}

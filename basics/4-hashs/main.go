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
	defer rdb.Close() // Close the Redis client connection when the main function exits

	// Handle interrupt signal to allow for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Flush all keys from the database (optional)
	err := rdb.FlushAll(ctx).Err()
	if err != nil {
		logger.Fatal("Failed to flush Redis database:", err)
	}

	// Example of using Redis Hash commands
	hashKey := "my-hash"

	// HSET: Sets the value of a field in a hash
	err = rdb.HSet(ctx, hashKey, "field1", "value1").Err()
	if err != nil {
		logger.Fatal("Failed to set hash field:", err)
	}
	logger.Println("Hash Set: field1 set to value1")

	// HGET: Retrieves the value of a specific field in a hash
	fieldValue, err := rdb.HGet(ctx, hashKey, "field1").Result()
	if err != nil {
		logger.Fatal("Failed to get hash field value:", err)
	}
	logger.Println("Hash Get: value for field1 -", fieldValue)

	// HMSET: Sets multiple field-value pairs in a hash in a single command
	err = rdb.HMSet(ctx, hashKey, map[string]interface{}{"field2": "value2", "field3": "value3"}).Err()
	if err != nil {
		logger.Fatal("Failed to set multiple hash fields:", err)
	}
	logger.Println("Hash MSet: field2 set to value2, field3 set to value3")

	// HMGET: Retrieves the values of multiple fields in a hash
	fields := []string{"field2", "field3"}
	fieldValues, err := rdb.HMGet(ctx, hashKey, fields...).Result()
	if err != nil {
		logger.Fatal("Failed to get multiple hash field values:", err)
	}
	logger.Println("Hash MGet:", fieldValues)

	// HDEL: Removes one or more fields from a hash
	err = rdb.HDel(ctx, hashKey, "field1").Err()
	if err != nil {
		logger.Fatal("Failed to delete hash field:", err)
	}
	logger.Println("Hash Del: field1 deleted")

	// HGETALL: Returns all fields and values for a hash
	allFieldValues, err := rdb.HGetAll(ctx, hashKey).Result()
	if err != nil {
		logger.Fatal("Failed to get all hash fields and values:", err)
	}
	logger.Println("Hash GetAll:", allFieldValues)

	// HEXISTS: Checks if a specific field exists in a hash
	fieldExists, err := rdb.HExists(ctx, hashKey, "field2").Result()
	if err != nil {
		logger.Fatal("Failed to check hash field existence:", err)
	}
	logger.Println("Hash Exists for field2:", fieldExists)

	// HINCRBY: Increments the numeric value of a field by a given integer
	err = rdb.HIncrBy(ctx, hashKey, "field4", 55).Err()
	if err != nil {
		logger.Fatal("Failed to increment hash field value:", err)
	}
	logger.Println("Hash IncrBy: field4 incremented by 55")

	// HINCRBYFLOAT: Increments the numeric value of a field by a given floating-point number
	err = rdb.HIncrByFloat(ctx, hashKey, "field5", 2.5).Err()
	if err != nil {
		logger.Fatal("Failed to increment hash field value by float:", err)
	}
	logger.Println("Hash IncrByFloat: field5 incremented by 2.5")

	// HKEYS: Returns all field names in a hash
	hashKeys, err := rdb.HKeys(ctx, hashKey).Result()
	if err != nil {
		logger.Fatal("Failed to get all hash field names:", err)
	}
	logger.Println("Hash Keys:", hashKeys)

	// HLEN: Returns the number of fields in a hash
	hashLen, err := rdb.HLen(ctx, hashKey).Result()
	if err != nil {
		logger.Fatal("Failed to get hash length:", err)
	}
	logger.Println("Hash Length:", hashLen)

	// HVALS: Returns all values in a hash
	hashValues, err := rdb.HVals(ctx, hashKey).Result()
	if err != nil {
		logger.Fatal("Failed to get all hash values:", err)
	}
	logger.Println("Hash Values:", hashValues)
}

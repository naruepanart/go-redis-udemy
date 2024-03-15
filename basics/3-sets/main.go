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

	// Add members to a set
	membersAddedCount, err := rdb.SAdd(ctx, "key-set", "member1", "member2", "member3").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Add:", membersAddedCount, "members added to key-set")

	// Get cardinality of a set
	cardinality, err := rdb.SCard(ctx, "key-set").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Cardinality:", cardinality)

	// Check if a member is present in a set
	isMember, err := rdb.SIsMember(ctx, "key-set", "member1").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Is member1 in set:", isMember)

	// Get all members of a set
	setMembers, err := rdb.SMembers(ctx, "key-set").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Members:", setMembers)

	// Remove members from a set
	membersRemovedCount, err := rdb.SRem(ctx, "key-set", "member2").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Remove:", membersRemovedCount, "members removed from key-set")

	// Remove and return a random member from a set
	randomMember, err := rdb.SPop(ctx, "key-set").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Pop:", "Removed random member:", randomMember)

	// Get one or more random members from a set
	randomMembers, err := rdb.SRandMember(ctx, "key-set").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set RandMember:", "Random members:", randomMembers)

	// Iterate over members of a set using a cursor
	var cursor uint64 = 0
	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.SScan(ctx, "key-set", cursor, "", 10).Result()
		if err != nil {
			logger.Fatal(err)
		}
		// Process keys
		logger.Println("Set Scan Result:", keys)
		if cursor == 0 {
			break
		}
	}

	// Perform set operations
	// Example: Union of multiple sets
	unionResult, err := rdb.SUnion(ctx, "key-set", "other-key-set").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Union Result:", unionResult)

	// Example: Intersection of multiple sets
	intersectionResult, err := rdb.SInter(ctx, "key-set", "other-key-set").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Intersection Result:", intersectionResult)

	// Example: Difference between multiple sets
	differenceResult, err := rdb.SDiff(ctx, "key-set", "other-key-set").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Difference Result:", differenceResult)

	// Example: Move a member from one set to another
	moveResult, err := rdb.SMove(ctx, "key-set", "other-key-set", "member1").Result()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Set Move Result:", moveResult)
}

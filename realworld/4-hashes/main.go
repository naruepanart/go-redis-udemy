package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()

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
		log.Fatal(err)
	}

	// Define routes
	app.Get("/like", func(c *fiber.Ctx) error {
		return findLikeCount(c, ctx, rdb)
	})
	app.Post("/like", func(c *fiber.Ctx) error {
		return updateLikeCount(c, ctx, rdb)
	})
	app.Put("/users-profile", func(c *fiber.Ctx) error {
		return updateUsersProfile(c, ctx, rdb)
	})

	// Start Fiber server
	log.Fatal(app.Listen(":3000"))
}

type UsersProfile struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	LikesCount    int    `json:"likes_count"`
	PostsCount    int    `json:"posts_count"`
	VisitorsCount int    `json:"visitors_count"`
}

type KeyValue struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}

const KEY_TESTER = "POST"

func findLikeCount(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	// Parse request body into KeyValue struct
	kv := KeyValue{}
	if err := c.BodyParser(&kv); err != nil {
		// Return error response if input JSON is invalid
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid input JSON"})
	}

	// Retrieve like count from Redis
	likeCount, err := rdb.HGet(ctx, KEY_TESTER+":"+kv.ID, "like_count").Result()
	if err != nil {
		// Return error response if Redis operation fails
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return like count as JSON response
	return c.JSON(likeCount)
}

func updateLikeCount(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	// Parse the request body JSON into a KeyValue struct
	kv := KeyValue{}
	if err := c.BodyParser(&kv); err != nil {
		// Return error response if JSON parsing fails
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid input JSON")
	}

	// Increment the like count for the given ID in Redis
	err := rdb.HIncrBy(ctx, KEY_TESTER+":"+kv.ID, "like_count", 1).Err()
	if err != nil {
		// Return error response if Redis operation fails
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return success status
	return c.SendStatus(fiber.StatusOK)
}

func updateUsersProfile(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	// Create a struct to hold user profile data
	userPro := UsersProfile{}

	// Parse the request body into the user profile struct
	if err := c.BodyParser(&userPro); err != nil {
		// If parsing fails, return an error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid input JSON"})
	}

	// Validate user ID
	if userPro.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// Prepare Redis pipeline
	pipe := rdb.TxPipeline()

	// Update user profile fields
	if userPro.Name != "" {
		pipe.HSet(ctx, KEY_TESTER+":"+userPro.ID, "name", userPro.Name)
	}
	if userPro.Email != "" {
		pipe.HSet(ctx, KEY_TESTER+":"+userPro.ID, "email", userPro.Email)
	}

	// Accumulate count increments
	counts := make(map[string]int)
	if userPro.LikesCount == 1 {
		counts["likes_count"]++
	}
	if userPro.PostsCount == 1 {
		counts["posts_count"]++
	}
	if userPro.VisitorsCount == 1 {
		counts["visitors_count"]++
	}

	// Execute batched Redis commands for count increments
	if len(counts) > 0 {
		for field, count := range counts {
			pipe.HIncrBy(ctx, KEY_TESTER+":"+userPro.ID, field, int64(count))
		}
	}

	// Execute the Redis pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		// If pipeline execution fails, return an internal server error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return a successful response
	return c.SendStatus(fiber.StatusOK)
}

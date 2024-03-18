package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
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
	app.Get("/votes", func(c *fiber.Ctx) error {
		return countVotes(c, ctx, rdb)
	})
	app.Post("/votes", func(c *fiber.Ctx) error {
		return createVotes(c, ctx, rdb)
	})

	// Start Fiber server
	log.Fatal(app.Listen(":3000"))
}

type KeyValue struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

const KEY_TESTER = "candidate:1"

func countVotes(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	// Retrieve the number of members (votes) in the set stored at the "KEY_TESTER" key in Redis.
	v, err := rdb.SCard(ctx, KEY_TESTER).Result()
	if err != nil {
		// If an error occurs while retrieving the count of votes from Redis, return an internal server error response.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Return the count of votes in the response body.
	return c.JSON(v)
}

func createVotes(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	// Define a struct to hold key-value pairs
	kv := KeyValue{}

	// Parse the request body into the KeyValue struct
	if err := c.BodyParser(&kv); err != nil {
		// If there's an error parsing the JSON body, return an error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid input JSON"})
	}

	// Check if the ID of the user already exists in the set stored in Redis
	exist, err := rdb.SIsMember(ctx, KEY_TESTER, kv.ID).Result()
	if err != nil {
		// If there's an error checking for existence in Redis, return an error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if exist {
		// If the user has already voted, return an error response
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "this user already voted"})
	}

	// If the user hasn't voted yet, add their ID to the set stored in Redis
	err = rdb.SAdd(ctx, KEY_TESTER, kv.ID).Err()
	if err != nil {
		// If there's an error adding the ID to Redis, return an error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return a successful response indicating that the vote was created
	return c.SendStatus(fiber.StatusCreated)
}

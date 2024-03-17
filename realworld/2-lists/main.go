package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
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
	app.Get("/posts", func(c *fiber.Ctx) error {
		return findPosts(c, ctx, rdb)
	})
	app.Post("/posts", func(c *fiber.Ctx) error {
		return createPosts(c, ctx, rdb)
	})
	app.Delete("/posts", func(c *fiber.Ctx) error {
		return deletePosts(c, ctx, rdb)
	})

	// Start Fiber server
	log.Fatal(app.Listen(":3000"))
}

type Posts struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Function to find posts
func findPosts(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	// Initialize variables to store a single post and an array of posts
	post := Posts{}
	posts := []Posts{}

	// Retrieve page number and count of posts per page from query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))   // Convert page parameter to integer, default to 1 if not provided
	count, _ := strconv.Atoi(c.Query("count", "5")) // Convert count parameter to integer, default to 5 if not provided

	// Calculate start and end indexes based on page and count
	start := (int64(page) - 1) * int64(count) // Calculate start index
	end := int64(page)*int64(count) - 1       // Calculate end index

	// Retrieve posts from Redis database within the specified range
	postJSON, err := rdb.LRange(ctx, "Keyyy", start, end).Result() // Retrieve posts from Redis
	if err != nil {
		// If there's an error, return an internal server error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Iterate over retrieved values
	for _, v := range postJSON {
		// Unmarshal each post from JSON format into 'post' struct
		if err := json.Unmarshal([]byte(v), &post); err != nil {
			// If there's an error during unmarshaling, return the error
			return err
		}
		// Append the unmarshaled post to the posts array
		posts = append(posts, post)
	}

	// Return the array of posts as a JSON response
	return c.JSON(posts)
}

func createPosts(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	// Create an empty struct to store posts.
	posts := Posts{}

	// Parse the request body and store it in the posts struct.
	if err := c.BodyParser(&posts); err != nil {
		// Return an error response if the input JSON is invalid.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid input JSON"})
	}

	// Convert the posts struct to JSON format.
	data, err := json.Marshal(posts)
	if err != nil {
		// Return an error response if there is an error in JSON marshaling.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if len(data) == 0 {
		// Return an error response if the JSON data is empty.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Empty JSON data"})
	}

	// Push the JSON data onto the left end of a Redis list with a specific key.
	if err := rdb.LPush(ctx, "Keyyy", data).Err(); err != nil {
		// Return an error response if there is an error in pushing data to Redis.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// If successful, return a 201 Created status without any additional data.
	return c.SendStatus(fiber.StatusCreated)
}

func deletePosts(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	// Create an empty struct to store posts.
	posts := Posts{}

	// Parse the request body and store it in the posts struct.
	if err := c.BodyParser(&posts); err != nil {
		// Return an error response if the input JSON is invalid.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid input JSON"})
	}

	// Convert the posts struct to JSON format.
	data, err := json.Marshal(posts)
	if err != nil {
		// Return an error response if there is an error in JSON marshaling.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Remove posts from the Redis list with the specified key.
	// LRem is a Redis command used to remove elements from a list.
	val, err := rdb.LRem(ctx, "Keyyy", 1, data).Result()
	if err != nil {
		// Return an error response if there is an error in deleting posts.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// If no posts were removed, return a 404 Not Found status.
	if val == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No posts found to delete"})
	}

	// If successful, return a 200 OK status without any additional data.
	return c.SendStatus(fiber.StatusOK)
}


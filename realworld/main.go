package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

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
	app.Get("/myip", myIPCache1)
	app.Get("/myip2", func(c *fiber.Ctx) error {
		return myIPCache2(c, ctx, rdb)
	})
	app.Post("/findip", findIP1)
	app.Post("/findip2", func(c *fiber.Ctx) error {
		return findIP2(c, ctx, rdb)
	})

	// Start Fiber server
	log.Fatal(app.Listen(":3000"))
}

// Struct for GeoIP data
type GeoIP struct {
	IP string `json:"ip"`
}

func myIPCache1(c *fiber.Ctx) error {
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch IP information")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read response body")
	}

	return c.Send(bodyBytes)
}

func findIP1(c *fiber.Ctx) error {
	geoIP := GeoIP{}

	if err := c.BodyParser(&geoIP); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input JSON"})
	}

	resp, err := http.Get("http://ip-api.com/json/" + geoIP.IP)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch IP information"})
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read response body"})
	}

	return c.Send(bodyBytes)
}

type IPAddress struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func myIPCache2(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	myIPKey := "myIP"
	ipAddress := IPAddress{}

	cachedIP, err := rdb.Get(ctx, myIPKey).Result()
	if err != nil && err != redis.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve cached IP information"})
	}

	if err := json.Unmarshal([]byte(cachedIP), &ipAddress); err == nil {
		return c.JSON(ipAddress)
	}

	req := fiber.Post("http://ip-api.com/json/").Body(c.Body())
	statusCode, body, _ := req.Bytes()

	if err := json.Unmarshal(body, &ipAddress); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unmarshal IP information"})
	}

	if err := rdb.Set(ctx, myIPKey, body, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to cache IP information"})
	}

	return c.Status(statusCode).JSON(ipAddress)
}

func findIP2(c *fiber.Ctx, ctx context.Context, rdb *redis.Client) error {
	geoIP := GeoIP{}

	if err := c.BodyParser(&geoIP); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input JSON"})
	}

	cachedIP, err := rdb.Get(ctx, geoIP.IP).Result()
	if err != nil && err != redis.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve cached IP information"})
	}

	ipAddress := IPAddress{}

	if err := json.Unmarshal([]byte(cachedIP), &ipAddress); err == nil {
		return c.JSON(ipAddress)
	}

	req := fiber.Post("http://ip-api.com/json/" + geoIP.IP).Body(c.Body())
	statusCode, body, _ := req.Bytes()

	if err := json.Unmarshal(body, &ipAddress); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unmarshal IP information"})
	}

	if err := rdb.Set(ctx, geoIP.IP, body, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to cache IP information"})
	}

	return c.Status(statusCode).JSON(ipAddress)
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// setupFiberApp creates a new Fiber app instance for benchmarking
func setupFiberApp() *fiber.App {
	// Create Fiber app with minimal config for better performance
	app := fiber.New(fiber.Config{
		Prefork:                   false,
		CaseSensitive:             true,
		StrictRouting:             true,
		ServerHeader:              "",
		AppName:                   "",
		DisableKeepalive:          true,
		DisableDefaultDate:        true,
		DisableDefaultContentType: true,
		DisableHeaderNormalizing:  true,
		DisableStartupMessage:     true,
	})

	// Simple GET handler
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Hello, Fiber Framework!")
	})

	// JSON response handler
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		user := User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}
		
		response := fiber.Map{
			"id":   id,
			"user": user,
		}
		
		return c.Status(http.StatusOK).JSON(response)
	})

	// POST handler with JSON body
	app.Post("/users", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}

		user.ID = 123
		response := fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}
		
		return c.Status(http.StatusCreated).JSON(response)
	})

	// Multiple route parameters
	app.Get("/users/:id/posts/:postId", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		postID := c.Params("postId")
		
		response := fiber.Map{
			"userId": userID,
			"postId": postID,
		}
		
		return c.Status(http.StatusOK).JSON(response)
	})

	// Query parameters
	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		limit := c.Query("limit")
		
		response := fiber.Map{
			"query": query,
			"limit": limit,
		}
		
		return c.Status(http.StatusOK).JSON(response)
	})

	return app
}

// Fiber Benchmark simple GET request
func BenchmarkFiberSimpleGET(b *testing.B) {
	app := setupFiberApp()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, err := app.Test(req, -1) // No timeout for benchmark consistency
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// Fiber Benchmark GET request with route parameter
func BenchmarkFiberGetWithParam(b *testing.B) {
	app := setupFiberApp()
	req := httptest.NewRequest("GET", "/user/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, 1)
		resp.Body.Close()
	}
}

// Fiber Benchmark GET request with multiple route parameters
func BenchmarkFiberGetWithMultipleParams(b *testing.B) {
	app := setupFiberApp()
	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}

// Fiber Benchmark GET request with query parameters
func BenchmarkFiberGetWithQueryParams(b *testing.B) {
	app := setupFiberApp()
	req := httptest.NewRequest("GET", "/search?q=golang&limit=10", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}

// Fiber Benchmark POST request with JSON payload
func BenchmarkFiberPostWithJSON(b *testing.B) {
	app := setupFiberApp()

	user := User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	jsonData, _ := json.Marshal(user)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}

// Fiber Benchmark routing performance with different number of routes
func BenchmarkFiberRouting10Routes(b *testing.B) {
	benchmarkFiberRouting(b, 10)
}

func BenchmarkFiberRouting100Routes(b *testing.B) {
	benchmarkFiberRouting(b, 100)
}

func BenchmarkFiberRouting1000Routes(b *testing.B) {
	benchmarkFiberRouting(b, 1000)
}

func benchmarkFiberRouting(b *testing.B, numRoutes int) {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "",
		AppName:       "",
	})

	// Create multiple routes
	for i := 0; i < numRoutes; i++ {
		route := "/route" + strconv.Itoa(i)
		routeIndex := i // Capture loop variable
		app.Get(route, func(c *fiber.Ctx) error {
			return c.Status(http.StatusOK).SendString("Route " + strconv.Itoa(routeIndex))
		})
	}

	// Test the last route (worst case scenario)
	req := httptest.NewRequest("GET", "/route"+strconv.Itoa(numRoutes-1), nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}

// Fiber Benchmark concurrent requests
func BenchmarkFiberConcurrentRequests(b *testing.B) {
	app := setupFiberApp()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, _ := app.Test(req, -1)
			resp.Body.Close()
		}
	})
}

// Fiber Benchmark memory allocations for different payload sizes
func BenchmarkFiberSmallPayload(b *testing.B) {
	benchmarkFiberPayloadSize(b, 100)
}

func BenchmarkFiberMediumPayload(b *testing.B) {
	benchmarkFiberPayloadSize(b, 1024)
}

func BenchmarkFiberLargePayload(b *testing.B) {
	benchmarkFiberPayloadSize(b, 10240)
}

func benchmarkFiberPayloadSize(b *testing.B, size int) {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "",
		AppName:       "",
	})

	app.Post("/data", func(c *fiber.Ctx) error {
		// Read and echo back the data
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}
		return c.Status(http.StatusOK).JSON(data)
	})

	// Create payload of specified size
	payload := make(map[string]string)
	for i := 0; i < size/10; i++ { // Approximate size control
		payload["key"+strconv.Itoa(i)] = "value" + strconv.Itoa(i)
	}

	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/data", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}

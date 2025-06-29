package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupGinApp creates a new Gin app instance for benchmarking
func setupGinApp() *gin.Engine {
	// Set Gin to release mode to avoid debug output affecting benchmarks
	gin.SetMode(gin.ReleaseMode)

	app := gin.New()

	// Simple GET handler
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin Framework!")
	})

	// JSON response handler
	app.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		user := User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}
		c.JSON(http.StatusOK, gin.H{
			"id":   id,
			"user": user,
		})
	})

	// POST handler with JSON body
	app.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		user.ID = 123
		c.JSON(http.StatusCreated, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		})
	})

	// Multiple route parameters
	app.GET("/users/:id/posts/:postId", func(c *gin.Context) {
		userID := c.Param("id")
		postID := c.Param("postId")

		c.JSON(http.StatusOK, gin.H{
			"userId": userID,
			"postId": postID,
		})
	})

	// Query parameters
	app.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		limit := c.Query("limit")

		c.JSON(http.StatusOK, gin.H{
			"query": query,
			"limit": limit,
		})
	})

	return app
}

// Gin Benchmark simple GET request
func BenchmarkGinSimpleGET(b *testing.B) {
	app := setupGinApp()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Gin Benchmark GET request with route parameter
func BenchmarkGinGetWithParam(b *testing.B) {
	app := setupGinApp()
	req := httptest.NewRequest("GET", "/user/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Gin Benchmark GET request with multiple route parameters
func BenchmarkGinGetWithMultipleParams(b *testing.B) {
	app := setupGinApp()
	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Gin Benchmark GET request with query parameters
func BenchmarkGinGetWithQueryParams(b *testing.B) {
	app := setupGinApp()
	req := httptest.NewRequest("GET", "/search?q=golang&limit=10", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Gin Benchmark POST request with JSON payload
func BenchmarkGinPostWithJSON(b *testing.B) {
	app := setupGinApp()

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
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Gin Benchmark routing performance with different number of routes
func BenchmarkGinRouting10Routes(b *testing.B) {
	benchmarkGinRouting(b, 10)
}

func BenchmarkGinRouting100Routes(b *testing.B) {
	benchmarkGinRouting(b, 100)
}

func BenchmarkGinRouting1000Routes(b *testing.B) {
	benchmarkGinRouting(b, 1000)
}

func benchmarkGinRouting(b *testing.B, numRoutes int) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()

	// Create multiple routes
	for i := 0; i < numRoutes; i++ {
		route := "/route" + strconv.Itoa(i)
		app.GET(route, func(c *gin.Context) {
			c.String(http.StatusOK, "Route "+strconv.Itoa(i))
		})
	}

	// Test the last route (worst case scenario)
	req := httptest.NewRequest("GET", "/route"+strconv.Itoa(numRoutes-1), nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Gin Benchmark concurrent requests
func BenchmarkGinConcurrentRequests(b *testing.B) {
	app := setupGinApp()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
		}
	})
}

// Gin Benchmark memory allocations for different payload sizes
func BenchmarkGinSmallPayload(b *testing.B) {
	benchmarkGinPayloadSize(b, 100)
}

func BenchmarkGinMediumPayload(b *testing.B) {
	benchmarkGinPayloadSize(b, 1024)
}

func BenchmarkGinLargePayload(b *testing.B) {
	benchmarkGinPayloadSize(b, 10240)
}

func benchmarkGinPayloadSize(b *testing.B, size int) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()

	app.POST("/data", func(c *gin.Context) {
		// Read and echo back the data
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		c.JSON(http.StatusOK, data)
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
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

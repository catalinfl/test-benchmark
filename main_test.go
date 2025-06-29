package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/catalinfl/tree-framework"
)

// setupApp creates a new app instance for benchmarking
func setupApp() *tree.Mux {
	app := tree.InitMux()

	// Simple GET handler
	app.GET("/", func(ctx *tree.Ctx) error {
		return ctx.SendString("Hello, Tree Framework!", http.StatusOK)
	})

	// JSON response handler
	app.GET("/user/:id", func(ctx *tree.Ctx) error {
		id, err := ctx.GetURLParam("id")
		if err != nil {
			return ctx.SendString("Invalid ID", http.StatusBadRequest)
		}
		user := User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}
		return ctx.SendJSON(tree.J{
			"id":   id,
			"user": user,
		}, http.StatusOK)
	})

	// POST handler with JSON body
	app.POST("/users", func(ctx *tree.Ctx) error {
		var user User
		if err := ctx.BindJSON(&user); err != nil {
			return ctx.SendJSON(tree.J{"error": "Invalid JSON"}, http.StatusBadRequest)
		}

		user.ID = 123
		return ctx.SendJSON(tree.J{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}, http.StatusCreated)
	})

	// Multiple route parameters
	app.GET("/users/:id/posts/:postId", func(ctx *tree.Ctx) error {
		userID, _ := ctx.GetURLParam("id")
		postID, _ := ctx.GetURLParam("postId")

		return ctx.SendJSON(tree.J{
			"userId": userID,
			"postId": postID,
		}, http.StatusOK)
	})

	// Query parameters
	app.GET("/search", func(ctx *tree.Ctx) error {
		query, _ := ctx.GetQuery("q")
		limit, _ := ctx.GetQuery("limit")

		return ctx.SendJSON(tree.J{
			"query": query,
			"limit": limit,
		}, http.StatusOK)
	})

	return app
}

// Benchmark simple GET request
func BenchmarkSimpleGET(b *testing.B) {
	app := setupApp()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Benchmark GET request with route parameter
func BenchmarkGetWithParam(b *testing.B) {
	app := setupApp()
	req := httptest.NewRequest("GET", "/user/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Benchmark GET request with multiple route parameters
func BenchmarkGetWithMultipleParams(b *testing.B) {
	app := setupApp()
	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Benchmark GET request with query parameters
func BenchmarkGetWithQueryParams(b *testing.B) {
	app := setupApp()
	req := httptest.NewRequest("GET", "/search?q=golang&limit=10", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Benchmark POST request with JSON payload
func BenchmarkPostWithJSON(b *testing.B) {
	app := setupApp()

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

// Benchmark routing performance with different number of routes
func BenchmarkRouting10Routes(b *testing.B) {
	benchmarkRouting(b, 10)
}

func BenchmarkRouting100Routes(b *testing.B) {
	benchmarkRouting(b, 100)
}

func BenchmarkRouting1000Routes(b *testing.B) {
	benchmarkRouting(b, 1000)
}

func benchmarkRouting(b *testing.B, numRoutes int) {
	app := tree.InitMux()

	// Create multiple routes
	for i := 0; i < numRoutes; i++ {
		route := "/route" + strconv.Itoa(i)
		app.GET(route, func(ctx *tree.Ctx) error {
			return ctx.SendString("Route "+strconv.Itoa(i), http.StatusOK)
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

// Benchmark concurrent requests
func BenchmarkConcurrentRequests(b *testing.B) {
	app := setupApp()
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

// Benchmark memory allocations for different payload sizes
func BenchmarkSmallPayload(b *testing.B) {
	benchmarkPayloadSize(b, 100)
}

func BenchmarkMediumPayload(b *testing.B) {
	benchmarkPayloadSize(b, 1024)
}

func BenchmarkLargePayload(b *testing.B) {
	benchmarkPayloadSize(b, 10240)
}

func benchmarkPayloadSize(b *testing.B, size int) {
	app := tree.InitMux()

	app.POST("/data", func(ctx *tree.Ctx) error {
		// Read and echo back the data
		var data map[string]interface{}
		if err := ctx.BindJSON(&data); err != nil {
			return ctx.SendJSON(tree.J{"error": "Invalid JSON"}, http.StatusBadRequest)
		}
		return ctx.SendJSON(tree.J(data), http.StatusOK)
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

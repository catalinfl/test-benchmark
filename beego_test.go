package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/beego/beego/v2/server/web"
)

// BeegoController for handling routes
type BeegoController struct {
	web.Controller
}

// Simple GET handler
func (c *BeegoController) Get() {
	c.Ctx.WriteString("Hello, Beego Framework!")
}

// User controller for user routes
type BeegoUserController struct {
	web.Controller
}

func (c *BeegoUserController) Get() {
	id := c.Ctx.Input.Param(":id")
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}
	
	response := map[string]interface{}{
		"id":   id,
		"user": user,
	}
	
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *BeegoUserController) Post() {
	var user User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	user.ID = 123
	response := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}
	
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Data["json"] = response
	c.ServeJSON()
}

// Multiple params controller
type BeegoMultiParamsController struct {
	web.Controller
}

func (c *BeegoMultiParamsController) Get() {
	userID := c.Ctx.Input.Param(":id")
	postID := c.Ctx.Input.Param(":postId")
	
	response := map[string]string{
		"userId": userID,
		"postId": postID,
	}
	
	c.Data["json"] = response
	c.ServeJSON()
}

// Search controller for query params
type BeegoSearchController struct {
	web.Controller
}

func (c *BeegoSearchController) Get() {
	query := c.GetString("q")
	limit := c.GetString("limit")
	
	response := map[string]string{
		"query": query,
		"limit": limit,
	}
	
	c.Data["json"] = response
	c.ServeJSON()
}

// setupBeegoApp creates a new Beego app instance for benchmarking
func setupBeegoApp() {
	// Disable logs for benchmarking
	web.BConfig.Log.AccessLogs = false
	web.BConfig.RunMode = "prod"
	
	// Simple GET handler
	web.Router("/", &BeegoController{})
	
	// JSON response handler
	web.Router("/user/:id", &BeegoUserController{})
	
	// POST handler
	web.Router("/users", &BeegoUserController{})
	
	// Multiple route parameters
	web.Router("/users/:id/posts/:postId", &BeegoMultiParamsController{})
	
	// Query parameters
	web.Router("/search", &BeegoSearchController{})
}

// Beego Benchmark simple GET request
func BenchmarkBeegoSimpleGET(b *testing.B) {
	setupBeegoApp()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		web.BeeApp.Handlers.ServeHTTP(w, req)
	}
}

// Beego Benchmark GET request with route parameter
func BenchmarkBeegoGetWithParam(b *testing.B) {
	setupBeegoApp()
	req := httptest.NewRequest("GET", "/user/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		web.BeeApp.Handlers.ServeHTTP(w, req)
	}
}

// Beego Benchmark GET request with multiple route parameters
func BenchmarkBeegoGetWithMultipleParams(b *testing.B) {
	setupBeegoApp()
	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		web.BeeApp.Handlers.ServeHTTP(w, req)
	}
}

// Beego Benchmark GET request with query parameters
func BenchmarkBeegoGetWithQueryParams(b *testing.B) {
	setupBeegoApp()
	req := httptest.NewRequest("GET", "/search?q=golang&limit=10", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		web.BeeApp.Handlers.ServeHTTP(w, req)
	}
}

// Beego Benchmark POST request with JSON payload
func BenchmarkBeegoPostWithJSON(b *testing.B) {
	setupBeegoApp()

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
		web.BeeApp.Handlers.ServeHTTP(w, req)
	}
}

// Dynamic route controllers for routing benchmarks
type BeegoRouteController struct {
	web.Controller
	RouteIndex int
}

func (c *BeegoRouteController) Get() {
	c.Ctx.WriteString("Route " + strconv.Itoa(c.RouteIndex))
}

// Beego Benchmark routing performance with different number of routes
func BenchmarkBeegoRouting10Routes(b *testing.B) {
	benchmarkBeegoRouting(b, 10)
}

func BenchmarkBeegoRouting100Routes(b *testing.B) {
	benchmarkBeegoRouting(b, 100)
}

func BenchmarkBeegoRouting1000Routes(b *testing.B) {
	benchmarkBeegoRouting(b, 1000)
}

func benchmarkBeegoRouting(b *testing.B, numRoutes int) {
	web.BConfig.Log.AccessLogs = false
	web.BConfig.RunMode = "prod"

	// Create multiple routes
	for i := 0; i < numRoutes; i++ {
		route := "/route" + strconv.Itoa(i)
		controller := &BeegoRouteController{RouteIndex: i}
		web.Router(route, controller)
	}

	// Test the last route (worst case scenario)
	req := httptest.NewRequest("GET", "/route"+strconv.Itoa(numRoutes-1), nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		web.BeeApp.Handlers.ServeHTTP(w, req)
	}
}

// Beego Benchmark concurrent requests
func BenchmarkBeegoConcurrentRequests(b *testing.B) {
	setupBeegoApp()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			web.BeeApp.Handlers.ServeHTTP(w, req)
		}
	})
}

// Data controller for payload benchmarks
type BeegoDataController struct {
	web.Controller
}

func (c *BeegoDataController) Post() {
	var data map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &data); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	
	c.Data["json"] = data
	c.ServeJSON()
}

// Beego Benchmark memory allocations for different payload sizes
func BenchmarkBeegoSmallPayload(b *testing.B) {
	benchmarkBeegoPayloadSize(b, 100)
}

func BenchmarkBeegoMediumPayload(b *testing.B) {
	benchmarkBeegoPayloadSize(b, 1024)
}

func BenchmarkBeegoLargePayload(b *testing.B) {
	benchmarkBeegoPayloadSize(b, 10240)
}

func benchmarkBeegoPayloadSize(b *testing.B, size int) {
	web.BConfig.Log.AccessLogs = false
	web.BConfig.RunMode = "prod"
	
	web.Router("/data", &BeegoDataController{})

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
		web.BeeApp.Handlers.ServeHTTP(w, req)
	}
}

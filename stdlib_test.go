package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// StandardHTTP Mux setup
func setupStandardHTTP() *http.ServeMux {
	mux := http.NewServeMux()

	// Simple GET handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Standard HTTP!"))
	})

	// JSON response handler with URL parameter parsing
	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// Extract ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/user/")
		if path == "" {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		
		user := User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}
		
		response := map[string]interface{}{
			"id":   path,
			"user": user,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	// POST handler with JSON body
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
			return
		}

		user.ID = 123
		response := map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	})

	// Multiple route parameters (simplified)
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// Simple parsing of /users/:id/posts/:postId
		path := strings.TrimPrefix(r.URL.Path, "/users/")
		parts := strings.Split(path, "/")
		
		if len(parts) >= 3 && parts[1] == "posts" {
			userID := parts[0]
			postID := parts[2]
			
			response := map[string]string{
				"userId": userID,
				"postId": postID,
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
		
		http.Error(w, "Not found", http.StatusNotFound)
	})

	// Query parameters
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		query := r.URL.Query().Get("q")
		limit := r.URL.Query().Get("limit")

		response := map[string]string{
			"query": query,
			"limit": limit,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	return mux
}

// Standard HTTP Benchmark simple GET request
func BenchmarkStandardHTTPSimpleGET(b *testing.B) {
	mux := setupStandardHTTP()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// Standard HTTP Benchmark GET request with route parameter
func BenchmarkStandardHTTPGetWithParam(b *testing.B) {
	mux := setupStandardHTTP()
	req := httptest.NewRequest("GET", "/user/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// Standard HTTP Benchmark GET request with multiple route parameters
func BenchmarkStandardHTTPGetWithMultipleParams(b *testing.B) {
	mux := setupStandardHTTP()
	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// Standard HTTP Benchmark GET request with query parameters
func BenchmarkStandardHTTPGetWithQueryParams(b *testing.B) {
	mux := setupStandardHTTP()
	req := httptest.NewRequest("GET", "/search?q=golang&limit=10", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// Standard HTTP Benchmark POST request with JSON payload
func BenchmarkStandardHTTPPostWithJSON(b *testing.B) {
	mux := setupStandardHTTP()

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
		mux.ServeHTTP(w, req)
	}
}

// Standard HTTP Benchmark routing performance with different number of routes
func BenchmarkStandardHTTPRouting10Routes(b *testing.B) {
	benchmarkStandardHTTPRouting(b, 10)
}

func BenchmarkStandardHTTPRouting100Routes(b *testing.B) {
	benchmarkStandardHTTPRouting(b, 100)
}

func BenchmarkStandardHTTPRouting1000Routes(b *testing.B) {
	benchmarkStandardHTTPRouting(b, 1000)
}

func benchmarkStandardHTTPRouting(b *testing.B, numRoutes int) {
	mux := http.NewServeMux()

	// Create multiple routes
	for i := 0; i < numRoutes; i++ {
		route := "/route" + strconv.Itoa(i)
		routeIndex := i // Capture loop variable
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Route " + strconv.Itoa(routeIndex)))
		})
	}

	// Test the last route (worst case scenario)
	req := httptest.NewRequest("GET", "/route"+strconv.Itoa(numRoutes-1), nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// Standard HTTP Benchmark concurrent requests
func BenchmarkStandardHTTPConcurrentRequests(b *testing.B) {
	mux := setupStandardHTTP()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
		}
	})
}

// Standard HTTP Benchmark memory allocations for different payload sizes
func BenchmarkStandardHTTPSmallPayload(b *testing.B) {
	benchmarkStandardHTTPPayloadSize(b, 100)
}

func BenchmarkStandardHTTPMediumPayload(b *testing.B) {
	benchmarkStandardHTTPPayloadSize(b, 1024)
}

func BenchmarkStandardHTTPLargePayload(b *testing.B) {
	benchmarkStandardHTTPPayloadSize(b, 10240)
}

func benchmarkStandardHTTPPayloadSize(b *testing.B, size int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// Read and echo back the data
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
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
		mux.ServeHTTP(w, req)
	}
}

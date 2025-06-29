# Tree Framework Benchmark Tests

This project contains comprehensive benchmark tests comparing the performance of `tree-framework` against popular Go web frameworks.

## Frameworks Tested

1. **Tree Framework** (`github.com/catalinfl/tree-framework`) - Your custom framework
2. **Gin** (`github.com/gin-gonic/gin`) - Popular high-performance HTTP web framework
3. **Fiber** (`github.com/gofiber/fiber/v2`) - Express.js inspired web framework built on Fasthttp
4. **Beego** (`github.com/beego/beego/v2`) - Full-featured MVC web framework
5. **Standard Library** (`net/http`) - Go's built-in HTTP package

## Test Categories

### Basic Operations
- Simple GET requests
- GET requests with route parameters
- GET requests with multiple route parameters
- GET requests with query parameters
- POST requests with JSON payload

### Performance Tests
- Routing performance with 10, 100, and 1000 routes
- Concurrent request handling
- Memory allocation tests with different payload sizes (small, medium, large)

## Running Benchmarks

### Run All Benchmarks
```powershell
go test -bench=. -benchmem
```

### Run Specific Framework Benchmarks
```powershell
# Tree Framework only
go test -bench=BenchmarkSimple -benchmem

# Gin only
go test -bench=BenchmarkGin -benchmem

# Fiber only
go test -bench=BenchmarkFiber -benchmem

# Beego only
go test -bench=BenchmarkBeego -benchmem

# Standard Library only
go test -bench=BenchmarkStandardHTTP -benchmem
```

### Run Specific Test Types
```powershell
# Routing performance tests
go test -bench=Routing -benchmem

# Concurrent tests
go test -bench=Concurrent -benchmem

# Payload size tests
go test -bench=Payload -benchmem
```

### Compare Specific Operations
```powershell
# Compare simple GET performance
go test -bench="SimpleGET|SimpleGet" -benchmem

# Compare JSON POST performance
go test -bench="PostWithJSON" -benchmem

# Compare routing performance
go test -bench="Routing100" -benchmem

# Compare all frameworks for a specific test
go test -bench="GetWithParam" -benchmem
```

## Understanding Results

Benchmark results show:
- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation (lower is better)
- **allocs/op**: Number of allocations per operation (lower is better)

Example output:
```
BenchmarkSimpleGET-8                    2000000    750 ns/op     96 B/op    3 allocs/op
BenchmarkGinSimpleGET-8                 1000000   1200 ns/op    144 B/op    5 allocs/op
BenchmarkFiberSimpleGET-8               3000000    400 ns/op     64 B/op    2 allocs/op
BenchmarkBeegoSimpleGET-8                800000   1500 ns/op    192 B/op    7 allocs/op
BenchmarkStandardHTTPSimpleGET-8        3000000    500 ns/op     48 B/op    2 allocs/op
```

## Files

- `main_test.go` - Tree Framework benchmarks
- `gin_test.go` - Gin framework benchmarks  
- `fiber_test.go` - Fiber framework benchmarks
- `beego_test.go` - Beego framework benchmarks
- `stdlib_test.go` - Standard library benchmarks
- `main.go` - Sample Tree Framework application

## Dependencies

- `github.com/catalinfl/tree-framework` - Your framework
- `github.com/gin-gonic/gin` - Gin framework for comparison
- `github.com/gofiber/fiber/v2` - Fiber framework for comparison
- `github.com/beego/beego/v2/server/web` - Beego framework for comparison

## Notes

- Gin is set to release mode (`gin.ReleaseMode`) to ensure fair performance comparison
- Fiber is configured with minimal settings for optimal performance
- Beego is set to production mode with access logs disabled
- All benchmarks use consistent testing methods for each framework
- Memory allocation tracking is enabled for all benchmarks
- Tests cover both CPU performance and memory efficiency
- Fiber uses `app.Test()` method which may have different overhead compared to direct `ServeHTTP` calls

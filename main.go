package main

import (
	"net/http"
	"regexp"

	"github.com/catalinfl/tree-framework"
)

// useRegex validates a string against a regex pattern
func useRegex(pattern, text string) bool {
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false
	}
	return matched
}

type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name" v:"required;minlen=3;maxlen=100;alphanumeric"`
	Description string   `json:"description" v:"required;"`
	Price       float64  `json:"price" v:"required;gt=0;lte=999999.99"`
	Category    string   `json:"category" v:"required;oneof=electronics clothing books home sports"`
	SKU         string   `json:"sku" v:"required;regex=^[A-Z]{3}\\d{5}$"`
	InStock     bool     `json:"in_stock" v:"required"`
	Tags        []string `json:"tags" v:"required;minlen=1;maxlen=5"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	app := tree.InitMux()

	// POST endpoint for creating a product with advanced validation
	app.POST("/product", func(c *tree.Ctx) error {
		var product Product

		// Bind JSON from request body to Product struct
		if err := c.BindJSON(&product); err != nil {
			return c.SendJSON(tree.J{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			}, http.StatusBadRequest)
		}

		// Simulate saving the product (in real app, save to database)
		product.ID = 12345 // Auto-generated ID

		return c.SendJSON(tree.J{
			"message": "Product created successfully",
			"product": product,
		}, http.StatusCreated)
	})

	// Additional regex validation endpoint for phone numbers
	app.GET("/validate/phone/:|^\\+?[1-9]\\d{1,14}$|", func(ctx *tree.Ctx) error {
		phone, err := ctx.RegexURLParam(1)
		if err != nil {
			return ctx.SendJSON(tree.J{
				"valid":   false,
				"error":   "Invalid phone format",
				"details": err.Error(),
			}, http.StatusBadRequest)
		}

		return ctx.SendJSON(tree.J{
			"valid":   true,
			"phone":   phone,
			"message": "Valid international phone format",
		}, http.StatusOK)
	})

	// Regex validation endpoint for email using RegexURLParam
	// Route pattern with regex: :|pattern| format
	app.GET("/validate/:|^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$|", func(ctx *tree.Ctx) error {
		email, err := ctx.RegexURLParam(1)
		if err != nil {
			return ctx.SendJSON(tree.J{
				"valid":   false,
				"error":   "Invalid email format - regex validation failed",
				"details": err.Error(),
			}, http.StatusBadRequest)
		}

		return ctx.SendJSON(tree.J{
			"valid":   true,
			"email":   email,
			"message": "Valid email format - passed regex validation",
		}, http.StatusOK)
	})

	// Multiple regex parameters example
	app.GET("/user/:|^[a-zA-Z0-9_]{3,20}$|/email/:|^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$|", func(ctx *tree.Ctx) error {
		// Get first regex param (username)
		username, err := ctx.RegexURLParam(1)
		if err != nil {
			return ctx.SendJSON(tree.J{
				"valid":   false,
				"error":   "Invalid username format",
				"details": err.Error(),
			}, http.StatusBadRequest)
		}

		// Get second regex param (email)
		email, err := ctx.RegexURLParam(2)
		if err != nil {
			return ctx.SendJSON(tree.J{
				"valid":   false,
				"error":   "Invalid email format",
				"details": err.Error(),
			}, http.StatusBadRequest)
		}

		return ctx.SendJSON(tree.J{
			"valid":    true,
			"username": username,
			"email":    email,
			"message":  "Both username and email are valid",
		}, http.StatusOK)
	})

	// GET endpoint for retrieving a product
	app.GET("/product/:id", func(c *tree.Ctx) error {
		id, err := c.GetURLParam("id")
		if err != nil {
			return c.SendString("Invalid ID", http.StatusBadRequest)
		}

		// Mock product data
		product := Product{
			ID:          123,
			Name:        "SampleProduct123",
			Description: "This is a sample product description with enough characters",
			Price:       299.99,
			Category:    "electronics",
			SKU:         "ABC12345",
			InStock:     true,
			Tags:        []string{"sample", "electronics", "gadget"},
		}

		return c.SendJSON(tree.J{
			"id":      id,
			"product": product,
		}, http.StatusOK)
	})

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

		// Simulate creating user
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
			"userId":  userID,
			"postId":  postID,
			"message": "Post retrieved successfully",
		}, http.StatusOK)
	})

	// Query parameters
	app.GET("/search", func(ctx *tree.Ctx) error {
		query, _ := ctx.GetQuery("q")
		limit, _ := ctx.GetQuery("limit")

		return ctx.SendJSON(tree.J{
			"query":   query,
			"limit":   limit,
			"results": "Sample search results",
		}, http.StatusOK)
	})

	app.StartExecuting()
}

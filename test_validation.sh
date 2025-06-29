#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🧪 Testing Tree Framework Product Validation${NC}"
echo "=============================================="

# Start the server in background
echo -e "${YELLOW}Starting server...${NC}"
go run main.go &
SERVER_PID=$!

# Wait for server to start
sleep 2

# Test valid product
echo -e "\n${GREEN}✅ Testing VALID product:${NC}"
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming123",
    "description": "High performance gaming laptop with advanced features and excellent build quality",
    "price": 1299.99,
    "category": "electronics",
    "sku": "ABC12345",
    "in_stock": true,
    "tags": ["gaming", "laptop", "electronics"]
  }' | jq .

# Test invalid name (too short)
echo -e "\n${RED}❌ Testing INVALID name (too short):${NC}"
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PC",
    "description": "Short name test - this should fail minlen validation",
    "price": 999.99,
    "category": "electronics",
    "sku": "DEF67890",
    "in_stock": true,
    "tags": ["test"]
  }' | jq .

# Test invalid name (non-alphanumeric)
echo -e "\n${RED}❌ Testing INVALID name (non-alphanumeric):${NC}"
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming@Laptop!",
    "description": "Name contains special characters which should fail alphanumeric validation",
    "price": 1299.99,
    "category": "electronics",
    "sku": "GHI11111",
    "in_stock": true,
    "tags": ["gaming"]
  }' | jq .

# Test invalid price (negative)
echo -e "\n${RED}❌ Testing INVALID price (negative):${NC}"
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "FreeProduct",
    "description": "This product has negative price which should fail gt validation",
    "price": -100.50,
    "category": "electronics",
    "sku": "MNO33333",
    "in_stock": true,
    "tags": ["free"]
  }' | jq .

# Test invalid category
echo -e "\n${RED}❌ Testing INVALID category:${NC}"
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "TestProduct",
    "description": "Product with invalid category that is not in the oneof list",
    "price": 299.99,
    "category": "invalidcategory",
    "sku": "PQR44444",
    "in_stock": true,
    "tags": ["test"]
  }' | jq .

# Test invalid SKU format
echo -e "\n${RED}❌ Testing INVALID SKU format:${NC}"
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Product",
    "description": "SKU does not match regex pattern - should be 3 letters + 5 digits",
    "price": 199.99,
    "category": "books",
    "sku": "INVALID1",
    "in_stock": true,
    "tags": ["book"]
  }' | jq .

# Test invalid tags (empty array)
echo -e "\n${RED}❌ Testing INVALID tags (empty):${NC}"
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "EmptyTags",
    "description": "This product has empty tags array which should fail minlen validation",
    "price": 99.99,
    "category": "home",
    "sku": "STU55555",
    "in_stock": true,
    "tags": []
  }' | jq .

# Test GET endpoint
echo -e "\n${GREEN}✅ Testing GET product by ID:${NC}"
curl -X GET http://localhost:8080/product/123 | jq .

# Kill the server
echo -e "\n${YELLOW}Stopping server...${NC}"
kill $SERVER_PID

echo -e "\n${GREEN}🎉 Testing completed!${NC}"

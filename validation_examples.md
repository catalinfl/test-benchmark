# Advanced Validation Examples for Tree Framework

## Test the validation endpoints with these curl commands

### 1. Valid User Registration Request
```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "email": "john.doe@example.com",
    "password": "MyP@ssw0rd123!",
    "first_name": "John",
    "last_name": "Doe",
    "age": 25,
    "phone": "+1234567890",
    "address": "123 Main Street, Anytown",
    "zip_code": "12345",
    "role": "user",
    "website": "https://www.johndoe.com",
    "bio": "Software developer with 5 years of experience",
    "accept_terms": true,
    "salary": 75000.50,
    "skills": ["JavaScript", "Go", "Python"],
    "country_code": "US",
    "ssn": "123-45-6789"
  }'
```

### 2. Invalid User Registration - Username too short
```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "usr",
    "email": "john.doe@example.com",
    "password": "MyP@ssw0rd123!",
    "first_name": "John",
    "last_name": "Doe",
    "age": 25,
    "phone": "+1234567890",
    "address": "123 Main Street, Anytown",
    "zip_code": "12345",
    "role": "user",
    "accept_terms": true,
    "salary": 75000.50,
    "skills": ["JavaScript"],
    "country_code": "US"
  }'
```

### 3. Invalid User Registration - Age under 18
```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "email": "john.doe@example.com",
    "password": "MyP@ssw0rd123!",
    "first_name": "John",
    "last_name": "Doe",
    "age": 16,
    "phone": "+1234567890",
    "address": "123 Main Street, Anytown",
    "zip_code": "12345",
    "role": "user",
    "accept_terms": true,
    "salary": 75000.50,
    "skills": ["JavaScript"],
    "country_code": "US"
  }'
```

### 4. Invalid User Registration - Invalid role
```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "email": "john.doe@example.com",
    "password": "MyP@ssw0rd123!",
    "first_name": "John",
    "last_name": "Doe",
    "age": 25,
    "phone": "+1234567890",
    "address": "123 Main Street, Anytown",
    "zip_code": "12345",
    "role": "superadmin",
    "accept_terms": true,
    "salary": 75000.50,
    "skills": ["JavaScript"],
    "country_code": "US"
  }'
```

### 5. Invalid User Registration - Terms not accepted
```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "email": "john.doe@example.com",
    "password": "MyP@ssw0rd123!",
    "first_name": "John",
    "last_name": "Doe",
    "age": 25,
    "phone": "+1234567890",
    "address": "123 Main Street, Anytown",
    "zip_code": "12345",
    "role": "user",
    "accept_terms": false,
    "salary": 75000.50,
    "skills": ["JavaScript"],
    "country_code": "US"
  }'
```

### 6. Valid Product Creation
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smartphone123",
    "description": "Latest smartphone with advanced features and great performance",
    "price": 699.99,
    "category": "electronics",
    "sku": "ELE12345",
    "in_stock": true,
    "tags": ["smartphone", "electronics", "mobile"]
  }'
```

### 7. Invalid Product - Price too low
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smartphone123",
    "description": "Latest smartphone with advanced features",
    "price": 0,
    "category": "electronics",
    "sku": "ELE12345",
    "in_stock": true,
    "tags": ["smartphone"]
  }'
```

### 8. Invalid Product - Wrong SKU format
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smartphone123",
    "description": "Latest smartphone with advanced features",
    "price": 699.99,
    "category": "electronics",
    "sku": "INVALID",
    "in_stock": true,
    "tags": ["smartphone"]
  }'
```

### 9. Test Validation Rules
```bash
curl -X POST http://localhost:8080/validate/test \
  -H "Content-Type: application/json" \
  -d '{
    "test_field": "any_value",
    "number_field": 123,
    "boolean_field": true
  }'
```

### 10. Health Check
```bash
curl -X GET http://localhost:8080/health
```

## Expected Validation Rules

### Username
- **required**: Must be present
- **len=5**: Must be exactly 5 characters
- **alphanumeric**: Only letters and numbers

### Email
- **required**: Must be present
- **email**: Must be valid email format
- **regex**: Must match email regex pattern

### Password
- **required**: Must be present
- **minlen=8**: Minimum 8 characters
- **maxlen=50**: Maximum 50 characters
- **contains=@#$%^&***: Must contain at least one special character

### Age
- **required**: Must be present
- **gte=18**: Must be 18 or older
- **lte=120**: Must be 120 or younger

### Phone
- **required**: Must be present
- **regex**: Must match international phone format (+1234567890)

### Role
- **required**: Must be present
- **oneof**: Must be one of: admin, user, moderator, guest

### Country Code
- **required**: Must be present
- **len=2**: Must be exactly 2 characters
- **regex**: Must be uppercase letters (e.g., US, UK, DE)

### SKU (Product)
- **required**: Must be present
- **len=8**: Must be exactly 8 characters
- **regex**: Must follow pattern: 3 uppercase letters + 5 digits (e.g., ELE12345)

## Running the Example

1. Start the server:
```bash
go run advanced_validator_example.go
```

2. Test the endpoints using the curl commands above

3. Observe the validation responses for both valid and invalid data

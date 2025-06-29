# Test Examples for RegexURLParam in Tree Framework

## How RegexURLParam Works

The `RegexURLParam` method extracts and validates URL segments based on regex patterns defined in the route.

### Route Pattern Format: `:|regex_pattern|`

## Test Cases

### 1. Email Validation

**Route:** `/validate/:|^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$|`

**Valid Tests:**
```bash
# Valid email addresses
curl "http://localhost:8080/validate/test@example.com"
curl "http://localhost:8080/validate/user.name+tag@domain.co.uk"
curl "http://localhost:8080/validate/john.doe123@company-name.org"
```

**Expected Response:**
```json
{
  "valid": true,
  "email": "test@example.com",
  "message": "Valid email format - passed regex validation"
}
```

**Invalid Tests:**
```bash
# Invalid email addresses
curl "http://localhost:8080/validate/invalid-email"
curl "http://localhost:8080/validate/test@"
curl "http://localhost:8080/validate/@domain.com"
curl "http://localhost:8080/validate/test.domain.com"
```

**Expected Response:**
```json
{
  "valid": false,
  "error": "Invalid email format - regex validation failed",
  "details": "regex not respected"
}
```

### 2. Phone Number Validation

**Route:** `/validate/phone/:|^\\+?[1-9]\\d{1,14}$|`

**Valid Tests:**
```bash
# Valid phone numbers
curl "http://localhost:8080/validate/phone/+1234567890"
curl "http://localhost:8080/validate/phone/1234567890"
curl "http://localhost:8080/validate/phone/+40123456789"
```

**Expected Response:**
```json
{
  "valid": true,
  "phone": "+1234567890",
  "message": "Valid international phone format"
}
```

**Invalid Tests:**
```bash
# Invalid phone numbers
curl "http://localhost:8080/validate/phone/+0123456789"  # starts with 0
curl "http://localhost:8080/validate/phone/123"         # too short
curl "http://localhost:8080/validate/phone/+abc123"     # contains letters
```

### 3. Multiple Regex Parameters

**Route:** `/user/:|^[a-zA-Z0-9_]{3,20}$|/email/:|^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$|`

**Valid Tests:**
```bash
# Valid username and email combination
curl "http://localhost:8080/user/john_doe123/email/john@example.com"
curl "http://localhost:8080/user/user_name/email/test.email@domain.org"
```

**Expected Response:**
```json
{
  "valid": true,
  "username": "john_doe123",
  "email": "john@example.com",
  "message": "Both username and email are valid"
}
```

**Invalid Tests:**
```bash
# Invalid username (too short)
curl "http://localhost:8080/user/ab/email/test@example.com"

# Invalid email
curl "http://localhost:8080/user/valid_user/email/invalid-email"

# Both invalid
curl "http://localhost:8080/user/ab/email/invalid"
```

## How RegexURLParam Index Works

The `RegexURLParam(index)` method works as follows:

1. **Index 1**: First regex pattern in the route
2. **Index 2**: Second regex pattern in the route
3. **Index N**: Nth regex pattern in the route

### Example Route Breakdown:
```
/user/:|^[a-zA-Z0-9_]{3,20}$|/email/:|^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$|
       ^                     ^        ^
       |                     |        |
   Index 1               Index 2   Index 2 continues
   (username)            (email)
```

## Testing with curl

### Start the server:
```bash
go run main.go
```

### Test all endpoints:
```bash
# Test valid email
curl "http://localhost:8080/validate/john.doe@example.com"

# Test invalid email
curl "http://localhost:8080/validate/not-an-email"

# Test valid phone
curl "http://localhost:8080/validate/phone/+1234567890"

# Test invalid phone
curl "http://localhost:8080/validate/phone/123"

# Test valid user/email combo
curl "http://localhost:8080/user/john_doe/email/john@example.com"

# Test invalid username
curl "http://localhost:8080/user/ab/email/john@example.com"
```

## Error Types

1. **Regex Pattern Not Matched**: When the URL segment doesn't match the regex
2. **Invalid Index**: When requesting a regex parameter that doesn't exist
3. **Malformed Route**: When the route pattern is incorrect

## Performance Considerations

- Regex compilation happens on each request
- Complex regex patterns may impact performance
- Consider caching compiled regex patterns for high-traffic endpoints

## Regex Patterns Used

1. **Email**: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$`
2. **Phone**: `^\\+?[1-9]\\d{1,14}$`
3. **Username**: `^[a-zA-Z0-9_]{3,20}$`

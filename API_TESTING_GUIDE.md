# Prothomuse Server - API Testing Guide

## Prerequisites
- Go server running on `http://localhost:8080`
- PostgreSQL database `postgres` with `users` table created
- Test data will be created during testing

---

## 1. REGISTER USER

**Endpoint:** `POST /api/auth/register`

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**cURL Command:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","email":"john@example.com","password":"SecurePass123!"}'
```

**Expected Response (201 Created):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "apiKey": "sk_live_abcd1234efgh5678ijkl...",
    "isActive": true
  },
  "message": "user registered successfully"
}
```

**Error Response (400 Bad Request - Email Already Exists):**
```json
{
  "status": "error",
  "message": "user with this email already exists"
}
```

---

## 2. LOGIN USER

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**cURL Command:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"SecurePass123!"}'
```

**Expected Response (200 OK):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "apiKey": "sk_live_abcd1234efgh5678ijkl...",
    "username": "john_doe"
  },
  "message": "user logged in successfully"
}
```

**Save the `token` and `apiKey` for use in the next requests!**

**Error Response (401 Unauthorized - Invalid Password):**
```json
{
  "status": "error",
  "message": "invalid password"
}
```

---

## 3. UPDATE USER

**Endpoint:** `PUT /api/auth/update`

**Authorization:** `Bearer <JWT_TOKEN>` (from login response)

**Request Body (update one or more fields):**
```json
{
  "username": "jane_doe",
  "email": "jane@example.com"
}
```

**cURL Command:**
```bash
# Replace TOKEN with the JWT token from login response
curl -X PUT http://localhost:8080/api/auth/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{"username":"jane_doe","email":"jane@example.com"}'
```

**Expected Response (200 OK):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "username": "jane_doe",
    "email": "jane@example.com",
    "isActive": true
  },
  "message": "user updated successfully"
}
```

**Error Response (401 Unauthorized - Missing/Invalid Token):**
```json
{
  "status": "error",
  "message": "invalid or expired token"
}
```

---

## 4. VALIDATE JWT TOKEN

**Endpoint:** `GET /api/auth/validate-jwt`

**Authorization:** `Bearer <JWT_TOKEN>`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/api/auth/validate-jwt \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Expected Response (200 OK):**
```json
{
  "status": "success",
  "data": {
    "userId": 1,
    "email": "jane@example.com"
  },
  "message": "JWT token is valid"
}
```

**Error Response (401 Unauthorized - Invalid Token):**
```json
{
  "status": "error",
  "message": "invalid or expired token"
}
```

---

## 5. VALIDATE API KEY

**Endpoint:** `GET /api/auth/validate-apikey`

**Authorization:** `ApiKey <API_KEY>` (from register or login response)

**cURL Command:**
```bash
# Replace API_KEY with the apiKey from register/login response
curl -X GET http://localhost:8080/api/auth/validate-apikey \
  -H "Authorization: ApiKey sk_live_abcd1234efgh5678ijkl..."
```

**Expected Response (200 OK):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "username": "jane_doe",
    "email": "jane@example.com",
    "isActive": true
  },
  "message": "API key is valid"
}
```

**Error Response (401 Unauthorized - Missing/Invalid API Key):**
```json
{
  "status": "error",
  "message": "API key is required"
}
```

---

## 6. HEALTH CHECK

**Endpoint:** `GET /health`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/health
```

**Expected Response (200 OK):**
```json
{
  "status": "healthy",
  "service": "prothomuse-health-server"
}
```

---

## COMPLETE TEST FLOW (Step-by-Step)

### Step 1: Register a user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"Test@1234"}'
```
**Save the `apiKey` from response**

---

### Step 2: Login with the same credentials
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test@1234"}'
```
**Save the `token` from response**

---

### Step 3: Validate JWT Token (replace TOKEN)
```bash
curl -X GET http://localhost:8080/api/auth/validate-jwt \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

### Step 4: Validate API Key (replace API_KEY)
```bash
curl -X GET http://localhost:8080/api/auth/validate-apikey \
  -H "Authorization: ApiKey YOUR_API_KEY_HERE"
```

---

### Step 5: Update User Profile (replace TOKEN)
```bash
curl -X PUT http://localhost:8080/api/auth/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"username":"updateduser","password":"NewPass@5678"}'
```

---

### Step 6: Health Check
```bash
curl -X GET http://localhost:8080/health
```

---

## PowerShell Testing Examples

If you prefer PowerShell instead of curl:

### Register User (PowerShell)
```powershell
$body = @{
    username = "testuser"
    email = "test@example.com"
    password = "Test@1234"
} | ConvertTo-Json

Invoke-WebRequest -Uri http://localhost:8080/api/auth/register `
  -Method POST `
  -Headers @{"Content-Type" = "application/json"} `
  -Body $body
```

### Login (PowerShell)
```powershell
$body = @{
    email = "test@example.com"
    password = "Test@1234"
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri http://localhost:8080/api/auth/login `
  -Method POST `
  -Headers @{"Content-Type" = "application/json"} `
  -Body $body

$token = ($response.Content | ConvertFrom-Json).data.token
Write-Output "JWT Token: $token"
```

### Update User (PowerShell)
```powershell
$body = @{
    username = "newusername"
} | ConvertTo-Json

Invoke-WebRequest -Uri http://localhost:8080/api/auth/update `
  -Method PUT `
  -Headers @{
      "Content-Type" = "application/json"
      "Authorization" = "Bearer $token"
  } `
  -Body $body
```

---

## Common Errors & Solutions

| Error | Cause | Solution |
|-------|-------|----------|
| `pq: column "id" does not exist` | `users` table not created or incorrect schema | Run SQL migration in PostgreSQL |
| `user not found` | Email doesn't exist in database | Register user first |
| `invalid password` | Wrong password provided | Check password spelling |
| `invalid or expired token` | JWT token is invalid or expired | Login again to get new token |
| `another user with this email already exists` | Email already registered | Use a different email |
| `connection refused` | Server not running | Start server: `go run ./cmd/server` |

---

## Important Notes

1. **JWT Token expiry:** Tokens expire after 24 hours
2. **API Key:** Regenerate by registering again (or update via database)
3. **Password:** Minimum 6 characters required
4. **Email:** Must be valid and unique
5. **Database:** Must have `users` table with correct schema

---

## Testing in Postman

If using Postman instead of curl:

1. Create new requests for each endpoint
2. Set method (GET, POST, PUT)
3. Enter URL: `http://localhost:8080/api/auth/register` (etc.)
4. For POST/PUT: Set Body → raw → JSON and paste request body
5. For auth: Set Headers tab → Add `Authorization: Bearer <token>` or `Authorization: ApiKey <key>`
6. Click Send

---

**Ready to test? Start with Step 1 above and share any error messages!**

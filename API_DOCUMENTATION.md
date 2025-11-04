# API Documentation Guide

This document explains how to access and use the Template Store API documentation.

## Overview

The Template Store API uses **OpenAPI 3.0** (Swagger) for API documentation. The documentation provides:

- Complete API endpoint reference
- Request/response schemas
- Authentication requirements
- Interactive API testing
- Code examples

## Accessing the Documentation

### Option 1: Interactive Swagger UI (Recommended)

Once the server is running, access the interactive documentation at:

```
http://localhost:8080/api-docs
```

This provides:
- ✅ Interactive API explorer
- ✅ Try-it-out functionality to test endpoints
- ✅ Request/response examples
- ✅ Schema definitions
- ✅ Authentication token management

### Option 2: Raw OpenAPI Specification

Access the raw OpenAPI YAML specification at:

```
http://localhost:8080/swagger.yaml
```

Use this file to:
- Import into API tools (Postman, Insomnia, etc.)
- Generate client SDKs
- Integrate with API gateways

## Using the Interactive Documentation

### 1. Start the Server

```bash
# Development mode (SQLite, no authentication)
GIN_MODE=debug go run cmd/server/main.go

# Or use the compiled binary
GIN_MODE=debug ./template-store-server
```

### 2. Open the Documentation

Navigate to `http://localhost:8080/api-docs` in your browser.

### 3. Explore API Endpoints

The documentation is organized by tags:

- **Authentication** - User registration and login
- **Profile** - User profile management
- **Preferences** - User settings and preferences
- **Security** - Password management and email verification
- **Dashboard** - User dashboard and statistics
- **Sessions** - Login history and session management
- **Activity** - Activity logging and audit trail
- **Admin** - Admin panel operations
- **Public** - Public endpoints (no authentication)

### 4. Test Endpoints

#### Without Authentication (Public Endpoints)

1. Find an endpoint (e.g., `GET /api/v1/users/{id}/profile`)
2. Click "Try it out"
3. Fill in required parameters
4. Click "Execute"
5. View the response

#### With Authentication (Protected Endpoints)

1. **Get an authentication token:**
   - In development mode, authentication is bypassed
   - In production, register/login to get a JWT token

2. **Set the Bearer token:**
   - Click the "Authorize" button at the top right
   - Enter: `Bearer YOUR_TOKEN_HERE`
   - Click "Authorize"
   - Click "Close"

3. **Test protected endpoints:**
   - All subsequent requests will include the token
   - Test any authenticated endpoint

## API Endpoint Categories

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |
| POST | `/api/v1/auth/forgot-password` | Request password reset |
| POST | `/api/v1/auth/reset-password` | Reset password with token |
| POST | `/api/v1/auth/verify-email` | Verify email address |
| POST | `/api/v1/auth/change-password` | Change password (auth required) |

### Profile Management Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/profile` | Get user profile |
| PUT | `/api/v1/profile` | Update profile |
| POST | `/api/v1/profile/avatar` | Upload avatar |
| DELETE | `/api/v1/profile/avatar` | Delete avatar |
| GET | `/api/v1/profile/preferences` | Get preferences |
| PUT | `/api/v1/profile/preferences` | Update preferences |
| POST | `/api/v1/profile/deactivate` | Deactivate account |
| DELETE | `/api/v1/profile` | Delete account |

### Dashboard Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/profile/dashboard` | Get dashboard data |
| GET | `/api/v1/profile/orders` | Get user orders |
| GET | `/api/v1/profile/orders/{id}` | Get order details |
| GET | `/api/v1/profile/purchased-templates` | Get purchased templates |
| GET | `/api/v1/profile/templates/{id}/download` | Download template |
| GET | `/api/v1/profile/blog-posts` | Get user blog posts |

### Security & Session Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/profile/login-history` | Get login history |
| GET | `/api/v1/profile/sessions` | Get active sessions |
| GET | `/api/v1/profile/activity` | Get activity log |
| POST | `/api/v1/auth/resend-verification` | Resend verification email |
| POST | `/api/v1/auth/logout-session/{id}` | Logout specific session |
| POST | `/api/v1/auth/logout-all` | Logout all sessions |

### Admin Endpoints (Requires Admin Role)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/dashboard` | Get admin dashboard |
| GET | `/api/v1/admin/users` | List all users |
| GET | `/api/v1/admin/users/{id}` | Get user details |
| PUT | `/api/v1/admin/users/{id}/role` | Update user role |
| POST | `/api/v1/admin/users/{id}/suspend` | Suspend user |
| POST | `/api/v1/admin/users/{id}/unsuspend` | Unsuspend user |
| DELETE | `/api/v1/admin/users/{id}` | Delete user |
| GET | `/api/v1/admin/users/{id}/activity` | Get user activity |

### Public Endpoints (No Authentication)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/{id}/profile` | Get public user profile |

## Common Request Examples

### Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "name": "John Doe"
  }'
```

### Get User Profile

```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Update Profile

```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Updated",
    "bio": "Software developer",
    "phone_number": "+1234567890"
  }'
```

### Get Dashboard

```bash
curl -X GET http://localhost:8080/api/v1/profile/dashboard \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### List Users (Admin Only)

```bash
curl -X GET "http://localhost:8080/api/v1/admin/users?limit=20&offset=0" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE"
```

## Response Format

### Success Response

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "status": "active",
  "email_verified": true,
  "created_at": "2025-10-28T12:00:00Z",
  "updated_at": "2025-10-28T12:00:00Z"
}
```

### Error Response

```json
{
  "error": "Invalid request parameters"
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error |

## Authentication

### JWT Bearer Token

Most endpoints require JWT authentication via AWS Cognito.

**Header Format:**
```
Authorization: Bearer <your_jwt_token>
```

### Development Mode

In development mode (`GIN_MODE=debug`), authentication is bypassed and all requests are treated as user ID 1.

### Production Mode

In production, you must:
1. Register or login to get a JWT token
2. Include the token in the Authorization header
3. Token expires after a configured time period

## Pagination

Endpoints that return lists support pagination:

**Query Parameters:**
- `limit` - Number of items per page (default: 10-20, max: 100)
- `offset` - Number of items to skip (default: 0)

**Example:**
```
GET /api/v1/profile/orders?limit=20&offset=40
```

**Response includes:**
```json
{
  "orders": [...],
  "total": 100,
  "limit": 20,
  "offset": 40
}
```

## Filtering

Admin user list endpoint supports filtering:

**Query Parameters:**
- `search` - Search by name or email
- `role` - Filter by role (user, admin, moderator)
- `status` - Filter by status (active, suspended, inactive)

**Example:**
```
GET /api/v1/admin/users?role=admin&status=active&search=john
```

## File Uploads

Avatar upload requires multipart/form-data:

```bash
curl -X POST http://localhost:8080/api/v1/profile/avatar \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "avatar=@/path/to/image.jpg"
```

**Constraints:**
- Max file size: 5MB
- Allowed formats: jpg, jpeg, png, webp

## Rate Limiting

(To be implemented)

Rate limits will be applied to prevent abuse:
- 100 requests per minute per IP
- 1000 requests per hour per user

## CORS

CORS is enabled for all origins in development mode.

For production, configure allowed origins in the server configuration.

## Import to API Tools

### Postman

1. Open Postman
2. Click "Import"
3. Select "Link" tab
4. Enter: `http://localhost:8080/swagger.yaml`
5. Click "Continue"
6. Click "Import"

### Insomnia

1. Open Insomnia
2. Click "Create" → "File"
3. Select "From URL"
4. Enter: `http://localhost:8080/swagger.yaml`
5. Click "Fetch and Import"

## Generate Client SDKs

Use the OpenAPI specification to generate client SDKs:

### JavaScript/TypeScript

```bash
npx @openapitools/openapi-generator-cli generate \
  -i http://localhost:8080/swagger.yaml \
  -g typescript-axios \
  -o ./client-sdk
```

### Python

```bash
openapi-generator-cli generate \
  -i http://localhost:8080/swagger.yaml \
  -g python \
  -o ./client-sdk
```

### Go

```bash
openapi-generator-cli generate \
  -i http://localhost:8080/swagger.yaml \
  -g go \
  -o ./client-sdk
```

## Additional Resources

- [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.3)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [AWS Cognito Documentation](https://docs.aws.amazon.com/cognito/)
- [Template Store GitHub Repository](https://github.com/yourusername/template-store)

## Support

For questions or issues:
- GitHub Issues: [Create an issue](https://github.com/yourusername/template-store/issues)
- Email: support@example.com

## Version History

- **v1.0.0** (2025-10-28) - Initial API documentation
  - Complete user account management endpoints
  - Profile, preferences, dashboard, security
  - Admin panel operations
  - Activity logging and audit trail

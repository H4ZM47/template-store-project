# Blog Management Guide

This guide provides comprehensive documentation for creating and managing blog posts in the Template Store platform.

## Table of Contents

1. [Overview](#overview)
2. [Blog Post Structure](#blog-post-structure)
3. [Authentication & Authorization](#authentication--authorization)
4. [API Endpoints](#api-endpoints)
5. [Creating Blog Posts](#creating-blog-posts)
6. [Managing Existing Blog Posts](#managing-existing-blog-posts)
7. [Querying Blog Posts](#querying-blog-posts)
8. [Best Practices](#best-practices)
9. [Code Examples](#code-examples)

---

## Overview

The Template Store platform includes a full-featured blog system that supports:

- **Markdown content** with automatic HTML conversion
- **Author attribution** with user relationships
- **Category organization** for easy filtering
- **SEO metadata** for search engine optimization
- **Full-text search** on titles and content
- **Pagination support** for listing endpoints
- **Automatic excerpt generation** for previews
- **Soft deletion** for content recovery

---

## Blog Post Structure

### Database Model

Blog posts are stored in the database with the following structure:

```go
type BlogPost struct {
    ID         uint           `json:"id"`           // Auto-generated
    Title      string         `json:"title"`        // Required
    Content    string         `json:"content"`      // Required (Markdown)
    AuthorID   uint           `json:"author_id"`    // Required
    Author     User           `json:"author"`       // Auto-populated
    CategoryID uint           `json:"category_id"`  // Optional
    Category   Category       `json:"category"`     // Auto-populated
    SEO        string         `json:"seo"`          // Optional
    CreatedAt  time.Time      `json:"created_at"`   // Auto-generated
    UpdatedAt  time.Time      `json:"updated_at"`   // Auto-maintained
    DeletedAt  gorm.DeletedAt `json:"-"`            // Soft delete
}
```

### Field Descriptions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | integer | Auto | Unique identifier for the blog post |
| `title` | string | Yes | The blog post title (used in search) |
| `content` | string | Yes | Blog content in **Markdown format** |
| `author_id` | integer | Yes | ID of the user who authored the post |
| `category_id` | integer | No | ID of the category for organization |
| `seo` | string | No | SEO keywords/metadata (comma-separated) |
| `created_at` | timestamp | Auto | When the post was created |
| `updated_at` | timestamp | Auto | When the post was last modified |

---

## Authentication & Authorization

### Authentication

All blog post creation, update, and deletion operations require authentication via AWS Cognito JWT tokens.

**How it works:**
- In **production mode**: Requests must include a valid JWT token in the `Authorization` header
- In **debug mode**: Authentication is bypassed with a test user ID

### Authorization Header Format

```
Authorization: Bearer <your-jwt-token>
```

### Debug Mode (Development)

When running in debug/development mode (`GIN_MODE=debug`):
- Authentication middleware is bypassed
- A test user ID (`1`) is automatically set
- Useful for local testing without AWS Cognito setup

### User Roles

The system supports the following user roles:
- `user` - Standard user (can create blog posts if granted author privileges)
- `author` - Can create and manage blog posts
- `admin` - Full administrative access

---

## API Endpoints

### Public Endpoints (No Authentication Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/blog` | List all blog posts with pagination and filtering |
| `GET` | `/api/v1/blog/:id` | Get a single blog post by ID |
| `GET` | `/api/v1/blog/category/:category_id` | Get blog posts by category |
| `GET` | `/api/v1/blog/author/:author_id` | Get blog posts by author |

### Authenticated Endpoints (Requires JWT Token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/blog` | Create a new blog post |
| `PUT` | `/api/v1/blog/:id` | Update an existing blog post |
| `DELETE` | `/api/v1/blog/:id` | Delete a blog post (soft delete) |

---

## Creating Blog Posts

### Step 1: Prepare Your Content

Write your blog post content in **Markdown format**. The system supports standard Markdown syntax including:

- Headers (`#`, `##`, `###`, etc.)
- Lists (ordered and unordered)
- Links `[text](url)`
- Images `![alt](url)`
- Code blocks (inline and fenced)
- Emphasis (`*italic*`, `**bold**`)
- Blockquotes

**Example Markdown Content:**

```markdown
# Understanding Security Governance

Security governance is the framework that ensures your organization's security policies align with business objectives.

## Key Components

1. **Policy Development** - Creating comprehensive security policies
2. **Risk Assessment** - Identifying and evaluating security risks
3. **Compliance Monitoring** - Ensuring adherence to regulations

For more information, check out our [Security Policy Template](/templates/security-policy).

## Best Practices

- Regular policy reviews
- Stakeholder engagement
- Continuous improvement
```

### Step 2: Make the API Request

**Endpoint:** `POST /api/v1/blog`

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <your-jwt-token>
```

**Request Body:**
```json
{
  "title": "Understanding Security Governance",
  "content": "# Understanding Security Governance\n\nSecurity governance is the framework...",
  "author_id": 1,
  "category_id": 1,
  "seo": "security, governance, compliance, policy"
}
```

**Success Response (201 Created):**
```json
{
  "message": "Blog post created successfully",
  "post": {
    "id": 10,
    "title": "Understanding Security Governance",
    "content": "# Understanding Security Governance\n\nSecurity governance...",
    "author_id": 1,
    "category_id": 1,
    "seo": "security, governance, compliance, policy",
    "created_at": "2025-11-04T10:30:00Z",
    "updated_at": "2025-11-04T10:30:00Z"
  }
}
```

### Step 3: Using cURL

```bash
curl -X POST http://localhost:8080/api/v1/blog \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Understanding Security Governance",
    "content": "# Understanding Security Governance\n\nYour markdown content here...",
    "author_id": 1,
    "category_id": 1,
    "seo": "security, governance, compliance"
  }'
```

### Step 4: Verify Creation

```bash
# Get the newly created blog post
curl http://localhost:8080/api/v1/blog/10
```

---

## Managing Existing Blog Posts

### Retrieving a Single Blog Post

**Endpoint:** `GET /api/v1/blog/:id`

**Example:**
```bash
curl http://localhost:8080/api/v1/blog/5
```

**Response:**
```json
{
  "post": {
    "id": 5,
    "title": "Understanding Security Governance",
    "content": "# Understanding Security Governance\n\n...",
    "html_content": "<h1>Understanding Security Governance</h1><p>...</p>",
    "author_id": 1,
    "author": {
      "id": 1,
      "username": "john_author",
      "email": "john@example.com"
    },
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Security & Compliance"
    },
    "seo": "security, governance, compliance",
    "created_at": "2025-11-04T10:30:00Z",
    "updated_at": "2025-11-04T10:30:00Z"
  }
}
```

**Key Features:**
- Returns both `content` (markdown) and `html_content` (processed HTML)
- Includes full author and category objects
- Preloads relationships for efficient querying

### Updating a Blog Post

**Endpoint:** `PUT /api/v1/blog/:id`

You can update any combination of fields. Only include the fields you want to change.

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <your-jwt-token>
```

**Example: Update Title and Content**
```bash
curl -X PUT http://localhost:8080/api/v1/blog/5 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Updated: Understanding Security Governance",
    "content": "# Updated Content\n\nNew markdown content..."
  }'
```

**Example: Update SEO Only**
```bash
curl -X PUT http://localhost:8080/api/v1/blog/5 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "seo": "security governance, risk management, compliance framework"
  }'
```

**Success Response (200 OK):**
```json
{
  "message": "Blog post updated successfully"
}
```

**Notes:**
- `updated_at` timestamp is automatically updated
- Only provided fields are modified
- Empty strings for required fields (`title`, `content`) will cause validation errors

### Deleting a Blog Post

**Endpoint:** `DELETE /api/v1/blog/:id`

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/blog/5 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response (200 OK):**
```json
{
  "message": "Blog post deleted successfully"
}
```

**Notes:**
- Uses **soft delete** (sets `deleted_at` timestamp)
- Blog post is not physically removed from database
- Can be recovered by database administrators if needed
- Will not appear in list queries after deletion

---

## Querying Blog Posts

### List All Blog Posts (with Pagination)

**Endpoint:** `GET /api/v1/blog`

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | integer | 10 | Number of posts per page |
| `offset` | integer | 0 | Number of posts to skip |
| `search` | string | - | Search in title and content |
| `category_id` | integer | - | Filter by category ID |

**Example: Get First 10 Posts**
```bash
curl "http://localhost:8080/api/v1/blog?limit=10&offset=0"
```

**Example: Search for "security"**
```bash
curl "http://localhost:8080/api/v1/blog?search=security"
```

**Example: Filter by Category**
```bash
curl "http://localhost:8080/api/v1/blog?category_id=1"
```

**Example: Pagination (Page 2 with 20 items per page)**
```bash
curl "http://localhost:8080/api/v1/blog?limit=20&offset=20"
```

**Response:**
```json
{
  "posts": [
    {
      "id": 1,
      "title": "Understanding Security Governance",
      "excerpt": "Security governance is the framework that ensures your organization's security policies align with business objectives...",
      "author_id": 1,
      "author": {
        "id": 1,
        "username": "john_author"
      },
      "category_id": 1,
      "category": {
        "id": 1,
        "name": "Security & Compliance"
      },
      "seo": "security, governance",
      "created_at": "2025-11-04T10:30:00Z",
      "updated_at": "2025-11-04T10:30:00Z"
    }
  ],
  "total": 45,
  "limit": 10,
  "offset": 0
}
```

**Key Features:**
- Automatic excerpt generation (150 characters by default)
- Posts ordered by `created_at` descending (newest first)
- Case-insensitive search using `ILIKE`
- Total count for pagination UI

### Get Blog Posts by Category

**Endpoint:** `GET /api/v1/blog/category/:category_id`

**Example:**
```bash
curl "http://localhost:8080/api/v1/blog/category/1"
```

**Response:**
```json
{
  "posts": [
    {
      "id": 1,
      "title": "Understanding Security Governance",
      "excerpt": "Security governance is the framework...",
      "author": { "id": 1, "username": "john_author" },
      "category": { "id": 1, "name": "Security & Compliance" },
      "created_at": "2025-11-04T10:30:00Z"
    }
  ]
}
```

### Get Blog Posts by Author

**Endpoint:** `GET /api/v1/blog/author/:author_id`

**Example:**
```bash
curl "http://localhost:8080/api/v1/blog/author/3"
```

**Response:** Same format as category query

---

## Best Practices

### Content Creation

1. **Write in Markdown**
   - Use proper heading hierarchy (`#`, `##`, `###`)
   - Include links to relevant templates or resources
   - Add code blocks for technical content
   - Use lists for better readability

2. **SEO Optimization**
   - Include relevant keywords in the `seo` field
   - Use comma-separated keywords
   - Keep titles clear and descriptive
   - Write engaging excerpts (first 150 characters are auto-generated)

3. **Content Structure**
   - Start with a clear introduction
   - Use headers to organize sections
   - Include practical examples
   - Link to related templates or resources

### API Usage

1. **Authentication**
   - Always include valid JWT tokens in production
   - Handle token expiration gracefully
   - Store tokens securely (never in client-side code)

2. **Error Handling**
   - Check HTTP status codes
   - Parse error messages from response body
   - Implement retry logic for network failures

3. **Pagination**
   - Use reasonable `limit` values (10-50 recommended)
   - Calculate `offset` as `page_number * limit`
   - Display total count for user navigation

4. **Search Performance**
   - Keep search queries focused
   - Consider caching frequently accessed posts
   - Use category filters to narrow results

### Security Considerations

1. **Input Validation**
   - Sanitize user-provided content
   - Validate author_id matches authenticated user
   - Check category_id exists before creating posts

2. **Authorization**
   - Only allow authors to create posts
   - Implement ownership checks for updates/deletes
   - Use admin middleware for sensitive operations

---

## Code Examples

### JavaScript/Fetch Example

```javascript
// Create a new blog post
async function createBlogPost(token, postData) {
  try {
    const response = await fetch('http://localhost:8080/api/v1/blog', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(postData)
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create blog post');
    }

    const result = await response.json();
    console.log('Blog post created:', result.post);
    return result.post;
  } catch (error) {
    console.error('Error creating blog post:', error);
    throw error;
  }
}

// Usage
const postData = {
  title: "Understanding Security Governance",
  content: "# Understanding Security Governance\n\n...",
  author_id: 1,
  category_id: 1,
  seo: "security, governance, compliance"
};

createBlogPost('your-jwt-token', postData);
```

### Python/Requests Example

```python
import requests

def create_blog_post(token, post_data):
    url = "http://localhost:8080/api/v1/blog"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}"
    }

    response = requests.post(url, json=post_data, headers=headers)

    if response.status_code == 201:
        print("Blog post created successfully")
        return response.json()["post"]
    else:
        print(f"Error: {response.json()['error']}")
        return None

# Usage
post_data = {
    "title": "Understanding Security Governance",
    "content": "# Understanding Security Governance\n\n...",
    "author_id": 1,
    "category_id": 1,
    "seo": "security, governance, compliance"
}

create_blog_post("your-jwt-token", post_data)
```

### Go Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type BlogPost struct {
    Title      string `json:"title"`
    Content    string `json:"content"`
    AuthorID   uint   `json:"author_id"`
    CategoryID uint   `json:"category_id"`
    SEO        string `json:"seo"`
}

func createBlogPost(token string, post BlogPost) error {
    url := "http://localhost:8080/api/v1/blog"

    jsonData, err := json.Marshal(post)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("failed to create blog post: %d", resp.StatusCode)
    }

    fmt.Println("Blog post created successfully")
    return nil
}

func main() {
    post := BlogPost{
        Title:      "Understanding Security Governance",
        Content:    "# Understanding Security Governance\n\n...",
        AuthorID:   1,
        CategoryID: 1,
        SEO:        "security, governance, compliance",
    }

    createBlogPost("your-jwt-token", post)
}
```

---

## Troubleshooting

### Common Errors

**401 Unauthorized**
```json
{
  "error": "Authorization header is missing"
}
```
**Solution:** Include `Authorization: Bearer <token>` header

**400 Bad Request - Empty Title**
```json
{
  "error": "blog post title is required"
}
```
**Solution:** Ensure `title` field is not empty

**400 Bad Request - Invalid JSON**
```json
{
  "error": "Invalid request body"
}
```
**Solution:** Validate JSON syntax before sending

**404 Not Found**
```json
{
  "error": "Blog post not found"
}
```
**Solution:** Verify the blog post ID exists and hasn't been deleted

### Debug Mode Testing

To test without authentication:

1. Set environment variable:
   ```bash
   export GIN_MODE=debug
   ```

2. Start the server:
   ```bash
   make run
   ```

3. Make requests without `Authorization` header:
   ```bash
   curl -X POST http://localhost:8080/api/v1/blog \
     -H "Content-Type: application/json" \
     -d '{"title": "Test", "content": "Test content", "author_id": 1}'
   ```

---

## Additional Resources

- [API Documentation](./API_DOCUMENTATION.md) - Complete API reference
- [Authentication Summary](./AUTHENTICATION_SUMMARY.md) - AWS Cognito setup
- [Swagger Documentation](http://localhost:8080/api-docs) - Interactive API explorer
- [README](./README.md) - Project overview and setup

---

## Support

For issues or questions:
1. Check the [GitHub Issues](https://github.com/H4ZM47/template-store-project/issues)
2. Review the API documentation at `/api-docs`
3. Check application logs for detailed error messages

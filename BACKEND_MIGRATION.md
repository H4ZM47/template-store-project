# Backend Migration: Go → Node.js/TypeScript

This document describes the migration of the Template Store backend from Go to Node.js with TypeScript.

## Overview

The backend has been completely refactored from Go (using Gin framework) to Node.js/TypeScript (using Express.js).

## Technology Stack

### Previous Stack (Go)
- **Language**: Go 1.24.0
- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: PostgreSQL (production) / SQLite (development)

### New Stack (Node.js/TypeScript)
- **Language**: TypeScript 5.7
- **Runtime**: Node.js ≥18.0.0
- **Framework**: Express.js 4.x
- **ORM**: TypeORM 0.3.x
- **Database**: PostgreSQL (production) / SQLite (development)

## Architecture

The new backend maintains the same clean architecture pattern:

```
src/
├── config/          # Configuration management
├── database/        # Database connection and setup
├── models/          # TypeORM entities (10 models)
├── services/        # Business logic layer (9 services)
├── routes/          # API route handlers
├── middleware/      # Authentication & RBAC middleware
├── utils/           # Utility functions (logger, etc.)
└── server.ts        # Main application entry point
```

## Database Models (TypeORM Entities)

All 10 database models have been migrated:

1. **User** - User accounts with authentication
2. **Category** - Content categories
3. **Template** - Digital templates for sale
4. **BlogPost** - Blog articles
5. **Order** - Purchase orders
6. **LoginHistory** - User login tracking
7. **ActivityLog** - User activity tracking
8. **PasswordResetToken** - Password reset tokens
9. **EmailVerificationToken** - Email verification tokens
10. **UserPreferences** - User settings and preferences

## Services Layer

The following services have been implemented:

1. **AuthService** - AWS Cognito integration
2. **StorageService** - AWS S3 file storage
3. **EmailService** - SendGrid email delivery
4. **StripeService** - Stripe payment processing
5. **UserService** - User management
6. **CategoryService** - Category CRUD operations
7. **TemplateService** - Template management
8. **BlogService** - Blog post management with Markdown rendering
9. **OrderService** - Order processing and tracking

## API Endpoints

All original API endpoints have been preserved:

### Authentication (`/api/v1/auth`)
- `POST /register` - User registration
- `POST /login` - User login
- `POST /forgot-password` - Password reset request
- `POST /reset-password` - Password reset confirmation
- `POST /change-password` - Change password (authenticated)

### Profile (`/api/v1/profile`)
- `GET /` - Get user profile
- `PUT /` - Update user profile
- `GET /orders` - List user orders
- `GET /purchased-templates` - List purchased templates
- `POST /deactivate` - Deactivate account

### Templates (`/api/v1/templates`)
- `GET /` - List all templates (public)
- `GET /:id` - Get template details (public)
- `GET /category/:category_id` - Get templates by category (public)
- `POST /` - Create template (author/admin)
- `PUT /:id` - Update template (author/admin)
- `DELETE /:id` - Delete template (author/admin)

### Blog (`/api/v1/blog`)
- `GET /` - List all blog posts (public)
- `GET /:id` - Get blog post (public)
- `GET /category/:category_id` - Get posts by category (public)
- `GET /author/:author_id` - Get posts by author (public)
- `POST /` - Create blog post (author/admin)
- `PUT /:id` - Update blog post (author/admin)
- `DELETE /:id` - Delete blog post (author/admin)

### Categories (`/api/v1/categories`)
- `GET /` - List categories (public)
- `GET /:id` - Get category (public)
- `POST /` - Create category (public)
- `PUT /:id` - Update category (public)
- `DELETE /:id` - Delete category (public)

### Payment (`/api/v1/payment`)
- `POST /checkout` - Create checkout session (authenticated)
- `GET /success` - Payment success callback (public)
- `GET /cancel` - Payment cancelled callback (public)
- `POST /webhooks/stripe` - Stripe webhook (public)

### Admin (`/api/v1/admin`)
- `GET /users` - List all users
- `GET /users/:id` - Get user details
- `PUT /users/:id/role` - Update user role
- `POST /users/:id/suspend` - Suspend user
- `POST /users/:id/unsuspend` - Unsuspend user
- `DELETE /users/:id` - Delete user

## External Integrations

All external service integrations have been migrated:

1. **AWS Cognito** - User authentication (@aws-sdk/client-cognito-identity-provider)
2. **AWS S3** - File storage (@aws-sdk/client-s3)
3. **Stripe** - Payment processing (stripe v17)
4. **SendGrid** - Email delivery (@sendgrid/mail)

## Setup Instructions

### 1. Install Dependencies

```bash
npm install
```

### 2. Configure Environment

Copy `.env.example` to `.env` and update with your configuration:

```bash
cp .env.example .env
```

### 3. Build the Project

```bash
npm run build
```

### 4. Run in Development Mode

```bash
npm run dev
```

### 5. Run in Production Mode

```bash
npm run build
npm start
```

## Development Scripts

- `npm run dev` - Start development server with hot reload
- `npm run build` - Compile TypeScript to JavaScript
- `npm start` - Start production server
- `npm run lint` - Run ESLint
- `npm run format` - Format code with Prettier
- `npm run typecheck` - Type-check without emitting files

## Database Migration

The new backend uses TypeORM with automatic synchronization enabled. The database schema will be automatically created/updated on server start.

For production, it's recommended to:
1. Disable `synchronize: true` in `src/database/connection.ts`
2. Use TypeORM migrations for schema changes

## Key Differences from Go Backend

1. **Language**: TypeScript with static typing instead of Go
2. **Framework**: Express.js instead of Gin
3. **ORM**: TypeORM instead of GORM
4. **Package Manager**: npm instead of Go modules
5. **Build Process**: TypeScript compilation instead of Go build

## Compatibility

The new backend maintains API compatibility with the existing frontend and any external clients. All endpoints, request/response formats, and authentication mechanisms remain the same.

## Testing

To test the backend:

1. Start the server: `npm run dev`
2. Health check: `curl http://localhost:8080/api/v1/health`
3. Test authentication endpoints with your frontend or API client

## Logging

Winston is used for structured logging with the following features:

- Console output with colors (development)
- File logging to `logs/` directory
- JSON format for easy parsing
- Separate error log file
- Configurable log levels

## Production Deployment

For production deployment:

1. Set `GIN_MODE=release` and `NODE_ENV=production`
2. Use PostgreSQL instead of SQLite
3. Configure proper environment variables
4. Use a process manager (PM2, systemd)
5. Set up proper logging and monitoring
6. Use HTTPS/TLS
7. Configure CORS properly
8. Use database migrations instead of auto-sync

## Migration Checklist

- [x] Project structure and configuration
- [x] Database models/entities
- [x] Service layer
- [x] Route handlers
- [x] Middleware (auth, RBAC)
- [x] External integrations (AWS, Stripe, SendGrid)
- [x] Error handling and logging
- [ ] Unit tests (future work)
- [ ] Integration tests (future work)
- [ ] API documentation (future work)

## Next Steps

1. Test all endpoints thoroughly
2. Add comprehensive unit and integration tests
3. Set up CI/CD pipeline
4. Generate API documentation (Swagger/OpenAPI)
5. Performance testing and optimization
6. Security audit

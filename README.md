# Template Store & Blog Platform

A modern web application for selling digital templates and hosting a company blog, built with Node.js, TypeScript, PostgreSQL, and AWS services.

## Architecture

This project follows a clean architecture pattern with:
- **Backend**: Node.js with TypeScript and Express.js framework
- **ORM**: TypeORM for database management
- **Database**: PostgreSQL (production) / SQLite (development)
- **Storage**: AWS S3 with CloudFront CDN
- **Authentication**: AWS Cognito
- **Payments**: Stripe integration
- **Email**: SendGrid for notifications
- **Frontend**: Vanilla HTML/CSS/JS with Tailwind CSS

## Prerequisites

- Node.js 18.0 or higher
- npm 9.0 or higher
- PostgreSQL 13 or higher (for production)
- Docker (optional, for containerized development)
- AWS CLI (for deployment)

## Quick Start

### Option 1: Using npm Scripts (Recommended)

The easiest way to get started:

```bash
# Clone the repository
git clone <repository-url>
cd template-store-project

# Install dependencies
npm install

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration values

# Start development server with hot reload
npm run dev
```

This will automatically:
- ✅ Start the backend API server on port 8080
- ✅ Connect to the database (PostgreSQL or SQLite)
- ✅ Enable hot reload for development
- ✅ Set up TypeORM with automatic synchronization

### Option 2: Manual Setup

If you prefer manual setup or need to customize the configuration:

```bash
# 1. Clone and setup
git clone <repository-url>
cd template-store-project

# 2. Install dependencies
npm install

# 3. Set up environment variables
cp .env.example .env
# Edit .env with your configuration values

# 4. Set up PostgreSQL (for production)
createdb template_store
# Or use SQLite for development (automatic)

# 5. Build the project
npm run build

# 6. Start the server
npm start
```

### Access Your Application

Once started, you can access:
- **Frontend**: http://localhost:3000 (served from web/ directory)
- **Backend API**: http://localhost:8080/api/v1/
- **Health Check**: http://localhost:8080/api/v1/health
- **API Documentation**: See [swagger.yaml](swagger.yaml)

## Development Scripts

Available npm scripts for development:

### Core Commands

```bash
npm run dev        # Start development server with hot reload
npm run build      # Compile TypeScript to JavaScript
npm start          # Start production server
npm run lint       # Run ESLint code analysis
npm run format     # Format code with Prettier
npm run typecheck  # Type-check without emitting files
```

### Development Workflow

```bash
# Initial setup
npm install        # Install all dependencies

# Daily development
npm run dev        # Start dev server (watches for changes)

# Before committing
npm run lint       # Check code style
npm run format     # Format code
npm run typecheck  # Verify types

# Production build
npm run build      # Build for production
npm start          # Run production server
```

## Development

### Project Structure
```
template-store-project/
├── src/                    # TypeScript source code
│   ├── config/            # Configuration management
│   ├── database/          # Database connection setup
│   ├── models/            # TypeORM entities (10 models)
│   ├── routes/            # API route handlers (7 modules)
│   ├── middleware/        # Authentication & RBAC middleware
│   ├── services/          # Business logic layer (9 services)
│   ├── utils/             # Utility functions (logger, etc.)
│   └── server.ts          # Main application entry point
├── dist/                  # Compiled JavaScript output
├── web/                   # Frontend assets (HTML/CSS/JS)
├── node_modules/          # Dependencies
├── logs/                  # Application logs
├── scripts/               # Utility scripts
├── package.json           # npm dependencies and scripts
├── tsconfig.json          # TypeScript configuration
├── .env.example           # Environment variables template
└── docs/                  # Documentation (*.md files)
```

### API Endpoints

The backend provides comprehensive REST APIs:

#### Authentication (`/api/v1/auth`)
- `POST /register` - User registration
- `POST /login` - User login
- `POST /forgot-password` - Password reset request
- `POST /reset-password` - Password reset confirmation
- `POST /change-password` - Change password (authenticated)

#### Profile (`/api/v1/profile`)
- `GET /` - Get user profile
- `PUT /` - Update user profile
- `GET /orders` - List user orders
- `GET /purchased-templates` - List purchased templates

#### Templates (`/api/v1/templates`)
- `GET /` - List all templates
- `GET /:id` - Get template details
- `POST /` - Create template (author/admin)
- `PUT /:id` - Update template (author/admin)
- `DELETE /:id` - Delete template (author/admin)

#### Blog (`/api/v1/blog`)
- `GET /` - List blog posts
- `GET /:id` - Get blog post
- `POST /` - Create post (author/admin)
- `PUT /:id` - Update post (author/admin)
- `DELETE /:id` - Delete post (author/admin)

#### Payment (`/api/v1/payment`)
- `POST /checkout` - Create Stripe checkout session
- `POST /webhooks/stripe` - Stripe webhook handler

For full API documentation, see [swagger.yaml](swagger.yaml).

### Environment Variables

See `.env.example` for all available configuration options including:
- Database connection settings
- AWS Cognito configuration
- AWS S3 credentials
- Stripe API keys
- SendGrid API key

## Testing

```bash
# Type checking
npm run typecheck

# Linting
npm run lint

# Format checking
npm run format

# Note: Unit and integration tests are planned for future implementation
```

## Building

```bash
# Build TypeScript to JavaScript
npm run build

# Output will be in dist/ directory
# Run with: npm start
```

## Deployment

### Local Development
```bash
# Run with hot reload (recommended)
npm run dev

# Run with Docker (if docker-compose.yml configured)
docker-compose up
```

### Production

The application is designed to be deployed on AWS with the following setup:

#### Infrastructure
- **Compute**: AWS ECS/EC2 or any Node.js hosting platform
- **Database**: AWS RDS PostgreSQL
- **Storage**: AWS S3 for file uploads
- **CDN**: CloudFront for content delivery
- **Authentication**: AWS Cognito user pools

#### Deployment Steps

1. **Build the application:**
   ```bash
   npm run build
   ```

2. **Set environment variables:**
   - Set `NODE_ENV=production`
   - Configure all AWS credentials
   - Set database connection strings
   - Configure Stripe and SendGrid API keys

3. **Run the production server:**
   ```bash
   npm start
   ```

4. **Process Management:**
   - Use PM2, systemd, or Docker for process management
   - Example with PM2:
     ```bash
     npm install -g pm2
     pm2 start dist/server.js --name template-store
     pm2 save
     ```

#### Production Checklist
- [ ] Set `NODE_ENV=production`
- [ ] Use PostgreSQL (not SQLite)
- [ ] Disable TypeORM `synchronize` (use migrations)
- [ ] Configure proper CORS settings
- [ ] Enable HTTPS/TLS
- [ ] Set up monitoring and logging
- [ ] Configure backup strategy
- [ ] Set up CI/CD pipeline

## Technology Stack

### Backend
- **Language**: TypeScript 5.7
- **Runtime**: Node.js 18+
- **Framework**: Express.js 4.x
- **ORM**: TypeORM 0.3.x
- **Database**: PostgreSQL / SQLite
- **Logging**: Winston

### External Services
- **AWS Cognito**: User authentication
- **AWS S3**: File storage
- **Stripe**: Payment processing
- **SendGrid**: Email notifications

### Frontend
- **HTML/CSS/JS**: Vanilla JavaScript
- **CSS Framework**: Tailwind CSS
- **Markdown**: Marked library for blog rendering

## Database Models

The application includes 10 TypeORM entities:
1. **User** - User accounts and profiles
2. **Category** - Content categories
3. **Template** - Digital templates for sale
4. **BlogPost** - Blog articles with Markdown support
5. **Order** - Purchase orders and transaction history
6. **LoginHistory** - User login tracking
7. **ActivityLog** - User activity tracking
8. **PasswordResetToken** - Password reset tokens
9. **EmailVerificationToken** - Email verification tokens
10. **UserPreferences** - User settings and preferences

## Migration History

This project was recently migrated from Go to Node.js/TypeScript. For details, see:
- [BACKEND_MIGRATION.md](BACKEND_MIGRATION.md) - Complete migration documentation
- [CLEANUP_SUMMARY.md](CLEANUP_SUMMARY.md) - Cleanup details

## Documentation

Additional documentation is available:
- [API_DOCUMENTATION.md](API_DOCUMENTATION.md) - API reference
- [AUTHENTICATION_SUMMARY.md](AUTHENTICATION_SUMMARY.md) - Authentication details
- [PAYMENT_INTEGRATION_SUMMARY.md](PAYMENT_INTEGRATION_SUMMARY.md) - Payment integration
- [BLOG_MANAGEMENT_GUIDE.md](BLOG_MANAGEMENT_GUIDE.md) - Blog management
- [swagger.yaml](swagger.yaml) - OpenAPI specification

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run linting and type checking:
   ```bash
   npm run lint
   npm run typecheck
   npm run format
   ```
5. Test your changes thoroughly
6. Submit a pull request

## License

MIT License

## Support

For support and questions:
- Create an issue on GitHub
- Check existing documentation in the [docs](.) folder
- Review the migration guide if you have questions about the architecture 
# Company Blog & Template Store Architecture

```mermaid
graph TB
    %% User Interface Layer
    subgraph "Frontend (Vanilla HTML/CSS/JS)"
        UI[User Interface]
        UI --> |API Calls| API
    end

    %% Backend API Layer
    subgraph "Backend (Go + Gin/Echo)"
        API[API Gateway]
        API --> Auth[Authentication Middleware]
        API --> BlogHandler[Blog Handler]
        API --> TemplateHandler[Template Handler]
        API --> PaymentHandler[Payment Handler]
        API --> AdminHandler[Admin Handler]
    end

    %% Database Layer
    subgraph "Database (PostgreSQL on AWS RDS)"
        DB[(PostgreSQL Database)]
        DB --> |Blog Posts| BlogPosts[(Blog Posts Table)]
        DB --> |Templates| Templates[(Templates Table)]
        DB --> |Orders| Orders[(Orders Table)]
        DB --> |Users| Users[(Users Table)]
        DB --> |Categories| Categories[(Categories Table)]
    end

    %% AWS Services
    subgraph "AWS Infrastructure"
        subgraph "Compute & Load Balancing"
            ALB[Application Load Balancer]
            ECS[ECS/Fargate Cluster]
            ECS --> |Go Backend Containers| Backend[Go Backend Service]
        end

        subgraph "Storage & CDN"
            S3[(S3 Bucket)]
            S3 --> |Template Files| TemplatesS3[Templates Storage]
            S3 --> |Blog Assets| BlogAssets[Blog Assets]
            CloudFront[CloudFront CDN]
            CloudFront --> S3
        end

        subgraph "Database & Monitoring"
            RDS[(RDS PostgreSQL)]
            CloudWatch[CloudWatch Logs]
            CloudWatch --> ECS
        end

        subgraph "DNS & Security"
            Route53[Route 53]
            Route53 --> ALB
            WAF[AWS WAF]
            WAF --> ALB
            Cognito[AWS Cognito]
            Cognito --> |User Authentication| Auth
        end
    end

    %% External Services
    subgraph "External Services"
        Stripe[Stripe Payment]
        Stripe --> |Webhooks| PaymentHandler
        SendGrid[SendGrid Email]
        SendGrid --> |Template Delivery| PaymentHandler
        SendGrid --> |Welcome Emails| Cognito
    end

    %% Data Flow Connections
    API --> DB
    BlogHandler --> BlogPosts
    TemplateHandler --> Templates
    PaymentHandler --> Orders
    AdminHandler --> Users
    AdminHandler --> Categories
    Cognito --> |User Data| Users

    %% AWS Service Connections
    Backend --> RDS
    Backend --> S3
    Backend --> CloudWatch
    ALB --> Backend

    %% Template Flow
    TemplatesS3 --> |Download Links| PaymentHandler
    PaymentHandler --> |Purchase Confirmation| SendGrid

    %% Blog Flow
    BlogHandler --> |Markdown Processing| BlogPosts
    BlogAssets --> |Images/Assets| BlogHandler

    %% Security & Monitoring
    Auth --> WAF
    CloudWatch --> |Logs & Metrics| Monitoring[Monitoring Dashboard]

    %% Styling
    classDef frontend fill:#e1f5fe
    classDef backend fill:#f3e5f5
    classDef database fill:#e8f5e8
    classDef aws fill:#fff3e0
    classDef external fill:#ffebee

    class UI frontend
    class API,Auth,BlogHandler,TemplateHandler,PaymentHandler,AdminHandler backend
    class DB,BlogPosts,Templates,Orders,Users,Categories database
    class ALB,ECS,Backend,S3,TemplatesS3,BlogAssets,CloudFront,RDS,CloudWatch,Route53,WAF,Cognito aws
    class Stripe,SendGrid external
```

## Component Details

### Frontend Components
- **User Interface**: Vanilla HTML/CSS/JS with Tailwind CSS
- **Pages**: Home, Blog, Template Catalog, Template Details, About, Contact
- **Features**: Responsive design, search, social sharing, breadcrumbs

### Backend Components (Go)
- **API Gateway**: Main entry point for all requests
- **Authentication Middleware**: AWS Cognito integration and JWT validation
- **Blog Handler**: Markdown processing and blog post management
- **Template Handler**: Template CRUD operations and preview generation
- **Payment Handler**: Stripe integration and order processing
- **Admin Handler**: Admin interface for content management

### Database Schema (PostgreSQL)
- **Blog Posts**: Content, metadata, SEO fields
- **Templates**: File info, categories, pricing, preview data
- **Orders**: Purchase history, delivery status, user info
- **Users**: Customer accounts and preferences (linked to Cognito)
- **Categories**: Template and blog categorization

### AWS Infrastructure
- **ECS/Fargate**: Containerized Go backend
- **RDS**: Managed PostgreSQL database
- **S3**: Template file storage with CloudFront CDN
- **ALB**: Load balancing and SSL termination
- **Route 53**: DNS management
- **CloudWatch**: Monitoring and logging
- **WAF**: Web application firewall
- **Cognito**: User authentication and management

### External Integrations
- **Stripe**: Payment processing and webhooks
- **SendGrid**: Reliable email delivery for templates and notifications

## Data Flow
1. **User Request** → ALB → ECS → Go Backend
2. **Authentication** → Cognito → JWT validation → API access
3. **API Calls** → Database queries → PostgreSQL
4. **Template Access** → S3 → CloudFront → User
5. **Payment Flow** → Stripe → Webhook → SendGrid → Email Delivery
6. **Blog Content** → Markdown processing → Database → Frontend 
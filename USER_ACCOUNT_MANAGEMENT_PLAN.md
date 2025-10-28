# User Account Management - Comprehensive Implementation Plan

## Overview
This document outlines the complete implementation plan for user account management functionality, including profile management, security features, account settings, user dashboard, and admin capabilities.

---

## Phase 1: Database Schema Enhancements

### 1.1 Extended User Model
**File:** `internal/models/user.go`

Add new fields to the User model:
```go
type User struct {
    // Existing fields
    ID             uint      `gorm:"primaryKey" json:"id"`
    CognitoSubject string    `gorm:"uniqueIndex" json:"cognito_subject"`
    Name           string    `json:"name"`
    Email          string    `gorm:"uniqueIndex" json:"email"`

    // New profile fields
    AvatarURL      string    `json:"avatar_url"`
    Bio            string    `gorm:"type:text" json:"bio"`
    PhoneNumber    string    `json:"phone_number"`

    // Address fields
    AddressLine1   string    `json:"address_line1"`
    AddressLine2   string    `json:"address_line2"`
    City           string    `json:"city"`
    State          string    `json:"state"`
    Country        string    `json:"country"`
    PostalCode     string    `json:"postal_code"`

    // Account status
    Role           string    `gorm:"default:'user'" json:"role"` // user, admin, author
    Status         string    `gorm:"default:'active'" json:"status"` // active, suspended, deleted
    EmailVerified  bool      `gorm:"default:false" json:"email_verified"`

    // Settings (JSON)
    Preferences    datatypes.JSON `json:"preferences"`

    // Metadata
    LastLoginAt    *time.Time `json:"last_login_at"`
    SuspendedAt    *time.Time `json:"suspended_at"`
    SuspendedBy    *uint      `json:"suspended_by"`
    SuspensionReason string   `json:"suspension_reason"`

    // Existing relationships
    BlogPosts      []BlogPost `gorm:"foreignKey:AuthorID" json:"blog_posts,omitempty"`
    Orders         []Order    `gorm:"foreignKey:UserID" json:"orders,omitempty"`

    // Timestamps
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
    DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### 1.2 New Models

#### UserPreferences Model
**File:** `internal/models/user_preferences.go`
```go
type UserPreferences struct {
    // Email preferences
    MarketingEmails      bool `json:"marketing_emails"`
    OrderNotifications   bool `json:"order_notifications"`
    BlogNotifications    bool `json:"blog_notifications"`

    // Display preferences
    Language             string `json:"language"` // en, es, fr, etc.
    Timezone             string `json:"timezone"` // UTC, America/New_York, etc.
    Theme                string `json:"theme"` // light, dark, auto

    // Privacy settings
    ProfileVisibility    string `json:"profile_visibility"` // public, private
    ShowEmail            bool   `json:"show_email"`
    ShowPurchaseHistory  bool   `json:"show_purchase_history"`
}
```

#### LoginHistory Model
**File:** `internal/models/login_history.go`
```go
type LoginHistory struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    UserID        uint      `gorm:"index" json:"user_id"`
    User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`

    IPAddress     string    `json:"ip_address"`
    UserAgent     string    `json:"user_agent"`
    Device        string    `json:"device"` // mobile, desktop, tablet
    Location      string    `json:"location"` // City, Country (from IP)

    LoginAt       time.Time `json:"login_at"`
    LogoutAt      *time.Time `json:"logout_at"`

    Success       bool      `json:"success"`
    FailureReason string    `json:"failure_reason"`

    CreatedAt     time.Time `json:"created_at"`
}
```

#### ActivityLog Model
**File:** `internal/models/activity_log.go`
```go
type ActivityLog struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    UserID      uint           `gorm:"index" json:"user_id"`
    User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`

    Action      string         `json:"action"` // profile_updated, password_changed, etc.
    Resource    string         `json:"resource"` // user, order, template, etc.
    ResourceID  *uint          `json:"resource_id"`
    Details     datatypes.JSON `json:"details"`

    IPAddress   string         `json:"ip_address"`
    UserAgent   string         `json:"user_agent"`

    CreatedAt   time.Time      `json:"created_at"`
}
```

#### PasswordResetToken Model
**File:** `internal/models/password_reset_token.go`
```go
type PasswordResetToken struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UserID      uint      `gorm:"index" json:"user_id"`
    User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`

    Token       string    `gorm:"uniqueIndex" json:"token"`
    ExpiresAt   time.Time `json:"expires_at"`
    UsedAt      *time.Time `json:"used_at"`

    CreatedAt   time.Time `json:"created_at"`
}
```

#### EmailVerificationToken Model
**File:** `internal/models/email_verification_token.go`
```go
type EmailVerificationToken struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UserID      uint      `gorm:"index" json:"user_id"`
    User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`

    Token       string    `gorm:"uniqueIndex" json:"token"`
    Email       string    `json:"email"` // For email change verification
    ExpiresAt   time.Time `json:"expires_at"`
    VerifiedAt  *time.Time `json:"verified_at"`

    CreatedAt   time.Time `json:"created_at"`
}
```

### 1.3 Migration Strategy
**File:** `internal/models/migrate.go`

Update the AutoMigrate function to include new models:
```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &User{},
        &Template{},
        &Category{},
        &BlogPost{},
        &Order{},
        &LoginHistory{},
        &ActivityLog{},
        &PasswordResetToken{},
        &EmailVerificationToken{},
    )
}
```

---

## Phase 2: API Endpoints Design

### 2.1 Profile Management Endpoints

#### GET /api/v1/profile
**Description:** Get current user's profile
**Auth:** Required
**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "avatar_url": "https://cdn.example.com/avatars/user1.jpg",
  "bio": "Template enthusiast",
  "phone_number": "+1234567890",
  "address_line1": "123 Main St",
  "city": "New York",
  "country": "USA",
  "role": "user",
  "status": "active",
  "email_verified": true,
  "preferences": {
    "marketing_emails": true,
    "language": "en",
    "theme": "dark"
  },
  "created_at": "2025-01-15T10:30:00Z"
}
```

#### PUT /api/v1/profile
**Description:** Update current user's profile
**Auth:** Required
**Request Body:**
```json
{
  "name": "John Doe",
  "bio": "Updated bio",
  "phone_number": "+1234567890",
  "address_line1": "123 Main St",
  "address_line2": "Apt 4B",
  "city": "New York",
  "state": "NY",
  "country": "USA",
  "postal_code": "10001"
}
```

#### POST /api/v1/profile/avatar
**Description:** Upload profile avatar
**Auth:** Required
**Content-Type:** multipart/form-data
**Request:** Form data with "avatar" file field
**Response:**
```json
{
  "avatar_url": "https://cdn.example.com/avatars/user1.jpg"
}
```

#### DELETE /api/v1/profile/avatar
**Description:** Remove profile avatar
**Auth:** Required
**Response:** 204 No Content

#### GET /api/v1/users/:id/profile
**Description:** Get public profile of any user
**Auth:** Optional
**Response:** Limited profile data based on privacy settings

### 2.2 Password & Security Endpoints

#### POST /api/v1/auth/forgot-password
**Description:** Request password reset
**Auth:** None
**Request Body:**
```json
{
  "email": "john@example.com"
}
```
**Response:**
```json
{
  "message": "If the email exists, a password reset link has been sent"
}
```

#### POST /api/v1/auth/reset-password
**Description:** Reset password with token
**Auth:** None
**Request Body:**
```json
{
  "token": "abc123...",
  "new_password": "newSecurePassword123"
}
```

#### POST /api/v1/auth/change-password
**Description:** Change password (authenticated)
**Auth:** Required
**Request Body:**
```json
{
  "current_password": "oldPassword123",
  "new_password": "newPassword123"
}
```

#### POST /api/v1/auth/verify-email
**Description:** Verify email with token
**Auth:** None
**Request Body:**
```json
{
  "token": "verification_token_here"
}
```

#### POST /api/v1/auth/resend-verification
**Description:** Resend verification email
**Auth:** Required
**Response:**
```json
{
  "message": "Verification email sent"
}
```

#### GET /api/v1/profile/sessions
**Description:** Get active login sessions
**Auth:** Required
**Response:**
```json
{
  "sessions": [
    {
      "id": 1,
      "device": "desktop",
      "ip_address": "192.168.1.1",
      "location": "New York, USA",
      "login_at": "2025-10-28T10:00:00Z",
      "user_agent": "Mozilla/5.0..."
    }
  ]
}
```

#### GET /api/v1/profile/login-history
**Description:** Get login history
**Auth:** Required
**Query Params:** ?limit=20&offset=0
**Response:**
```json
{
  "history": [
    {
      "id": 1,
      "ip_address": "192.168.1.1",
      "device": "desktop",
      "location": "New York, USA",
      "login_at": "2025-10-28T10:00:00Z",
      "success": true
    }
  ],
  "total": 50
}
```

#### POST /api/v1/auth/logout-session/:id
**Description:** Logout specific session
**Auth:** Required

#### POST /api/v1/auth/logout-all
**Description:** Logout all sessions
**Auth:** Required

### 2.3 Account Settings Endpoints

#### GET /api/v1/profile/preferences
**Description:** Get user preferences
**Auth:** Required
**Response:**
```json
{
  "marketing_emails": true,
  "order_notifications": true,
  "blog_notifications": false,
  "language": "en",
  "timezone": "America/New_York",
  "theme": "dark",
  "profile_visibility": "public",
  "show_email": false
}
```

#### PUT /api/v1/profile/preferences
**Description:** Update user preferences
**Auth:** Required
**Request Body:**
```json
{
  "marketing_emails": false,
  "language": "es",
  "theme": "light"
}
```

#### POST /api/v1/profile/change-email
**Description:** Request email change
**Auth:** Required
**Request Body:**
```json
{
  "new_email": "newemail@example.com",
  "password": "currentPassword123"
}
```
**Response:**
```json
{
  "message": "Verification email sent to new address"
}
```

#### POST /api/v1/profile/deactivate
**Description:** Deactivate account
**Auth:** Required
**Request Body:**
```json
{
  "password": "currentPassword123",
  "reason": "No longer need the service"
}
```

#### DELETE /api/v1/profile
**Description:** Delete account permanently
**Auth:** Required
**Request Body:**
```json
{
  "password": "currentPassword123",
  "confirmation": "DELETE MY ACCOUNT"
}
```

### 2.4 User Dashboard Endpoints

#### GET /api/v1/profile/dashboard
**Description:** Get user dashboard summary
**Auth:** Required
**Response:**
```json
{
  "stats": {
    "total_orders": 15,
    "total_spent": 299.85,
    "templates_purchased": 12,
    "blog_posts_authored": 5,
    "account_age_days": 180
  },
  "recent_orders": [...],
  "recent_blog_posts": [...]
}
```

#### GET /api/v1/profile/orders
**Description:** Get user's order history
**Auth:** Required
**Query Params:** ?limit=10&offset=0&status=completed
**Response:**
```json
{
  "orders": [
    {
      "id": 1,
      "template": {
        "id": 5,
        "name": "Invoice Template Pro"
      },
      "amount": 29.99,
      "status": "completed",
      "purchased_at": "2025-10-15T14:30:00Z",
      "download_url": "https://..."
    }
  ],
  "total": 15,
  "total_spent": 299.85
}
```

#### GET /api/v1/profile/orders/:id
**Description:** Get specific order details
**Auth:** Required
**Response:**
```json
{
  "id": 1,
  "template": {...},
  "amount": 29.99,
  "status": "completed",
  "payment_method": "card",
  "stripe_session_id": "cs_...",
  "purchased_at": "2025-10-15T14:30:00Z",
  "download_url": "https://...",
  "invoice_url": "https://..."
}
```

#### GET /api/v1/profile/purchased-templates
**Description:** Get all purchased templates
**Auth:** Required
**Response:**
```json
{
  "templates": [
    {
      "id": 5,
      "name": "Invoice Template Pro",
      "category": "Business",
      "purchased_at": "2025-10-15T14:30:00Z",
      "download_url": "https://...",
      "download_count": 3
    }
  ]
}
```

#### GET /api/v1/profile/blog-posts
**Description:** Get user's authored blog posts
**Auth:** Required
**Query Params:** ?status=published&limit=10&offset=0
**Response:**
```json
{
  "posts": [
    {
      "id": 1,
      "title": "10 Best Invoice Templates",
      "status": "published",
      "views": 1250,
      "created_at": "2025-09-01T10:00:00Z"
    }
  ],
  "total": 5
}
```

#### GET /api/v1/profile/activity
**Description:** Get user activity log
**Auth:** Required
**Query Params:** ?limit=20&offset=0
**Response:**
```json
{
  "activities": [
    {
      "id": 1,
      "action": "profile_updated",
      "details": "Updated profile picture",
      "created_at": "2025-10-28T10:00:00Z"
    }
  ],
  "total": 150
}
```

### 2.5 Admin Endpoints

#### GET /api/v1/admin/users
**Description:** List all users (admin only)
**Auth:** Required (admin role)
**Query Params:** ?search=john&role=user&status=active&limit=50&offset=0&sort=created_at&order=desc
**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "status": "active",
      "email_verified": true,
      "total_orders": 5,
      "total_spent": 149.95,
      "last_login_at": "2025-10-28T09:00:00Z",
      "created_at": "2025-08-01T10:00:00Z"
    }
  ],
  "total": 1250,
  "page": 1,
  "per_page": 50
}
```

#### GET /api/v1/admin/users/:id
**Description:** Get detailed user info (admin only)
**Auth:** Required (admin role)
**Response:** Full user object with all relationships

#### PUT /api/v1/admin/users/:id/role
**Description:** Update user role
**Auth:** Required (admin role)
**Request Body:**
```json
{
  "role": "admin"
}
```

#### POST /api/v1/admin/users/:id/suspend
**Description:** Suspend user account
**Auth:** Required (admin role)
**Request Body:**
```json
{
  "reason": "Violation of terms of service",
  "duration_days": 30
}
```

#### POST /api/v1/admin/users/:id/unsuspend
**Description:** Unsuspend user account
**Auth:** Required (admin role)

#### DELETE /api/v1/admin/users/:id
**Description:** Delete user (soft delete)
**Auth:** Required (admin role)

#### GET /api/v1/admin/dashboard
**Description:** Admin dashboard statistics
**Auth:** Required (admin role)
**Response:**
```json
{
  "stats": {
    "total_users": 1250,
    "active_users": 1100,
    "suspended_users": 10,
    "new_users_this_month": 85,
    "total_orders": 5430,
    "total_revenue": 54299.50,
    "revenue_this_month": 4250.00
  },
  "recent_users": [...],
  "recent_orders": [...]
}
```

#### GET /api/v1/admin/users/:id/activity
**Description:** View user's activity log
**Auth:** Required (admin role)

#### POST /api/v1/admin/users/:id/impersonate
**Description:** Impersonate user (for support)
**Auth:** Required (admin role)
**Response:**
```json
{
  "token": "impersonation_jwt_token",
  "expires_in": 3600
}
```

---

## Phase 3: Service Layer Implementation

### 3.1 Profile Service
**File:** `internal/services/profile.go`

```go
type ProfileService interface {
    // Profile operations
    GetProfile(userID uint) (*models.User, error)
    UpdateProfile(userID uint, updates map[string]interface{}) error
    UploadAvatar(userID uint, file multipart.File, filename string) (string, error)
    DeleteAvatar(userID uint) error
    GetPublicProfile(userID uint, viewerID *uint) (*models.User, error)

    // Preferences
    GetPreferences(userID uint) (*models.UserPreferences, error)
    UpdatePreferences(userID uint, prefs *models.UserPreferences) error

    // Account operations
    ChangeEmail(userID uint, newEmail, password string) error
    DeactivateAccount(userID uint, reason string) error
    DeleteAccount(userID uint, password string) error
}
```

### 3.2 Security Service
**File:** `internal/services/security.go`

```go
type SecurityService interface {
    // Password operations
    ChangePassword(userID uint, currentPassword, newPassword string) error
    RequestPasswordReset(email string) error
    ResetPasswordWithToken(token, newPassword string) error
    ValidatePasswordResetToken(token string) (*models.PasswordResetToken, error)

    // Email verification
    SendVerificationEmail(userID uint) error
    VerifyEmail(token string) error
    RequestEmailChange(userID uint, newEmail string) error

    // Login tracking
    RecordLogin(userID uint, ipAddress, userAgent string, success bool) error
    GetLoginHistory(userID uint, limit, offset int) ([]models.LoginHistory, int, error)
    GetActiveSessions(userID uint) ([]models.LoginHistory, error)
    LogoutSession(sessionID uint) error
    LogoutAllSessions(userID uint) error

    // Activity logging
    LogActivity(userID uint, action, resource string, details map[string]interface{}) error
    GetActivityLog(userID uint, limit, offset int) ([]models.ActivityLog, int, error)
}
```

### 3.3 Dashboard Service
**File:** `internal/services/dashboard.go`

```go
type DashboardService interface {
    // User dashboard
    GetDashboardStats(userID uint) (map[string]interface{}, error)
    GetOrderHistory(userID uint, limit, offset int) ([]models.Order, int, error)
    GetPurchasedTemplates(userID uint) ([]models.Template, error)
    GetAuthoredBlogPosts(userID uint, limit, offset int) ([]models.BlogPost, int, error)
}
```

### 3.4 Admin Service
**File:** `internal/services/admin.go`

```go
type AdminService interface {
    // User management
    ListUsers(filters map[string]interface{}, limit, offset int) ([]models.User, int, error)
    GetUserDetails(userID uint) (*models.User, error)
    UpdateUserRole(userID uint, role string) error
    SuspendUser(userID uint, adminID uint, reason string, durationDays int) error
    UnsuspendUser(userID uint) error
    DeleteUser(userID uint) error

    // Statistics
    GetAdminDashboard() (map[string]interface{}, error)
    GetUserActivity(userID uint, limit, offset int) ([]models.ActivityLog, int, error)

    // Impersonation
    CreateImpersonationToken(userID uint, adminID uint) (string, error)
}
```

### 3.5 Storage Service Extension
**File:** `internal/services/storage.go`

Add methods for avatar uploads:
```go
func (s *StorageService) UploadAvatar(file multipart.File, userID uint, filename string) (string, error)
func (s *StorageService) DeleteAvatar(avatarURL string) error
```

### 3.6 Email Service
**File:** `internal/services/email.go`

```go
type EmailService interface {
    // Template emails
    SendWelcomeEmail(user *models.User) error
    SendPasswordResetEmail(user *models.User, resetLink string) error
    SendEmailVerificationEmail(user *models.User, verificationLink string) error
    SendEmailChangeConfirmation(user *models.User, newEmail, verificationLink string) error
    SendAccountSuspensionEmail(user *models.User, reason string) error
    SendAccountDeletionEmail(user *models.User) error

    // Notification emails
    SendOrderConfirmation(user *models.User, order *models.Order) error
    SendPasswordChangedNotification(user *models.User) error
}
```

---

## Phase 4: Handlers Implementation

### 4.1 Profile Handler
**File:** `internal/handlers/profile.go`

```go
type ProfileHandler struct {
    profileService  services.ProfileService
    securityService services.SecurityService
}

// Implement all profile endpoint handlers
func (h *ProfileHandler) GetProfile(c *gin.Context)
func (h *ProfileHandler) UpdateProfile(c *gin.Context)
func (h *ProfileHandler) UploadAvatar(c *gin.Context)
func (h *ProfileHandler) DeleteAvatar(c *gin.Context)
func (h *ProfileHandler) GetPublicProfile(c *gin.Context)
func (h *ProfileHandler) GetPreferences(c *gin.Context)
func (h *ProfileHandler) UpdatePreferences(c *gin.Context)
func (h *ProfileHandler) ChangeEmail(c *gin.Context)
func (h *ProfileHandler) DeactivateAccount(c *gin.Context)
func (h *ProfileHandler) DeleteAccount(c *gin.Context)
```

### 4.2 Security Handler
**File:** `internal/handlers/security.go`

```go
type SecurityHandler struct {
    securityService services.SecurityService
    authService     services.AuthService
}

func (h *SecurityHandler) ForgotPassword(c *gin.Context)
func (h *SecurityHandler) ResetPassword(c *gin.Context)
func (h *SecurityHandler) ChangePassword(c *gin.Context)
func (h *SecurityHandler) VerifyEmail(c *gin.Context)
func (h *SecurityHandler) ResendVerification(c *gin.Context)
func (h *SecurityHandler) GetLoginHistory(c *gin.Context)
func (h *SecurityHandler) GetActiveSessions(c *gin.Context)
func (h *SecurityHandler) LogoutSession(c *gin.Context)
func (h *SecurityHandler) LogoutAll(c *gin.Context)
```

### 4.3 Dashboard Handler
**File:** `internal/handlers/dashboard.go`

```go
type DashboardHandler struct {
    dashboardService services.DashboardService
}

func (h *DashboardHandler) GetDashboard(c *gin.Context)
func (h *DashboardHandler) GetOrders(c *gin.Context)
func (h *DashboardHandler) GetOrder(c *gin.Context)
func (h *DashboardHandler) GetPurchasedTemplates(c *gin.Context)
func (h *DashboardHandler) GetBlogPosts(c *gin.Context)
func (h *DashboardHandler) GetActivity(c *gin.Context)
```

### 4.4 Admin Handler
**File:** `internal/handlers/admin.go`

```go
type AdminHandler struct {
    adminService services.AdminService
}

func (h *AdminHandler) ListUsers(c *gin.Context)
func (h *AdminHandler) GetUser(c *gin.Context)
func (h *AdminHandler) UpdateUserRole(c *gin.Context)
func (h *AdminHandler) SuspendUser(c *gin.Context)
func (h *AdminHandler) UnsuspendUser(c *gin.Context)
func (h *AdminHandler) DeleteUser(c *gin.Context)
func (h *AdminHandler) GetDashboard(c *gin.Context)
func (h *AdminHandler) GetUserActivity(c *gin.Context)
func (h *AdminHandler) ImpersonateUser(c *gin.Context)
```

---

## Phase 5: Middleware Enhancements

### 5.1 Role-Based Authorization Middleware
**File:** `internal/middleware/rbac.go`

```go
// RequireRole ensures user has specified role
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := getUserRoleFromContext(c)

        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{
            "error": "Insufficient permissions",
        })
        c.Abort()
    }
}

// RequireAdmin is a shorthand for RequireRole("admin")
func RequireAdmin() gin.HandlerFunc {
    return RequireRole("admin")
}
```

### 5.2 Activity Logging Middleware
**File:** `internal/middleware/activity.go`

```go
// LogActivity logs user actions
func LogActivity(securityService services.SecurityService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Track request details
        c.Next()

        // After request, log if successful
        if c.Writer.Status() < 400 {
            userID := getUserIDFromContext(c)
            action := determineAction(c)

            securityService.LogActivity(userID, action, c.Request.URL.Path, nil)
        }
    }
}
```

### 5.3 Rate Limiting Middleware
**File:** `internal/middleware/rate_limit.go`

```go
// RateLimitByUser limits requests per user
func RateLimitByUser(maxRequests int, window time.Duration) gin.HandlerFunc {
    // Implementation using in-memory cache or Redis
}

// RateLimitByIP limits requests per IP
func RateLimitByIP(maxRequests int, window time.Duration) gin.HandlerFunc {
    // Implementation using in-memory cache or Redis
}
```

---

## Phase 6: Frontend Implementation

### 6.1 New Pages/Views

#### Profile Page
**File:** `web/profile.html`
- View and edit profile information
- Upload/change avatar
- Update address details

#### Account Settings Page
**File:** `web/settings.html`
- Email preferences
- Privacy settings
- Language/timezone selection
- Theme selection
- Account deletion

#### Security Page
**File:** `web/security.html`
- Change password form
- Login history table
- Active sessions management
- Two-factor authentication (future)

#### Dashboard Page
**File:** `web/dashboard.html`
- User statistics cards
- Recent orders table
- Purchased templates grid
- Authored blog posts list
- Activity timeline

#### Admin Panel
**File:** `web/admin.html`
- User management table with search/filters
- User details modal
- Role management
- Suspend/unsuspend actions
- Admin dashboard statistics

### 6.2 JavaScript Components

#### API Client Extension
**File:** `web/js/api.js`
```javascript
class APIClient {
    // Profile methods
    async getProfile() { ... }
    async updateProfile(data) { ... }
    async uploadAvatar(file) { ... }
    async deleteAvatar() { ... }

    // Security methods
    async changePassword(currentPassword, newPassword) { ... }
    async forgotPassword(email) { ... }
    async resetPassword(token, newPassword) { ... }
    async getLoginHistory(limit, offset) { ... }
    async logoutSession(sessionId) { ... }

    // Settings methods
    async getPreferences() { ... }
    async updatePreferences(prefs) { ... }

    // Dashboard methods
    async getDashboard() { ... }
    async getOrders(limit, offset) { ... }
    async getPurchasedTemplates() { ... }

    // Admin methods
    async getUsers(filters) { ... }
    async getUserDetails(userId) { ... }
    async updateUserRole(userId, role) { ... }
    async suspendUser(userId, reason) { ... }
}
```

#### Profile Component
**File:** `web/js/components/profile.js`
```javascript
class ProfileComponent {
    constructor() {
        this.api = new APIClient();
    }

    async loadProfile() { ... }
    async saveProfile(formData) { ... }
    handleAvatarUpload(file) { ... }
    validateForm(data) { ... }
}
```

#### Settings Component
**File:** `web/js/components/settings.js`
```javascript
class SettingsComponent {
    async loadPreferences() { ... }
    async savePreferences(prefs) { ... }
    handleThemeChange(theme) { ... }
}
```

#### Dashboard Component
**File:** `web/js/components/dashboard.js`
```javascript
class DashboardComponent {
    async loadDashboard() { ... }
    renderStats(stats) { ... }
    renderOrders(orders) { ... }
    renderTemplates(templates) { ... }
}
```

#### Admin Component
**File:** `web/js/components/admin.js`
```javascript
class AdminComponent {
    async loadUsers(filters) { ... }
    async suspendUser(userId, reason) { ... }
    async updateRole(userId, role) { ... }
    renderUserTable(users) { ... }
}
```

### 6.3 UI/UX Components

#### Avatar Upload Component
- Drag & drop zone
- Image preview
- Crop functionality (optional)
- File type/size validation

#### Form Validation
- Real-time field validation
- Password strength indicator
- Email format validation
- Phone number formatting

#### Data Tables
- Sortable columns
- Pagination
- Search functionality
- Export to CSV (optional)

#### Modals/Dialogs
- Confirmation dialogs (delete account, etc.)
- User details modal (admin)
- Image upload modal

#### Toast Notifications
- Success messages
- Error messages
- Info messages
- Loading states

---

## Phase 7: Email Templates

### 7.1 Email Template Files
**Directory:** `internal/templates/emails/`

Create HTML email templates:

1. **welcome.html** - Welcome new users
2. **password_reset.html** - Password reset link
3. **email_verification.html** - Email verification link
4. **email_change.html** - Confirm email change
5. **password_changed.html** - Password change notification
6. **account_suspended.html** - Account suspension notice
7. **account_deleted.html** - Account deletion confirmation
8. **order_confirmation.html** - Order confirmation (enhance existing)

### 7.2 Email Template Service
**File:** `internal/services/email_template.go`

```go
type EmailTemplateService struct {
    sendgridClient *sendgrid.Client
}

func (s *EmailTemplateService) RenderTemplate(templateName string, data map[string]interface{}) (string, error)
func (s *EmailTemplateService) SendEmail(to, subject, htmlContent string) error
```

---

## Phase 8: Security Considerations

### 8.1 Authentication & Authorization
- ✅ JWT token validation on all protected routes
- ✅ Role-based access control (RBAC)
- ✅ Admin-only endpoints protected
- ✅ User can only access their own data
- ✅ Password hashing via Cognito

### 8.2 Data Protection
- Validate all user inputs
- Sanitize HTML content in bio/profiles
- SQL injection prevention (GORM handles this)
- XSS protection in frontend
- CSRF protection for state-changing operations

### 8.3 Privacy
- Soft delete accounts (keep for audit)
- Anonymize data after account deletion
- GDPR compliance considerations
- Privacy settings enforcement
- Public profile visibility controls

### 8.4 Rate Limiting
- Login attempts: 5 per 15 minutes
- Password reset: 3 per hour
- Profile updates: 10 per hour
- API requests: 100 per minute (general)

### 8.5 Audit Trail
- Log all sensitive operations
- Track admin actions
- Record login attempts
- Store IP addresses and user agents

---

## Phase 9: Testing Strategy

### 9.1 Unit Tests
**Files:** `*_test.go` files for each service/handler

Test coverage for:
- Profile service methods
- Security service methods
- Admin service methods
- Dashboard service methods
- Email service methods

### 9.2 Integration Tests
- API endpoint tests
- Database transaction tests
- Email sending tests (mock SendGrid)
- File upload tests (mock S3)

### 9.3 End-to-End Tests
- User registration flow
- Profile update flow
- Password reset flow
- Email verification flow
- Admin user management flow

### 9.4 Manual Testing Checklist
- [ ] User can register and login
- [ ] User can update profile
- [ ] User can upload avatar
- [ ] User can change password
- [ ] User can reset forgotten password
- [ ] User receives verification email
- [ ] User can view order history
- [ ] User can view purchased templates
- [ ] Admin can view all users
- [ ] Admin can suspend users
- [ ] Admin can change user roles
- [ ] Email preferences work
- [ ] Privacy settings work
- [ ] Account deletion works

---

## Phase 10: Deployment Considerations

### 10.1 Environment Variables
Add to `.env`:
```bash
# Email settings
SENDGRID_API_KEY=SG.xxx
SENDGRID_FROM_EMAIL=noreply@templatestore.com
SENDGRID_FROM_NAME=Template Store

# File upload settings
MAX_AVATAR_SIZE_MB=5
ALLOWED_AVATAR_TYPES=image/jpeg,image/png,image/webp

# Security settings
PASSWORD_RESET_TOKEN_EXPIRY_HOURS=24
EMAIL_VERIFICATION_TOKEN_EXPIRY_HOURS=72
SESSION_TIMEOUT_HOURS=24

# Rate limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# Frontend URL (for email links)
FRONTEND_URL=https://templatestore.com
```

### 10.2 Database Migration
- Run `AutoMigrate` to add new tables
- Create indexes for performance:
  - `users.email` (unique)
  - `users.role`
  - `users.status`
  - `login_history.user_id`
  - `activity_log.user_id`
  - `password_reset_token.token`

### 10.3 CDN Configuration
- Configure S3 bucket for avatar uploads
- Set up CloudFront for avatar delivery
- Configure CORS for file uploads

### 10.4 Monitoring
- Track failed login attempts
- Monitor password reset requests
- Alert on unusual admin activity
- Track API error rates

---

## Implementation Timeline Estimate

### Week 1-2: Database & Backend Foundation
- Day 1-2: Update database models and run migrations
- Day 3-4: Implement profile service
- Day 5-6: Implement security service
- Day 7-8: Implement dashboard service
- Day 9-10: Implement admin service

### Week 3: API Endpoints
- Day 11-12: Profile handlers and routes
- Day 13-14: Security handlers and routes
- Day 15-16: Dashboard handlers and routes
- Day 17: Admin handlers and routes

### Week 4: Email & Storage
- Day 18-19: Email templates and service
- Day 20: Avatar upload/storage integration
- Day 21: Testing email flows

### Week 5-6: Frontend Development
- Day 22-24: Profile page
- Day 25-26: Settings page
- Day 27-28: Security page
- Day 29-31: Dashboard page
- Day 32-34: Admin panel

### Week 7: Testing & Polish
- Day 35-37: Unit and integration tests
- Day 38-39: End-to-end testing
- Day 40-41: Bug fixes and polish
- Day 42: Documentation

### Week 8: Deployment & Launch
- Day 43-44: Staging deployment and testing
- Day 45: Production deployment
- Day 46-47: Monitoring and hotfixes
- Day 48: Post-launch review

---

## Priority Order for Implementation

### Phase 1 (High Priority - Core Features)
1. Database schema updates
2. Profile management (view/edit)
3. Avatar upload
4. Password change
5. Password reset (forgot password)
6. Email verification
7. User dashboard (orders, templates)

### Phase 2 (Medium Priority - Enhanced Features)
8. Account settings (preferences)
9. Login history
10. Activity log
11. Email change
12. Admin user list
13. Admin user management

### Phase 3 (Low Priority - Advanced Features)
14. Session management
15. Account deactivation
16. Account deletion
17. Admin dashboard
18. User impersonation
19. Advanced analytics

---

## Success Metrics

### User Engagement
- Profile completion rate
- Avatar upload rate
- Settings customization rate
- Password changes per month

### Security
- Failed login attempts
- Password reset requests
- Email verification rate
- Suspicious activity detected

### Business
- User retention rate
- Admin efficiency (time to resolve user issues)
- Support ticket reduction

---

## Future Enhancements (Post-MVP)

1. **Two-Factor Authentication (2FA)**
   - SMS verification
   - Authenticator app (TOTP)
   - Backup codes

2. **Social Login**
   - Google OAuth
   - GitHub OAuth
   - Apple Sign In

3. **Advanced Privacy**
   - Data export (GDPR)
   - Privacy dashboard
   - Cookie consent management

4. **Notifications System**
   - In-app notifications
   - Push notifications
   - Notification preferences

5. **User Roles & Permissions**
   - Custom roles
   - Granular permissions
   - Team/organization accounts

6. **Advanced Analytics**
   - User behavior tracking
   - Conversion funnels
   - A/B testing

---

## Questions to Answer Before Starting

1. **Email Provider:** Confirm SendGrid setup and templates?
2. **Avatar Storage:** S3 bucket name and CDN configuration?
3. **Cognito Integration:** How to handle password changes with Cognito?
4. **Admin Roles:** Who gets initial admin access?
5. **Data Retention:** How long to keep deleted user data?
6. **Email Frequency:** Limits on verification/reset emails?
7. **Testing:** Need test user accounts in Cognito?

---

## Documentation to Create

1. **API Documentation**
   - Swagger/OpenAPI spec for all new endpoints
   - Authentication guide
   - Rate limiting documentation

2. **User Guides**
   - Profile management guide
   - Security best practices
   - Admin panel guide

3. **Developer Documentation**
   - Service architecture
   - Database schema
   - Email template customization
   - Testing guide

---

This plan provides a complete roadmap for implementing comprehensive user account management. We can start with Phase 1 (database schema) and proceed incrementally, or focus on specific high-priority features first. Which approach would you prefer?

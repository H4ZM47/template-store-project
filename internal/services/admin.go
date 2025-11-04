package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"template-store/internal/models"
)

var (
	ErrUnauthorized        = errors.New("unauthorized: admin access required")
	ErrCannotSuspendSelf   = errors.New("cannot suspend your own account")
	ErrCannotDeleteSelf    = errors.New("cannot delete your own account")
	ErrInvalidRole         = errors.New("invalid role specified")
	ErrUserAlreadySuspended = errors.New("user is already suspended")
)

// AdminDashboardStats represents admin dashboard statistics
type AdminDashboardStats struct {
	TotalUsers        int64   `json:"total_users"`
	ActiveUsers       int64   `json:"active_users"`
	SuspendedUsers    int64   `json:"suspended_users"`
	NewUsersThisMonth int64   `json:"new_users_this_month"`
	TotalOrders       int64   `json:"total_orders"`
	TotalRevenue      float64 `json:"total_revenue"`
	RevenueThisMonth  float64 `json:"revenue_this_month"`
}

// UserFilters represents filters for user listing
type UserFilters struct {
	Search string // Search by name or email
	Role   string // Filter by role
	Status string // Filter by status
}

// AdminService defines the interface for admin operations
type AdminService interface {
	// User management
	ListUsers(filters UserFilters, sortBy, order string, limit, offset int) ([]models.User, int64, error)
	GetUserDetails(userID uint) (*models.User, error)
	UpdateUserRole(adminID uint, userID uint, role string) error
	SuspendUser(adminID uint, userID uint, reason string, durationDays int) error
	UnsuspendUser(adminID uint, userID uint) error
	DeleteUser(adminID uint, userID uint) error

	// Statistics
	GetAdminDashboard() (*AdminDashboardStats, map[string]interface{}, error)
	GetUserActivity(userID uint, limit, offset int) ([]models.ActivityLog, int64, error)
	GetUserStats(userID uint) (map[string]interface{}, error)

	// Verification
	IsAdmin(userID uint) (bool, error)
}

// AdminServiceImpl implements the AdminService interface
type AdminServiceImpl struct {
	db              *gorm.DB
	securityService SecurityService
}

// NewAdminService creates a new AdminService instance
func NewAdminService(db *gorm.DB, securityService SecurityService) AdminService {
	return &AdminServiceImpl{
		db:              db,
		securityService: securityService,
	}
}

// IsAdmin checks if a user has admin role
func (s *AdminServiceImpl) IsAdmin(userID uint) (bool, error) {
	var user models.User
	if err := s.db.Select("role").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrUserNotFound
		}
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	return user.Role == "admin", nil
}

// ListUsers retrieves a paginated, filtered list of users
func (s *AdminServiceImpl) ListUsers(filters UserFilters, sortBy, order string, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})

	// Apply filters
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", searchTerm, searchTerm)
	}

	if filters.Role != "" {
		query = query.Where("role = ?", filters.Role)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Apply sorting
	validSortFields := map[string]bool{
		"id": true, "name": true, "email": true, "created_at": true,
		"last_login_at": true, "role": true, "status": true,
	}
	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	orderClause := fmt.Sprintf("%s %s", sortBy, order)

	// Get paginated records
	if err := query.
		Preload("Orders").
		Order(orderClause).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Calculate additional stats for each user
	for i := range users {
		// Count orders
		var orderCount int64
		s.db.Model(&models.Order{}).Where("user_id = ?", users[i].ID).Count(&orderCount)
		// Note: We could add this to a custom struct if needed
	}

	return users, total, nil
}

// GetUserDetails retrieves detailed information about a user
func (s *AdminServiceImpl) GetUserDetails(userID uint) (*models.User, error) {
	var user models.User

	if err := s.db.
		Preload("Orders").
		Preload("Orders.Template").
		Preload("BlogPosts").
		First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	return &user, nil
}

// UpdateUserRole updates a user's role
func (s *AdminServiceImpl) UpdateUserRole(adminID uint, userID uint, role string) error {
	// Verify admin privileges
	isAdmin, err := s.IsAdmin(adminID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrUnauthorized
	}

	// Validate role
	validRoles := map[string]bool{
		"user": true, "admin": true, "author": true,
	}
	if !validRoles[role] {
		return ErrInvalidRole
	}

	// Check if user exists
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Update role
	if err := s.db.Model(&user).Update("role", role).Error; err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	// Log activity
	details := map[string]interface{}{
		"old_role": user.Role,
		"new_role": role,
		"admin_id": adminID,
	}
	_ = s.securityService.LogActivity(userID, models.ActivityRoleChanged, "user", &userID, details, "", "")

	return nil
}

// SuspendUser suspends a user account
func (s *AdminServiceImpl) SuspendUser(adminID uint, userID uint, reason string, durationDays int) error {
	// Verify admin privileges
	isAdmin, err := s.IsAdmin(adminID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrUnauthorized
	}

	// Prevent self-suspension
	if adminID == userID {
		return ErrCannotSuspendSelf
	}

	// Check if user exists
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if already suspended
	if user.Status == "suspended" {
		return ErrUserAlreadySuspended
	}

	// Calculate suspension end date if duration is specified
	now := time.Now()
	updates := map[string]interface{}{
		"status":            "suspended",
		"suspended_at":      now,
		"suspended_by":      adminID,
		"suspension_reason": reason,
	}

	// Update user
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to suspend user: %w", err)
	}

	// Log activity
	details := map[string]interface{}{
		"reason":        reason,
		"duration_days": durationDays,
		"admin_id":      adminID,
	}
	_ = s.securityService.LogActivity(userID, models.ActivityUserSuspended, "user", &userID, details, "", "")

	// TODO: Send suspension notification email

	return nil
}

// UnsuspendUser unsuspends a user account
func (s *AdminServiceImpl) UnsuspendUser(adminID uint, userID uint) error {
	// Verify admin privileges
	isAdmin, err := s.IsAdmin(adminID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrUnauthorized
	}

	// Check if user exists
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Update user status
	updates := map[string]interface{}{
		"status":            "active",
		"suspended_at":      nil,
		"suspended_by":      nil,
		"suspension_reason": "",
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to unsuspend user: %w", err)
	}

	// Log activity
	details := map[string]interface{}{
		"admin_id": adminID,
	}
	_ = s.securityService.LogActivity(userID, models.ActivityUserUnsuspended, "user", &userID, details, "", "")

	// TODO: Send unsuspension notification email

	return nil
}

// DeleteUser deletes a user (soft delete)
func (s *AdminServiceImpl) DeleteUser(adminID uint, userID uint) error {
	// Verify admin privileges
	isAdmin, err := s.IsAdmin(adminID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrUnauthorized
	}

	// Prevent self-deletion
	if adminID == userID {
		return ErrCannotDeleteSelf
	}

	// Check if user exists
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Log activity before deletion
	details := map[string]interface{}{
		"admin_id": adminID,
		"user_name": user.Name,
		"user_email": user.Email,
	}
	_ = s.securityService.LogActivity(adminID, models.ActivityUserDeleted, "user", &userID, details, "", "")

	// Soft delete
	if err := s.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetAdminDashboard retrieves admin dashboard statistics
func (s *AdminServiceImpl) GetAdminDashboard() (*AdminDashboardStats, map[string]interface{}, error) {
	stats := &AdminDashboardStats{}

	// Count total users
	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Count active users
	if err := s.db.Model(&models.User{}).Where("status = ?", "active").Count(&stats.ActiveUsers).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count active users: %w", err)
	}

	// Count suspended users
	if err := s.db.Model(&models.User{}).Where("status = ?", "suspended").Count(&stats.SuspendedUsers).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count suspended users: %w", err)
	}

	// Count new users this month
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Truncate(24 * time.Hour)
	if err := s.db.Model(&models.User{}).Where("created_at >= ?", startOfMonth).Count(&stats.NewUsersThisMonth).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count new users: %w", err)
	}

	// Count total orders
	if err := s.db.Model(&models.Order{}).Count(&stats.TotalOrders).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count orders: %w", err)
	}

	// Calculate revenue (TODO: implement when Order model has amount field)
	stats.TotalRevenue = 0
	stats.RevenueThisMonth = 0

	// Get recent users (last 10)
	var recentUsers []models.User
	if err := s.db.Order("created_at DESC").Limit(10).Find(&recentUsers).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get recent users: %w", err)
	}

	// Get recent orders (last 10)
	var recentOrders []models.Order
	if err := s.db.Preload("User").Preload("Template").Order("created_at DESC").Limit(10).Find(&recentOrders).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get recent orders: %w", err)
	}

	additional := map[string]interface{}{
		"recent_users":  recentUsers,
		"recent_orders": recentOrders,
	}

	return stats, additional, nil
}

// GetUserActivity retrieves a user's activity log (admin view)
func (s *AdminServiceImpl) GetUserActivity(userID uint, limit, offset int) ([]models.ActivityLog, int64, error) {
	return s.securityService.GetActivityLog(userID, limit, offset)
}

// GetUserStats retrieves statistics for a specific user
func (s *AdminServiceImpl) GetUserStats(userID uint) (map[string]interface{}, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Count orders
	var orderCount int64
	s.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&orderCount)

	// Count blog posts
	var blogPostCount int64
	s.db.Model(&models.BlogPost{}).Where("author_id = ?", userID).Count(&blogPostCount)

	// Count login attempts
	var loginCount int64
	s.db.Model(&models.LoginHistory{}).Where("user_id = ? AND success = ?", userID, true).Count(&loginCount)

	// Get last login
	var lastLogin models.LoginHistory
	s.db.Where("user_id = ? AND success = ?", userID, true).Order("login_at DESC").First(&lastLogin)

	stats := map[string]interface{}{
		"total_orders":      orderCount,
		"total_blog_posts":  blogPostCount,
		"total_logins":      loginCount,
		"last_login_at":     user.LastLoginAt,
		"last_login_ip":     lastLogin.IPAddress,
		"account_age_days":  int(time.Since(user.CreatedAt).Hours() / 24),
		"status":            user.Status,
		"role":              user.Role,
		"email_verified":    user.EmailVerified,
	}

	return stats, nil
}

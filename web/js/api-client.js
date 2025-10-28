// API Client for Template Store User Account Management
class APIClient {
    constructor() {
        this.apiBase = 'http://localhost:8080/api/v1';
        this.token = localStorage.getItem('auth_token');
    }

    // Helper method to get auth headers
    getHeaders(includeAuth = true) {
        const headers = {
            'Content-Type': 'application/json'
        };

        if (includeAuth && this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        return headers;
    }

    // Helper method for multipart/form-data
    getMultipartHeaders(includeAuth = true) {
        const headers = {};

        if (includeAuth && this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        return headers;
    }

    // Set authentication token
    setToken(token) {
        this.token = token;
        localStorage.setItem('auth_token', token);
    }

    // Clear authentication token
    clearToken() {
        this.token = null;
        localStorage.removeItem('auth_token');
    }

    // Generic request handler
    async request(url, options = {}) {
        try {
            const response = await fetch(url, options);

            // Handle non-JSON responses
            const contentType = response.headers.get('content-type');
            if (!contentType || !contentType.includes('application/json')) {
                if (response.ok) {
                    return { success: true };
                }
                throw new Error('Server error');
            }

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || `HTTP error ${response.status}`);
            }

            return data;
        } catch (error) {
            console.error('API Request failed:', error);
            throw error;
        }
    }

    // ========== AUTH ENDPOINTS ==========

    async register(name, email, password) {
        return this.request(`${this.apiBase}/auth/register`, {
            method: 'POST',
            headers: this.getHeaders(false),
            body: JSON.stringify({ name, email, password })
        });
    }

    async login(email, password) {
        const data = await this.request(`${this.apiBase}/auth/login`, {
            method: 'POST',
            headers: this.getHeaders(false),
            body: JSON.stringify({ email, password })
        });

        if (data.token) {
            this.setToken(data.token);
        }

        return data;
    }

    async forgotPassword(email) {
        return this.request(`${this.apiBase}/auth/forgot-password`, {
            method: 'POST',
            headers: this.getHeaders(false),
            body: JSON.stringify({ email })
        });
    }

    async resetPassword(token, newPassword) {
        return this.request(`${this.apiBase}/auth/reset-password`, {
            method: 'POST',
            headers: this.getHeaders(false),
            body: JSON.stringify({ token, new_password: newPassword })
        });
    }

    async changePassword(currentPassword, newPassword) {
        return this.request(`${this.apiBase}/auth/change-password`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify({ current_password: currentPassword, new_password: newPassword })
        });
    }

    async verifyEmail(token) {
        return this.request(`${this.apiBase}/auth/verify-email`, {
            method: 'POST',
            headers: this.getHeaders(false),
            body: JSON.stringify({ token })
        });
    }

    async resendVerification() {
        return this.request(`${this.apiBase}/auth/resend-verification`, {
            method: 'POST',
            headers: this.getHeaders()
        });
    }

    // ========== PROFILE ENDPOINTS ==========

    async getProfile() {
        return this.request(`${this.apiBase}/profile`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async updateProfile(updates) {
        return this.request(`${this.apiBase}/profile`, {
            method: 'PUT',
            headers: this.getHeaders(),
            body: JSON.stringify(updates)
        });
    }

    async uploadAvatar(file) {
        const formData = new FormData();
        formData.append('avatar', file);

        return this.request(`${this.apiBase}/profile/avatar`, {
            method: 'POST',
            headers: this.getMultipartHeaders(),
            body: formData
        });
    }

    async deleteAvatar() {
        return this.request(`${this.apiBase}/profile/avatar`, {
            method: 'DELETE',
            headers: this.getHeaders()
        });
    }

    async getPublicProfile(userId) {
        return this.request(`${this.apiBase}/users/${userId}/profile`, {
            method: 'GET',
            headers: this.getHeaders(false)
        });
    }

    // ========== PREFERENCES ENDPOINTS ==========

    async getPreferences() {
        return this.request(`${this.apiBase}/profile/preferences`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async updatePreferences(preferences) {
        return this.request(`${this.apiBase}/profile/preferences`, {
            method: 'PUT',
            headers: this.getHeaders(),
            body: JSON.stringify(preferences)
        });
    }

    // ========== ACCOUNT MANAGEMENT ENDPOINTS ==========

    async deactivateAccount(password, reason) {
        return this.request(`${this.apiBase}/profile/deactivate`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify({ password, reason })
        });
    }

    async deleteAccount(password, confirmation) {
        return this.request(`${this.apiBase}/profile`, {
            method: 'DELETE',
            headers: this.getHeaders(),
            body: JSON.stringify({ password, confirmation })
        });
    }

    // ========== DASHBOARD ENDPOINTS ==========

    async getDashboard() {
        return this.request(`${this.apiBase}/profile/dashboard`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async getOrders(limit = 10, offset = 0, status = '') {
        let url = `${this.apiBase}/profile/orders?limit=${limit}&offset=${offset}`;
        if (status) url += `&status=${status}`;

        return this.request(url, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async getOrder(orderId) {
        return this.request(`${this.apiBase}/profile/orders/${orderId}`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async getPurchasedTemplates() {
        return this.request(`${this.apiBase}/profile/purchased-templates`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async getTemplateDownload(templateId) {
        return this.request(`${this.apiBase}/profile/templates/${templateId}/download`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async getBlogPosts(limit = 10, offset = 0, status = '') {
        let url = `${this.apiBase}/profile/blog-posts?limit=${limit}&offset=${offset}`;
        if (status) url += `&status=${status}`;

        return this.request(url, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    // ========== SECURITY ENDPOINTS ==========

    async getLoginHistory(limit = 20, offset = 0) {
        return this.request(`${this.apiBase}/profile/login-history?limit=${limit}&offset=${offset}`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async getActiveSessions() {
        return this.request(`${this.apiBase}/profile/sessions`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async logoutSession(sessionId) {
        return this.request(`${this.apiBase}/auth/logout-session/${sessionId}`, {
            method: 'POST',
            headers: this.getHeaders()
        });
    }

    async logoutAll() {
        return this.request(`${this.apiBase}/auth/logout-all`, {
            method: 'POST',
            headers: this.getHeaders()
        });
    }

    async getActivityLog(limit = 20, offset = 0) {
        return this.request(`${this.apiBase}/profile/activity?limit=${limit}&offset=${offset}`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    // ========== ADMIN ENDPOINTS ==========

    async getAdminDashboard() {
        return this.request(`${this.apiBase}/admin/dashboard`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async listUsers(filters = {}, limit = 50, offset = 0) {
        let url = `${this.apiBase}/admin/users?limit=${limit}&offset=${offset}`;

        if (filters.search) url += `&search=${encodeURIComponent(filters.search)}`;
        if (filters.role) url += `&role=${filters.role}`;
        if (filters.status) url += `&status=${filters.status}`;
        if (filters.sort) url += `&sort=${filters.sort}`;
        if (filters.order) url += `&order=${filters.order}`;

        return this.request(url, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async getUser(userId) {
        return this.request(`${this.apiBase}/admin/users/${userId}`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }

    async updateUserRole(userId, role) {
        return this.request(`${this.apiBase}/admin/users/${userId}/role`, {
            method: 'PUT',
            headers: this.getHeaders(),
            body: JSON.stringify({ role })
        });
    }

    async suspendUser(userId, reason, durationDays = 0) {
        return this.request(`${this.apiBase}/admin/users/${userId}/suspend`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify({ reason, duration_days: durationDays })
        });
    }

    async unsuspendUser(userId) {
        return this.request(`${this.apiBase}/admin/users/${userId}/unsuspend`, {
            method: 'POST',
            headers: this.getHeaders()
        });
    }

    async deleteUser(userId) {
        return this.request(`${this.apiBase}/admin/users/${userId}`, {
            method: 'DELETE',
            headers: this.getHeaders()
        });
    }

    async getUserActivity(userId, limit = 20, offset = 0) {
        return this.request(`${this.apiBase}/admin/users/${userId}/activity?limit=${limit}&offset=${offset}`, {
            method: 'GET',
            headers: this.getHeaders()
        });
    }
}

// Create singleton instance
const apiClient = new APIClient();

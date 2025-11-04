// Security Page JavaScript
class SecurityPage {
    constructor() {
        this.api = apiClient;
        this.historyPage = 0;
        this.activityPage = 0;
        this.pageSize = 10;
        this.init();
    }

    async init() {
        // Check if user is authenticated
        if (!this.api.token) {
            window.location.href = 'index.html';
            return;
        }

        this.setupEventListeners();
        await Promise.all([
            this.loadActiveSessions(),
            this.loadLoginHistory(),
            this.loadActivity()
        ]);
    }

    setupEventListeners() {
        // Change password form
        document.getElementById('change-password-form').addEventListener('submit', (e) => this.handleChangePassword(e));

        // Logout all sessions
        document.getElementById('logout-all-btn').addEventListener('click', () => this.handleLogoutAll());

        // Load more buttons
        document.getElementById('load-more-history').addEventListener('click', () => this.loadMoreHistory());
        document.getElementById('load-more-activity').addEventListener('click', () => this.loadMoreActivity());

        // Logout
        document.getElementById('logout-btn').addEventListener('click', () => this.handleLogout());
    }

    async handleChangePassword(e) {
        e.preventDefault();

        const currentPassword = document.getElementById('current_password').value;
        const newPassword = document.getElementById('new_password').value;
        const confirmPassword = document.getElementById('confirm_password').value;

        // Validate passwords match
        if (newPassword !== confirmPassword) {
            this.showMessage('New passwords do not match', 'error');
            return;
        }

        // Validate password strength
        if (newPassword.length < 8) {
            this.showMessage('Password must be at least 8 characters', 'error');
            return;
        }

        try {
            await this.api.changePassword(currentPassword, newPassword);
            this.showMessage('Password changed successfully!', 'success');
            document.getElementById('change-password-form').reset();
        } catch (error) {
            this.showMessage('Failed to change password: ' + error.message, 'error');
        }
    }

    async loadActiveSessions() {
        try {
            const data = await this.api.getActiveSessions();
            const sessions = data.sessions || [];
            this.displaySessions(sessions);
        } catch (error) {
            this.showMessage('Failed to load active sessions: ' + error.message, 'error');
            document.getElementById('sessions-container').innerHTML = '<div class="text-center text-gray-500 py-4">Failed to load sessions</div>';
        }
    }

    displaySessions(sessions) {
        const container = document.getElementById('sessions-container');

        if (sessions.length === 0) {
            container.innerHTML = '<div class="text-center text-gray-500 py-4">No active sessions</div>';
            return;
        }

        container.innerHTML = sessions.map(session => `
            <div class="border border-gray-200 rounded-lg p-4">
                <div class="flex justify-between items-start">
                    <div class="flex-1">
                        <div class="flex items-center mb-2">
                            <i class="fas fa-${this.getDeviceIcon(session.device)} text-gray-600 mr-2"></i>
                            <span class="font-medium text-gray-900">${this.formatDevice(session.device)}</span>
                        </div>
                        <p class="text-sm text-gray-600">
                            <i class="fas fa-map-marker-alt text-gray-400 mr-1"></i>
                            ${session.location || 'Unknown location'}
                        </p>
                        <p class="text-sm text-gray-600">
                            <i class="fas fa-globe text-gray-400 mr-1"></i>
                            ${session.ip_address || 'Unknown IP'}
                        </p>
                        <p class="text-xs text-gray-500 mt-1">
                            Active since ${this.formatDateTime(session.login_at)}
                        </p>
                    </div>
                    <button class="text-red-600 hover:text-red-800 text-sm" onclick="securityPage.handleLogoutSession(${session.id})">
                        <i class="fas fa-sign-out-alt"></i> Logout
                    </button>
                </div>
            </div>
        `).join('');
    }

    async loadLoginHistory() {
        try {
            const data = await this.api.getLoginHistory(this.pageSize, this.historyPage * this.pageSize);
            const history = data.history || [];
            this.displayLoginHistory(history, this.historyPage === 0);
        } catch (error) {
            this.showMessage('Failed to load login history: ' + error.message, 'error');
            document.getElementById('history-container').innerHTML = '<div class="text-center text-gray-500 py-4">Failed to load history</div>';
        }
    }

    displayLoginHistory(history, replace = true) {
        const container = document.getElementById('history-container');

        if (history.length === 0 && replace) {
            container.innerHTML = '<div class="text-center text-gray-500 py-4">No login history</div>';
            return;
        }

        const html = history.map(item => `
            <div class="border border-gray-200 rounded-lg p-4">
                <div class="flex items-start justify-between">
                    <div class="flex-1">
                        <div class="flex items-center mb-1">
                            <i class="fas fa-${item.success ? 'check-circle text-green-600' : 'times-circle text-red-600'} mr-2"></i>
                            <span class="font-medium text-gray-900">${item.success ? 'Successful Login' : 'Failed Login'}</span>
                        </div>
                        <p class="text-sm text-gray-600">
                            <i class="fas fa-${this.getDeviceIcon(item.device)} text-gray-400 mr-1"></i>
                            ${this.formatDevice(item.device)} • ${item.location || 'Unknown'}
                        </p>
                        <p class="text-sm text-gray-600">
                            <i class="fas fa-globe text-gray-400 mr-1"></i>
                            ${item.ip_address}
                        </p>
                        ${!item.success && item.failure_reason ? `
                            <p class="text-sm text-red-600 mt-1">
                                <i class="fas fa-exclamation-circle mr-1"></i>
                                ${item.failure_reason}
                            </p>
                        ` : ''}
                        <p class="text-xs text-gray-500 mt-1">
                            ${this.formatDateTime(item.login_at)}
                        </p>
                    </div>
                </div>
            </div>
        `).join('');

        if (replace) {
            container.innerHTML = html;
        } else {
            container.innerHTML += html;
        }
    }

    async loadActivity() {
        try {
            const data = await this.api.getActivityLog(this.pageSize, this.activityPage * this.pageSize);
            const activities = data.activities || [];
            this.displayActivity(activities, this.activityPage === 0);
        } catch (error) {
            this.showMessage('Failed to load activity: ' + error.message, 'error');
            document.getElementById('activity-container').innerHTML = '<div class="text-center text-gray-500 py-4">Failed to load activity</div>';
        }
    }

    displayActivity(activities, replace = true) {
        const container = document.getElementById('activity-container');

        if (activities.length === 0 && replace) {
            container.innerHTML = '<div class="text-center text-gray-500 py-4">No recent activity</div>';
            return;
        }

        const html = activities.map(activity => `
            <div class="border-l-4 ${this.getActivityColor(activity.action)} pl-4 py-2">
                <div class="flex items-start justify-between">
                    <div class="flex-1">
                        <p class="font-medium text-gray-900">${this.formatActivityAction(activity.action)}</p>
                        <p class="text-sm text-gray-600">
                            ${activity.resource}${activity.resource_id ? ` #${activity.resource_id}` : ''}
                        </p>
                        <p class="text-xs text-gray-500 mt-1">
                            ${this.formatDateTime(activity.created_at)}
                        </p>
                    </div>
                    <i class="fas fa-${this.getActivityIcon(activity.action)} text-gray-400"></i>
                </div>
            </div>
        `).join('');

        if (replace) {
            container.innerHTML = html;
        } else {
            container.innerHTML += html;
        }
    }

    async loadMoreHistory() {
        this.historyPage++;
        await this.loadLoginHistory();
    }

    async loadMoreActivity() {
        this.activityPage++;
        await this.loadActivity();
    }

    async handleLogoutSession(sessionId) {
        if (!confirm('Are you sure you want to logout this session?')) {
            return;
        }

        try {
            await this.api.logoutSession(sessionId);
            this.showMessage('Session logged out successfully', 'success');
            await this.loadActiveSessions();
        } catch (error) {
            this.showMessage('Failed to logout session: ' + error.message, 'error');
        }
    }

    async handleLogoutAll() {
        if (!confirm('Are you sure you want to logout from all devices? You will need to login again.')) {
            return;
        }

        try {
            await this.api.logoutAll();
            this.showMessage('All sessions logged out. Redirecting...', 'success');
            setTimeout(() => {
                this.api.clearToken();
                window.location.href = 'index.html';
            }, 2000);
        } catch (error) {
            this.showMessage('Failed to logout all sessions: ' + error.message, 'error');
        }
    }

    handleLogout() {
        if (confirm('Are you sure you want to logout?')) {
            this.api.clearToken();
            window.location.href = 'index.html';
        }
    }

    // Helper methods
    getDeviceIcon(device) {
        switch (device) {
            case 'mobile': return 'mobile-alt';
            case 'tablet': return 'tablet-alt';
            case 'desktop': return 'desktop';
            default: return 'question-circle';
        }
    }

    formatDevice(device) {
        if (!device) return 'Unknown Device';
        return device.charAt(0).toUpperCase() + device.slice(1);
    }

    getActivityColor(action) {
        if (action.includes('delete') || action.includes('suspend')) return 'border-red-400';
        if (action.includes('create') || action.includes('upload')) return 'border-green-400';
        if (action.includes('update') || action.includes('change')) return 'border-blue-400';
        return 'border-gray-400';
    }

    getActivityIcon(action) {
        if (action.includes('password')) return 'key';
        if (action.includes('email')) return 'envelope';
        if (action.includes('profile')) return 'user';
        if (action.includes('avatar')) return 'image';
        if (action.includes('login')) return 'sign-in-alt';
        if (action.includes('logout')) return 'sign-out-alt';
        return 'circle';
    }

    formatActivityAction(action) {
        return action.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
    }

    formatDateTime(dateString) {
        if (!dateString) return 'Unknown';
        const date = new Date(dateString);
        return date.toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    showMessage(message, type) {
        const container = document.getElementById('message-container');
        const bgColor = type === 'success' ? 'bg-green-100 border-green-400 text-green-700' :
                        type === 'error' ? 'bg-red-100 border-red-400 text-red-700' :
                        'bg-blue-100 border-blue-400 text-blue-700';

        container.className = `border-l-4 p-4 ${bgColor}`;
        container.innerHTML = `
            <div class="flex">
                <div class="flex-shrink-0">
                    <i class="fas fa-${type === 'success' ? 'check-circle' : type === 'error' ? 'exclamation-circle' : 'info-circle'}"></i>
                </div>
                <div class="ml-3">
                    <p class="text-sm">${message}</p>
                </div>
            </div>
        `;
        container.classList.remove('hidden');

        // Auto-hide after 5 seconds
        setTimeout(() => {
            container.classList.add('hidden');
        }, 5000);
    }
}

// Initialize the page when DOM is ready
let securityPage;
document.addEventListener('DOMContentLoaded', () => {
    securityPage = new SecurityPage();
});

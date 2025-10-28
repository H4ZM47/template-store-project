// Admin Panel JavaScript
class AdminPanel {
    constructor() {
        this.api = apiClient;
        this.currentPage = 0;
        this.pageSize = 50;
        this.totalUsers = 0;
        this.filters = {};
        this.init();
    }

    async init() {
        // Check if user is authenticated
        if (!this.api.token) {
            window.location.href = 'index.html';
            return;
        }

        this.setupEventListeners();
        await this.loadAdminDashboard();
        await this.loadUsers();
    }

    setupEventListeners() {
        // Apply filters
        document.getElementById('apply-filters').addEventListener('click', () => this.applyFilters());
        document.getElementById('search-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.applyFilters();
        });

        // Pagination
        document.getElementById('prev-page').addEventListener('click', () => this.prevPage());
        document.getElementById('next-page').addEventListener('click', () => this.nextPage());

        // Modal close
        document.getElementById('close-modal').addEventListener('click', () => this.closeModal());
        document.getElementById('cancel-suspend').addEventListener('click', () => this.closeSuspendModal());

        // Suspend form
        document.getElementById('suspend-form').addEventListener('submit', (e) => this.handleSuspend(e));

        // Logout
        document.getElementById('logout-btn').addEventListener('click', () => this.handleLogout());
    }

    async loadAdminDashboard() {
        try {
            const data = await this.api.getAdminDashboard();
            const stats = data.stats || {};

            document.getElementById('total-users').textContent = stats.total_users || 0;
            document.getElementById('active-users').textContent = stats.active_users || 0;
            document.getElementById('new-users').textContent = stats.new_users_this_month || 0;
            document.getElementById('total-revenue').textContent = `$${(stats.total_revenue || 0).toFixed(2)}`;
        } catch (error) {
            console.error('Failed to load admin dashboard:', error);
            alert('You do not have admin access');
            window.location.href = 'dashboard.html';
        }
    }

    async loadUsers() {
        try {
            const data = await this.api.listUsers(this.filters, this.pageSize, this.currentPage * this.pageSize);
            const users = data.users || [];
            this.totalUsers = data.total || 0;

            this.displayUsers(users);
            this.updatePagination();
        } catch (error) {
            console.error('Failed to load users:', error);
            document.getElementById('users-table').innerHTML = '<tr><td colspan="6" class="px-6 py-8 text-center text-red-500">Failed to load users</td></tr>';
        }
    }

    displayUsers(users) {
        const tbody = document.getElementById('users-table');

        if (users.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="px-6 py-8 text-center text-gray-500">No users found</td></tr>';
            return;
        }

        tbody.innerHTML = users.map(user => `
            <tr class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10">
                            <img class="h-10 w-10 rounded-full" src="${user.avatar_url || 'https://via.placeholder.com/40'}" alt="">
                        </div>
                        <div class="ml-4">
                            <div class="text-sm font-medium text-gray-900">${user.name || 'No Name'}</div>
                            <div class="text-sm text-gray-500">${user.email}</div>
                        </div>
                    </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${this.getRoleBadgeColor(user.role)}">
                        ${this.formatRole(user.role)}
                    </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${this.getStatusBadgeColor(user.status)}">
                        ${this.formatStatus(user.status)}
                    </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    ${user.Orders?.length || 0}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    ${this.formatDate(user.created_at)}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                    <button onclick="adminPanel.viewUser(${user.id})" class="text-blue-600 hover:text-blue-900">
                        <i class="fas fa-eye"></i>
                    </button>
                    <button onclick="adminPanel.changeRole(${user.id})" class="text-green-600 hover:text-green-900">
                        <i class="fas fa-user-tag"></i>
                    </button>
                    ${user.status === 'active' ? `
                        <button onclick="adminPanel.suspendUser(${user.id})" class="text-yellow-600 hover:text-yellow-900">
                            <i class="fas fa-ban"></i>
                        </button>
                    ` : user.status === 'suspended' ? `
                        <button onclick="adminPanel.unsuspendUser(${user.id})" class="text-green-600 hover:text-green-900">
                            <i class="fas fa-check-circle"></i>
                        </button>
                    ` : ''}
                    <button onclick="adminPanel.deleteUser(${user.id})" class="text-red-600 hover:text-red-900">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            </tr>
        `).join('');
    }

    async viewUser(userId) {
        try {
            const data = await this.api.getUser(userId);
            const user = data.user || {};
            const stats = data.stats || {};

            const content = `
                <div class="space-y-4">
                    <div class="flex items-center space-x-4 pb-4 border-b">
                        <img src="${user.avatar_url || 'https://via.placeholder.com/100'}" class="w-20 h-20 rounded-full" alt="">
                        <div>
                            <h4 class="text-xl font-bold">${user.name}</h4>
                            <p class="text-gray-600">${user.email}</p>
                            <div class="flex space-x-2 mt-2">
                                <span class="px-2 py-1 text-xs rounded ${this.getRoleBadgeColor(user.role)}">${this.formatRole(user.role)}</span>
                                <span class="px-2 py-1 text-xs rounded ${this.getStatusBadgeColor(user.status)}">${this.formatStatus(user.status)}</span>
                            </div>
                        </div>
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <p class="text-sm text-gray-600">Phone</p>
                            <p class="font-medium">${user.phone_number || 'Not provided'}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-600">Joined</p>
                            <p class="font-medium">${this.formatDate(user.created_at)}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-600">Total Orders</p>
                            <p class="font-medium">${stats.total_orders || 0}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-600">Blog Posts</p>
                            <p class="font-medium">${stats.total_blog_posts || 0}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-600">Last Login</p>
                            <p class="font-medium">${this.formatDateTime(user.last_login_at)}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-600">Email Verified</p>
                            <p class="font-medium">${user.email_verified ? 'Yes' : 'No'}</p>
                        </div>
                    </div>

                    ${user.bio ? `
                        <div>
                            <p class="text-sm text-gray-600">Bio</p>
                            <p class="font-medium">${user.bio}</p>
                        </div>
                    ` : ''}

                    ${user.suspended_at ? `
                        <div class="bg-red-50 border border-red-200 rounded p-4">
                            <p class="text-sm font-medium text-red-900">Suspension Details</p>
                            <p class="text-sm text-red-700">Reason: ${user.suspension_reason || 'No reason provided'}</p>
                            <p class="text-sm text-red-700">Suspended on: ${this.formatDateTime(user.suspended_at)}</p>
                        </div>
                    ` : ''}
                </div>
            `;

            document.getElementById('user-details-content').innerHTML = content;
            document.getElementById('user-modal').classList.remove('hidden');
        } catch (error) {
            alert('Failed to load user details: ' + error.message);
        }
    }

    async changeRole(userId) {
        const newRole = prompt('Enter new role (user, admin, author):');
        if (!newRole) return;

        if (!['user', 'admin', 'author'].includes(newRole.toLowerCase())) {
            alert('Invalid role. Must be user, admin, or author.');
            return;
        }

        try {
            await this.api.updateUserRole(userId, newRole.toLowerCase());
            alert('Role updated successfully');
            await this.loadUsers();
        } catch (error) {
            alert('Failed to update role: ' + error.message);
        }
    }

    suspendUser(userId) {
        document.getElementById('suspend-user-id').value = userId;
        document.getElementById('suspend-modal').classList.remove('hidden');
    }

    async handleSuspend(e) {
        e.preventDefault();

        const userId = parseInt(document.getElementById('suspend-user-id').value);
        const reason = document.getElementById('suspend-reason').value;
        const duration = parseInt(document.getElementById('suspend-duration').value) || 0;

        try {
            await this.api.suspendUser(userId, reason, duration);
            alert('User suspended successfully');
            this.closeSuspendModal();
            await this.loadUsers();
        } catch (error) {
            alert('Failed to suspend user: ' + error.message);
        }
    }

    async unsuspendUser(userId) {
        if (!confirm('Are you sure you want to unsuspend this user?')) return;

        try {
            await this.api.unsuspendUser(userId);
            alert('User unsuspended successfully');
            await this.loadUsers();
        } catch (error) {
            alert('Failed to unsuspend user: ' + error.message);
        }
    }

    async deleteUser(userId) {
        if (!confirm('Are you sure you want to delete this user? This action cannot be undone.')) return;

        try {
            await this.api.deleteUser(userId);
            alert('User deleted successfully');
            await this.loadUsers();
        } catch (error) {
            alert('Failed to delete user: ' + error.message);
        }
    }

    applyFilters() {
        this.filters = {
            search: document.getElementById('search-input').value,
            role: document.getElementById('role-filter').value,
            status: document.getElementById('status-filter').value
        };
        this.currentPage = 0;
        this.loadUsers();
    }

    prevPage() {
        if (this.currentPage > 0) {
            this.currentPage--;
            this.loadUsers();
        }
    }

    nextPage() {
        if ((this.currentPage + 1) * this.pageSize < this.totalUsers) {
            this.currentPage++;
            this.loadUsers();
        }
    }

    updatePagination() {
        const showing = Math.min((this.currentPage + 1) * this.pageSize, this.totalUsers);
        document.getElementById('showing-count').textContent = showing;
        document.getElementById('total-count').textContent = this.totalUsers;

        document.getElementById('prev-page').disabled = this.currentPage === 0;
        document.getElementById('next-page').disabled = (this.currentPage + 1) * this.pageSize >= this.totalUsers;
    }

    closeModal() {
        document.getElementById('user-modal').classList.add('hidden');
    }

    closeSuspendModal() {
        document.getElementById('suspend-modal').classList.add('hidden');
        document.getElementById('suspend-form').reset();
    }

    handleLogout() {
        if (confirm('Are you sure you want to logout?')) {
            this.api.clearToken();
            window.location.href = 'index.html';
        }
    }

    // Helper methods
    formatRole(role) {
        if (!role) return 'User';
        return role.charAt(0).toUpperCase() + role.slice(1);
    }

    formatStatus(status) {
        if (!status) return 'Unknown';
        return status.charAt(0).toUpperCase() + status.slice(1);
    }

    getRoleBadgeColor(role) {
        switch (role) {
            case 'admin': return 'bg-red-100 text-red-800';
            case 'author': return 'bg-purple-100 text-purple-800';
            default: return 'bg-gray-100 text-gray-800';
        }
    }

    getStatusBadgeColor(status) {
        switch (status) {
            case 'active': return 'bg-green-100 text-green-800';
            case 'suspended': return 'bg-red-100 text-red-800';
            case 'inactive': return 'bg-gray-100 text-gray-800';
            default: return 'bg-gray-100 text-gray-800';
        }
    }

    formatDate(dateString) {
        if (!dateString) return 'Unknown';
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
    }

    formatDateTime(dateString) {
        if (!dateString) return 'Never';
        const date = new Date(dateString);
        return date.toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }
}

// Initialize the page when DOM is ready
let adminPanel;
document.addEventListener('DOMContentLoaded', () => {
    adminPanel = new AdminPanel();
});

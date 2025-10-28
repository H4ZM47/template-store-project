// Dashboard Page JavaScript
class DashboardPage {
    constructor() {
        this.api = apiClient;
        this.ordersPage = 0;
        this.blogPage = 0;
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
        await this.loadDashboardStats();
        await this.loadOrders();
    }

    setupEventListeners() {
        // Tab switching
        document.getElementById('tab-orders').addEventListener('click', () => this.switchTab('orders'));
        document.getElementById('tab-templates').addEventListener('click', () => this.switchTab('templates'));
        document.getElementById('tab-blog').addEventListener('click', () => this.switchTab('blog'));

        // Load more buttons
        document.getElementById('load-more-orders').addEventListener('click', () => this.loadMoreOrders());
        document.getElementById('load-more-blog').addEventListener('click', () => this.loadMoreBlog());

        // Logout
        document.getElementById('logout-btn').addEventListener('click', () => this.handleLogout());
    }

    async switchTab(tab) {
        // Update tab buttons
        document.querySelectorAll('.tab-button').forEach(btn => {
            btn.classList.remove('active', 'border-blue-600', 'text-blue-600');
            btn.classList.add('border-transparent', 'text-gray-600');
        });

        const activeBtn = document.getElementById(`tab-${tab}`);
        activeBtn.classList.add('active', 'border-blue-600', 'text-blue-600');
        activeBtn.classList.remove('border-transparent', 'text-gray-600');

        // Update content
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.add('hidden');
        });

        document.getElementById(`content-${tab}`).classList.remove('hidden');

        // Load data if not already loaded
        if (tab === 'templates') {
            await this.loadPurchasedTemplates();
        } else if (tab === 'blog') {
            await this.loadBlogPosts();
        }
    }

    async loadDashboardStats() {
        try {
            const data = await this.api.getDashboard();
            const stats = data.stats || {};

            document.getElementById('total-orders').textContent = stats.total_orders || 0;
            document.getElementById('total-spent').textContent = `$${(stats.total_spent || 0).toFixed(2)}`;
            document.getElementById('templates-purchased').textContent = stats.templates_purchased || 0;
            document.getElementById('blog-posts-authored').textContent = stats.blog_posts_authored || 0;
        } catch (error) {
            console.error('Failed to load dashboard stats:', error);
        }
    }

    async loadOrders() {
        try {
            const data = await this.api.getOrders(this.pageSize, this.ordersPage * this.pageSize);
            const orders = data.orders || [];
            this.displayOrders(orders, this.ordersPage === 0);

            // Show/hide load more button
            if (orders.length >= this.pageSize) {
                document.getElementById('load-more-orders').classList.remove('hidden');
            }
        } catch (error) {
            console.error('Failed to load orders:', error);
            document.getElementById('orders-container').innerHTML = '<div class="text-center text-red-500 py-8">Failed to load orders</div>';
        }
    }

    displayOrders(orders, replace = true) {
        const container = document.getElementById('orders-container');

        if (orders.length === 0 && replace) {
            container.innerHTML = '<div class="text-center text-gray-500 py-8"><i class="fas fa-inbox text-4xl mb-2"></i><p>No orders yet</p></div>';
            return;
        }

        const html = orders.map(order => `
            <div class="border border-gray-200 rounded-lg p-6 hover:shadow-md transition">
                <div class="flex items-center justify-between mb-4">
                    <div>
                        <h3 class="text-lg font-semibold text-gray-900">Order #${order.id}</h3>
                        <p class="text-sm text-gray-600">${this.formatDate(order.created_at)}</p>
                    </div>
                    <span class="px-3 py-1 rounded-full text-sm font-medium ${this.getOrderStatusColor(order.delivery_status)}">
                        ${this.formatStatus(order.delivery_status)}
                    </span>
                </div>
                ${order.Template ? `
                    <div class="flex items-center space-x-4">
                        <div class="flex-1">
                            <p class="font-medium text-gray-900">${order.Template.name || 'Template'}</p>
                            <p class="text-sm text-gray-600">Category: ${order.Template.Category?.name || 'Uncategorized'}</p>
                        </div>
                        <div class="text-right">
                            <p class="text-lg font-bold text-gray-900">$${(order.Template.price || 0).toFixed(2)}</p>
                            ${order.delivery_status === 'completed' ? `
                                <button class="mt-2 text-blue-600 hover:text-blue-800 text-sm" onclick="dashboardPage.downloadTemplate(${order.Template.id})">
                                    <i class="fas fa-download mr-1"></i>Download
                                </button>
                            ` : ''}
                        </div>
                    </div>
                ` : ''}
            </div>
        `).join('');

        if (replace) {
            container.innerHTML = html;
        } else {
            container.innerHTML += html;
        }
    }

    async loadPurchasedTemplates() {
        try {
            const data = await this.api.getPurchasedTemplates();
            const templates = data.templates || [];
            this.displayTemplates(templates);
        } catch (error) {
            console.error('Failed to load templates:', error);
            document.getElementById('templates-container').innerHTML = '<div class="col-span-full text-center text-red-500 py-8">Failed to load templates</div>';
        }
    }

    displayTemplates(templates) {
        const container = document.getElementById('templates-container');

        if (templates.length === 0) {
            container.innerHTML = '<div class="col-span-full text-center text-gray-500 py-8"><i class="fas fa-file-alt text-4xl mb-2"></i><p>No templates purchased yet</p></div>';
            return;
        }

        container.innerHTML = templates.map(item => `
            <div class="border border-gray-200 rounded-lg overflow-hidden hover:shadow-lg transition">
                <div class="p-6">
                    <div class="flex items-center justify-between mb-2">
                        <i class="fas fa-file-alt text-blue-600 text-3xl"></i>
                        <span class="text-xs text-gray-500">Purchased ${this.formatDate(item.purchased_at)}</span>
                    </div>
                    <h3 class="text-lg font-semibold text-gray-900 mb-2">${item.template.name || 'Template'}</h3>
                    <p class="text-sm text-gray-600 mb-4">${item.template.Category?.name || 'Uncategorized'}</p>
                    <button class="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" onclick="dashboardPage.downloadTemplate(${item.template.id})">
                        <i class="fas fa-download mr-2"></i>Download
                    </button>
                </div>
            </div>
        `).join('');
    }

    async loadBlogPosts() {
        try {
            const data = await this.api.getBlogPosts(this.pageSize, this.blogPage * this.pageSize);
            const posts = data.posts || [];
            this.displayBlogPosts(posts, this.blogPage === 0);

            // Show/hide load more button
            if (posts.length >= this.pageSize) {
                document.getElementById('load-more-blog').classList.remove('hidden');
            }
        } catch (error) {
            console.error('Failed to load blog posts:', error);
            document.getElementById('blog-container').innerHTML = '<div class="text-center text-red-500 py-8">Failed to load blog posts</div>';
        }
    }

    displayBlogPosts(posts, replace = true) {
        const container = document.getElementById('blog-container');

        if (posts.length === 0 && replace) {
            container.innerHTML = '<div class="text-center text-gray-500 py-8"><i class="fas fa-pen text-4xl mb-2"></i><p>No blog posts authored yet</p></div>';
            return;
        }

        const html = posts.map(post => `
            <div class="border border-gray-200 rounded-lg p-6 hover:shadow-md transition">
                <div class="flex items-center justify-between mb-2">
                    <h3 class="text-xl font-semibold text-gray-900">${post.title}</h3>
                    <span class="px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        Published
                    </span>
                </div>
                <p class="text-gray-600 mb-2">${post.excerpt || post.content?.substring(0, 150) + '...' || ''}</p>
                <div class="flex items-center justify-between text-sm text-gray-500">
                    <span><i class="fas fa-folder mr-1"></i>${post.Category?.name || 'Uncategorized'}</span>
                    <span><i class="fas fa-calendar mr-1"></i>${this.formatDate(post.created_at)}</span>
                </div>
            </div>
        `).join('');

        if (replace) {
            container.innerHTML = html;
        } else {
            container.innerHTML += html;
        }
    }

    async loadMoreOrders() {
        this.ordersPage++;
        await this.loadOrders();
    }

    async loadMoreBlog() {
        this.blogPage++;
        await this.loadBlogPosts();
    }

    async downloadTemplate(templateId) {
        try {
            const data = await this.api.getTemplateDownload(templateId);
            if (data.download_url) {
                window.open(data.download_url, '_blank');
            } else {
                alert('Download link not available');
            }
        } catch (error) {
            alert('Failed to get download link: ' + error.message);
        }
    }

    handleLogout() {
        if (confirm('Are you sure you want to logout?')) {
            this.api.clearToken();
            window.location.href = 'index.html';
        }
    }

    // Helper methods
    formatDate(dateString) {
        if (!dateString) return 'Unknown';
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
    }

    formatStatus(status) {
        if (!status) return 'Unknown';
        return status.charAt(0).toUpperCase() + status.slice(1);
    }

    getOrderStatusColor(status) {
        switch (status) {
            case 'completed': return 'bg-green-100 text-green-800';
            case 'pending': return 'bg-yellow-100 text-yellow-800';
            case 'cancelled': return 'bg-red-100 text-red-800';
            default: return 'bg-gray-100 text-gray-800';
        }
    }
}

// Initialize the page when DOM is ready
let dashboardPage;
document.addEventListener('DOMContentLoaded', () => {
    dashboardPage = new DashboardPage();
});

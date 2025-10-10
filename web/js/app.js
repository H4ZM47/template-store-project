// TemplateStore Frontend Application

// Data models
class Category {
    constructor({ id, name, created_at, updated_at }) {
        this.id = id ?? null;
        this.name = name ?? '';
        this.created_at = created_at ?? null;
        this.updated_at = updated_at ?? null;
    }
}

class TemplateModel {
    constructor(obj = {}) {
        // Normalize keys coming from API (Category vs category)
        const category = obj.category || obj.Category || null;
        this.id = obj.id ?? null;
        this.name = obj.name ?? '';
        this.file_info = obj.file_info ?? '';
        this.category_id = obj.category_id ?? (category ? category.id : null);
        this.category = category ? new Category(category) : null;
        this.price = typeof obj.price === 'number' ? obj.price : 0;
        this.preview_data = obj.preview_data ?? '';
        this.created_at = obj.created_at ?? null;
        this.updated_at = obj.updated_at ?? null;
    }
}

class BlogPostModel {
    constructor(obj = {}) {
        // Normalize keys (author/category provided by API handlers)
        this.id = obj.id ?? null;
        this.title = obj.title ?? '';
        this.content = obj.content ?? '';
        this.html_content = obj.html_content ?? '';
        this.author_id = obj.author_id ?? (obj.author ? obj.author.id : null);
        this.author = obj.author || null;
        this.category_id = obj.category_id ?? (obj.category ? obj.category.id : null);
        this.category = obj.category ? new Category(obj.category) : null;
        this.seo = obj.seo ?? '';
        this.created_at = obj.created_at ?? null;
        this.updated_at = obj.updated_at ?? null;
        this.excerpt = obj.excerpt ?? '';
    }
}

class TemplateStoreApp {
    constructor() {
        this.apiBase = 'http://localhost:8080/api/v1';
        this.templates = [];
        this.blogPosts = [];
        this.categories = [];
        this.currentTemplatePage = 0;
        this.currentBlogPage = 0;
        this.templateLimit = 8;
        this.blogLimit = 6;

        this.init();
    }

    async init() {
        this.setupEventListeners();
        await this.loadCategories();
        await this.loadTemplates();
        await this.loadBlogPosts();
        this.setupMobileMenu();
    }

    setupEventListeners() {
        // Template search and filter
        document.getElementById('search-btn').addEventListener('click', () => this.searchTemplates());
        document.getElementById('template-search').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.searchTemplates();
        });
        document.getElementById('category-filter').addEventListener('change', () => this.filterTemplatesByCategory());
        
        // Blog search
        document.getElementById('blog-search-btn').addEventListener('click', () => this.searchBlogPosts());
        document.getElementById('blog-search').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.searchBlogPosts();
        });
        
        // Load more buttons
        document.getElementById('load-more').addEventListener('click', () => this.loadMoreTemplates());
        document.getElementById('load-more-blog').addEventListener('click', () => this.loadMoreBlogPosts());
        
        // Modal close
        document.getElementById('close-modal').addEventListener('click', () => this.closeModal());
        document.getElementById('template-modal').addEventListener('click', (e) => {
            if (e.target === e.currentTarget) this.closeModal();
        });
        
        // Contact form
        document.querySelector('#contact form').addEventListener('submit', (e) => this.handleContactForm(e));
        
        // Smooth scrolling for navigation links
        document.querySelectorAll('a[href^="#"]').forEach(anchor => {
            anchor.addEventListener('click', (e) => {
                e.preventDefault();
                const target = document.querySelector(anchor.getAttribute('href'));
                if (target) {
                    target.scrollIntoView({ behavior: 'smooth' });
                }
            });
        });
    }

    setupMobileMenu() {
        const mobileMenuButton = document.getElementById('mobile-menu-button');
        const mobileMenu = document.getElementById('mobile-menu');
        
        mobileMenuButton.addEventListener('click', () => {
            mobileMenu.classList.toggle('hidden');
        });
        
        // Close mobile menu when clicking on a link
        mobileMenu.querySelectorAll('a').forEach(link => {
            link.addEventListener('click', () => {
                mobileMenu.classList.add('hidden');
            });
        });
    }

    async loadCategories() {
        try {
            const response = await fetch(`${this.apiBase}/categories`);
            const data = await response.json();
            this.categories = (data.categories || []).map(c => new Category(c));
            this.populateCategoryFilter();
        } catch (error) {
            console.error('Error loading categories:', error);
        }
    }

    populateCategoryFilter() {
        const categoryFilter = document.getElementById('category-filter');
        categoryFilter.innerHTML = '<option value="">All Categories</option>';
        
        this.categories.forEach(category => {
            const option = document.createElement('option');
            option.value = category.id;
            option.textContent = category.name;
            categoryFilter.appendChild(option);
        });
    }

    async loadTemplates(page = 0, search = '', categoryId = '') {
        try {
            let url = `${this.apiBase}/templates?limit=${this.templateLimit}&offset=${page * this.templateLimit}`;
            if (search) url += `&search=${encodeURIComponent(search)}`;
            if (categoryId) url += `&category_id=${categoryId}`;
            
            const response = await fetch(url);
            const data = await response.json();
            const newItems = (data.templates || []).map(t => new TemplateModel(t));
            
            if (page === 0) {
                this.templates = newItems;
            } else {
                this.templates = [...this.templates, ...newItems];
            }
            
            this.renderTemplates();
            this.currentTemplatePage = page;
            
            // Hide load more button if no more templates
            const loadMoreBtn = document.getElementById('load-more');
            loadMoreBtn.style.display = newItems.length < this.templateLimit ? 'none' : 'block';
            
        } catch (error) {
            console.error('Error loading templates:', error);
            this.showError('Failed to load templates');
        }
    }

    async loadBlogPosts(page = 0, search = '') {
        try {
            let url = `${this.apiBase}/blog?limit=${this.blogLimit}&offset=${page * this.blogLimit}`;
            if (search) url += `&search=${encodeURIComponent(search)}`;
            
            const response = await fetch(url);
            const data = await response.json();
            const newItems = (data.posts || []).map(p => new BlogPostModel(p));
            
            if (page === 0) {
                this.blogPosts = newItems;
            } else {
                this.blogPosts = [...this.blogPosts, ...newItems];
            }
            
            this.renderBlogPosts();
            this.currentBlogPage = page;
            
            // Hide load more button if no more posts
            const loadMoreBtn = document.getElementById('load-more-blog');
            loadMoreBtn.style.display = newItems.length < this.blogLimit ? 'none' : 'block';
            
        } catch (error) {
            console.error('Error loading blog posts:', error);
            this.showError('Failed to load blog posts');
        }
    }

    renderTemplates() {
        const templatesGrid = document.getElementById('templates-grid');
        templatesGrid.innerHTML = '';
        
        this.templates.forEach(template => {
            const templateCard = this.createTemplateCard(template);
            templatesGrid.appendChild(templateCard);
        });
        
        // Add animation class
        setTimeout(() => {
            document.querySelectorAll('.template-card').forEach((card, index) => {
                card.style.animationDelay = `${index * 0.1}s`;
                card.classList.add('animate-fade-in-up');
            });
        }, 100);
    }

    createTemplateCard(template) {
        const card = document.createElement('div');
        card.className = 'template-card bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow duration-300';

        // Use thumbnail if available, otherwise use default icon
        const thumbnailUrl = `${this.apiBase}/templates/${template.id}/thumbnail`;
        const thumbnailHTML = `
            <div class="w-full h-48 overflow-hidden bg-gray-100">
                <img src="${thumbnailUrl}"
                     alt="${template.name}"
                     class="w-full h-full object-cover"
                     onerror="this.onerror=null; this.parentElement.innerHTML='<div class=\\'w-full h-full bg-gradient-to-br from-blue-100 to-blue-200 flex items-center justify-center\\'><svg class=\\'w-16 h-16 text-blue-400\\' fill=\\'none\\' stroke=\\'currentColor\\' viewBox=\\'0 0 24 24\\'><path stroke-linecap=\\'round\\' stroke-linejoin=\\'round\\' stroke-width=\\'2\\' d=\\'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z\\'></path></svg></div>';">
            </div>
        `;

        card.innerHTML = `
            ${thumbnailHTML}
            <div class="p-6">
                <h4 class="text-lg font-semibold text-gray-900 mb-2">${template.name}</h4>
                <p class="text-gray-600 text-sm mb-4 line-clamp-2">${template.preview_data || template.file_info || 'Professional template'}</p>
                <div class="flex items-center justify-between">
                    <span class="text-2xl font-bold text-primary">$${template.price.toFixed(2)}</span>
                    <button class="bg-primary text-white px-4 py-2 rounded-lg hover:bg-secondary transition-colors"
                            onclick="app.previewTemplate(${template.id})">
                        Preview
                    </button>
                </div>
                ${template.category ? `<div class="mt-3"><span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">${template.category.name}</span></div>` : ''}
            </div>
        `;

        return card;
    }

    renderBlogPosts() {
        const blogGrid = document.getElementById('blog-grid');
        blogGrid.innerHTML = '';
        
        this.blogPosts.forEach(post => {
            const blogCard = this.createBlogCard(post);
            blogGrid.appendChild(blogCard);
        });
        
        // Add animation class
        setTimeout(() => {
            document.querySelectorAll('.blog-card').forEach((card, index) => {
                card.style.animationDelay = `${index * 0.1}s`;
                card.classList.add('animate-fade-in-up');
            });
        }, 100);
    }

    createBlogCard(post) {
        const card = document.createElement('div');
        card.className = 'blog-card bg-white rounded-lg shadow-md overflow-hidden';
        
        const excerpt = post.excerpt || 'No preview available';
        const authorName = post.author ? (post.author.name || 'Unknown Author') : 'Unknown Author';
        const categoryName = post.category ? post.category.name : 'Uncategorized';
        
        card.innerHTML = `
            <div class="p-6">
                <div class="flex items-center mb-3">
                    <span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">${categoryName}</span>
                </div>
                <h4 class="text-xl font-semibold text-gray-900 mb-3">${post.title}</h4>
                <p class="text-gray-600 mb-4 line-clamp-3">${excerpt}</p>
                <div class="flex items-center justify-between text-sm text-gray-500">
                    <span>By ${authorName}</span>
                    <span>${post.created_at ? new Date(post.created_at).toLocaleDateString() : ''}</span>
                </div>
                <button class="mt-4 w-full bg-primary text-white px-4 py-2 rounded-lg hover:bg-secondary transition-colors"
                        onclick="app.readBlogPost(${post.id})">
                    Read More
                </button>
            </div>
        `;
        
        return card;
    }

    async previewTemplate(templateId) {
        try {
            const response = await fetch(`${this.apiBase}/templates/${templateId}`);
            const data = await response.json();
            const template = new TemplateModel(data.template || {});
            
            this.showTemplateModal(template);
        } catch (error) {
            console.error('Error loading template:', error);
            this.showError('Failed to load template preview');
        }
    }

    showTemplateModal(template) {
        const modal = document.getElementById('template-modal');
        const title = document.getElementById('modal-title');
        const content = document.getElementById('modal-content');

        title.textContent = template.name;
        content.innerHTML = `
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                    <div class="bg-gradient-to-br from-blue-100 to-blue-200 rounded-lg p-8 text-center">
                        <svg class="w-24 h-24 text-blue-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                        </svg>
                        <p class="text-blue-600 font-medium">Template Preview</p>
                    </div>
                </div>
                <div class="space-y-4">
                    <div>
                        <h4 class="font-semibold text-gray-900">Description</h4>
                        <p class="text-gray-600">${template.preview_data || template.file_info || 'Professional template for your needs'}</p>
                    </div>
                    <div>
                        <h4 class="font-semibold text-gray-900">Category</h4>
                        <p class="text-gray-600">${template.category ? template.category.name : 'Uncategorized'}</p>
                    </div>
                    <div>
                        <h4 class="font-semibold text-gray-900">Price</h4>
                        <p class="text-3xl font-bold text-primary">$${template.price}</p>
                    </div>
                    <div class="pt-4 space-y-3">
                        <button onclick="app.viewTemplate(${template.id})" class="w-full bg-green-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-green-700 transition-colors flex items-center justify-center gap-2">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                            </svg>
                            View Template
                        </button>
                        <button onclick="app.customizeAndDownload(${template.id}, '${template.name}')" class="w-full bg-purple-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-purple-700 transition-colors flex items-center justify-center gap-2">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                            </svg>
                            Customize & Download
                        </button>
                        ${this.getStripeBuyButton(template.name)}
                    </div>
                </div>
            </div>
        `;

        modal.classList.remove('hidden');
        modal.classList.add('modal-enter');

        // Style Stripe buy button after it loads
        setTimeout(() => {
            this.styleStripeBuyButton();
        }, 500);
    }

    styleStripeBuyButton() {
        const stripeButtons = document.querySelectorAll('stripe-buy-button');
        stripeButtons.forEach(button => {
            if (button.shadowRoot) {
                // Inject custom styles into the shadow DOM
                const style = document.createElement('style');
                style.textContent = `
                    * {
                        box-sizing: border-box !important;
                    }
                    button {
                        background-color: #3B82F6 !important;
                        color: white !important;
                        padding: 0.75rem 1.5rem !important;
                        border-radius: 0.5rem !important;
                        font-weight: 600 !important;
                        font-size: 1rem !important;
                        width: 100% !important;
                        min-width: 100% !important;
                        max-width: 100% !important;
                        display: block !important;
                        border: none !important;
                        cursor: pointer !important;
                        transition: background-color 0.2s ease !important;
                        height: 48px !important;
                        margin: 0 !important;
                    }
                    button:hover {
                        background-color: #1E40AF !important;
                    }
                `;
                button.shadowRoot.appendChild(style);
            }
        });
    }

    closeModal() {
        const modal = document.getElementById('template-modal');
        modal.classList.add('modal-exit');

        setTimeout(() => {
            modal.classList.add('hidden');
            modal.classList.remove('modal-enter', 'modal-exit');
        }, 300);
    }

    getStripeBuyButton(templateName) {
        const buyButtons = {
            'Vulnerability Management Standard': {
                buyButtonId: 'buy_btn_1SGmXsAH3GeIobtC0xuoNOCX',
                publishableKey: 'pk_test_51JsbhNAH3GeIobtCuxkeFNfWeiUZavWH6d4OpQfDxGZS5qeKLnRzmbrMbHkxvCB0AA6CMFAi4PaWDOnFPNEAnpdD00pRBBFjSD'
            },
            'Data Classification Standard': {
                buyButtonId: 'buy_btn_1SGmfhAH3GeIobtCYpvCC5FU',
                publishableKey: 'pk_test_51JsbhNAH3GeIobtCuxkeFNfWeiUZavWH6d4OpQfDxGZS5qeKLnRzmbrMbHkxvCB0AA6CMFAi4PaWDOnFPNEAnpdD00pRBBFjSD'
            }
        };

        const buttonConfig = buyButtons[templateName];

        if (buttonConfig) {
            return `<div class="stripe-button-wrapper">
                      <stripe-buy-button
                        buy-button-id="${buttonConfig.buyButtonId}"
                        publishable-key="${buttonConfig.publishableKey}"
                      >
                      </stripe-buy-button>
                    </div>`;
        }

        // Default buy button for templates without specific configuration
        return `<button class="w-full bg-primary text-white px-6 py-3 rounded-lg font-semibold hover:bg-secondary transition-colors">
                    Purchase Template
                </button>`;
    }

    async readBlogPost(postId) {
        try {
            const response = await fetch(`${this.apiBase}/blog/${postId}`);
            const data = await response.json();
            const post = new BlogPostModel(data.post || {});
            
            this.showBlogPostModal(post);
        } catch (error) {
            console.error('Error loading blog post:', error);
            this.showError('Failed to load blog post');
        }
    }

    showBlogPostModal(post) {
        const modal = document.getElementById('template-modal');
        const title = document.getElementById('modal-title');
        const content = document.getElementById('modal-content');
        
        title.textContent = post.title;
        content.innerHTML = `
            <div class="space-y-6">
                <div class="flex items-center space-x-4 text-sm text-gray-500">
                    <span>By ${post.author ? (post.author.name || 'Unknown Author') : 'Unknown Author'}</span>
                    <span>•</span>
                    <span>${post.created_at ? new Date(post.created_at).toLocaleDateString() : ''}</span>
                    ${post.category ? `<span>•</span><span class="bg-blue-100 text-blue-800 px-2 py-1 rounded-full">${post.category.name}</span>` : ''}
                </div>
                <div class="prose max-w-none">
                    ${post.html_content || post.content}
                </div>
            </div>
        `;
        
        modal.classList.remove('hidden');
        modal.classList.add('modal-enter');
    }

    searchTemplates() {
        const searchTerm = document.getElementById('template-search').value.trim();
        const categoryId = document.getElementById('category-filter').value;
        
        this.currentTemplatePage = 0;
        this.loadTemplates(0, searchTerm, categoryId);
    }

    filterTemplatesByCategory() {
        const categoryId = document.getElementById('category-filter').value;
        const searchTerm = document.getElementById('template-search').value.trim();
        
        this.currentTemplatePage = 0;
        this.loadTemplates(0, searchTerm, categoryId);
    }

    searchBlogPosts() {
        const searchTerm = document.getElementById('blog-search').value.trim();
        
        this.currentBlogPage = 0;
        this.loadBlogPosts(0, searchTerm);
    }

    loadMoreTemplates() {
        const searchTerm = document.getElementById('template-search').value.trim();
        const categoryId = document.getElementById('category-filter').value;
        
        this.loadTemplates(this.currentTemplatePage + 1, searchTerm, categoryId);
    }

    loadMoreBlogPosts() {
        const searchTerm = document.getElementById('blog-search').value.trim();
        
        this.loadBlogPosts(this.currentBlogPage + 1, searchTerm);
    }

    async handleContactForm(event) {
        event.preventDefault();
        
        const formData = new FormData(event.target);
        const name = formData.get('name');
        const email = formData.get('email');
        const message = formData.get('message');
        
        // Here you would typically send the form data to your backend
        // For now, we'll just show a success message
        this.showSuccess('Thank you for your message! We\'ll get back to you soon.');
        event.target.reset();
    }

    showError(message) {
        this.showNotification(message, 'error');
    }

    showSuccess(message) {
        this.showNotification(message, 'success');
    }

    showNotification(message, type = 'info') {
        // Create notification element
        const notification = document.createElement('div');
        notification.className = `fixed top-4 right-4 z-50 px-6 py-3 rounded-lg shadow-lg transition-all duration-300 transform translate-x-full`;

        const bgColor = type === 'error' ? 'bg-red-500' : type === 'success' ? 'bg-green-500' : 'bg-blue-500';
        notification.className += ` ${bgColor} text-white`;

        notification.textContent = message;

        document.body.appendChild(notification);

        // Animate in
        setTimeout(() => {
            notification.classList.remove('translate-x-full');
        }, 100);

        // Remove after 5 seconds
        setTimeout(() => {
            notification.classList.add('translate-x-full');
            setTimeout(() => {
                document.body.removeChild(notification);
            }, 300);
        }, 5000);
    }

    viewTemplate(templateId) {
        // Open template in new window/tab
        window.open(`${this.apiBase}/templates/${templateId}/view`, '_blank');
    }

    async customizeAndDownload(templateId, templateName) {
        try {
            // Fetch template variables
            const response = await fetch(`${this.apiBase}/templates/${templateId}/variables`);
            const data = await response.json();
            const variables = data.variables || [];

            if (variables.length === 0) {
                // No variables, download directly
                this.downloadTemplatePDF(templateId, templateName);
                return;
            }

            // Show customization modal
            this.showCustomizationModal(templateId, templateName, variables);
        } catch (error) {
            console.error('Error loading template variables:', error);
            this.showError('Failed to load template customization options');
        }
    }

    showCustomizationModal(templateId, templateName, variables) {
        const modal = document.getElementById('template-modal');
        const title = document.getElementById('modal-title');
        const content = document.getElementById('modal-content');

        title.textContent = `Customize: ${templateName}`;

        // Build form HTML
        const formFields = variables.map(variable => {
            const required = variable.required ? 'required' : '';
            const requiredLabel = variable.required ? '<span class="text-red-500">*</span>' : '';

            if (variable.type === 'date') {
                return `
                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-1">
                            ${variable.label} ${requiredLabel}
                        </label>
                        <input type="date" name="${variable.name}" value="${variable.default || ''}" ${required}
                               class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary">
                    </div>
                `;
            } else if (variable.type === 'select' && variable.options) {
                const options = variable.options.map(opt =>
                    `<option value="${opt}" ${opt === variable.default ? 'selected' : ''}>${opt}</option>`
                ).join('');
                return `
                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-1">
                            ${variable.label} ${requiredLabel}
                        </label>
                        <select name="${variable.name}" ${required}
                                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary">
                            ${options}
                        </select>
                    </div>
                `;
            } else {
                // Default to text input
                return `
                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-1">
                            ${variable.label} ${requiredLabel}
                        </label>
                        <input type="text" name="${variable.name}" value="${variable.default || ''}" ${required}
                               class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                               placeholder="${variable.label}">
                    </div>
                `;
            }
        }).join('');

        content.innerHTML = `
            <form id="customize-form" class="space-y-4">
                <p class="text-gray-600 mb-4">Fill in the following information to customize your template:</p>
                ${formFields}
                <div class="flex gap-3 pt-4">
                    <button type="submit" class="flex-1 bg-purple-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-purple-700 transition-colors">
                        Generate & Download PDF
                    </button>
                    <button type="button" onclick="app.closeModal()" class="px-6 py-3 border border-gray-300 rounded-lg font-semibold hover:bg-gray-50 transition-colors">
                        Cancel
                    </button>
                </div>
            </form>
        `;

        // Add form submit handler
        document.getElementById('customize-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            await this.generateCustomPDF(templateId, templateName, e.target);
        });

        modal.classList.remove('hidden');
        modal.classList.add('modal-enter');
    }

    async generateCustomPDF(templateId, templateName, form) {
        try {
            // Collect form data
            const formData = new FormData(form);
            const variables = {};
            for (const [key, value] of formData.entries()) {
                variables[key] = value;
            }

            // Call generate endpoint
            const response = await fetch(`${this.apiBase}/templates/${templateId}/generate`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(variables)
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to generate PDF');
            }

            // Download the PDF
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = `${templateName.toLowerCase().replace(/\s+/g, '_')}_custom.pdf`;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            window.URL.revokeObjectURL(url);

            this.closeModal();
            this.showSuccess('Customized PDF downloaded successfully!');
        } catch (error) {
            console.error('Error generating custom PDF:', error);
            this.showError(error.message || 'Failed to generate custom PDF');
        }
    }

    downloadTemplatePDF(templateId, templateName) {
        // Download PDF directly (for templates without variables)
        const downloadUrl = `${this.apiBase}/templates/${templateId}/download`;

        // Create a temporary link and trigger download
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = `${templateName.toLowerCase().replace(/\s+/g, '_')}.pdf`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

        this.showSuccess('PDF download started...');
    }
}

// Initialize the application when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.app = new TemplateStoreApp();
});

// Add some utility functions to the global scope
window.utils = {
    // Format price
    formatPrice: (price) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(price);
    },
    
    // Format date
    formatDate: (dateString) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    },
    
    // Debounce function for search
    debounce: (func, wait) => {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }
};

// Settings Page JavaScript
class SettingsPage {
    constructor() {
        this.api = apiClient;
        this.currentPreferences = null;
        this.init();
    }

    async init() {
        // Check if user is authenticated
        if (!this.api.token) {
            window.location.href = 'index.html';
            return;
        }

        this.setupEventListeners();
        await this.loadPreferences();
    }

    setupEventListeners() {
        // Save preferences
        document.getElementById('save-btn').addEventListener('click', () => this.handleSavePreferences());

        // Reset preferences
        document.getElementById('reset-btn').addEventListener('click', () => this.loadPreferences());

        // Logout
        document.getElementById('logout-btn').addEventListener('click', () => this.handleLogout());

        // Deactivate account
        document.getElementById('deactivate-btn').addEventListener('click', () => this.showDeactivateModal());
        document.getElementById('cancel-deactivate').addEventListener('click', () => this.hideDeactivateModal());
        document.getElementById('deactivate-form').addEventListener('submit', (e) => this.handleDeactivate(e));

        // Delete account
        document.getElementById('delete-btn').addEventListener('click', () => this.showDeleteModal());
        document.getElementById('cancel-delete').addEventListener('click', () => this.hideDeleteModal());
        document.getElementById('delete-form').addEventListener('submit', (e) => this.handleDelete(e));
    }

    async loadPreferences() {
        try {
            const data = await this.api.getPreferences();
            this.currentPreferences = data;
            this.populateForm(data);
        } catch (error) {
            this.showMessage('Failed to load preferences: ' + error.message, 'error');
        }
    }

    populateForm(prefs) {
        // Email preferences
        document.getElementById('marketing_emails').checked = prefs.marketing_emails ?? true;
        document.getElementById('order_notifications').checked = prefs.order_notifications ?? true;
        document.getElementById('blog_notifications').checked = prefs.blog_notifications ?? true;

        // Display preferences
        document.getElementById('language').value = prefs.language || 'en';
        document.getElementById('timezone').value = prefs.timezone || 'UTC';
        document.getElementById('theme').value = prefs.theme || 'light';

        // Privacy settings
        document.getElementById('profile_visibility').value = prefs.profile_visibility || 'public';
        document.getElementById('show_email').checked = prefs.show_email ?? false;
        document.getElementById('show_purchase_history').checked = prefs.show_purchase_history ?? false;
    }

    async handleSavePreferences() {
        const saveBtn = document.getElementById('save-btn');
        const saveText = document.getElementById('save-text');
        const saveSpinner = document.getElementById('save-spinner');

        // Show loading state
        saveBtn.disabled = true;
        saveText.classList.add('hidden');
        saveSpinner.classList.remove('hidden');

        try {
            const preferences = {
                marketing_emails: document.getElementById('marketing_emails').checked,
                order_notifications: document.getElementById('order_notifications').checked,
                blog_notifications: document.getElementById('blog_notifications').checked,
                language: document.getElementById('language').value,
                timezone: document.getElementById('timezone').value,
                theme: document.getElementById('theme').value,
                profile_visibility: document.getElementById('profile_visibility').value,
                show_email: document.getElementById('show_email').checked,
                show_purchase_history: document.getElementById('show_purchase_history').checked
            };

            await this.api.updatePreferences(preferences);
            this.currentPreferences = preferences;
            this.showMessage('Preferences saved successfully!', 'success');
        } catch (error) {
            this.showMessage('Failed to save preferences: ' + error.message, 'error');
        } finally {
            saveBtn.disabled = false;
            saveText.classList.remove('hidden');
            saveSpinner.classList.add('hidden');
        }
    }

    showDeactivateModal() {
        document.getElementById('deactivate-modal').classList.remove('hidden');
    }

    hideDeactivateModal() {
        document.getElementById('deactivate-modal').classList.add('hidden');
        document.getElementById('deactivate-form').reset();
    }

    async handleDeactivate(e) {
        e.preventDefault();

        const password = document.getElementById('deactivate-password').value;
        const reason = document.getElementById('deactivate-reason').value;

        try {
            await this.api.deactivateAccount(password, reason);
            alert('Account deactivated successfully. You will be logged out.');
            this.api.clearToken();
            window.location.href = 'index.html';
        } catch (error) {
            this.hideDeactivateModal();
            this.showMessage('Failed to deactivate account: ' + error.message, 'error');
        }
    }

    showDeleteModal() {
        document.getElementById('delete-modal').classList.remove('hidden');
    }

    hideDeleteModal() {
        document.getElementById('delete-modal').classList.add('hidden');
        document.getElementById('delete-form').reset();
    }

    async handleDelete(e) {
        e.preventDefault();

        const password = document.getElementById('delete-password').value;
        const confirmation = document.getElementById('delete-confirmation').value;

        if (confirmation !== 'DELETE MY ACCOUNT') {
            alert('Please type "DELETE MY ACCOUNT" exactly to confirm.');
            return;
        }

        if (!confirm('This is your last chance. Are you absolutely sure you want to delete your account?')) {
            return;
        }

        try {
            await this.api.deleteAccount(password, confirmation);
            alert('Account deleted successfully. We\'re sorry to see you go.');
            this.api.clearToken();
            window.location.href = 'index.html';
        } catch (error) {
            this.hideDeleteModal();
            this.showMessage('Failed to delete account: ' + error.message, 'error');
        }
    }

    handleLogout() {
        if (confirm('Are you sure you want to logout?')) {
            this.api.clearToken();
            window.location.href = 'index.html';
        }
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
document.addEventListener('DOMContentLoaded', () => {
    new SettingsPage();
});

// Profile Page JavaScript
class ProfilePage {
    constructor() {
        this.api = apiClient;
        this.currentProfile = null;
        this.init();
    }

    async init() {
        // Check if user is authenticated
        if (!this.api.token) {
            window.location.href = 'index.html';
            return;
        }

        this.setupEventListeners();
        await this.loadProfile();
    }

    setupEventListeners() {
        // Profile form submission
        document.getElementById('profile-form').addEventListener('submit', (e) => this.handleSaveProfile(e));

        // Cancel button
        document.getElementById('cancel-btn').addEventListener('click', () => this.loadProfile());

        // Avatar upload
        document.getElementById('avatar-upload').addEventListener('change', (e) => this.handleAvatarUpload(e));

        // Remove avatar
        document.getElementById('remove-avatar-btn').addEventListener('click', () => this.handleRemoveAvatar());

        // Logout
        document.getElementById('logout-btn').addEventListener('click', () => this.handleLogout());
    }

    async loadProfile() {
        try {
            const data = await this.api.getProfile();
            this.currentProfile = data;
            this.populateForm(data);
            this.showMessage('Profile loaded successfully', 'success');
        } catch (error) {
            this.showMessage('Failed to load profile: ' + error.message, 'error');
        }
    }

    populateForm(profile) {
        // Basic info
        document.getElementById('name').value = profile.name || '';
        document.getElementById('email').value = profile.email || '';
        document.getElementById('phone_number').value = profile.phone_number || '';
        document.getElementById('bio').value = profile.bio || '';

        // Address
        document.getElementById('address_line1').value = profile.address_line1 || '';
        document.getElementById('address_line2').value = profile.address_line2 || '';
        document.getElementById('city').value = profile.city || '';
        document.getElementById('state').value = profile.state || '';
        document.getElementById('country').value = profile.country || '';
        document.getElementById('postal_code').value = profile.postal_code || '';

        // Avatar
        if (profile.avatar_url) {
            document.getElementById('avatar-preview').src = profile.avatar_url;
        }

        // Display info
        document.getElementById('user-name').textContent = profile.name || 'No Name';
        document.getElementById('user-email').textContent = profile.email || '';

        // Account info
        document.getElementById('account-status').textContent = this.formatStatus(profile.status);
        document.getElementById('account-status').className = `font-medium ${this.getStatusColor(profile.status)}`;
        document.getElementById('email-verified').textContent = profile.email_verified ? 'Yes' : 'No';
        document.getElementById('email-verified').className = `font-medium ${profile.email_verified ? 'text-green-600' : 'text-red-600'}`;
        document.getElementById('user-role').textContent = this.formatRole(profile.role);
        document.getElementById('member-since').textContent = this.formatDate(profile.created_at);
    }

    async handleSaveProfile(e) {
        e.preventDefault();

        const saveBtn = document.getElementById('save-btn');
        const saveText = document.getElementById('save-text');
        const saveSpinner = document.getElementById('save-spinner');

        // Show loading state
        saveBtn.disabled = true;
        saveText.classList.add('hidden');
        saveSpinner.classList.remove('hidden');

        try {
            const updates = {
                name: document.getElementById('name').value,
                phone_number: document.getElementById('phone_number').value,
                bio: document.getElementById('bio').value,
                address_line1: document.getElementById('address_line1').value,
                address_line2: document.getElementById('address_line2').value,
                city: document.getElementById('city').value,
                state: document.getElementById('state').value,
                country: document.getElementById('country').value,
                postal_code: document.getElementById('postal_code').value
            };

            const data = await this.api.updateProfile(updates);
            this.currentProfile = data;
            this.showMessage('Profile updated successfully!', 'success');
        } catch (error) {
            this.showMessage('Failed to update profile: ' + error.message, 'error');
        } finally {
            saveBtn.disabled = false;
            saveText.classList.remove('hidden');
            saveSpinner.classList.add('hidden');
        }
    }

    async handleAvatarUpload(e) {
        const file = e.target.files[0];
        if (!file) return;

        // Validate file size (5MB max)
        if (file.size > 5 * 1024 * 1024) {
            this.showMessage('Avatar file is too large. Maximum size is 5MB.', 'error');
            return;
        }

        // Validate file type
        if (!file.type.startsWith('image/')) {
            this.showMessage('Please select a valid image file.', 'error');
            return;
        }

        try {
            // Show preview immediately
            const reader = new FileReader();
            reader.onload = (event) => {
                document.getElementById('avatar-preview').src = event.target.result;
            };
            reader.readAsDataURL(file);

            // Upload to server
            const data = await this.api.uploadAvatar(file);
            this.showMessage('Avatar uploaded successfully!', 'success');

            // Update preview with server URL
            if (data.avatar_url) {
                document.getElementById('avatar-preview').src = data.avatar_url;
            }
        } catch (error) {
            this.showMessage('Failed to upload avatar: ' + error.message, 'error');
            // Reload profile to reset avatar
            await this.loadProfile();
        }
    }

    async handleRemoveAvatar() {
        if (!confirm('Are you sure you want to remove your avatar?')) {
            return;
        }

        try {
            await this.api.deleteAvatar();
            document.getElementById('avatar-preview').src = 'https://via.placeholder.com/150';
            this.showMessage('Avatar removed successfully!', 'success');
        } catch (error) {
            this.showMessage('Failed to remove avatar: ' + error.message, 'error');
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

    formatStatus(status) {
        if (!status) return 'Unknown';
        return status.charAt(0).toUpperCase() + status.slice(1);
    }

    getStatusColor(status) {
        switch (status) {
            case 'active': return 'text-green-600';
            case 'suspended': return 'text-red-600';
            case 'inactive': return 'text-gray-600';
            default: return 'text-gray-600';
        }
    }

    formatRole(role) {
        if (!role) return 'User';
        return role.charAt(0).toUpperCase() + role.slice(1);
    }

    formatDate(dateString) {
        if (!dateString) return 'Unknown';
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
    }
}

// Initialize the page when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new ProfilePage();
});

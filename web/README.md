# TemplateStore Frontend

This is the frontend application for the TemplateStore & Blog Platform, built with vanilla HTML, CSS, and JavaScript using Tailwind CSS for styling.

## Features

### 🎨 **Modern Design**
- Responsive design that works on all devices
- Beautiful gradient hero section
- Smooth animations and transitions
- Professional color scheme

### 🔍 **Template Store**
- Template catalog with grid layout
- Search functionality
- Category filtering
- Template preview modal
- Responsive card design

### 📝 **Blog System**
- Blog post listing with excerpts
- Search functionality
- Full blog post reader
- Author and category information
- Markdown rendering support

### 📱 **Mobile Responsive**
- Mobile-first design approach
- Collapsible navigation menu
- Touch-friendly interactions
- Optimized for all screen sizes

### ⚡ **Performance**
- Vanilla JavaScript for fast loading
- Optimized CSS animations
- Efficient DOM manipulation
- Minimal dependencies

## File Structure

```
web/
├── index.html          # Main HTML template
├── css/
│   └── styles.css      # Custom CSS styles
├── js/
│   └── app.js         # Main JavaScript application
└── README.md           # This file
```

## Getting Started

### Prerequisites
- Go backend server running on port 8080
- Modern web browser
- Local development environment

### Running the Frontend

1. **Start the Go backend server:**
   ```bash
   go run cmd/server/main.go
   ```

2. **Start the frontend server:**
   ```bash
   go run cmd/web/main.go
   ```

3. **Open your browser:**
   Navigate to `http://localhost:3000`

### Alternative: Direct File Access
You can also open `web/index.html` directly in your browser, but you'll need to:
- Ensure CORS is properly configured
- Have the backend running on port 8080

## API Integration

The frontend communicates with the Go backend API at `http://localhost:8080/api/v1`:

### Templates
- `GET /api/v1/templates` - List templates
- `GET /api/v1/templates/:id` - Get template details
- `GET /api/v1/templates?search=query` - Search templates
- `GET /api/v1/templates?category_id=1` - Filter by category

### Blog Posts
- `GET /api/v1/blog` - List blog posts
- `GET /api/v1/blog/:id` - Get blog post details
- `GET /api/v1/blog?search=query` - Search blog posts

### Categories
- `GET /api/v1/categories` - List categories

## Customization

### Colors
The color scheme is defined in the Tailwind config within `index.html`:

```javascript
tailwind.config = {
    theme: {
        extend: {
            colors: {
                primary: '#3B82F6',    // Blue
                secondary: '#1E40AF',  // Dark Blue
                accent: '#F59E0B'      // Orange
            }
        }
    }
}
```

### Styling
Custom CSS is in `css/styles.css`:
- Smooth scrolling
- Custom animations
- Hover effects
- Responsive adjustments
- Accessibility features

### JavaScript
The main application logic is in `js/app.js`:
- Template management
- Blog post handling
- Search and filtering
- Modal functionality
- API communication

## Browser Support

- **Modern Browsers:** Chrome 80+, Firefox 75+, Safari 13+, Edge 80+
- **Mobile:** iOS Safari 13+, Chrome Mobile 80+
- **Features:** ES6+, CSS Grid, Flexbox, CSS Animations

## Development

### Adding New Features
1. **HTML:** Add new sections to `index.html`
2. **CSS:** Add styles to `css/styles.css`
3. **JavaScript:** Add functionality to `js/app.js`

### Testing
- Test on multiple devices and screen sizes
- Verify API integration
- Check accessibility features
- Test performance

### Building for Production
1. Minify CSS and JavaScript
2. Optimize images
3. Enable compression
4. Set up CDN for static assets

## Troubleshooting

### Common Issues

**Frontend not loading templates/blog posts:**
- Ensure backend server is running on port 8080
- Check browser console for API errors
- Verify CORS configuration

**Styles not applying:**
- Check if Tailwind CSS CDN is accessible
- Verify `css/styles.css` is loading
- Clear browser cache

**Mobile menu not working:**
- Check JavaScript console for errors
- Verify event listeners are properly attached
- Test on different mobile devices

### Debug Mode
Open browser developer tools to:
- View console logs
- Inspect network requests
- Debug JavaScript errors
- Test responsive design

## Performance Tips

1. **Optimize Images:** Use appropriate formats and sizes
2. **Minimize HTTP Requests:** Combine CSS/JS files
3. **Enable Caching:** Set proper cache headers
4. **Lazy Loading:** Implement for images and content
5. **Code Splitting:** Separate concerns in JavaScript

## Accessibility

- **Keyboard Navigation:** All interactive elements are keyboard accessible
- **Screen Readers:** Proper ARIA labels and semantic HTML
- **Color Contrast:** Meets WCAG AA standards
- **Focus Management:** Clear focus indicators
- **Reduced Motion:** Respects user preferences

## Future Enhancements

- [ ] Dark mode toggle
- [ ] Advanced filtering options
- [ ] User authentication
- [ ] Shopping cart functionality
- [ ] Template favorites
- [ ] Social sharing
- [ ] Progressive Web App (PWA) features
- [ ] Internationalization (i18n)

## Contributing

1. Follow the existing code style
2. Test changes on multiple devices
3. Ensure accessibility compliance
4. Update documentation as needed
5. Test API integration thoroughly

## License

This project is part of the TemplateStore & Blog Platform.

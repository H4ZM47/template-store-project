import { Router } from 'express';
import authRoutes from './auth.routes';
import profileRoutes from './profile.routes';
import templateRoutes from './template.routes';
import blogRoutes from './blog.routes';
import categoryRoutes from './category.routes';
import paymentRoutes from './payment.routes';
import adminRoutes from './admin.routes';

const router = Router();

// Mount routes
router.use('/auth', authRoutes);
router.use('/profile', profileRoutes);
router.use('/templates', templateRoutes);
router.use('/blog', blogRoutes);
router.use('/categories', categoryRoutes);
router.use('/payment', paymentRoutes);
router.use('/admin', adminRoutes);

// Health check
router.get('/health', (req, res) => {
  res.json({
    status: 'ok',
    timestamp: new Date().toISOString(),
    service: 'template-store-api',
  });
});

export default router;

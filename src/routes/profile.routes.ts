import { Router, Response } from 'express';
import { UserService, OrderService } from '../services';
import { authMiddleware, AuthRequest } from '../middleware';
import { logger } from '../utils/logger';

const router = Router();
const userService = new UserService();
const orderService = new OrderService();

// Get profile (authenticated)
router.get('/', authMiddleware, async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const userId = req.userId!;

    const user = await userService.getUserById(userId);

    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    res.json({ user });
  } catch (error: any) {
    logger.error('Get profile error:', error);
    res.status(500).json({ error: 'Failed to get profile' });
  }
});

// Update profile (authenticated)
router.put('/', authMiddleware, async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const userId = req.userId!;
    const updates = req.body;

    // Remove fields that shouldn't be updated via this endpoint
    delete updates.id;
    delete updates.email;
    delete updates.cognitoSubject;
    delete updates.role;
    delete updates.password;

    const user = await userService.updateUser(userId, updates);

    res.json({
      message: 'Profile updated successfully',
      user,
    });
  } catch (error: any) {
    logger.error('Update profile error:', error);
    res.status(500).json({ error: 'Failed to update profile' });
  }
});

// Get user orders (authenticated)
router.get('/orders', authMiddleware, async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const userId = req.userId!;
    const limit = parseInt(req.query.limit as string) || 50;

    const orders = await orderService.getOrdersByUser(userId, limit);

    res.json({
      orders,
      count: orders.length,
    });
  } catch (error: any) {
    logger.error('Get user orders error:', error);
    res.status(500).json({ error: 'Failed to get orders' });
  }
});

// Get purchased templates (authenticated)
router.get('/purchased-templates', authMiddleware, async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const userId = req.userId!;

    const orders = await orderService.getUserPurchasedTemplates(userId);

    res.json({
      templates: orders.map(order => order.template),
      count: orders.length,
    });
  } catch (error: any) {
    logger.error('Get purchased templates error:', error);
    res.status(500).json({ error: 'Failed to get purchased templates' });
  }
});

// Deactivate account (authenticated)
router.post('/deactivate', authMiddleware, async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const userId = req.userId!;

    await userService.updateUser(userId, { status: 'deactivated' });

    res.json({ message: 'Account deactivated successfully' });
  } catch (error: any) {
    logger.error('Deactivate account error:', error);
    res.status(500).json({ error: 'Failed to deactivate account' });
  }
});

export default router;

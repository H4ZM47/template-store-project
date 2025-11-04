import { Router, Response } from 'express';
import { UserService } from '../services';
import { authMiddleware, AuthRequest, requireAdmin } from '../middleware';
import { logger } from '../utils/logger';

const router = Router();
const userService = new UserService();

// All admin routes require authentication and admin role
router.use(authMiddleware);
router.use(requireAdmin);

// List all users
router.get('/users', async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const limit = parseInt(req.query.limit as string) || 50;
    const offset = parseInt(req.query.offset as string) || 0;

    const users = await userService.listUsers(limit, offset);

    res.json({
      users,
      count: users.length,
      limit,
      offset,
    });
  } catch (error: any) {
    logger.error('Admin list users error:', error);
    res.status(500).json({ error: 'Failed to list users' });
  }
});

// Get user by ID
router.get('/users/:id', async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const { id } = req.params;

    const user = await userService.getUserById(id);

    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    res.json({ user });
  } catch (error: any) {
    logger.error('Admin get user error:', error);
    res.status(500).json({ error: 'Failed to get user' });
  }
});

// Update user role
router.put('/users/:id/role', async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const { id } = req.params;
    const { role } = req.body;

    if (!role || !['user', 'admin', 'author'].includes(role)) {
      res.status(400).json({ error: 'Invalid role' });
      return;
    }

    const user = await userService.updateUser(id, { role });

    res.json({
      message: 'User role updated successfully',
      user,
    });
  } catch (error: any) {
    logger.error('Admin update user role error:', error);
    res.status(500).json({ error: 'Failed to update user role' });
  }
});

// Suspend user
router.post('/users/:id/suspend', async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const { id } = req.params;

    const user = await userService.updateUser(id, { status: 'suspended' });

    res.json({
      message: 'User suspended successfully',
      user,
    });
  } catch (error: any) {
    logger.error('Admin suspend user error:', error);
    res.status(500).json({ error: 'Failed to suspend user' });
  }
});

// Unsuspend user
router.post('/users/:id/unsuspend', async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const { id } = req.params;

    const user = await userService.updateUser(id, { status: 'active' });

    res.json({
      message: 'User unsuspended successfully',
      user,
    });
  } catch (error: any) {
    logger.error('Admin unsuspend user error:', error);
    res.status(500).json({ error: 'Failed to unsuspend user' });
  }
});

// Delete user
router.delete('/users/:id', async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const { id } = req.params;

    await userService.deleteUser(id);

    res.json({ message: 'User deleted successfully' });
  } catch (error: any) {
    logger.error('Admin delete user error:', error);
    res.status(500).json({ error: 'Failed to delete user' });
  }
});

export default router;

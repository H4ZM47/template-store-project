import { Router, Request, Response } from 'express';
import { AuthService, UserService, EmailService } from '../services';
import { authMiddleware, AuthRequest } from '../middleware';
import { logger } from '../utils/logger';
import bcrypt from 'bcrypt';

const router = Router();
const authService = new AuthService();
const userService = new UserService();
const emailService = new EmailService();

// Register
router.post('/register', async (req: Request, res: Response): Promise<void> => {
  try {
    const { email, password, name } = req.body;

    if (!email || !password || !name) {
      res.status(400).json({ error: 'Missing required fields' });
      return;
    }

    // Check if user already exists
    const existingUser = await userService.getUserByEmail(email);
    if (existingUser) {
      res.status(400).json({ error: 'User already exists' });
      return;
    }

    // Sign up with Cognito
    const cognitoResult = await authService.signUp({ email, password, name });

    // Create user in database
    const user = await userService.createUser({
      email,
      name,
      cognitoSubject: cognitoResult.userSub,
    });

    // Send welcome email
    await emailService.sendWelcomeEmail(email, name);

    res.status(201).json({
      message: 'User registered successfully',
      userId: user.id,
      emailVerificationRequired: cognitoResult.emailVerificationRequired,
    });
  } catch (error: any) {
    logger.error('Register error:', error);
    res.status(500).json({ error: error.message || 'Registration failed' });
  }
});

// Login
router.post('/login', async (req: Request, res: Response): Promise<void> => {
  try {
    const { email, password } = req.body;

    if (!email || !password) {
      res.status(400).json({ error: 'Missing email or password' });
      return;
    }

    // Sign in with Cognito
    const authResult = await authService.signIn({ email, password });

    // Get user from database
    const user = await userService.getUserByEmail(email);
    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    // Update last login
    await userService.updateUser(user.id, { lastLogin: new Date() });

    res.json({
      message: 'Login successful',
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
        role: user.role,
      },
      tokens: {
        accessToken: authResult.accessToken,
        idToken: authResult.idToken,
        refreshToken: authResult.refreshToken,
        expiresIn: authResult.expiresIn,
      },
    });
  } catch (error: any) {
    logger.error('Login error:', error);
    res.status(401).json({ error: error.message || 'Login failed' });
  }
});

// Forgot password
router.post('/forgot-password', async (req: Request, res: Response): Promise<void> => {
  try {
    const { email } = req.body;

    if (!email) {
      res.status(400).json({ error: 'Email is required' });
      return;
    }

    await authService.forgotPassword(email);

    res.json({ message: 'Password reset email sent' });
  } catch (error: any) {
    logger.error('Forgot password error:', error);
    res.status(500).json({ error: error.message || 'Failed to send reset email' });
  }
});

// Reset password
router.post('/reset-password', async (req: Request, res: Response): Promise<void> => {
  try {
    const { email, code, newPassword } = req.body;

    if (!email || !code || !newPassword) {
      res.status(400).json({ error: 'Missing required fields' });
      return;
    }

    await authService.confirmForgotPassword(email, code, newPassword);

    res.json({ message: 'Password reset successfully' });
  } catch (error: any) {
    logger.error('Reset password error:', error);
    res.status(500).json({ error: error.message || 'Password reset failed' });
  }
});

// Change password (authenticated)
router.post('/change-password', authMiddleware, async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const { currentPassword, newPassword } = req.body;
    const userId = req.userId!;

    if (!currentPassword || !newPassword) {
      res.status(400).json({ error: 'Missing required fields' });
      return;
    }

    const user = await userService.getUserById(userId);
    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    // Verify current password
    if (user.password) {
      const isValid = await bcrypt.compare(currentPassword, user.password);
      if (!isValid) {
        res.status(401).json({ error: 'Current password is incorrect' });
        return;
      }
    }

    // Hash new password
    const hashedPassword = await bcrypt.hash(newPassword, 10);
    await userService.updateUser(userId, { password: hashedPassword });

    res.json({ message: 'Password changed successfully' });
  } catch (error: any) {
    logger.error('Change password error:', error);
    res.status(500).json({ error: error.message || 'Password change failed' });
  }
});

export default router;

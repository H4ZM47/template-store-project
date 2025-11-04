import { Response, NextFunction } from 'express';
import { AuthRequest } from './auth';
import { logger } from '../utils/logger';

export const requireAdmin = (req: AuthRequest, res: Response, next: NextFunction): void => {
  try {
    if (!req.user) {
      res.status(401).json({ error: 'Authentication required' });
      return;
    }

    if (req.user.role !== 'admin') {
      logger.warn('Unauthorized admin access attempt', {
        userId: req.userId,
        userRole: req.user.role,
      });
      res.status(403).json({ error: 'Admin access required' });
      return;
    }

    next();
  } catch (error) {
    logger.error('RBAC middleware error:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
};

export const requireRole = (allowedRoles: string[]) => {
  return (req: AuthRequest, res: Response, next: NextFunction): void => {
    try {
      if (!req.user) {
        res.status(401).json({ error: 'Authentication required' });
        return;
      }

      if (!allowedRoles.includes(req.user.role)) {
        logger.warn('Unauthorized role access attempt', {
          userId: req.userId,
          userRole: req.user.role,
          requiredRoles: allowedRoles,
        });
        res.status(403).json({ error: 'Insufficient permissions' });
        return;
      }

      next();
    } catch (error) {
      logger.error('RBAC middleware error:', error);
      res.status(500).json({ error: 'Internal server error' });
    }
  };
};

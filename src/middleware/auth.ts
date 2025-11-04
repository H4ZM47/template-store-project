import { Request, Response, NextFunction } from 'express';
import { CognitoJwtVerifier } from 'aws-jwt-verify';
import { AppDataSource } from '../database/connection';
import { User } from '../models';
import { config } from '../config/config';
import { logger } from '../utils/logger';

export interface AuthRequest extends Request {
  userId?: string;
  user?: User;
  cognitoSub?: string;
}

const jwtVerifier = CognitoJwtVerifier.create({
  userPoolId: config.aws.cognitoUserPoolId,
  tokenUse: 'access',
  clientId: config.aws.cognitoClientId,
});

export const authMiddleware = async (
  req: AuthRequest,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    // Debug mode bypass (for testing)
    if (config.server.mode === 'debug' && process.env.DEBUG_USER_ID) {
      const userRepository = AppDataSource.getRepository(User);
      const user = await userRepository.findOne({
        where: { id: process.env.DEBUG_USER_ID },
      });

      if (user) {
        req.userId = user.id;
        req.user = user;
        logger.debug('Debug mode: User authenticated', { userId: user.id });
        return next();
      }
    }

    // Extract token from Authorization header
    const authHeader = req.headers.authorization;
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      res.status(401).json({ error: 'No authorization token provided' });
      return;
    }

    const token = authHeader.substring(7);

    // Verify JWT token with Cognito
    const payload = await jwtVerifier.verify(token);
    const cognitoSub = payload.sub;

    if (!cognitoSub) {
      res.status(401).json({ error: 'Invalid token: missing subject' });
      return;
    }

    // Find user by Cognito subject
    const userRepository = AppDataSource.getRepository(User);
    const user = await userRepository.findOne({
      where: { cognitoSubject: cognitoSub },
    });

    if (!user) {
      res.status(401).json({ error: 'User not found' });
      return;
    }

    if (user.status === 'suspended') {
      res.status(403).json({ error: 'Account suspended' });
      return;
    }

    if (user.status === 'deactivated') {
      res.status(403).json({ error: 'Account deactivated' });
      return;
    }

    // Attach user information to request
    req.userId = user.id;
    req.user = user;
    req.cognitoSub = cognitoSub;

    logger.debug('User authenticated', { userId: user.id, email: user.email });
    next();
  } catch (error) {
    logger.error('Authentication error:', error);
    res.status(401).json({ error: 'Invalid or expired token' });
  }
};

export const optionalAuth = async (
  req: AuthRequest,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const authHeader = req.headers.authorization;
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return next();
    }

    const token = authHeader.substring(7);
    const payload = await jwtVerifier.verify(token);
    const cognitoSub = payload.sub;

    if (cognitoSub) {
      const userRepository = AppDataSource.getRepository(User);
      const user = await userRepository.findOne({
        where: { cognitoSubject: cognitoSub },
      });

      if (user && user.status === 'active') {
        req.userId = user.id;
        req.user = user;
        req.cognitoSub = cognitoSub;
      }
    }

    next();
  } catch (error) {
    // Fail silently for optional auth
    next();
  }
};

import 'reflect-metadata';
import express, { Request, Response, NextFunction } from 'express';
import cors from 'cors';
import { config } from './config/config';
import { initializeDatabase } from './database/connection';
import { logger } from './utils/logger';
import routes from './routes';

const app = express();

// Middleware
app.use(cors());
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true, limit: '10mb' }));

// Request logging middleware
app.use((req: Request, res: Response, next: NextFunction) => {
  logger.info(`${req.method} ${req.path}`, {
    method: req.method,
    path: req.path,
    query: req.query,
    ip: req.ip,
  });
  next();
});

// API routes
app.use('/api/v1', routes);

// Root endpoint
app.get('/', (req: Request, res: Response) => {
  res.json({
    message: 'Template Store API',
    version: '1.0.0',
    documentation: '/api-docs',
  });
});

// 404 handler
app.use((req: Request, res: Response) => {
  res.status(404).json({
    error: 'Not found',
    path: req.path,
  });
});

// Error handling middleware
app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
  logger.error('Unhandled error:', err);
  res.status(500).json({
    error: 'Internal server error',
    message: config.server.mode === 'debug' ? err.message : undefined,
  });
});

// Start server
const startServer = async () => {
  try {
    // Initialize database
    await initializeDatabase();

    // Start listening
    app.listen(config.server.port, () => {
      logger.info(`Server started on port ${config.server.port}`, {
        port: config.server.port,
        mode: config.server.mode,
        nodeEnv: process.env.NODE_ENV,
      });
      logger.info(`API available at http://localhost:${config.server.port}/api/v1`);
      logger.info(`Health check at http://localhost:${config.server.port}/api/v1/health`);
    });
  } catch (error) {
    logger.error('Failed to start server:', error);
    process.exit(1);
  }
};

// Handle uncaught errors
process.on('uncaughtException', (error: Error) => {
  logger.error('Uncaught exception:', error);
  process.exit(1);
});

process.on('unhandledRejection', (reason: any) => {
  logger.error('Unhandled rejection:', reason);
  process.exit(1);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.info('SIGTERM received, shutting down gracefully');
  process.exit(0);
});

process.on('SIGINT', () => {
  logger.info('SIGINT received, shutting down gracefully');
  process.exit(0);
});

// Start the server
startServer();

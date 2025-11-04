import 'reflect-metadata';
import { DataSource } from 'typeorm';
import { config } from '../config/config';
import {
  User,
  Category,
  Template,
  BlogPost,
  Order,
  LoginHistory,
  ActivityLog,
  PasswordResetToken,
  EmailVerificationToken,
  UserPreferences,
} from '../models';

const entities = [
  User,
  Category,
  Template,
  BlogPost,
  Order,
  LoginHistory,
  ActivityLog,
  PasswordResetToken,
  EmailVerificationToken,
  UserPreferences,
];

export const AppDataSource = new DataSource(
  config.database.type === 'sqlite'
    ? {
        type: 'sqlite',
        database: config.database.database,
        entities,
        synchronize: true,
        logging: config.server.mode === 'debug',
      }
    : {
        type: 'postgres',
        host: config.database.host,
        port: config.database.port,
        username: config.database.username,
        password: config.database.password,
        database: config.database.database,
        ssl: config.database.ssl ? { rejectUnauthorized: false } : false,
        entities,
        synchronize: true,
        logging: config.server.mode === 'debug',
      }
);

export const initializeDatabase = async (): Promise<void> => {
  try {
    await AppDataSource.initialize();
    console.log('Database connection established successfully');
  } catch (error) {
    console.error('Error initializing database:', error);
    throw error;
  }
};

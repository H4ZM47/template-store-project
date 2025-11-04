import { Repository } from 'typeorm';
import { AppDataSource } from '../database/connection';
import { User, UserRole } from '../models/User';
import { logger } from '../utils/logger';
import bcrypt from 'bcrypt';

export interface CreateUserParams {
  email: string;
  name: string;
  password?: string;
  cognitoSubject?: string;
  role?: UserRole;
}

export class UserService {
  private userRepository: Repository<User>;

  constructor() {
    this.userRepository = AppDataSource.getRepository(User);
  }

  async createUser(params: CreateUserParams): Promise<User> {
    try {
      let hashedPassword: string | undefined;
      if (params.password) {
        hashedPassword = await bcrypt.hash(params.password, 10);
      }

      const user = this.userRepository.create({
        email: params.email,
        name: params.name,
        password: hashedPassword,
        cognitoSubject: params.cognitoSubject,
        role: params.role || 'user',
        status: 'active',
      });

      await this.userRepository.save(user);
      logger.info('User created', { userId: user.id, email: user.email });
      return user;
    } catch (error) {
      logger.error('Create user error:', error);
      throw error;
    }
  }

  async getUserById(id: string): Promise<User | null> {
    try {
      const user = await this.userRepository.findOne({ where: { id } });
      return user;
    } catch (error) {
      logger.error('Get user by ID error:', error);
      throw error;
    }
  }

  async getUserByEmail(email: string): Promise<User | null> {
    try {
      const user = await this.userRepository.findOne({ where: { email } });
      return user;
    } catch (error) {
      logger.error('Get user by email error:', error);
      throw error;
    }
  }

  async getUserByCognitoSubject(cognitoSubject: string): Promise<User | null> {
    try {
      const user = await this.userRepository.findOne({ where: { cognitoSubject } });
      return user;
    } catch (error) {
      logger.error('Get user by Cognito subject error:', error);
      throw error;
    }
  }

  async updateUser(id: string, updates: Partial<User>): Promise<User> {
    try {
      await this.userRepository.update(id, updates);
      const user = await this.getUserById(id);
      if (!user) throw new Error('User not found after update');
      logger.info('User updated', { userId: id });
      return user;
    } catch (error) {
      logger.error('Update user error:', error);
      throw error;
    }
  }

  async deleteUser(id: string): Promise<void> {
    try {
      await this.userRepository.softDelete(id);
      logger.info('User deleted', { userId: id });
    } catch (error) {
      logger.error('Delete user error:', error);
      throw error;
    }
  }

  async listUsers(limit: number = 50, offset: number = 0): Promise<User[]> {
    try {
      const users = await this.userRepository.find({
        take: limit,
        skip: offset,
        order: { createdAt: 'DESC' },
      });
      return users;
    } catch (error) {
      logger.error('List users error:', error);
      throw error;
    }
  }

  async verifyPassword(user: User, password: string): Promise<boolean> {
    if (!user.password) return false;
    return bcrypt.compare(password, user.password);
  }
}

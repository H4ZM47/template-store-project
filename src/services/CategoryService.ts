import { Repository } from 'typeorm';
import { AppDataSource } from '../database/connection';
import { Category } from '../models/Category';
import { logger } from '../utils/logger';

export interface CreateCategoryParams {
  name: string;
  description?: string;
}

export class CategoryService {
  private categoryRepository: Repository<Category>;

  constructor() {
    this.categoryRepository = AppDataSource.getRepository(Category);
  }

  async createCategory(params: CreateCategoryParams): Promise<Category> {
    try {
      const category = this.categoryRepository.create(params);
      await this.categoryRepository.save(category);
      logger.info('Category created', { categoryId: category.id, name: category.name });
      return category;
    } catch (error) {
      logger.error('Create category error:', error);
      throw error;
    }
  }

  async getCategoryById(id: string): Promise<Category | null> {
    try {
      const category = await this.categoryRepository.findOne({ where: { id } });
      return category;
    } catch (error) {
      logger.error('Get category by ID error:', error);
      throw error;
    }
  }

  async getCategoryByName(name: string): Promise<Category | null> {
    try {
      const category = await this.categoryRepository.findOne({ where: { name } });
      return category;
    } catch (error) {
      logger.error('Get category by name error:', error);
      throw error;
    }
  }

  async updateCategory(id: string, updates: Partial<Category>): Promise<Category> {
    try {
      await this.categoryRepository.update(id, updates);
      const category = await this.getCategoryById(id);
      if (!category) throw new Error('Category not found after update');
      logger.info('Category updated', { categoryId: id });
      return category;
    } catch (error) {
      logger.error('Update category error:', error);
      throw error;
    }
  }

  async deleteCategory(id: string): Promise<void> {
    try {
      await this.categoryRepository.softDelete(id);
      logger.info('Category deleted', { categoryId: id });
    } catch (error) {
      logger.error('Delete category error:', error);
      throw error;
    }
  }

  async listCategories(): Promise<Category[]> {
    try {
      const categories = await this.categoryRepository.find({
        order: { name: 'ASC' },
      });
      return categories;
    } catch (error) {
      logger.error('List categories error:', error);
      throw error;
    }
  }
}

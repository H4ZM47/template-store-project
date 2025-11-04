import { Repository } from 'typeorm';
import { AppDataSource } from '../database/connection';
import { Template } from '../models/Template';
import { logger } from '../utils/logger';

export interface CreateTemplateParams {
  name: string;
  description?: string;
  categoryId?: string;
  price: number;
  fileUrl?: string;
  fileSize?: number;
  fileType?: string;
  previewUrl?: string;
  thumbnailUrl?: string;
  variables?: Record<string, any>;
}

export class TemplateService {
  private templateRepository: Repository<Template>;

  constructor() {
    this.templateRepository = AppDataSource.getRepository(Template);
  }

  async createTemplate(params: CreateTemplateParams): Promise<Template> {
    try {
      const template = this.templateRepository.create(params);
      await this.templateRepository.save(template);
      logger.info('Template created', { templateId: template.id, name: template.name });
      return template;
    } catch (error) {
      logger.error('Create template error:', error);
      throw error;
    }
  }

  async getTemplateById(id: string): Promise<Template | null> {
    try {
      const template = await this.templateRepository.findOne({
        where: { id },
        relations: ['category'],
      });
      return template;
    } catch (error) {
      logger.error('Get template by ID error:', error);
      throw error;
    }
  }

  async updateTemplate(id: string, updates: Partial<Template>): Promise<Template> {
    try {
      await this.templateRepository.update(id, updates);
      const template = await this.getTemplateById(id);
      if (!template) throw new Error('Template not found after update');
      logger.info('Template updated', { templateId: id });
      return template;
    } catch (error) {
      logger.error('Update template error:', error);
      throw error;
    }
  }

  async deleteTemplate(id: string): Promise<void> {
    try {
      await this.templateRepository.softDelete(id);
      logger.info('Template deleted', { templateId: id });
    } catch (error) {
      logger.error('Delete template error:', error);
      throw error;
    }
  }

  async listTemplates(limit: number = 50, offset: number = 0): Promise<Template[]> {
    try {
      const templates = await this.templateRepository.find({
        where: { active: true },
        relations: ['category'],
        take: limit,
        skip: offset,
        order: { createdAt: 'DESC' },
      });
      return templates;
    } catch (error) {
      logger.error('List templates error:', error);
      throw error;
    }
  }

  async getTemplatesByCategory(categoryId: string, limit: number = 50): Promise<Template[]> {
    try {
      const templates = await this.templateRepository.find({
        where: { categoryId, active: true },
        relations: ['category'],
        take: limit,
        order: { createdAt: 'DESC' },
      });
      return templates;
    } catch (error) {
      logger.error('Get templates by category error:', error);
      throw error;
    }
  }

  async incrementDownloadCount(id: string): Promise<void> {
    try {
      await this.templateRepository.increment({ id }, 'downloads', 1);
      logger.debug('Template download count incremented', { templateId: id });
    } catch (error) {
      logger.error('Increment download count error:', error);
      throw error;
    }
  }
}

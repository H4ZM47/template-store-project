import { Router, Request, Response } from 'express';
import { TemplateService, StorageService } from '../services';
import { authMiddleware, AuthRequest, requireRole } from '../middleware';
import { logger } from '../utils/logger';

const router = Router();
const templateService = new TemplateService();
const storageService = new StorageService();

// List all templates (public)
router.get('/', async (req: Request, res: Response): Promise<void> => {
  try {
    const limit = parseInt(req.query.limit as string) || 50;
    const offset = parseInt(req.query.offset as string) || 0;

    const templates = await templateService.listTemplates(limit, offset);

    res.json({
      templates,
      count: templates.length,
      limit,
      offset,
    });
  } catch (error: any) {
    logger.error('List templates error:', error);
    res.status(500).json({ error: 'Failed to list templates' });
  }
});

// Get template by ID (public)
router.get('/:id', async (req: Request, res: Response): Promise<void> => {
  try {
    const { id } = req.params;

    const template = await templateService.getTemplateById(id);

    if (!template) {
      res.status(404).json({ error: 'Template not found' });
      return;
    }

    res.json({ template });
  } catch (error: any) {
    logger.error('Get template error:', error);
    res.status(500).json({ error: 'Failed to get template' });
  }
});

// Get templates by category (public)
router.get('/category/:category_id', async (req: Request, res: Response): Promise<void> => {
  try {
    const { category_id } = req.params;
    const limit = parseInt(req.query.limit as string) || 50;

    const templates = await templateService.getTemplatesByCategory(category_id, limit);

    res.json({
      templates,
      count: templates.length,
      categoryId: category_id,
    });
  } catch (error: any) {
    logger.error('Get templates by category error:', error);
    res.status(500).json({ error: 'Failed to get templates' });
  }
});

// Create template (authenticated, author/admin only)
router.post(
  '/',
  authMiddleware,
  requireRole(['admin', 'author']),
  async (req: AuthRequest, res: Response): Promise<void> => {
    try {
      const templateData = req.body;

      const template = await templateService.createTemplate(templateData);

      res.status(201).json({
        message: 'Template created successfully',
        template,
      });
    } catch (error: any) {
      logger.error('Create template error:', error);
      res.status(500).json({ error: 'Failed to create template' });
    }
  }
);

// Update template (authenticated, author/admin only)
router.put(
  '/:id',
  authMiddleware,
  requireRole(['admin', 'author']),
  async (req: AuthRequest, res: Response): Promise<void> => {
    try {
      const { id } = req.params;
      const updates = req.body;

      const template = await templateService.updateTemplate(id, updates);

      res.json({
        message: 'Template updated successfully',
        template,
      });
    } catch (error: any) {
      logger.error('Update template error:', error);
      res.status(500).json({ error: 'Failed to update template' });
    }
  }
);

// Delete template (authenticated, author/admin only)
router.delete(
  '/:id',
  authMiddleware,
  requireRole(['admin', 'author']),
  async (req: AuthRequest, res: Response): Promise<void> => {
    try {
      const { id } = req.params;

      await templateService.deleteTemplate(id);

      res.json({ message: 'Template deleted successfully' });
    } catch (error: any) {
      logger.error('Delete template error:', error);
      res.status(500).json({ error: 'Failed to delete template' });
    }
  }
);

export default router;

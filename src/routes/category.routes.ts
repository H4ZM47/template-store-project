import { Router, Request, Response } from 'express';
import { CategoryService } from '../services';
import { authMiddleware, requireAdmin } from '../middleware';
import { logger } from '../utils/logger';

const router = Router();
const categoryService = new CategoryService();

// List all categories (public)
router.get('/', async (req: Request, res: Response): Promise<void> => {
  try {
    const categories = await categoryService.listCategories();

    res.json({
      categories,
      count: categories.length,
    });
  } catch (error: any) {
    logger.error('List categories error:', error);
    res.status(500).json({ error: 'Failed to list categories' });
  }
});

// Get category by ID (public)
router.get('/:id', async (req: Request, res: Response): Promise<void> => {
  try {
    const { id } = req.params;

    const category = await categoryService.getCategoryById(id);

    if (!category) {
      res.status(404).json({ error: 'Category not found' });
      return;
    }

    res.json({ category });
  } catch (error: any) {
    logger.error('Get category error:', error);
    res.status(500).json({ error: 'Failed to get category' });
  }
});

// Create category (public for now, can be restricted)
router.post('/', async (req: Request, res: Response): Promise<void> => {
  try {
    const { name, description } = req.body;

    if (!name) {
      res.status(400).json({ error: 'Name is required' });
      return;
    }

    const category = await categoryService.createCategory({ name, description });

    res.status(201).json({
      message: 'Category created successfully',
      category,
    });
  } catch (error: any) {
    logger.error('Create category error:', error);
    res.status(500).json({ error: 'Failed to create category' });
  }
});

// Update category (public for now, can be restricted)
router.put('/:id', async (req: Request, res: Response): Promise<void> => {
  try {
    const { id } = req.params;
    const updates = req.body;

    const category = await categoryService.updateCategory(id, updates);

    res.json({
      message: 'Category updated successfully',
      category,
    });
  } catch (error: any) {
    logger.error('Update category error:', error);
    res.status(500).json({ error: 'Failed to update category' });
  }
});

// Delete category (public for now, can be restricted)
router.delete('/:id', async (req: Request, res: Response): Promise<void> => {
  try {
    const { id } = req.params;

    await categoryService.deleteCategory(id);

    res.json({ message: 'Category deleted successfully' });
  } catch (error: any) {
    logger.error('Delete category error:', error);
    res.status(500).json({ error: 'Failed to delete category' });
  }
});

export default router;

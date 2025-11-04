import { Router, Request, Response } from 'express';
import { BlogService } from '../services';
import { authMiddleware, AuthRequest, requireRole } from '../middleware';
import { logger } from '../utils/logger';

const router = Router();
const blogService = new BlogService();

// List all blog posts (public)
router.get('/', async (req: Request, res: Response): Promise<void> => {
  try {
    const limit = parseInt(req.query.limit as string) || 50;
    const offset = parseInt(req.query.offset as string) || 0;

    const blogPosts = await blogService.listBlogPosts(limit, offset, true);

    res.json({
      blogPosts,
      count: blogPosts.length,
      limit,
      offset,
    });
  } catch (error: any) {
    logger.error('List blog posts error:', error);
    res.status(500).json({ error: 'Failed to list blog posts' });
  }
});

// Get blog post by ID (public)
router.get('/:id', async (req: Request, res: Response): Promise<void> => {
  try {
    const { id } = req.params;

    const blogPost = await blogService.getBlogPostById(id);

    if (!blogPost) {
      res.status(404).json({ error: 'Blog post not found' });
      return;
    }

    // Increment view count
    await blogService.incrementViewCount(id);

    // Render markdown content
    const renderedContent = blogService.renderMarkdown(blogPost.content);

    res.json({
      blogPost: {
        ...blogPost,
        renderedContent,
      },
    });
  } catch (error: any) {
    logger.error('Get blog post error:', error);
    res.status(500).json({ error: 'Failed to get blog post' });
  }
});

// Get blog posts by category (public)
router.get('/category/:category_id', async (req: Request, res: Response): Promise<void> => {
  try {
    const { category_id } = req.params;
    const limit = parseInt(req.query.limit as string) || 50;

    const blogPosts = await blogService.getBlogPostsByCategory(category_id, limit);

    res.json({
      blogPosts,
      count: blogPosts.length,
      categoryId: category_id,
    });
  } catch (error: any) {
    logger.error('Get blog posts by category error:', error);
    res.status(500).json({ error: 'Failed to get blog posts' });
  }
});

// Get blog posts by author (public)
router.get('/author/:author_id', async (req: Request, res: Response): Promise<void> => {
  try {
    const { author_id } = req.params;
    const limit = parseInt(req.query.limit as string) || 50;

    const blogPosts = await blogService.getBlogPostsByAuthor(author_id, limit);

    res.json({
      blogPosts,
      count: blogPosts.length,
      authorId: author_id,
    });
  } catch (error: any) {
    logger.error('Get blog posts by author error:', error);
    res.status(500).json({ error: 'Failed to get blog posts' });
  }
});

// Create blog post (authenticated, author/admin only)
router.post(
  '/',
  authMiddleware,
  requireRole(['admin', 'author']),
  async (req: AuthRequest, res: Response): Promise<void> => {
    try {
      const blogPostData = {
        ...req.body,
        authorId: req.userId!,
      };

      const blogPost = await blogService.createBlogPost(blogPostData);

      res.status(201).json({
        message: 'Blog post created successfully',
        blogPost,
      });
    } catch (error: any) {
      logger.error('Create blog post error:', error);
      res.status(500).json({ error: 'Failed to create blog post' });
    }
  }
);

// Update blog post (authenticated, author/admin only)
router.put(
  '/:id',
  authMiddleware,
  requireRole(['admin', 'author']),
  async (req: AuthRequest, res: Response): Promise<void> => {
    try {
      const { id } = req.params;
      const updates = req.body;

      const blogPost = await blogService.updateBlogPost(id, updates);

      res.json({
        message: 'Blog post updated successfully',
        blogPost,
      });
    } catch (error: any) {
      logger.error('Update blog post error:', error);
      res.status(500).json({ error: 'Failed to update blog post' });
    }
  }
);

// Delete blog post (authenticated, author/admin only)
router.delete(
  '/:id',
  authMiddleware,
  requireRole(['admin', 'author']),
  async (req: AuthRequest, res: Response): Promise<void> => {
    try {
      const { id } = req.params;

      await blogService.deleteBlogPost(id);

      res.json({ message: 'Blog post deleted successfully' });
    } catch (error: any) {
      logger.error('Delete blog post error:', error);
      res.status(500).json({ error: 'Failed to delete blog post' });
    }
  }
);

export default router;

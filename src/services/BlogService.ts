import { Repository } from 'typeorm';
import { AppDataSource } from '../database/connection';
import { BlogPost } from '../models/BlogPost';
import { logger } from '../utils/logger';
import { marked } from 'marked';

export interface CreateBlogPostParams {
  title: string;
  content: string;
  authorId: string;
  categoryId?: string;
  excerpt?: string;
  featuredImage?: string;
  tags?: string[];
  published?: boolean;
  slug?: string;
  metaTitle?: string;
  metaDescription?: string;
}

export class BlogService {
  private blogRepository: Repository<BlogPost>;

  constructor() {
    this.blogRepository = AppDataSource.getRepository(BlogPost);
  }

  async createBlogPost(params: CreateBlogPostParams): Promise<BlogPost> {
    try {
      const blogPost = this.blogRepository.create({
        ...params,
        published: params.published || false,
        publishedAt: params.published ? new Date() : undefined,
      });

      await this.blogRepository.save(blogPost);
      logger.info('Blog post created', { blogPostId: blogPost.id, title: blogPost.title });
      return blogPost;
    } catch (error) {
      logger.error('Create blog post error:', error);
      throw error;
    }
  }

  async getBlogPostById(id: string): Promise<BlogPost | null> {
    try {
      const blogPost = await this.blogRepository.findOne({
        where: { id },
        relations: ['author', 'category'],
      });
      return blogPost;
    } catch (error) {
      logger.error('Get blog post by ID error:', error);
      throw error;
    }
  }

  async getBlogPostBySlug(slug: string): Promise<BlogPost | null> {
    try {
      const blogPost = await this.blogRepository.findOne({
        where: { slug },
        relations: ['author', 'category'],
      });
      return blogPost;
    } catch (error) {
      logger.error('Get blog post by slug error:', error);
      throw error;
    }
  }

  async updateBlogPost(id: string, updates: Partial<BlogPost>): Promise<BlogPost> {
    try {
      // If publishing, set publishedAt
      if (updates.published && !updates.publishedAt) {
        updates.publishedAt = new Date();
      }

      await this.blogRepository.update(id, updates);
      const blogPost = await this.getBlogPostById(id);
      if (!blogPost) throw new Error('Blog post not found after update');
      logger.info('Blog post updated', { blogPostId: id });
      return blogPost;
    } catch (error) {
      logger.error('Update blog post error:', error);
      throw error;
    }
  }

  async deleteBlogPost(id: string): Promise<void> {
    try {
      await this.blogRepository.softDelete(id);
      logger.info('Blog post deleted', { blogPostId: id });
    } catch (error) {
      logger.error('Delete blog post error:', error);
      throw error;
    }
  }

  async listBlogPosts(
    limit: number = 50,
    offset: number = 0,
    publishedOnly: boolean = true
  ): Promise<BlogPost[]> {
    try {
      const where = publishedOnly ? { published: true } : {};
      const blogPosts = await this.blogRepository.find({
        where,
        relations: ['author', 'category'],
        take: limit,
        skip: offset,
        order: { publishedAt: 'DESC', createdAt: 'DESC' },
      });
      return blogPosts;
    } catch (error) {
      logger.error('List blog posts error:', error);
      throw error;
    }
  }

  async getBlogPostsByCategory(categoryId: string, limit: number = 50): Promise<BlogPost[]> {
    try {
      const blogPosts = await this.blogRepository.find({
        where: { categoryId, published: true },
        relations: ['author', 'category'],
        take: limit,
        order: { publishedAt: 'DESC' },
      });
      return blogPosts;
    } catch (error) {
      logger.error('Get blog posts by category error:', error);
      throw error;
    }
  }

  async getBlogPostsByAuthor(authorId: string, limit: number = 50): Promise<BlogPost[]> {
    try {
      const blogPosts = await this.blogRepository.find({
        where: { authorId },
        relations: ['category'],
        take: limit,
        order: { createdAt: 'DESC' },
      });
      return blogPosts;
    } catch (error) {
      logger.error('Get blog posts by author error:', error);
      throw error;
    }
  }

  async incrementViewCount(id: string): Promise<void> {
    try {
      await this.blogRepository.increment({ id }, 'viewCount', 1);
      logger.debug('Blog post view count incremented', { blogPostId: id });
    } catch (error) {
      logger.error('Increment view count error:', error);
      throw error;
    }
  }

  renderMarkdown(content: string): string {
    return marked(content) as string;
  }
}

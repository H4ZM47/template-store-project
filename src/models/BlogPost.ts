import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  DeleteDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { User } from './User';
import { Category } from './Category';

@Entity('blog_posts')
export class BlogPost {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column()
  title!: string;

  @Column({ type: 'text' })
  content!: string;

  @Column({ name: 'author_id' })
  authorId!: string;

  @Column({ name: 'category_id', nullable: true })
  categoryId?: string;

  @Column({ nullable: true })
  excerpt?: string;

  @Column({ name: 'featured_image', nullable: true })
  featuredImage?: string;

  @Column({ type: 'varchar', array: true, nullable: true })
  tags?: string[];

  @Column({ default: false })
  published!: boolean;

  @Column({ name: 'published_at', type: 'timestamp', nullable: true })
  publishedAt?: Date;

  @Column({ nullable: true })
  slug?: string;

  @Column({ name: 'meta_title', nullable: true })
  metaTitle?: string;

  @Column({ name: 'meta_description', nullable: true })
  metaDescription?: string;

  @Column({ name: 'view_count', type: 'int', default: 0 })
  viewCount!: number;

  @CreateDateColumn({ name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt!: Date;

  @DeleteDateColumn({ name: 'deleted_at' })
  deletedAt?: Date;

  // Relationships
  @ManyToOne(() => User, (user) => user.blogPosts)
  @JoinColumn({ name: 'author_id' })
  author!: User;

  @ManyToOne(() => Category, (category) => category.blogPosts)
  @JoinColumn({ name: 'category_id' })
  category?: Category;
}

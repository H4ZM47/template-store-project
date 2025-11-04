import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  DeleteDateColumn,
  ManyToOne,
  JoinColumn,
  OneToMany,
} from 'typeorm';
import { Category } from './Category';
import { Order } from './Order';

@Entity('templates')
export class Template {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column()
  name!: string;

  @Column({ nullable: true })
  description?: string;

  @Column({ name: 'category_id', nullable: true })
  categoryId?: string;

  @Column({ type: 'decimal', precision: 10, scale: 2, default: 0 })
  price!: number;

  @Column({ name: 'file_url', nullable: true })
  fileUrl?: string;

  @Column({ name: 'file_size', type: 'bigint', nullable: true })
  fileSize?: number;

  @Column({ name: 'file_type', nullable: true })
  fileType?: string;

  @Column({ name: 'preview_url', nullable: true })
  previewUrl?: string;

  @Column({ name: 'thumbnail_url', nullable: true })
  thumbnailUrl?: string;

  @Column({ type: 'jsonb', nullable: true })
  variables?: Record<string, any>;

  @Column({ type: 'int', default: 0 })
  downloads!: number;

  @Column({ default: true })
  active!: boolean;

  @CreateDateColumn({ name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt!: Date;

  @DeleteDateColumn({ name: 'deleted_at' })
  deletedAt?: Date;

  // Relationships
  @ManyToOne(() => Category, (category) => category.templates)
  @JoinColumn({ name: 'category_id' })
  category?: Category;

  @OneToMany(() => Order, (order) => order.template)
  orders!: Order[];
}

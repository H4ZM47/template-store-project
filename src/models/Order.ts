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
import { Template } from './Template';

export type OrderStatus = 'pending' | 'completed' | 'failed' | 'refunded';
export type DeliveryStatus = 'pending' | 'delivered' | 'failed';

@Entity('orders')
export class Order {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column({ name: 'user_id' })
  userId!: string;

  @Column({ name: 'template_id' })
  templateId!: string;

  @Column({ type: 'decimal', precision: 10, scale: 2 })
  amount!: number;

  @Column({ type: 'varchar', default: 'pending' })
  status!: OrderStatus;

  @Column({ name: 'delivery_status', type: 'varchar', default: 'pending' })
  deliveryStatus!: DeliveryStatus;

  @Column({ name: 'stripe_payment_intent_id', nullable: true })
  stripePaymentIntentId?: string;

  @Column({ name: 'stripe_session_id', nullable: true })
  stripeSessionId?: string;

  @Column({ name: 'download_url', nullable: true })
  downloadUrl?: string;

  @Column({ name: 'download_expires_at', type: 'timestamp', nullable: true })
  downloadExpiresAt?: Date;

  @Column({ type: 'jsonb', nullable: true })
  metadata?: Record<string, any>;

  @CreateDateColumn({ name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt!: Date;

  @DeleteDateColumn({ name: 'deleted_at' })
  deletedAt?: Date;

  // Relationships
  @ManyToOne(() => User, (user) => user.orders)
  @JoinColumn({ name: 'user_id' })
  user!: User;

  @ManyToOne(() => Template, (template) => template.orders)
  @JoinColumn({ name: 'template_id' })
  template!: Template;
}

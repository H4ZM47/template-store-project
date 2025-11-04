import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  OneToOne,
  JoinColumn,
} from 'typeorm';
import { User } from './User';

@Entity('user_preferences')
export class UserPreferences {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column({ name: 'user_id', unique: true })
  userId!: string;

  @Column({ type: 'jsonb', default: '{}' })
  preferences!: Record<string, any>;

  @Column({ name: 'email_notifications', default: true })
  emailNotifications!: boolean;

  @Column({ name: 'marketing_emails', default: false })
  marketingEmails!: boolean;

  @Column({ nullable: true })
  language?: string;

  @Column({ nullable: true })
  timezone?: string;

  @CreateDateColumn({ name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt!: Date;

  // Relationships
  @OneToOne(() => User, (user) => user.preferences)
  @JoinColumn({ name: 'user_id' })
  user!: User;
}

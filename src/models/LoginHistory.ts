import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { User } from './User';

@Entity('login_history')
export class LoginHistory {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column({ name: 'user_id' })
  userId!: string;

  @Column({ name: 'ip_address', nullable: true })
  ipAddress?: string;

  @Column({ name: 'user_agent', nullable: true })
  userAgent?: string;

  @Column({ name: 'login_method', nullable: true })
  loginMethod?: string;

  @Column({ default: true })
  successful!: boolean;

  @Column({ name: 'failure_reason', nullable: true })
  failureReason?: string;

  @CreateDateColumn({ name: 'created_at' })
  createdAt!: Date;

  // Relationships
  @ManyToOne(() => User, (user) => user.loginHistory)
  @JoinColumn({ name: 'user_id' })
  user!: User;
}

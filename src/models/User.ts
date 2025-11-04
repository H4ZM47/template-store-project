import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  DeleteDateColumn,
  OneToMany,
  OneToOne,
} from 'typeorm';
import { BlogPost } from './BlogPost';
import { Order } from './Order';
import { LoginHistory } from './LoginHistory';
import { ActivityLog } from './ActivityLog';
import { UserPreferences } from './UserPreferences';

export type UserRole = 'user' | 'admin' | 'author';
export type UserStatus = 'active' | 'suspended' | 'deactivated';

@Entity('users')
export class User {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column({ unique: true })
  email!: string;

  @Column({ name: 'cognito_subject', nullable: true, unique: true })
  cognitoSubject?: string;

  @Column({ nullable: true })
  password?: string;

  @Column()
  name!: string;

  @Column({ type: 'varchar', default: 'user' })
  role!: UserRole;

  @Column({ type: 'varchar', default: 'active' })
  status!: UserStatus;

  @Column({ name: 'avatar_url', nullable: true })
  avatarUrl?: string;

  @Column({ name: 'phone_number', nullable: true })
  phoneNumber?: string;

  @Column({ nullable: true })
  address?: string;

  @Column({ nullable: true })
  city?: string;

  @Column({ nullable: true })
  state?: string;

  @Column({ name: 'postal_code', nullable: true })
  postalCode?: string;

  @Column({ nullable: true })
  country?: string;

  @Column({ name: 'email_verified', default: false })
  emailVerified!: boolean;

  @Column({ name: 'last_login', type: 'timestamp', nullable: true })
  lastLogin?: Date;

  @CreateDateColumn({ name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt!: Date;

  @DeleteDateColumn({ name: 'deleted_at' })
  deletedAt?: Date;

  // Relationships
  @OneToMany(() => BlogPost, (blogPost) => blogPost.author)
  blogPosts!: BlogPost[];

  @OneToMany(() => Order, (order) => order.user)
  orders!: Order[];

  @OneToMany(() => LoginHistory, (loginHistory) => loginHistory.user)
  loginHistory!: LoginHistory[];

  @OneToMany(() => ActivityLog, (activityLog) => activityLog.user)
  activityLogs!: ActivityLog[];

  @OneToOne(() => UserPreferences, (preferences) => preferences.user)
  preferences?: UserPreferences;
}

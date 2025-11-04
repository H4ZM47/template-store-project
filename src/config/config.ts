import dotenv from 'dotenv';

dotenv.config();

export interface ServerConfig {
  port: number;
  mode: 'debug' | 'release';
}

export interface DatabaseConfig {
  host: string;
  port: number;
  username: string;
  password: string;
  database: string;
  ssl: boolean;
  type: 'postgres' | 'sqlite';
}

export interface AWSConfig {
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
  s3Bucket: string;
  cognitoUserPoolId: string;
  cognitoClientId: string;
  cognitoRegion: string;
}

export interface StripeConfig {
  apiKey: string;
  webhookSecret: string;
  successUrl: string;
  cancelUrl: string;
}

export interface SendGridConfig {
  apiKey: string;
  fromEmail: string;
}

export interface AppConfig {
  server: ServerConfig;
  database: DatabaseConfig;
  aws: AWSConfig;
  stripe: StripeConfig;
  sendgrid: SendGridConfig;
  jwtSecret: string;
}

const getConfig = (): AppConfig => {
  const mode = (process.env.GIN_MODE || 'debug') as 'debug' | 'release';

  return {
    server: {
      port: parseInt(process.env.PORT || '8080', 10),
      mode,
    },
    database: {
      type: mode === 'debug' ? 'sqlite' : 'postgres',
      host: process.env.DB_HOST || 'localhost',
      port: parseInt(process.env.DB_PORT || '5432', 10),
      username: process.env.DB_USER || 'postgres',
      password: process.env.DB_PASSWORD || '',
      database: process.env.DB_NAME || mode === 'debug' ? 'gorm.db' : 'template_store',
      ssl: process.env.DB_SSLMODE === 'require',
    },
    aws: {
      region: process.env.AWS_REGION || 'us-east-1',
      accessKeyId: process.env.AWS_ACCESS_KEY_ID || '',
      secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY || '',
      s3Bucket: process.env.AWS_S3_BUCKET || '',
      cognitoUserPoolId: process.env.AWS_COGNITO_USER_POOL_ID || '',
      cognitoClientId: process.env.AWS_COGNITO_CLIENT_ID || '',
      cognitoRegion: process.env.AWS_COGNITO_REGION || 'us-east-1',
    },
    stripe: {
      apiKey: process.env.STRIPE_API_KEY || '',
      webhookSecret: process.env.STRIPE_WEBHOOK_SECRET || '',
      successUrl: process.env.STRIPE_SUCCESS_URL || 'http://localhost:3000/payment/success',
      cancelUrl: process.env.STRIPE_CANCEL_URL || 'http://localhost:3000/payment/cancel',
    },
    sendgrid: {
      apiKey: process.env.SENDGRID_API_KEY || '',
      fromEmail: process.env.SENDGRID_FROM_EMAIL || 'noreply@templatestore.com',
    },
    jwtSecret: process.env.JWT_SECRET || 'your-secret-key',
  };
};

export const config = getConfig();

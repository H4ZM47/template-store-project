import sgMail from '@sendgrid/mail';
import { config } from '../config/config';
import { logger } from '../utils/logger';

export interface EmailParams {
  to: string;
  subject: string;
  text?: string;
  html?: string;
}

export class EmailService {
  constructor() {
    sgMail.setApiKey(config.sendgrid.apiKey);
  }

  async sendEmail(params: EmailParams): Promise<void> {
    try {
      const mailData: any = {
        from: config.sendgrid.fromEmail,
        to: params.to,
        subject: params.subject,
      };

      if (params.text) {
        mailData.text = params.text;
      }

      if (params.html) {
        mailData.html = params.html;
      }

      await sgMail.send(mailData);

      logger.info('Email sent successfully', { to: params.to, subject: params.subject });
    } catch (error) {
      logger.error('Email send error:', error);
      throw error;
    }
  }

  async sendWelcomeEmail(email: string, name: string): Promise<void> {
    const subject = 'Welcome to Template Store';
    const html = `
      <h1>Welcome ${name}!</h1>
      <p>Thank you for joining Template Store. We're excited to have you on board.</p>
      <p>Start exploring our collection of templates and find the perfect one for your needs.</p>
    `;

    await this.sendEmail({ to: email, subject, html });
  }

  async sendPasswordResetEmail(email: string, resetToken: string): Promise<void> {
    const resetUrl = `${process.env.FRONTEND_URL || 'http://localhost:3000'}/reset-password?token=${resetToken}`;
    const subject = 'Password Reset Request';
    const html = `
      <h1>Password Reset Request</h1>
      <p>You requested to reset your password. Click the link below to reset it:</p>
      <p><a href="${resetUrl}">Reset Password</a></p>
      <p>This link will expire in 1 hour.</p>
      <p>If you didn't request this, please ignore this email.</p>
    `;

    await this.sendEmail({ to: email, subject, html });
  }

  async sendEmailVerification(email: string, verificationToken: string): Promise<void> {
    const verificationUrl = `${process.env.FRONTEND_URL || 'http://localhost:3000'}/verify-email?token=${verificationToken}`;
    const subject = 'Email Verification';
    const html = `
      <h1>Email Verification</h1>
      <p>Please verify your email address by clicking the link below:</p>
      <p><a href="${verificationUrl}">Verify Email</a></p>
      <p>This link will expire in 24 hours.</p>
    `;

    await this.sendEmail({ to: email, subject, html });
  }

  async sendOrderConfirmation(email: string, orderId: string, templateName: string): Promise<void> {
    const subject = 'Order Confirmation';
    const html = `
      <h1>Order Confirmation</h1>
      <p>Thank you for your purchase!</p>
      <p><strong>Order ID:</strong> ${orderId}</p>
      <p><strong>Template:</strong> ${templateName}</p>
      <p>You can download your template from your dashboard.</p>
    `;

    await this.sendEmail({ to: email, subject, html });
  }
}

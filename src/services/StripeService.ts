import Stripe from 'stripe';
import { config } from '../config/config';
import { logger } from '../utils/logger';

export interface CheckoutSessionParams {
  templateId: string;
  templateName: string;
  price: number;
  userEmail: string;
  userId: string;
}

export class StripeService {
  private stripe: Stripe;

  constructor() {
    this.stripe = new Stripe(config.stripe.apiKey, {
      apiVersion: '2025-02-24.acacia',
    });
  }

  async createCheckoutSession(params: CheckoutSessionParams): Promise<string> {
    try {
      const session = await this.stripe.checkout.sessions.create({
        payment_method_types: ['card'],
        line_items: [
          {
            price_data: {
              currency: 'usd',
              product_data: {
                name: params.templateName,
                description: `Template: ${params.templateName}`,
              },
              unit_amount: Math.round(params.price * 100), // Convert to cents
            },
            quantity: 1,
          },
        ],
        mode: 'payment',
        success_url: `${config.stripe.successUrl}?session_id={CHECKOUT_SESSION_ID}`,
        cancel_url: config.stripe.cancelUrl,
        customer_email: params.userEmail,
        metadata: {
          userId: params.userId,
          templateId: params.templateId,
        },
      });

      logger.info('Checkout session created', { sessionId: session.id, userId: params.userId });
      return session.url!;
    } catch (error) {
      logger.error('Stripe checkout session creation error:', error);
      throw error;
    }
  }

  async getSession(sessionId: string): Promise<Stripe.Checkout.Session> {
    try {
      const session = await this.stripe.checkout.sessions.retrieve(sessionId);
      return session;
    } catch (error) {
      logger.error('Stripe get session error:', error);
      throw error;
    }
  }

  async constructWebhookEvent(
    payload: string | Buffer,
    signature: string
  ): Promise<Stripe.Event> {
    try {
      const event = this.stripe.webhooks.constructEvent(
        payload,
        signature,
        config.stripe.webhookSecret
      );
      return event;
    } catch (error) {
      logger.error('Stripe webhook verification error:', error);
      throw error;
    }
  }

  async createRefund(paymentIntentId: string): Promise<Stripe.Refund> {
    try {
      const refund = await this.stripe.refunds.create({
        payment_intent: paymentIntentId,
      });

      logger.info('Refund created', { refundId: refund.id, paymentIntentId });
      return refund;
    } catch (error) {
      logger.error('Stripe refund creation error:', error);
      throw error;
    }
  }
}

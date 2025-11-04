import { Router, Request, Response } from 'express';
import { StripeService, OrderService, TemplateService, EmailService, UserService } from '../services';
import { authMiddleware, AuthRequest } from '../middleware';
import { logger } from '../utils/logger';
import { config } from '../config/config';

const router = Router();
const stripeService = new StripeService();
const orderService = new OrderService();
const templateService = new TemplateService();
const emailService = new EmailService();
const userService = new UserService();

// Create checkout session (authenticated)
router.post('/checkout', authMiddleware, async (req: AuthRequest, res: Response): Promise<void> => {
  try {
    const { templateId } = req.body;
    const userId = req.userId!;

    if (!templateId) {
      res.status(400).json({ error: 'Template ID is required' });
      return;
    }

    // Get template
    const template = await templateService.getTemplateById(templateId);
    if (!template) {
      res.status(404).json({ error: 'Template not found' });
      return;
    }

    // Get user
    const user = await userService.getUserById(userId);
    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    // Create Stripe checkout session
    const checkoutUrl = await stripeService.createCheckoutSession({
      templateId: template.id,
      templateName: template.name,
      price: Number(template.price),
      userEmail: user.email,
      userId: user.id,
    });

    res.json({
      checkoutUrl,
      message: 'Checkout session created',
    });
  } catch (error: any) {
    logger.error('Create checkout session error:', error);
    res.status(500).json({ error: 'Failed to create checkout session' });
  }
});

// Payment success callback (public)
router.get('/success', async (req: Request, res: Response): Promise<void> => {
  try {
    const { session_id } = req.query;

    if (!session_id) {
      res.status(400).json({ error: 'Session ID is required' });
      return;
    }

    // Get session from Stripe
    const session = await stripeService.getSession(session_id as string);

    if (session.payment_status === 'paid') {
      // Find or create order
      let order = await orderService.getOrderByStripeSessionId(session_id as string);

      if (!order) {
        // Create order
        order = await orderService.createOrder({
          userId: session.metadata!.userId,
          templateId: session.metadata!.templateId,
          amount: session.amount_total! / 100,
          stripeSessionId: session.id,
          stripePaymentIntentId: session.payment_intent as string,
        });
      }

      // Update order status
      await orderService.updateOrderStatus(order.id, 'completed');

      res.json({
        message: 'Payment successful',
        orderId: order.id,
      });
    } else {
      res.status(400).json({ error: 'Payment not completed' });
    }
  } catch (error: any) {
    logger.error('Payment success error:', error);
    res.status(500).json({ error: 'Failed to process payment' });
  }
});

// Payment cancel callback (public)
router.get('/cancel', async (req: Request, res: Response): Promise<void> => {
  res.json({ message: 'Payment cancelled' });
});

// Stripe webhook (public)
router.post('/webhooks/stripe', async (req: Request, res: Response): Promise<void> => {
  try {
    const signature = req.headers['stripe-signature'] as string;
    const payload = req.body;

    const event = await stripeService.constructWebhookEvent(payload, signature);

    // Handle the event
    switch (event.type) {
      case 'checkout.session.completed': {
        const session = event.data.object as any;
        logger.info('Checkout session completed', { sessionId: session.id });

        // Create order
        const order = await orderService.createOrder({
          userId: session.metadata.userId,
          templateId: session.metadata.templateId,
          amount: session.amount_total / 100,
          stripeSessionId: session.id,
          stripePaymentIntentId: session.payment_intent,
        });

        await orderService.updateOrderStatus(order.id, 'completed');
        break;
      }
      case 'payment_intent.succeeded': {
        const paymentIntent = event.data.object as any;
        logger.info('Payment intent succeeded', { paymentIntentId: paymentIntent.id });
        break;
      }
      default:
        logger.debug('Unhandled event type', { type: event.type });
    }

    res.json({ received: true });
  } catch (error: any) {
    logger.error('Stripe webhook error:', error);
    res.status(400).json({ error: 'Webhook processing failed' });
  }
});

export default router;

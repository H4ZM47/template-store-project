import { Repository } from 'typeorm';
import { AppDataSource } from '../database/connection';
import { Order, OrderStatus, DeliveryStatus } from '../models/Order';
import { logger } from '../utils/logger';

export interface CreateOrderParams {
  userId: string;
  templateId: string;
  amount: number;
  stripeSessionId?: string;
  stripePaymentIntentId?: string;
}

export class OrderService {
  private orderRepository: Repository<Order>;

  constructor() {
    this.orderRepository = AppDataSource.getRepository(Order);
  }

  async createOrder(params: CreateOrderParams): Promise<Order> {
    try {
      const order = this.orderRepository.create({
        ...params,
        status: 'pending',
        deliveryStatus: 'pending',
      });

      await this.orderRepository.save(order);
      logger.info('Order created', { orderId: order.id, userId: params.userId });
      return order;
    } catch (error) {
      logger.error('Create order error:', error);
      throw error;
    }
  }

  async getOrderById(id: string): Promise<Order | null> {
    try {
      const order = await this.orderRepository.findOne({
        where: { id },
        relations: ['user', 'template'],
      });
      return order;
    } catch (error) {
      logger.error('Get order by ID error:', error);
      throw error;
    }
  }

  async getOrderByStripeSessionId(sessionId: string): Promise<Order | null> {
    try {
      const order = await this.orderRepository.findOne({
        where: { stripeSessionId: sessionId },
        relations: ['user', 'template'],
      });
      return order;
    } catch (error) {
      logger.error('Get order by Stripe session ID error:', error);
      throw error;
    }
  }

  async updateOrderStatus(id: string, status: OrderStatus): Promise<Order> {
    try {
      await this.orderRepository.update(id, { status });
      const order = await this.getOrderById(id);
      if (!order) throw new Error('Order not found after update');
      logger.info('Order status updated', { orderId: id, status });
      return order;
    } catch (error) {
      logger.error('Update order status error:', error);
      throw error;
    }
  }

  async updateDeliveryStatus(id: string, deliveryStatus: DeliveryStatus): Promise<Order> {
    try {
      await this.orderRepository.update(id, { deliveryStatus });
      const order = await this.getOrderById(id);
      if (!order) throw new Error('Order not found after update');
      logger.info('Order delivery status updated', { orderId: id, deliveryStatus });
      return order;
    } catch (error) {
      logger.error('Update delivery status error:', error);
      throw error;
    }
  }

  async getOrdersByUser(userId: string, limit: number = 50): Promise<Order[]> {
    try {
      const orders = await this.orderRepository.find({
        where: { userId },
        relations: ['template'],
        take: limit,
        order: { createdAt: 'DESC' },
      });
      return orders;
    } catch (error) {
      logger.error('Get orders by user error:', error);
      throw error;
    }
  }

  async getUserPurchasedTemplates(userId: string): Promise<Order[]> {
    try {
      const orders = await this.orderRepository.find({
        where: { userId, status: 'completed' },
        relations: ['template'],
        order: { createdAt: 'DESC' },
      });
      return orders;
    } catch (error) {
      logger.error('Get user purchased templates error:', error);
      throw error;
    }
  }

  async setDownloadUrl(id: string, downloadUrl: string, expiresAt: Date): Promise<Order> {
    try {
      await this.orderRepository.update(id, { downloadUrl, downloadExpiresAt: expiresAt });
      const order = await this.getOrderById(id);
      if (!order) throw new Error('Order not found after update');
      logger.info('Order download URL set', { orderId: id });
      return order;
    } catch (error) {
      logger.error('Set download URL error:', error);
      throw error;
    }
  }
}

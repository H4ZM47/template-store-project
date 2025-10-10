package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"template-store/internal/services"

	"github.com/gin-gonic/gin"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	stripeService   *services.StripeService
	orderService    *services.OrderService
	templateService *services.TemplateService
	userService     *services.UserService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(stripeService *services.StripeService, orderService *services.OrderService, templateService *services.TemplateService, userService *services.UserService) *PaymentHandler {
	return &PaymentHandler{
		stripeService:   stripeService,
		orderService:    orderService,
		templateService: templateService,
		userService:     userService,
	}
}

// CreateCheckoutSessionRequest represents the request to create a checkout session
type CreateCheckoutSessionRequest struct {
	TemplateID uint   `json:"template_id" binding:"required"`
	UserID     uint   `json:"user_id" binding:"required"`
	UserEmail  string `json:"user_email" binding:"required,email"`
}

// CreateCheckoutSession creates a Stripe checkout session
func (h *PaymentHandler) CreateCheckoutSession(c *gin.Context) {
	var req CreateCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get template details
	template, err := h.templateService.GetTemplate(req.TemplateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Verify user exists
	user, err := h.userService.GetUser(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Use user email from database if not provided
	email := req.UserEmail
	if email == "" {
		email = user.Email
	}

	// Create checkout session
	session, err := h.stripeService.CreateCheckoutSession(
		template.ID,
		template.Name,
		template.Price,
		email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session"})
		return
	}

	// Create pending order
	order, err := h.orderService.CreateOrder(req.UserID, req.TemplateID, session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":  session.ID,
		"session_url": session.URL,
		"order_id":    order.ID,
	})
}

// CreatePaymentIntentRequest represents the request to create a payment intent
type CreatePaymentIntentRequest struct {
	TemplateID uint   `json:"template_id" binding:"required"`
	UserID     uint   `json:"user_id" binding:"required"`
	Currency   string `json:"currency"`
}

// CreatePaymentIntent creates a Stripe payment intent for direct payments
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get template details
	template, err := h.templateService.GetTemplate(req.TemplateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Verify user exists
	_, err = h.userService.GetUser(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Default currency to USD
	currency := req.Currency
	if currency == "" {
		currency = "usd"
	}

	// Create payment intent
	metadata := map[string]string{
		"template_id": fmt.Sprintf("%d", req.TemplateID),
		"user_id":     fmt.Sprintf("%d", req.UserID),
	}

	pi, err := h.stripeService.CreatePaymentIntent(template.Price, currency, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client_secret": pi.ClientSecret,
		"payment_intent_id": pi.ID,
	})
}

// GetCheckoutSessionSuccess handles successful checkout redirects
func (h *PaymentHandler) GetCheckoutSessionSuccess(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session_id"})
		return
	}

	// Get checkout session
	session, err := h.stripeService.GetCheckoutSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}

	// Get order by session ID
	order, err := h.orderService.GetOrderByStripeSession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Payment successful",
		"session_status": session.PaymentStatus,
		"order":          services.ToOrderResponse(order),
	})
}

// HandleWebhook handles Stripe webhook events
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// Read the request body
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Get the Stripe signature
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe signature"})
		return
	}

	// Verify webhook signature
	event, err := h.stripeService.VerifyWebhookSignature(payload, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// Handle the event
	switch event.Type {
	case "checkout.session.completed":
		session := event.Data.Object
		sessionID, _ := session["id"].(string)

		// Handle successful checkout
		if err := h.handleCheckoutSessionCompletedByID(sessionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process checkout"})
			return
		}

	case "payment_intent.succeeded":
		pi := event.Data.Object
		piID, _ := pi["id"].(string)
		metadata, _ := pi["metadata"].(map[string]interface{})
		amount, _ := pi["amount"].(float64)
		status, _ := pi["status"].(string)

		// Handle successful payment intent
		if err := h.handlePaymentIntentSucceededData(piID, metadata, amount, status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
			return
		}

	case "payment_intent.payment_failed":
		pi := event.Data.Object
		piID, _ := pi["id"].(string)
		metadata, _ := pi["metadata"].(map[string]interface{})

		// Handle failed payment
		if err := h.handlePaymentIntentFailedData(piID, metadata); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment failure"})
			return
		}

	default:
		// Unhandled event type
		c.JSON(http.StatusOK, gin.H{"message": "Event received but not handled"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

// handleCheckoutSessionCompletedByID processes completed checkout sessions
func (h *PaymentHandler) handleCheckoutSessionCompletedByID(sessionID string) error {
	// Get full session details
	session, err := h.stripeService.GetCheckoutSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Find order by session ID
	order, err := h.orderService.GetOrderByStripeSession(session.ID)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}

	// Update order status to paid
	paymentDetails := fmt.Sprintf("Stripe Session: %s | Payment Status: %s | Amount: $%.2f",
		session.ID,
		session.PaymentStatus,
		float64(session.AmountTotal)/100,
	)

	if err := h.orderService.MarkOrderAsPaid(order.ID, paymentDetails); err != nil {
		return fmt.Errorf("failed to mark order as paid: %w", err)
	}

	// TODO: Send delivery email with template download link
	// TODO: Mark as delivered after successful email

	return nil
}

// handlePaymentIntentSucceededData processes successful payment intents
func (h *PaymentHandler) handlePaymentIntentSucceededData(piID string, metadata map[string]interface{}, amount float64, status string) error {
	// Extract metadata
	templateIDStr, ok := metadata["template_id"].(string)
	if !ok {
		return fmt.Errorf("missing template_id in metadata")
	}
	userIDStr, ok := metadata["user_id"].(string)
	if !ok {
		return fmt.Errorf("missing user_id in metadata")
	}

	templateID, _ := strconv.ParseUint(templateIDStr, 10, 32)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	// Create order
	order, err := h.orderService.CreateOrder(uint(userID), uint(templateID), piID)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Mark as paid
	paymentDetails := fmt.Sprintf("Stripe Payment Intent: %s | Status: %s | Amount: $%.2f",
		piID,
		status,
		amount/100,
	)

	if err := h.orderService.MarkOrderAsPaid(order.ID, paymentDetails); err != nil {
		return fmt.Errorf("failed to mark order as paid: %w", err)
	}

	return nil
}

// handlePaymentIntentFailedData processes failed payment intents
func (h *PaymentHandler) handlePaymentIntentFailedData(piID string, metadata map[string]interface{}) error {
	// Try to find existing order
	templateIDStr, ok := metadata["template_id"].(string)
	if !ok {
		return nil // No order to update
	}
	userIDStr, ok := metadata["user_id"].(string)
	if !ok {
		return nil
	}

	templateID, _ := strconv.ParseUint(templateIDStr, 10, 32)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	// Create failed order record
	order, err := h.orderService.CreateOrder(uint(userID), uint(templateID), piID)
	if err != nil {
		return fmt.Errorf("failed to create order record: %w", err)
	}

	// Mark as failed
	reason := fmt.Sprintf("Payment failed: %s", piID)
	if err := h.orderService.MarkOrderAsFailed(order.ID, reason); err != nil {
		return fmt.Errorf("failed to mark order as failed: %w", err)
	}

	return nil
}

// GetOrderByID retrieves an order by ID
func (h *PaymentHandler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, services.ToOrderResponse(order))
}

// GetUserOrders retrieves all orders for a user
func (h *PaymentHandler) GetUserOrders(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	orders, err := h.orderService.GetOrdersByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	// Convert to response format
	responses := make([]services.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = services.ToOrderResponse(&order)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": responses,
		"count":  len(responses),
	})
}

// GetAllOrders retrieves all orders with pagination
func (h *PaymentHandler) GetAllOrders(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	orders, err := h.orderService.GetAllOrders(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	// Convert to response format
	responses := make([]services.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = services.ToOrderResponse(&order)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": responses,
		"count":  len(responses),
		"limit":  limit,
		"offset": offset,
	})
}

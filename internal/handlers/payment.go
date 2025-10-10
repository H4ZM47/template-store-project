package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"template-store/internal/models"
	"template-store/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

// PaymentHandler handles payment-related HTTP requests.
type PaymentHandler struct {
	paymentService      services.PaymentService
	templateService     services.TemplateService
	orderService        services.OrderService
	userService         services.UserService
	stripeWebhookSecret string
}

// NewPaymentHandler creates a new PaymentHandler.
func NewPaymentHandler(paymentService services.PaymentService, templateService services.TemplateService, orderService services.OrderService, userService services.UserService, stripeWebhookSecret string) *PaymentHandler {
	return &PaymentHandler{
		paymentService:      paymentService,
		templateService:     templateService,
		orderService:        orderService,
		userService:         userService,
		stripeWebhookSecret: stripeWebhookSecret,
	}
}

// CheckoutRequest defines the request body for creating a checkout session.
type CheckoutRequest struct {
	TemplateID uint `json:"template_id" binding:"required"`
}

// CreateCheckoutSession handles the creation of a payment intent for a template.
func (h *PaymentHandler) CreateCheckoutSession(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user ID from claims
	claims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User claims not found"})
		return
	}

	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse user claims"})
		return
	}

	cognitoSub, ok := claimsMap["sub"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user ID from claims"})
		return
	}

	// Get the template details to find the price
	template, err := h.templateService.GetTemplate(req.TemplateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Convert price to cents
	amountInCents := int64(template.Price * 100)

	// Prepare metadata
	metadata := map[string]string{
		"cognito_sub": cognitoSub,
		"template_id": strconv.FormatUint(uint64(req.TemplateID), 10),
	}

	// Define success and cancel URLs
	// In a real application, these would be dynamically generated
	successURL := "http://localhost:8080/payment/success?session_id={CHECKOUT_SESSION_ID}"
	cancelURL := "http://localhost:8080/payment/cancel"

	// Create a checkout session with Stripe
	session, err := h.paymentService.CreateCheckoutSession(amountInCents, "usd", metadata, successURL, cancelURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkout_url": session.URL,
	})
}

// StripeWebhook handles incoming webhooks from Stripe.
func (h *PaymentHandler) StripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf("Error reading webhook request body: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading request body"})
		return
	}

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), h.stripeWebhookSecret)
	if err != nil {
		logrus.Errorf("Webhook signature verification failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook signature verification failed"})
		return
	}

	// Handle the event
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			logrus.Errorf("Error parsing webhook JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing webhook JSON"})
			return
		}

		// Extract metadata
		cognitoSub := session.Metadata["cognito_sub"]
		templateIDStr := session.Metadata["template_id"]

		// Find user by Cognito sub
		user, err := h.userService.GetUserByCognitoSub(cognitoSub)
		if err != nil {
			logrus.Errorf("Webhook error: user with cognito_sub %s not found: %v", cognitoSub, err)
			// Still return 200 to Stripe to acknowledge receipt of the event
			c.Status(http.StatusOK)
			return
		}

		templateID, err := strconv.ParseUint(templateIDStr, 10, 32)
		if err != nil {
			logrus.Errorf("Webhook error: invalid template_id %s: %v", templateIDStr, err)
			c.Status(http.StatusOK)
			return
		}

		// Create the order
		order := &models.Order{
			UserID:         user.ID,
			TemplateID:     uint(templateID),
			PurchaseHistory: session.PaymentIntent.ID, // Store payment intent ID for reference
			DeliveryStatus: "Completed",
		}

		if err := h.orderService.CreateOrder(order); err != nil {
			logrus.Errorf("Webhook error: failed to create order: %v", err)
			c.Status(http.StatusOK)
			return
		}

		logrus.Infof("Successfully created order for user %d and template %d", user.ID, templateID)

	default:
		logrus.Warnf("Unhandled Stripe event type: %s", event.Type)
	}

	c.Status(http.StatusOK)
}

// PaymentSuccess handles the success response from stripe
func (h *PaymentHandler) PaymentSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Payment successful!"})
}

// PaymentCancel handles the cancel response from stripe
func (h *PaymentHandler) PaymentCancel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Payment canceled."})
}

package services

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripeService handles Stripe payment operations
type StripeService struct {
	apiKey          string
	webhookSecret   string
	successURL      string
	cancelURL       string
}

// NewStripeService creates a new Stripe service instance
func NewStripeService(apiKey, webhookSecret, successURL, cancelURL string) *StripeService {
	stripe.Key = apiKey
	return &StripeService{
		apiKey:        apiKey,
		webhookSecret: webhookSecret,
		successURL:    successURL,
		cancelURL:     cancelURL,
	}
}

// CreateCheckoutSession creates a Stripe checkout session for template purchase
func (s *StripeService) CreateCheckoutSession(templateID uint, templateName string, price float64, userEmail string) (*stripe.CheckoutSession, error) {
	// Convert price to cents (Stripe uses smallest currency unit)
	priceInCents := int64(price * 100)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(templateName),
						Description: stripe.String(fmt.Sprintf("Digital Template: %s", templateName)),
					},
					UnitAmount: stripe.Int64(priceInCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.cancelURL),
		CustomerEmail: stripe.String(userEmail),
		Metadata: map[string]string{
			"template_id": fmt.Sprintf("%d", templateID),
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess, nil
}

// CreatePaymentIntent creates a payment intent for direct payment
func (s *StripeService) CreatePaymentIntent(amount float64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	// Convert to smallest currency unit (cents for USD)
	amountInCents := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(currency),
		Metadata: metadata,
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return pi, nil
}

// GetPaymentIntent retrieves a payment intent by ID
func (s *StripeService) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}
	return pi, nil
}

// VerifyWebhookSignature verifies the Stripe webhook signature
func (s *StripeService) VerifyWebhookSignature(payload []byte, signature string) (stripe.Event, error) {
	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		return stripe.Event{}, fmt.Errorf("failed to verify webhook signature: %w", err)
	}
	return event, nil
}

// GetCheckoutSession retrieves a checkout session by ID
func (s *StripeService) GetCheckoutSession(sessionID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")

	sess, err := session.Get(sessionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkout session: %w", err)
	}
	return sess, nil
}

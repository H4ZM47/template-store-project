package services

import (
	"errors"
	appconfig "template-store/internal/config"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

// PaymentService defines the interface for payment operations.
type PaymentService interface {
	CreateCheckoutSession(amount int64, currency string, metadata map[string]string, successURL string, cancelURL string) (*stripe.CheckoutSession, error)
}

// stripePaymentService is the concrete implementation of the PaymentService interface.
type stripePaymentService struct {
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(cfg *appconfig.Config) (PaymentService, error) {
	if cfg.Stripe.SecretKey == "" {
		return nil, errors.New("Stripe secret key is not configured")
	}
	stripe.Key = cfg.Stripe.SecretKey
	return &stripePaymentService{}, nil
}

// CreateCheckoutSession creates a new Stripe Checkout Session.
// Amount should be in the smallest currency unit (e.g., cents for USD).
func (s *stripePaymentService) CreateCheckoutSession(amount int64, currency string, metadata map[string]string, successURL string, cancelURL string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					UnitAmount: stripe.Int64(amount),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Template"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: metadata,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

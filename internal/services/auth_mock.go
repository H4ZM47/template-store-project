package services

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	appconfig "template-store/internal/config"
)

// MockAuthService is a mock implementation of AuthService for development.
type MockAuthService struct{}

// NewMockAuthService creates a new MockAuthService.
func NewMockAuthService(cfg *appconfig.Config) (AuthService, error) {
	return &MockAuthService{}, nil
}

// SignUp is a mock implementation of the SignUp method.
func (s *MockAuthService) SignUp(ctx context.Context, email, password, name string) (*cognitoidentityprovider.SignUpOutput, error) {
	return nil, errors.New("SignUp not implemented in mock service")
}

// SignIn is a mock implementation of the SignIn method.
func (s *MockAuthService) SignIn(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return nil, errors.New("SignIn not implemented in mock service")
}

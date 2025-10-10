package services

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	appconfig "template-store/internal/config"
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	SignUp(ctx context.Context, email, password, name string) (*cognitoidentityprovider.SignUpOutput, error)
	SignIn(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error)
}

// CognitoService provides authentication services using AWS Cognito.
type CognitoService struct {
	client   *cognitoidentityprovider.Client
	poolID   string
	clientID string
}

// NewAuthService creates a new CognitoService.
func NewAuthService(cfg *appconfig.Config) (*CognitoService, error) {
	if cfg.AWS.Region == "" || cfg.AWS.CognitoPoolID == "" || cfg.AWS.CognitoAppClientID == "" {
		return nil, errors.New("AWS Cognito configuration is missing")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.AWS.Region))
	if err != nil {
		return nil, err
	}

	return &CognitoService{
		client:   cognitoidentityprovider.NewFromConfig(awsCfg),
		poolID:   cfg.AWS.CognitoPoolID,
		clientID: cfg.AWS.CognitoAppClientID,
	}, nil
}

// SignUp handles user registration in Cognito.
func (s *CognitoService) SignUp(ctx context.Context, email, password, name string) (*cognitoidentityprovider.SignUpOutput, error) {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: &s.clientID,
		Password: &password,
		Username: &email,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("name"),
				Value: &name,
			},
			{
				Name:  aws.String("email"),
				Value: &email,
			},
		},
	}

	result, err := s.client.SignUp(ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SignIn handles user authentication in Cognito.
func (s *CognitoService) SignIn(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserSrpAuth,
		ClientId: &s.clientID,
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}

	result, err := s.client.InitiateAuth(ctx, authInput)
	if err != nil {
		return nil, err
	}

	return result, nil
}

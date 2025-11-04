import {
  CognitoIdentityProviderClient,
  SignUpCommand,
  InitiateAuthCommand,
  AuthFlowType,
  ConfirmSignUpCommand,
  ForgotPasswordCommand,
  ConfirmForgotPasswordCommand,
} from '@aws-sdk/client-cognito-identity-provider';
import { config } from '../config/config';
import { logger } from '../utils/logger';

export interface SignUpParams {
  email: string;
  password: string;
  name: string;
}

export interface SignInParams {
  email: string;
  password: string;
}

export interface AuthResponse {
  accessToken: string;
  idToken: string;
  refreshToken: string;
  expiresIn: number;
}

export class AuthService {
  private client: CognitoIdentityProviderClient;

  constructor() {
    this.client = new CognitoIdentityProviderClient({
      region: config.aws.cognitoRegion,
      credentials: {
        accessKeyId: config.aws.accessKeyId,
        secretAccessKey: config.aws.secretAccessKey,
      },
    });
  }

  async signUp(params: SignUpParams): Promise<{ userSub: string; emailVerificationRequired: boolean }> {
    try {
      const command = new SignUpCommand({
        ClientId: config.aws.cognitoClientId,
        Username: params.email,
        Password: params.password,
        UserAttributes: [
          { Name: 'email', Value: params.email },
          { Name: 'name', Value: params.name },
        ],
      });

      const response = await this.client.send(command);
      logger.info('User signed up successfully', { email: params.email });

      return {
        userSub: response.UserSub!,
        emailVerificationRequired: !response.UserConfirmed,
      };
    } catch (error) {
      logger.error('Sign up error:', error);
      throw error;
    }
  }

  async signIn(params: SignInParams): Promise<AuthResponse> {
    try {
      const command = new InitiateAuthCommand({
        ClientId: config.aws.cognitoClientId,
        AuthFlow: AuthFlowType.USER_PASSWORD_AUTH,
        AuthParameters: {
          USERNAME: params.email,
          PASSWORD: params.password,
        },
      });

      const response = await this.client.send(command);

      if (!response.AuthenticationResult) {
        throw new Error('Authentication failed');
      }

      logger.info('User signed in successfully', { email: params.email });

      return {
        accessToken: response.AuthenticationResult.AccessToken!,
        idToken: response.AuthenticationResult.IdToken!,
        refreshToken: response.AuthenticationResult.RefreshToken!,
        expiresIn: response.AuthenticationResult.ExpiresIn!,
      };
    } catch (error) {
      logger.error('Sign in error:', error);
      throw error;
    }
  }

  async confirmSignUp(email: string, code: string): Promise<void> {
    try {
      const command = new ConfirmSignUpCommand({
        ClientId: config.aws.cognitoClientId,
        Username: email,
        ConfirmationCode: code,
      });

      await this.client.send(command);
      logger.info('User email confirmed', { email });
    } catch (error) {
      logger.error('Confirm sign up error:', error);
      throw error;
    }
  }

  async forgotPassword(email: string): Promise<void> {
    try {
      const command = new ForgotPasswordCommand({
        ClientId: config.aws.cognitoClientId,
        Username: email,
      });

      await this.client.send(command);
      logger.info('Password reset initiated', { email });
    } catch (error) {
      logger.error('Forgot password error:', error);
      throw error;
    }
  }

  async confirmForgotPassword(email: string, code: string, newPassword: string): Promise<void> {
    try {
      const command = new ConfirmForgotPasswordCommand({
        ClientId: config.aws.cognitoClientId,
        Username: email,
        ConfirmationCode: code,
        Password: newPassword,
      });

      await this.client.send(command);
      logger.info('Password reset confirmed', { email });
    } catch (error) {
      logger.error('Confirm forgot password error:', error);
      throw error;
    }
  }
}

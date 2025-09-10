variable "app_name" {
  description = "The name of the application."
  type        = string
  default     = "template-store"
}

variable "image_uri" {
  description = "The URI of the Docker image in ECR."
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC."
  type        = string
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs for the ECS service."
  type        = list(string)
}

variable "ecs_security_group_id" {
  description = "The ID of the security group for the ECS tasks."
  type        = string
}

variable "alb_target_group_arn" {
  description = "The ARN of the Application Load Balancer target group."
  type        = string
}

variable "db_endpoint" {
  description = "The endpoint of the RDS database."
  type        = string
}

variable "db_port" {
  description = "The port of the RDS database."
  type        = number
}

variable "db_credentials_secret_arn" {
  description = "ARN of the Secrets Manager secret for the database credentials."
  type        = string
}

variable "s3_bucket_arn" {
  description = "ARN of the S3 bucket for application assets."
  type        = string
}

variable "s3_bucket_id" {
  description = "ID (name) of the S3 bucket for application assets."
  type        = string
}

variable "aws_region" {
  description = "The AWS region."
  type        = string
}

variable "cognito_user_pool_id" {
  description = "The ID of the Cognito User Pool."
  type        = string
}

variable "cognito_user_pool_client_id" {
  description = "The ID of the Cognito User Pool Client."
  type        = string
}

variable "stripe_secret_key_arn" {
  description = "ARN of the Secrets Manager secret for the Stripe Secret Key."
  type        = string
}

variable "sendgrid_api_key_arn" {
  description = "ARN of the Secrets Manager secret for the SendGrid API Key."
  type        = string
}

variable "container_port" {
  description = "The port the container listens on."
  type        = number
  default     = 8080
}

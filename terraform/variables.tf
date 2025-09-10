variable "aws_region" {
  description = "The AWS region to deploy the resources in."
  type        = string
  default     = "us-east-1"
}

variable "stripe_secret_key" {
  description = "The Stripe secret key."
  type        = string
  sensitive   = true
}

variable "sendgrid_api_key" {
  description = "The SendGrid API key."
  type        = string
  sensitive   = true
}

variable "certificate_arn" {
  description = "The ARN of the ACM certificate for the ALB."
  type        = string
}

variable "domain_name" {
  description = "The domain name for the application."
  type        = string
}

variable "image_tag" {
  description = "The tag for the Docker image in ECR."
  type        = string
  default     = "latest"
}

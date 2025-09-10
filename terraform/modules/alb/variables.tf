variable "app_name" {
  description = "The name of the application."
  type        = string
  default     = "template-store"
}

variable "vpc_id" {
  description = "The ID of the VPC."
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs for the ALB."
  type        = list(string)
}

variable "certificate_arn" {
  description = "The ARN of the ACM certificate for HTTPS."
  type        = string
}

variable "health_check_path" {
  description = "The path for the health check."
  type        = string
  default     = "/health"
}

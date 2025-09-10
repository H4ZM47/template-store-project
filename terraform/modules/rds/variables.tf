variable "vpc_id" {
  description = "The ID of the VPC where the RDS instance will be deployed."
  type        = string
}

variable "private_subnet_ids" {
  description = "A list of private subnet IDs for the RDS DB subnet group."
  type        = list(string)
}

variable "ecs_security_group_id" {
  description = "The ID of the security group for the ECS service that needs access to the RDS instance."
  type        = string
}

variable "db_name" {
  description = "The name of the database."
  type        = string
  default     = "template_store"
}

variable "db_username" {
  description = "The username for the database."
  type        = string
  default     = "postgres"
}

variable "db_instance_class" {
  description = "The instance class for the RDS instance."
  type        = string
  default     = "db.t3.micro"
}

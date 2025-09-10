output "db_instance_endpoint" {
  description = "The endpoint of the RDS instance."
  value       = aws_db_instance.main.endpoint
}

output "db_instance_port" {
  description = "The port of the RDS instance."
  value       = aws_db_instance.main.port
}

output "db_credentials_secret_arn" {
  description = "The ARN of the Secrets Manager secret for the database credentials."
  value       = aws_secretsmanager_secret.db_credentials.arn
}

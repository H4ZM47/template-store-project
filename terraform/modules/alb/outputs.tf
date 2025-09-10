output "alb_dns_name" {
  description = "The DNS name of the ALB."
  value       = aws_lb.main.dns_name
}

output "target_group_arn" {
  description = "The ARN of the target group."
  value       = aws_lb_target_group.main.arn
}

output "https_listener_arn" {
  description = "The ARN of the HTTPS listener."
  value       = aws_lb_listener.https.arn
}

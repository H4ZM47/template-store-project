output "ecr_repository_url" {
  description = "The URL of the ECR repository."
  value       = aws_ecr_repository.app.repository_url
}

output "alb_dns_name" {
  description = "The DNS name of the ALB."
  value       = module.alb.alb_dns_name
}

output "cloudfront_domain_name" {
  description = "The domain name of the CloudFront distribution."
  value       = module.cloudfront.distribution_domain_name
}

output "s3_bucket_name" {
  description = "The name of the S3 bucket."
  value       = module.s3.bucket_id
}

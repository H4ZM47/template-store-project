variable "s3_bucket_domain_name" {
  description = "The domain name of the S3 bucket."
  type        = string
}

variable "s3_bucket_id" {
  description = "The ID of the S3 bucket."
  type        = string
}

variable "cloudfront_oai_path" {
  description = "The path of the CloudFront Origin Access Identity."
  type        = string
}

variable "certificate_arn" {
  description = "The ARN of the ACM certificate for the custom domain."
  type        = string
}

variable "aliases" {
  description = "A list of CNAMEs for the distribution."
  type        = list(string)
  default     = []
}

variable "bucket_name" {
  description = "The name of the S3 bucket. Must be globally unique."
  type        = string
}

variable "tags" {
  description = "A map of tags to assign to the bucket."
  type        = map(string)
  default     = {}
}

variable "cloudfront_oai_arn" {
  description = "The ARN of the CloudFront Origin Access Identity."
  type        = string
  default     = ""
}

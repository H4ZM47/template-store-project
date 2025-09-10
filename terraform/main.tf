terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

module "vpc" {
  source = "./modules/vpc"

  aws_region = var.aws_region
}

module "rds" {
  source = "./modules/rds"

  vpc_id                = module.vpc.vpc_id
  private_subnet_ids    = [module.vpc.private_subnet_id]
  ecs_security_group_id = aws_security_group.ecs.id
}

module "s3" {
  source = "./modules/s3"

  bucket_name        = "template-store-assets-${random_id.main.hex}"
  cloudfront_oai_arn = aws_cloudfront_origin_access_identity.main.iam_arn
  tags = {
    Name = "template-store-assets"
  }
}

module "cloudfront" {
  source = "./modules/cloudfront"

  s3_bucket_domain_name = module.s3.bucket_domain_name
  s3_bucket_id          = module.s3.bucket_id
  certificate_arn       = var.certificate_arn
  cloudfront_oai_path   = aws_cloudfront_origin_access_identity.main.cloudfront_access_identity_path
  aliases               = [var.domain_name]
}

module "alb" {
  source = "./modules/alb"

  app_name          = "template-store"
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = [module.vpc.public_subnet_id]
  certificate_arn   = var.certificate_arn
}

module "ecs" {
  source = "./modules/ecs"

  app_name                      = "template-store"
  image_uri                     = "${aws_ecr_repository.app.repository_url}:${var.image_tag}"
  vpc_id                        = module.vpc.vpc_id
  private_subnet_ids            = [module.vpc.private_subnet_id]
  ecs_security_group_id         = aws_security_group.ecs.id
  alb_target_group_arn          = module.alb.target_group_arn
  db_endpoint                   = module.rds.db_instance_endpoint
  db_port                       = module.rds.db_instance_port
  db_credentials_secret_arn     = module.rds.db_credentials_secret_arn
  s3_bucket_arn                 = module.s3.bucket_arn
  s3_bucket_id                  = module.s3.bucket_id
  aws_region                    = var.aws_region
  cognito_user_pool_id          = module.cognito.user_pool_id
  cognito_user_pool_client_id   = module.cognito.user_pool_client_id
  stripe_secret_key_arn         = aws_secretsmanager_secret.stripe_secret_key.arn
  sendgrid_api_key_arn          = aws_secretsmanager_secret.sendgrid_api_key.arn
}

module "cognito" {
  source = "./modules/cognito"

  app_name = "template-store"
}

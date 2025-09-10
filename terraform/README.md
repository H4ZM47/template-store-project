# Terraform for Template Store

This directory contains the Terraform code to deploy the Template Store application to AWS.

## Prerequisites

- Terraform v1.0 or higher
- AWS CLI
- An AWS account with the necessary permissions

## How to Use

1. **Initialize Terraform:**
   ```bash
   terraform init
   ```

2. **Build and Push the Docker Image to ECR:**

   Before you can apply the Terraform configuration, you need to build the Docker image for the application and push it to the ECR repository that Terraform will create.

   a. **Authenticate Docker with ECR:**
      Replace `your-aws-account-id` and `your-aws-region` with your actual AWS account ID and region.
      ```bash
      aws ecr get-login-password --region your-aws-region | docker login --username AWS --password-stdin your-aws-account-id.dkr.ecr.your-aws-region.amazonaws.com
      ```

   b. **Build the Docker image:**
      From the root of the project, run:
      ```bash
      docker build -t template-store:latest .
      ```

   c. **Tag the Docker image:**
      Replace `your-aws-account-id` and `your-aws-region` with your actual AWS account ID and region. You can change `latest` to a different tag if you want.
      ```bash
      docker tag template-store:latest your-aws-account-id.dkr.ecr.your-aws-region.amazonaws.com/template-store:latest
      ```

   d. **Push the Docker image to ECR:**
      ```bash
      docker push your-aws-account-id.dkr.ecr.your-aws-region.amazonaws.com/template-store:latest
      ```

3. **Apply the Terraform Configuration:**

   You will need to provide values for the following variables. You can do this by creating a `terraform.tfvars` file or by passing them as command-line arguments.

   - `stripe_secret_key`: Your Stripe secret key.
   - `sendgrid_api_key`: Your SendGrid API key.
   - `certificate_arn`: The ARN of the ACM certificate for your domain.
   - `domain_name`: Your domain name (e.g., `www.example.com`).

   ```bash
   terraform apply
   ```

   This will provision all the necessary AWS resources.

4. **Destroy the Infrastructure:**
   When you are done, you can destroy all the resources created by Terraform:
   ```bash
   terraform destroy
   ```

resource "aws_ecr_repository" "app" {
  name                 = "template-store"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "template-store-ecr"
  }
}

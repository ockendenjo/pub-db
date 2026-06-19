resource "aws_s3_bucket" "static_files" {
  bucket_prefix = "pub-db-data-${var.env}-"
  force_destroy = true
}

resource "aws_s3_bucket_policy" "static_read_ew1" {
  bucket = aws_s3_bucket.static_files.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowCloudFrontRead"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action = "s3:GetObject"
        Resource = [
          aws_s3_bucket.static_files.arn,
          "${aws_s3_bucket.static_files.arn}/*"
        ]
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = [for d in var.allowed_cfront_distros : "arn:aws:cloudfront::${var.aws_account_id}:distribution/${d}"]
          }
        }
      }
    ]
  })
}

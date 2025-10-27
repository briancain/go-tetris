#!/bin/bash

# Bootstrap script to create Terraform remote state infrastructure
# Run this once before any terraform commands

set -e

# Use configured AWS profile and region
PROFILE=${AWS_PROFILE:-$(aws configure get profile 2>/dev/null || echo "default")}
REGION=$(aws configure get region --profile "$PROFILE" 2>/dev/null || echo "us-west-2")
BUCKET_NAME="tetris-terraform-state-$(date +%s)"
DYNAMODB_TABLE="tetris-terraform-locks"

echo "🚀 Bootstrapping Terraform remote state infrastructure..."
echo "Profile: $PROFILE"
echo "Region: $REGION"
echo "Bucket: $BUCKET_NAME"
echo "DynamoDB Table: $DYNAMODB_TABLE"

# Create S3 bucket for Terraform state
echo "📦 Creating S3 bucket for Terraform state..."
aws s3api create-bucket \
    --bucket "$BUCKET_NAME" \
    --region "$REGION" \
    --create-bucket-configuration LocationConstraint="$REGION" \
    --profile "$PROFILE"

# Enable versioning on the bucket
echo "🔄 Enabling versioning on S3 bucket..."
aws s3api put-bucket-versioning \
    --bucket "$BUCKET_NAME" \
    --versioning-configuration Status=Enabled \
    --profile "$PROFILE"

# Enable server-side encryption
echo "🔒 Enabling server-side encryption on S3 bucket..."
aws s3api put-bucket-encryption \
    --bucket "$BUCKET_NAME" \
    --server-side-encryption-configuration '{
        "Rules": [
            {
                "ApplyServerSideEncryptionByDefault": {
                    "SSEAlgorithm": "AES256"
                }
            }
        ]
    }' \
    --profile "$PROFILE"

# Block public access
echo "🚫 Blocking public access on S3 bucket..."
aws s3api put-public-access-block \
    --bucket "$BUCKET_NAME" \
    --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true" \
    --profile "$PROFILE"

# Create DynamoDB table for state locking
echo "🔐 Creating DynamoDB table for state locking..."
aws dynamodb create-table \
    --table-name "$DYNAMODB_TABLE" \
    --attribute-definitions AttributeName=LockID,AttributeType=S \
    --key-schema AttributeName=LockID,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region "$REGION" \
    --profile "$PROFILE"

# Wait for table to be active
echo "⏳ Waiting for DynamoDB table to be active..."
aws dynamodb wait table-exists \
    --table-name "$DYNAMODB_TABLE" \
    --region "$REGION" \
    --profile "$PROFILE"

# Create terraform backend configuration file
echo "📝 Creating Terraform backend configuration..."
cat > ../backend.tf << EOF
terraform {
  backend "s3" {
    bucket         = "$BUCKET_NAME"
    key            = "terraform.tfstate"
    region         = "$REGION"
    dynamodb_table = "$DYNAMODB_TABLE"
    encrypt        = true
    profile        = "$PROFILE"
  }
}
EOF

echo "✅ Bootstrap complete!"
echo ""
echo "📋 Next steps:"
echo "1. Run 'terraform init' to initialize the backend"
echo "2. Proceed with your Terraform deployment"
echo ""
echo "🗑️  To cleanup later, run:"
echo "   aws s3 rb s3://$BUCKET_NAME --force --profile $PROFILE"
echo "   aws dynamodb delete-table --table-name $DYNAMODB_TABLE --profile $PROFILE"

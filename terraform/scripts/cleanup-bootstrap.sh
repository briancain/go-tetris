#!/bin/bash

# Cleanup script to remove Terraform remote state infrastructure
# Use this to rollback the bootstrap if needed

set -e

# Use configured AWS profile and region
PROFILE=${AWS_PROFILE:-$(aws configure get profile 2>/dev/null || echo "default")}
REGION=$(aws configure get region --profile "$PROFILE" 2>/dev/null || echo "us-west-2")
DYNAMODB_TABLE="tetris-terraform-locks"

echo "🧹 Cleaning up Terraform remote state infrastructure..."
echo "Profile: $PROFILE"
echo "Region: $REGION"

# Find the S3 bucket (it has a timestamp suffix)
echo "🔍 Finding Terraform state bucket..."
BUCKET_NAME=$(aws s3api list-buckets --profile "$PROFILE" --query "Buckets[?starts_with(Name, 'tetris-terraform-state-')].Name" --output text)

if [ -z "$BUCKET_NAME" ]; then
    echo "❌ No Terraform state bucket found"
else
    echo "📦 Found bucket: $BUCKET_NAME"
    
    # Delete all objects in the bucket
    echo "🗑️  Deleting all objects in bucket..."
    aws s3 rm "s3://$BUCKET_NAME" --recursive --profile "$PROFILE" || true
    
    # Delete all object versions
    echo "🗑️  Deleting all object versions..."
    aws s3api delete-objects --bucket "$BUCKET_NAME" --delete "$(aws s3api list-object-versions --bucket "$BUCKET_NAME" --profile "$PROFILE" --output json --query '{Objects: Versions[].{Key:Key,VersionId:VersionId}}')" --profile "$PROFILE" 2>/dev/null || true
    
    # Delete all delete markers
    echo "🗑️  Deleting all delete markers..."
    aws s3api delete-objects --bucket "$BUCKET_NAME" --delete "$(aws s3api list-object-versions --bucket "$BUCKET_NAME" --profile "$PROFILE" --output json --query '{Objects: DeleteMarkers[].{Key:Key,VersionId:VersionId}}')" --profile "$PROFILE" 2>/dev/null || true
    
    # Delete the bucket
    echo "🗑️  Deleting S3 bucket..."
    aws s3api delete-bucket --bucket "$BUCKET_NAME" --profile "$PROFILE"
    echo "✅ S3 bucket deleted"
fi

# Delete DynamoDB table
echo "🗑️  Deleting DynamoDB table..."
aws dynamodb delete-table --table-name "$DYNAMODB_TABLE" --region "$REGION" --profile "$PROFILE" 2>/dev/null || echo "❌ DynamoDB table not found or already deleted"

# Remove backend.tf file
if [ -f "backend.tf" ]; then
    echo "🗑️  Removing backend.tf..."
    rm backend.tf
    echo "✅ backend.tf removed"
fi

echo "✅ Cleanup complete!"

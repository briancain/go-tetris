#!/bin/bash

# Deploy WebAssembly client to S3 and invalidate CloudFront cache

set -e

# Get bucket name and CloudFront distribution ID from Terraform outputs
BUCKET_NAME=$(terraform output -raw website_bucket_name)
DISTRIBUTION_ID=$(terraform output -raw cloudfront_distribution_id)

echo "🚀 Deploying WebAssembly client..."
echo "Bucket: $BUCKET_NAME"
echo "Distribution: $DISTRIBUTION_ID"

# Build the WebAssembly client
echo "🔨 Building WebAssembly client..."
cd ..
make build-web

# Upload to S3
echo "📦 Uploading to S3..."
aws s3 sync bin/web/ "s3://$BUCKET_NAME/" --delete

# Invalidate CloudFront cache for immediate updates
echo "🔄 Invalidating CloudFront cache..."
aws cloudfront create-invalidation \
    --distribution-id "$DISTRIBUTION_ID" \
    --paths "/*" \
    --query 'Invalidation.Id' \
    --output text

echo "✅ Deployment complete!"
echo "🌐 Website URL: $(cd terraform && terraform output -raw website_url)"

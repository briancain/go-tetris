#!/bin/bash

# Deploy WebAssembly client to S3 and invalidate CloudFront cache

set -e

# Change to terraform directory to get outputs
cd "$(dirname "$0")/.."

# Get bucket name, CloudFront distribution ID, and ALB URL from Terraform outputs
BUCKET_NAME=$(terraform output -raw website_bucket_name)
DISTRIBUTION_ID=$(terraform output -raw cloudfront_distribution_id)
ALB_DNS_NAME=$(terraform output -raw alb_dns_name)
SERVER_URL="http://${ALB_DNS_NAME}"

echo "🚀 Deploying WebAssembly client..."
echo "Bucket: $BUCKET_NAME"
echo "Distribution: $DISTRIBUTION_ID"
echo "Server URL: $SERVER_URL"

# Build the WebAssembly client
echo "🔨 Building WebAssembly client..."
cd ..
make build-web

# Replace server URL template in index.html
echo "🔧 Configuring server URL..."
sed -i.bak "s|{{SERVER_URL}}|${SERVER_URL}|g" bin/web/index.html
rm bin/web/index.html.bak

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
echo "🔗 Server URL: $SERVER_URL"

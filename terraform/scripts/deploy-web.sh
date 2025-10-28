#!/bin/bash

# Deploy WebAssembly client to S3 and invalidate CloudFront cache

set -e

# Change to terraform directory to get outputs
cd "$(dirname "$0")/.."

# Get bucket name, CloudFront distribution ID, and server URL from Terraform outputs
BUCKET_NAME=$(terraform output -raw website_bucket_name)
DISTRIBUTION_ID=$(terraform output -raw cloudfront_distribution_id)
CLOUDFRONT_URL=$(terraform output -raw website_url)
SERVER_URL=$(terraform output -raw api_url)
SSL_ENABLED=$(terraform output -raw ssl_enabled)

echo "🚀 Deploying WebAssembly client..."
echo "Bucket: $BUCKET_NAME"
echo "Distribution: $DISTRIBUTION_ID"
echo "CloudFront URL: $CLOUDFRONT_URL"
echo "Server URL: $SERVER_URL"
echo "SSL Enabled: $SSL_ENABLED"

# Check if SSL is required for production deployment
if [ "$SSL_ENABLED" = "false" ]; then
    echo "⚠️  WARNING: SSL is disabled!"
    echo "   CloudFront serves HTTPS but ALB uses HTTP"
    echo "   WebSocket connections will fail due to mixed content policy"
    echo "   This deployment is for testing HTTP API calls only"
    echo ""
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Deployment cancelled"
        exit 1
    fi
fi

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
echo "🌐 Website URL: $CLOUDFRONT_URL"
echo "🔗 API URL: $SERVER_URL"
echo "🔌 WebSocket URL: $(terraform output -raw websocket_url)"

if [ "$SSL_ENABLED" = "true" ]; then
    echo "🔒 SSL/HTTPS: Enabled - Full functionality available"
else
    echo "⚠️  SSL/HTTPS: Disabled - WebSocket connections will fail from HTTPS site"
fi

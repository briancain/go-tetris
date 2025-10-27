#!/bin/bash

# Build and push tetris-server image to ECR
# Usage: ./push-image.sh [tag]

set -e

# Configuration
REGION="us-west-2"
TAG=${1:-latest}

echo "🚀 Building and pushing tetris-server image to ECR..."
echo "Region: $REGION"
echo "Tag: $TAG"

# Get ECR repository URL from Terraform
cd "$(dirname "$0")/.."
ECR_URL=$(terraform output -raw ecr_repository_url)

if [ -z "$ECR_URL" ]; then
    echo "❌ Error: Could not get ECR repository URL from Terraform"
    exit 1
fi

echo "ECR Repository: $ECR_URL"

# Authenticate with ECR
echo "🔐 Authenticating with ECR..."
aws ecr get-login-password --region "$REGION" | finch login --username AWS --password-stdin "$ECR_URL"

# Build image
echo "🔨 Building Docker image..."
cd ..
finch build -t "go-tetris-server:$TAG" .

# Tag for ECR
echo "🏷️  Tagging image for ECR..."
finch tag "go-tetris-server:$TAG" "$ECR_URL:$TAG"

# Push to ECR
echo "📤 Pushing image to ECR..."
finch push "$ECR_URL:$TAG"

echo "✅ Image pushed successfully!"
echo "Image: $ECR_URL:$TAG"

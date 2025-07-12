#!/bin/bash

# Heroku Deployment Script for Skill Swap Backend
# Make sure you have heroku CLI installed and are logged in

set -e

echo "🚀 Starting Heroku deployment for Skill Swap Backend..."

# Check if Heroku CLI is installed
if ! command -v heroku &> /dev/null; then
    echo "❌ Heroku CLI is not installed. Please install it first."
    exit 1
fi

# Check if logged in to Heroku
if ! heroku auth:whoami &> /dev/null; then
    echo "❌ You are not logged in to Heroku. Please run 'heroku login' first."
    exit 1
fi

# Prompt for app name
read -p "Enter your Heroku app name (or press Enter to create a new one): " APP_NAME

if [ -z "$APP_NAME" ]; then
    echo "📝 Creating a new Heroku app..."
    heroku create --region us
    APP_NAME=$(heroku info --json | grep -o '"name":"[^"]*' | grep -o '[^"]*$')
    echo "✅ Created app: $APP_NAME"
else
    echo "📝 Using existing app: $APP_NAME"
fi

# Set the stack to container
echo "🐳 Setting stack to container..."
heroku stack:set container -a $APP_NAME

# Add PostgreSQL addon
echo "🗄️ Adding PostgreSQL addon..."
heroku addons:create heroku-postgresql:essential-0 -a $APP_NAME || echo "PostgreSQL addon might already exist"

# Generate a secure JWT secret
JWT_SECRET=$(openssl rand -base64 32)
echo "🔐 Setting JWT secret..."
heroku config:set JWT_SECRET="$JWT_SECRET" -a $APP_NAME

# Set production BASE_URL
echo "🌐 Setting BASE_URL..."
heroku config:set BASE_URL="https://$APP_NAME.herokuapp.com" -a $APP_NAME

# Set GIN_MODE to release
echo "⚙️ Setting GIN_MODE to release..."
heroku config:set GIN_MODE=release -a $APP_NAME

# Deploy the application
echo "🚀 Deploying to Heroku..."
git add .
git commit -m "Deploy to Heroku" || echo "No changes to commit"
heroku git:remote -a $APP_NAME
git push heroku main || git push heroku master

echo "✅ Deployment completed!"
echo "🌐 Your app is available at: https://$APP_NAME.herokuapp.com"
echo "📋 API Base URL: https://$APP_NAME.herokuapp.com/api/v1"
echo "❤️ Health Check: https://$APP_NAME.herokuapp.com/health"

# Show app info
echo "📊 App Information:"
heroku info -a $APP_NAME

echo "🎉 Deployment successful! Your Skill Swap backend is now live on Heroku!"

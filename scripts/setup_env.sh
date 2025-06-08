#!/bin/bash

# Environment Setup Script for NutrieDiet Go App
# This file helps you set up environment variables securely

echo "🔧 Setting up SMTP Environment Variables"
echo "========================================"

# Check if .env file exists
if [ -f ".env" ]; then
    echo "⚠️  .env file already exists. Backing up to .env.backup"
    cp .env .env.backup
fi

# Create .env file
cat > .env << 'EOF'
# SMTP Configuration for Password Reset
# IMPORTANT: Keep this file secure and never commit to version control

# Gmail SMTP Settings (recommended)
SMTP_EMAIL=your_app_email@gmail.com
SMTP_PASSWORD=your_16_character_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# Alternative: SendGrid (for production)
# SMTP_EMAIL=apikey
# SMTP_PASSWORD=your_sendgrid_api_key
# SMTP_HOST=smtp.sendgrid.net
# SMTP_PORT=587

# Alternative: AWS SES (for production)
# SMTP_EMAIL=your_aws_access_key
# SMTP_PASSWORD=your_aws_secret_key
# SMTP_HOST=email-smtp.us-east-1.amazonaws.com
# SMTP_PORT=587
EOF

echo "✅ Created .env file with template"
echo ""
echo "📝 Next Steps:"
echo "1. Edit .env file with your actual credentials"
echo "2. Ensure .env is in your .gitignore file"
echo "3. For production, use a secure secret management system"
echo ""
echo "🔐 Security Reminder:"
echo "- Never commit .env to version control"
echo "- Use a dedicated email account for your app"
echo "- Generate Gmail App Password (not regular password)" 
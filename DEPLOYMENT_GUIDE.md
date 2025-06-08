# 🚀 Production Deployment Guide

## 🔐 **Environment Variables in Production**

**Important**: `.env` files are for development only! Never deploy them to production.

## 📡 **Netlify Deployment**

### **Step 1: Set Environment Variables**
1. Go to [Netlify Dashboard](https://app.netlify.com)
2. Select your site
3. Navigate to **Site Settings** → **Environment Variables**
4. Add these variables:

```
SMTP_EMAIL=nutriediet.help@gmail.com
SMTP_PASSWORD=beti ourj bntk hkmy
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

### **Step 2: Build Settings**
```toml
# netlify.toml
[build]
  command = "go build -o main ."
  publish = "."

[build.environment]
  GO_VERSION = "1.20"
```

### **Step 3: Deploy**
```bash
# Push to GitHub (without .env file)
git add .
git commit -m "Deploy with environment variables"
git push origin main
```

## 🏗️ **Other Platforms**

### **Heroku**
```bash
# Set environment variables
heroku config:set SMTP_EMAIL=nutriediet.help@gmail.com
heroku config:set SMTP_PASSWORD="beti ourj bntk hkmy"
heroku config:set SMTP_HOST=smtp.gmail.com
heroku config:set SMTP_PORT=587

# Deploy
git push heroku main
```

### **Vercel**
```bash
# Set environment variables
vercel env add SMTP_EMAIL
vercel env add SMTP_PASSWORD
vercel env add SMTP_HOST
vercel env add SMTP_PORT

# Deploy
vercel --prod
```

### **Railway**
```bash
# Set environment variables
railway variables set SMTP_EMAIL=nutriediet.help@gmail.com
railway variables set SMTP_PASSWORD="beti ourj bntk hkmy"

# Deploy
railway up
```

## 🎯 **Production Recommendations**

### **1. Use Professional Email Service**
For production, upgrade from Gmail to a professional service:

#### **SendGrid (Recommended)**
```bash
SMTP_EMAIL=apikey
SMTP_PASSWORD=SG.your_sendgrid_api_key
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
```

#### **AWS SES**
```bash
SMTP_EMAIL=your_access_key
SMTP_PASSWORD=your_secret_key
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
```

#### **Mailgun**
```bash
SMTP_EMAIL=postmaster@mg.yourdomain.com
SMTP_PASSWORD=your_mailgun_password
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
```

### **2. Environment-Specific Configs**

#### **Development**
```bash
SMTP_EMAIL=nutriediet-dev@gmail.com
SMTP_PASSWORD=dev_app_password
```

#### **Staging**
```bash
SMTP_EMAIL=nutriediet-staging@gmail.com
SMTP_PASSWORD=staging_app_password
```

#### **Production**
```bash
SMTP_EMAIL=apikey
SMTP_PASSWORD=SG.live_sendgrid_key
SMTP_HOST=smtp.sendgrid.net
```

## 🔒 **Security Best Practices**

### **✅ Do's**
- Set environment variables in hosting platform
- Use different credentials for each environment
- Rotate credentials regularly
- Monitor failed authentication attempts
- Use professional email services for production

### **❌ Don'ts**
- Never commit `.env` files
- Never put secrets in source code
- Never use development credentials in production
- Never share credentials in chat/email

## 🧪 **Testing Production Environment**

### **1. Verify Environment Variables**
Add a health check endpoint:

```go
router.GET("/health", func(c *gin.Context) {
    smtpConfigured := os.Getenv("SMTP_EMAIL") != ""
    c.JSON(200, gin.H{
        "status": "healthy",
        "smtp_configured": smtpConfigured,
    })
})
```

### **2. Test Email Functionality**
After deployment, test the forgot password flow with a real email.

## 📊 **Monitoring & Alerts**

### **1. Email Delivery Monitoring**
```go
// Add logging for email sends
log.Printf("OTP email sent to %s at %v", email, time.Now())
```

### **2. Failed Authentication Tracking**
```go
// Log failed OTP attempts
log.Printf("Invalid OTP attempt for %s", email)
```

## 🚨 **Troubleshooting**

### **Environment Variables Not Loading**
```bash
# Check if variables are set
curl https://your-app.netlify.app/health
```

### **Email Not Sending in Production**
1. Verify environment variables are set
2. Check SMTP credentials are correct
3. Ensure production email service is configured
4. Check application logs for errors

### **CORS Issues**
```go
// Update CORS for production domain
config := cors.Config{
    AllowOrigins: []string{
        "https://your-frontend.netlify.app",
        "https://nutriediet.netlify.app",
    },
}
``` 
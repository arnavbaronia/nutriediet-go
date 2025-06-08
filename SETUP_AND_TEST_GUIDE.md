# 🚀 Complete Setup & Testing Guide

## 📧 **Step 1: Create Dedicated Gmail Account**

### **🧠 Why a Separate Account?**
- **Security**: Isolates app access from personal email
- **Control**: Easy to revoke without affecting personal communications
- **Monitoring**: Track app-related emails separately
- **Professional**: Clean separation of concerns

### **Setup Process:**
1. **Create New Gmail Account**
   - Go to [accounts.google.com](https://accounts.google.com)
   - Create account like: `nutriediet-app@gmail.com`
   - Use a strong, unique password

2. **Enable 2-Factor Authentication**
   - Go to Google Account → Security
   - Turn on 2-Step Verification
   - Add your phone number

3. **Generate App Password**
   - Go to Security → 2-Step Verification → App passwords
   - Select "Mail" as the app
   - Copy the 16-character password (e.g., `abcd efgh ijkl mnop`)

## 🔧 **Step 2: Configure Environment Variables**

### **Option A: Using Our Script (Recommended)**
```bash
# Run the setup script
./scripts/setup_env.sh

# Edit the .env file with your credentials
nano .env
```

### **Option B: Manual Setup**
```bash
# Create .env file
touch .env

# Add your credentials (replace with actual values)
echo "SMTP_EMAIL=nutriediet-app@gmail.com" >> .env
echo "SMTP_PASSWORD=abcd efgh ijkl mnop" >> .env
echo "SMTP_HOST=smtp.gmail.com" >> .env
echo "SMTP_PORT=587" >> .env
```

### **Option C: System Environment Variables**
```bash
# For temporary testing
export SMTP_EMAIL="nutriediet-app@gmail.com"
export SMTP_PASSWORD="abcd efgh ijkl mnop"
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"
```

## 🏗️ **Step 3: Build and Run**

```bash
# Build the application
go build -o nutriediet-go .

# Run the application
./nutriediet-go
```

You should see:
```
✅ Environment variables loaded from .env file
```

## 🧪 **Step 4: Test the Functionality**

### **Option A: Using Our Test Script (Recommended)**
```bash
# Run the interactive test script
./scripts/test_password_reset.sh
```

### **Option B: Manual cURL Testing**

#### Test 1: Request OTP
```bash
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "your-test-user@example.com"}'
```

#### Test 2: Reset Password
```bash
curl -X POST http://localhost:8080/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-test-user@example.com",
    "otp": "123456",
    "new_password": "newStrongPassword123"
  }'
```

## 🔐 **Production Security Recommendations**

### **For Production Deployment:**

1. **Use Professional Email Service**
   ```bash
   # SendGrid (recommended)
   SMTP_EMAIL=apikey
   SMTP_PASSWORD=your_sendgrid_api_key
   SMTP_HOST=smtp.sendgrid.net
   SMTP_PORT=587
   ```

2. **Use Secret Management**
   - AWS Secrets Manager
   - HashiCorp Vault
   - Kubernetes Secrets
   - Docker Secrets

3. **Environment-Specific Configs**
   ```bash
   # Development
   .env.development
   
   # Staging
   .env.staging
   
   # Production
   .env.production
   ```

## 🚨 **Troubleshooting**

### **Email Not Sending**
```bash
# Check environment variables
echo $SMTP_EMAIL
echo $SMTP_PASSWORD

# Test SMTP connection
telnet smtp.gmail.com 587
```

### **Common Errors & Solutions**

| Error | Solution |
|-------|----------|
| "SMTP credentials not configured" | Check .env file exists and has correct values |
| "Authentication failed" | Verify App Password (not regular password) |
| "User not found" | Ensure user exists in database |
| "OTP expired" | Request new OTP (5-minute expiry) |

### **Database Issues**
```bash
# Check if table exists
mysql -u your_user -p -e "DESCRIBE defaultdb.password_otps;"

# Run migration if needed
go run migrate/migrate.go
```

## 📊 **Testing Checklist**

- [ ] Environment variables loaded correctly
- [ ] Server starts without errors
- [ ] Forgot password endpoint accepts valid email
- [ ] Email received with 6-digit OTP
- [ ] Reset password endpoint accepts valid OTP
- [ ] Password successfully updated in database
- [ ] Login works with new password
- [ ] OTP expires after 5 minutes
- [ ] Invalid OTP rejected
- [ ] Non-existent user handled gracefully

## 🎯 **Next Steps**

1. **Rate Limiting**: Add middleware to prevent abuse
2. **Email Templates**: Create branded HTML templates
3. **Monitoring**: Add logging and metrics
4. **SMS Backup**: Alternative OTP delivery method
5. **Admin Dashboard**: Monitor password reset requests

## 🔒 **Security Reminders**

- ✅ Never commit `.env` files to version control
- ✅ Use dedicated email account for app
- ✅ Generate App Passwords, not regular passwords
- ✅ Monitor failed authentication attempts
- ✅ Consider rate limiting in production
- ✅ Use HTTPS in production
- ✅ Regularly rotate SMTP credentials 
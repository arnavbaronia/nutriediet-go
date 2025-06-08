# Testing Password Reset Functionality

## Prerequisites

1. Set up SMTP environment variables (see SMTP_SETUP.md)
2. Start your Go server
3. Ensure you have a user account in the database

## API Endpoints

### 1. Request Password Reset OTP

**Endpoint**: `POST /auth/forgot-password`

**Request Body**:
```json
{
  "email": "user@example.com"
}
```

**Success Response** (200):
```json
{
  "message": "OTP sent successfully to your email address",
  "email": "user@example.com"
}
```

**Error Responses**:
- 400: Invalid email format or user not found
- 500: Server error (check SMTP configuration)

### 2. Reset Password with OTP

**Endpoint**: `POST /auth/reset-password`

**Request Body**:
```json
{
  "email": "user@example.com",
  "otp": "123456",
  "new_password": "newStrongPassword123"
}
```

**Success Response** (200):
```json
{
  "message": "Password reset successfully",
  "email": "user@example.com"
}
```

**Error Responses**:
- 400: Invalid OTP, expired OTP, or validation errors
- 500: Server error

## cURL Examples

### Request OTP:
```bash
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

### Reset Password:
```bash
curl -X POST http://localhost:8080/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "123456",
    "new_password": "newStrongPassword123"
  }'
```

## Security Features Implemented

✅ **OTP Hashing**: OTPs are stored as bcrypt hashes, never in plain text
✅ **Time Expiry**: OTPs expire after 5 minutes
✅ **User Validation**: Ensures user exists before sending OTP
✅ **Password Security**: New passwords are hashed with bcrypt
✅ **Cleanup**: Used/expired OTPs are automatically removed
✅ **Rate Limiting**: UPSERT prevents spam (one active OTP per email)
✅ **Input Validation**: Email format and password strength requirements

## Common Issues & Solutions

### Email Not Sending
- Check SMTP environment variables
- Verify Gmail app password is correct
- Ensure 2FA is enabled for Gmail

### OTP Expired Error
- OTPs expire in 5 minutes
- Request a new OTP if expired

### User Not Found
- Ensure the email address has a registered account
- Check the correct user_type if applicable 
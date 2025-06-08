# SMTP Configuration for Password Reset

## Environment Variables Required

Set the following environment variables for email functionality:

```bash
export SMTP_EMAIL=nutriediet.help@gmail.com
export SMTP_PASSWORD=beti ourj bntk hkmy
export SMTP_HOST=smtp.gmail.com
export SMTP_PORT=587
```

## Gmail Setup (Recommended)

1. **Enable 2-Factor Authentication** on your Gmail account
2. **Generate an App Password**:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate a password for "Mail"
   - Use this app password (not your regular password) in `SMTP_PASSWORD`

## Alternative SMTP Providers

- **Outlook**: `smtp-mail.outlook.com:587`
- **Yahoo**: `smtp.mail.yahoo.com:587`
- **Custom SMTP**: Use your provider's settings

## Testing the Setup

You can test the email functionality by:
1. Setting the environment variables
2. Starting your Go server
3. Making a POST request to `/auth/forgot-password`

## Security Notes

- Never commit SMTP credentials to version control
- Use environment variables or a secure config management system
- Consider using services like SendGrid or AWS SES for production 
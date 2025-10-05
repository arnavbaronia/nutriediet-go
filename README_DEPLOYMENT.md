# NutrieDiet Go - Production Deployment Documentation

## 📚 Documentation Overview

This repository contains comprehensive guides for deploying your NutrieDiet Go API to production on Digital Ocean.

---

## 📖 Available Guides

### 1. **DO_QUICK_START.md** ⚡ 
**Start here if you want to deploy quickly (~2 hours)**

Step-by-step walkthrough covering:
- Creating and configuring Digital Ocean droplet
- Installing all required software (Go, MySQL, Nginx)
- Setting up SSL certificates with Let's Encrypt
- Configuring backups and monitoring
- Testing your deployed application

**Perfect for:** First-time deployment, getting your app live quickly

---

### 2. **DIGITAL_OCEAN_DEPLOYMENT.md** 📘
**Comprehensive production deployment guide**

In-depth documentation covering:
- Server setup and security hardening
- MySQL optimization for production
- Nginx configuration with rate limiting
- SSL/TLS best practices
- systemd service configuration
- Monitoring and logging setup
- Backup strategies
- Performance optimization
- Troubleshooting common issues

**Perfect for:** Understanding the complete deployment architecture, advanced configuration

---

### 3. **PRODUCTION_IMPROVEMENTS.md** 🔒
**Security and code improvements analysis**

Complete analysis of 30+ improvements needed:
- 🔴 8 Critical security issues
- 🟠 10 High priority issues
- 🟡 10 Medium priority issues
- 🟢 2 Low priority issues

Each issue includes:
- Risk assessment
- Detailed explanation
- Complete code fixes
- Testing strategies
- Implementation priority

**Perfect for:** Security review, code quality improvements, ongoing maintenance

---

### 4. **SECURITY_QUICK_FIXES.md** 🚀
**Top 5 critical fixes (can be done in 1 hour)**

Addresses ~70% of critical security vulnerabilities:
1. Database credentials (environment variables)
2. JWT secret key (secure generation)
3. Rate limiting (brute force protection)
4. Security headers (XSS, clickjacking protection)
5. Unprotected admin routes (proper authorization)

**Perfect for:** Quick security improvements before soft launch

---

### 5. **env.example** ⚙️
**Environment variables template**

Complete `.env` configuration file with:
- Development settings
- Production settings (Digital Ocean)
- All required environment variables
- Detailed comments explaining each setting
- Security best practices

**Perfect for:** Setting up your environment correctly

---

## 🚀 Quick Start Guide

### For First-Time Deployment:

1. **Read this order:**
   ```
   1. DO_QUICK_START.md          (Deploy in 2 hours)
   2. SECURITY_QUICK_FIXES.md    (Secure critical issues)
   3. PRODUCTION_IMPROVEMENTS.md (Plan remaining improvements)
   ```

2. **What you need:**
   - Digital Ocean account
   - Domain name
   - SSH key pair
   - Gmail account (for SMTP)
   - ~2-3 hours of time

3. **Estimated costs:**
   - Droplet: $12/month
   - Backups: $2.40/month
   - Domain: $10-15/year
   - **Total: ~$14.40/month**

---

## 🔧 Implementation Status

### ✅ Completed
- [x] Database credentials moved to environment variables
- [x] Connection pooling configured
- [x] Support for both local and cloud databases
- [x] TLS configuration based on environment
- [x] Comprehensive deployment documentation

### 🔄 In Progress (Follow SECURITY_QUICK_FIXES.md)
- [ ] JWT secret key from environment
- [ ] Rate limiting middleware
- [ ] Security headers middleware
- [ ] Protected admin routes
- [ ] Strong password requirements

### 📋 Planned (Follow PRODUCTION_IMPROVEMENTS.md)
- [ ] OTP attempt limiting
- [ ] File upload validation
- [ ] Structured logging
- [ ] Health check endpoints
- [ ] Monitoring and alerts

---

## 🏗️ Architecture Overview

### Digital Ocean Setup:
```
Internet
    ↓
Domain (yourdomain.com)
    ↓
Digital Ocean Droplet (Ubuntu 22.04)
    ├── Nginx (Port 80/443)
    │   ├── SSL/TLS (Let's Encrypt)
    │   ├── Rate Limiting
    │   └── Reverse Proxy
    │       ↓
    ├── Go Application (Port 8080)
    │   ├── JWT Authentication
    │   ├── Password Reset (OTP)
    │   └── API Endpoints
    │       ↓
    └── MySQL Database (localhost:3306)
        ├── Local connection (no TLS)
        └── Daily backups
```

### Security Layers:
1. **UFW Firewall** - Blocks all except 22, 80, 443
2. **Fail2Ban** - Blocks brute force attempts
3. **Nginx Rate Limiting** - 5 req/min for auth, 100 req/min for API
4. **SSL/TLS** - All traffic encrypted
5. **JWT Authentication** - Token-based auth
6. **bcrypt** - Password hashing
7. **OTP** - Email-based password reset

---

## 📊 Deployment Timeline

### Week 1: Initial Deployment
- [ ] Set up Digital Ocean droplet
- [ ] Configure MySQL database
- [ ] Deploy application
- [ ] Set up SSL certificates
- [ ] Configure firewall
- [ ] Implement critical security fixes (#1-5)

### Week 2: Security Hardening
- [ ] JWT token improvements
- [ ] OTP attempt limiting
- [ ] File upload security
- [ ] Error handling improvements
- [ ] Security testing

### Week 3: Infrastructure
- [ ] Structured logging
- [ ] Health checks
- [ ] Monitoring setup
- [ ] Backup automation
- [ ] Performance optimization

### Week 4: Polish & Testing
- [ ] Pagination implementation
- [ ] API versioning
- [ ] Load testing
- [ ] Security audit
- [ ] Documentation updates

---

## 🔐 Security Checklist

### Before Going Live:
- [ ] All database credentials in environment variables
- [ ] JWT secret key configured (64+ characters)
- [ ] Rate limiting enabled
- [ ] Security headers configured
- [ ] Admin routes protected
- [ ] SSL certificate installed
- [ ] Firewall configured
- [ ] Strong passwords enforced
- [ ] Error messages don't leak info
- [ ] Backups configured and tested

### After Going Live:
- [ ] Monitor error rates
- [ ] Check rate limiting effectiveness
- [ ] Review authentication logs
- [ ] Test all critical flows
- [ ] Monitor database performance
- [ ] Set up alerts
- [ ] Document incident response plan

---

## 🧪 Testing Checklist

### Local Testing:
```bash
# Test database connection
go run main.go

# Test migrations
cd migrate && go run migrate.go

# Test API endpoints
curl http://localhost:8080/health

# Test rate limiting (should block after 5)
for i in {1..10}; do curl -X POST http://localhost:8080/login; done
```

### Production Testing:
```bash
# Test SSL
curl -I https://yourdomain.com

# Test health check
curl https://yourdomain.com/health

# Test authentication flow
curl -X POST https://yourdomain.com/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!@#","first_name":"Test","last_name":"User","user_type":"CLIENT"}'

# Test password reset flow
curl -X POST https://yourdomain.com/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# Check logs
ssh nutriediet@your_droplet_ip
tail -f /opt/nutriediet/logs/app.log
```

---

## 📞 Getting Help

### If something breaks:

1. **Application won't start:**
   ```bash
   sudo systemctl status nutriediet
   sudo journalctl -u nutriediet -n 50
   ```

2. **Database issues:**
   ```bash
   sudo systemctl status mysql
   mysql -u nutriediet_app -p nutriediet_production
   ```

3. **SSL certificate issues:**
   ```bash
   sudo certbot certificates
   sudo certbot renew
   ```

4. **Can't access via domain:**
   ```bash
   nslookup yourdomain.com
   sudo systemctl status nginx
   ```

### Common Fixes:
- **Port already in use:** `sudo lsof -i :8080`
- **Permission denied:** Check file ownership and chmod
- **Out of disk space:** `df -h` and clean up old logs/backups
- **Database connection failed:** Verify credentials in .env

---

## 📈 Monitoring

### Key Metrics to Watch:
- Response time (should be < 200ms)
- Error rate (should be < 1%)
- Database connections (should be < 50)
- Disk usage (alert at 80%)
- Memory usage (alert at 90%)
- SSL certificate expiry (auto-renews, but verify)

### Log Locations:
- Application: `/opt/nutriediet/logs/app.log`
- Nginx Access: `/var/log/nginx/nutriediet_access.log`
- Nginx Error: `/var/log/nginx/nutriediet_error.log`
- MySQL: `/var/log/mysql/error.log`
- systemd: `sudo journalctl -u nutriediet`

---

## 🔄 Update Process

### When you make code changes:
```bash
# On your local machine
git add .
git commit -m "Your changes"
git push origin main

# On the server
ssh nutriediet@your_droplet_ip
cd /opt/nutriediet
git pull origin main
go build -a -installsuffix cgo -ldflags="-w -s" -o nutriediet-go .
sudo systemctl restart nutriediet
sudo systemctl status nutriediet
```

### For database schema changes:
```bash
# Update migration files first
# Then on server:
cd /opt/nutriediet/migrate
go run migrate.go
sudo systemctl restart nutriediet
```

---

## 💡 Best Practices

1. **Never commit .env files** - Already in .gitignore
2. **Use strong passwords** - 12+ characters, mixed case, numbers, symbols
3. **Rotate secrets regularly** - JWT keys, database passwords
4. **Monitor logs daily** - Check for unusual activity
5. **Test backups monthly** - Verify you can restore
6. **Keep system updated** - `sudo apt update && sudo apt upgrade`
7. **Document changes** - Update README when adding features
8. **Use branches** - Don't push directly to main
9. **Test locally first** - Don't debug in production
10. **Have a rollback plan** - Keep previous binary

---

## 🎯 Success Criteria

Your deployment is production-ready when:
- ✅ All services start automatically on boot
- ✅ SSL certificate is valid and auto-renewing
- ✅ Rate limiting blocks excessive requests
- ✅ Backups run daily without errors
- ✅ Health checks return 200 OK
- ✅ No sensitive data in logs
- ✅ Firewall blocks unnecessary ports
- ✅ All environment variables properly set
- ✅ Application recovers from crashes
- ✅ You can deploy updates without downtime

---

## 📝 Additional Resources

### Official Documentation:
- [Digital Ocean Droplet Docs](https://docs.digitalocean.com/products/droplets/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM](https://gorm.io/docs/)
- [Let's Encrypt](https://letsencrypt.org/docs/)
- [MySQL Performance](https://dev.mysql.com/doc/refman/8.0/en/optimization.html)

### Security Resources:
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://github.com/OWASP/Go-SCP)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)

---

## 🙏 Need Help?

1. Check the relevant guide above
2. Review the troubleshooting section
3. Check application logs
4. Search GitHub issues
5. Review Digital Ocean community tutorials

---

## 📅 Maintenance Schedule

### Daily:
- Check application logs for errors
- Verify backups completed
- Monitor disk space

### Weekly:
- Review security logs
- Check for system updates
- Test critical endpoints

### Monthly:
- Test backup restoration
- Review and rotate logs
- Update dependencies
- Security audit

### Quarterly:
- Rotate JWT secret keys
- Update SSL certificates (auto, but verify)
- Performance review
- Cost optimization review

---

**Last Updated:** 2025-10-05  
**Version:** 1.0  
**Status:** Ready for Production Deployment

---

## 🚀 You're Ready to Deploy!

Follow the guides in order, take your time, and don't hesitate to refer back to this documentation. Good luck! 🎉


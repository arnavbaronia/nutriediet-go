# Digital Ocean Deployment Guide - NutrieDiet Go

## Overview
This guide covers deploying your Go application on a Digital Ocean droplet with:
- Linux server (Ubuntu 22.04 LTS recommended)
- MySQL database on the same machine
- Nginx as reverse proxy
- SSL/TLS with Let's Encrypt
- systemd for process management
- Domain pointing to your droplet

---

## Table of Contents
1. [Server Setup](#server-setup)
2. [Database Configuration](#database-configuration)
3. [Application Deployment](#application-deployment)
4. [Nginx Configuration](#nginx-configuration)
5. [SSL/TLS Setup](#ssltls-setup)
6. [Security Hardening](#security-hardening)
7. [Monitoring & Logging](#monitoring--logging)
8. [Backup Strategy](#backup-strategy)

---

## Server Setup

### 1. Create Digital Ocean Droplet

**Recommended Specs:**
- **Size:** Basic - $12/month (2GB RAM, 1 CPU, 50GB SSD)
- **OS:** Ubuntu 22.04 LTS x64
- **Datacenter:** Choose closest to your users
- **Add-ons:** 
  - ✅ Monitoring (free)
  - ✅ Backups ($2.40/month - highly recommended)

### 2. Initial Server Configuration

```bash
# SSH into your droplet
ssh root@your_droplet_ip

# Update system
apt update && apt upgrade -y

# Set timezone
timedatectl set-timezone Asia/Kolkata  # or your timezone

# Create a non-root user
adduser nutriediet
usermod -aG sudo nutriediet

# Set up SSH key authentication for new user
mkdir -p /home/nutriediet/.ssh
cp ~/.ssh/authorized_keys /home/nutriediet/.ssh/
chown -R nutriediet:nutriediet /home/nutriediet/.ssh
chmod 700 /home/nutriediet/.ssh
chmod 600 /home/nutriediet/.ssh/authorized_keys

# Switch to new user
su - nutriediet
```

### 3. Install Required Software

```bash
# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version

# Install MySQL
sudo apt install mysql-server -y

# Install Nginx
sudo apt install nginx -y

# Install Certbot (for SSL)
sudo apt install certbot python3-certbot-nginx -y

# Install git
sudo apt install git -y

# Install build essentials
sudo apt install build-essential -y
```

---

## Database Configuration

### 1. Secure MySQL Installation

```bash
sudo mysql_secure_installation
```

**Recommended answers:**
- Set root password: **YES**
- Remove anonymous users: **YES**
- Disallow root login remotely: **YES**
- Remove test database: **YES**
- Reload privilege tables: **YES**

### 2. Create Application Database

```bash
# Login to MySQL
sudo mysql -u root -p

# Create database and user
CREATE DATABASE nutriediet_production CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'nutriediet_app'@'localhost' IDENTIFIED BY 'your_strong_password_here';

GRANT ALL PRIVILEGES ON nutriediet_production.* TO 'nutriediet_app'@'localhost';

FLUSH PRIVILEGES;

# Verify
SHOW DATABASES;
SELECT user, host FROM mysql.user;

EXIT;
```

### 3. Optimize MySQL for Production

```bash
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
```

Add/modify these settings:
```ini
[mysqld]
# Basic Settings
max_connections = 200
connect_timeout = 10
wait_timeout = 600
max_allowed_packet = 64M
thread_cache_size = 128
sort_buffer_size = 4M
bulk_insert_buffer_size = 16M
tmp_table_size = 64M
max_heap_table_size = 64M

# InnoDB Settings
innodb_buffer_pool_size = 512M  # 70% of available RAM
innodb_log_file_size = 128M
innodb_file_per_table = 1
innodb_flush_method = O_DIRECT

# Query Cache (for MySQL 5.7, skip for 8.0+)
# query_cache_type = 1
# query_cache_size = 32M

# Logging
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow-query.log
long_query_time = 2
log_error = /var/log/mysql/error.log

# Character Set
character_set_server = utf8mb4
collation_server = utf8mb4_unicode_ci
```

Restart MySQL:
```bash
sudo systemctl restart mysql
sudo systemctl enable mysql
```

---

## Application Deployment

### 1. Clone and Build Application

```bash
# Create application directory
sudo mkdir -p /opt/nutriediet
sudo chown nutriediet:nutriediet /opt/nutriediet
cd /opt/nutriediet

# Clone repository
git clone https://github.com/cd-Ishita/nutriediet-go.git .

# Create necessary directories
mkdir -p images logs

# Install dependencies
go mod download
go mod verify
```

### 2. Create Production Environment File

```bash
nano /opt/nutriediet/.env
```

**Production .env file:**
```bash
# Application
ENVIRONMENT=production
PORT=8080
GIN_MODE=release

# Database (LOCAL - same machine)
DB_USER=nutriediet_app
DB_PASSWORD=your_strong_password_here
DB_HOST=localhost
DB_PORT=3306
DB_NAME=nutriediet_production

# JWT Secret (generate with: openssl rand -base64 64)
JWT_SECRET_KEY=your_very_long_secure_random_secret_key_minimum_64_characters_generated_with_openssl

# SMTP Configuration
SMTP_EMAIL=nutriediet.help@gmail.com
SMTP_PASSWORD=your_16_character_gmail_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# CORS - Your Domain
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Rate Limiting
RATE_LIMIT_LOGIN=5
RATE_LIMIT_WINDOW=1m

# Logging
LOG_LEVEL=info
LOG_FILE=/opt/nutriediet/logs/app.log
```

**Secure the .env file:**
```bash
chmod 600 /opt/nutriediet/.env
```

### 3. Build Production Binary

```bash
cd /opt/nutriediet

# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o nutriediet-go .

# Make it executable
chmod +x nutriediet-go

# Test run
./nutriediet-go
# Press Ctrl+C after verifying it starts correctly
```

### 4. Run Database Migrations

```bash
cd /opt/nutriediet/migrate
go run migrate.go
```

### 5. Create systemd Service

```bash
sudo nano /etc/systemd/system/nutriediet.service
```

**Service configuration:**
```ini
[Unit]
Description=NutrieDiet Go API Service
After=network.target mysql.service
Requires=mysql.service

[Service]
Type=simple
User=nutriediet
Group=nutriediet
WorkingDirectory=/opt/nutriediet
ExecStart=/opt/nutriediet/nutriediet-go
Restart=always
RestartSec=5
StandardOutput=append:/opt/nutriediet/logs/app.log
StandardError=append:/opt/nutriediet/logs/error.log

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/nutriediet/images /opt/nutriediet/logs

# Environment
Environment="GIN_MODE=release"
EnvironmentFile=/opt/nutriediet/.env

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

**Enable and start service:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable nutriediet
sudo systemctl start nutriediet

# Check status
sudo systemctl status nutriediet

# View logs
sudo journalctl -u nutriediet -f
```

---

## Nginx Configuration

### 1. Basic Nginx Setup

```bash
sudo nano /etc/nginx/sites-available/nutriediet
```

**Nginx configuration:**
```nginx
# Rate limiting zones
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=5r/m;
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/m;
limit_req_zone $binary_remote_addr zone=general_limit:10m rate=200r/m;

# Upstream Go application
upstream nutriediet_backend {
    server 127.0.0.1:8080 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

# HTTP server - redirect to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name yourdomain.com www.yourdomain.com;

    # Let's Encrypt validation
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    # Redirect to HTTPS
    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS server
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL certificates (will be added by certbot)
    # ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    # ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_stapling on;
    ssl_stapling_verify on;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self'; img-src 'self' data: https:; script-src 'self'; style-src 'self' 'unsafe-inline'" always;

    # Logging
    access_log /var/log/nginx/nutriediet_access.log;
    error_log /var/log/nginx/nutriediet_error.log;

    # Max body size for file uploads
    client_max_body_size 10M;
    client_body_buffer_size 128k;

    # Timeouts
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;

    # Serve static files directly
    location /images/ {
        alias /opt/nutriediet/images/;
        expires 30d;
        add_header Cache-Control "public, immutable";
        access_log off;
    }

    # Health check (no rate limit)
    location /health {
        proxy_pass http://nutriediet_backend;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        access_log off;
    }

    # Auth endpoints - strict rate limiting
    location ~ ^/(signup|login|auth/forgot-password|auth/reset-password) {
        limit_req zone=auth_limit burst=2 nodelay;
        limit_req_status 429;

        proxy_pass http://nutriediet_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        proxy_cache_bypass $http_upgrade;
    }

    # API endpoints - moderate rate limiting
    location /api/ {
        limit_req zone=api_limit burst=20 nodelay;
        limit_req_status 429;

        proxy_pass http://nutriediet_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        proxy_cache_bypass $http_upgrade;
    }

    # All other endpoints
    location / {
        limit_req zone=general_limit burst=50 nodelay;
        limit_req_status 429;

        proxy_pass http://nutriediet_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;
        proxy_cache_bypass $http_upgrade;
    }
}
```

**Enable site:**
```bash
# Test configuration
sudo nginx -t

# Create symlink
sudo ln -s /etc/nginx/sites-available/nutriediet /etc/nginx/sites-enabled/

# Remove default site
sudo rm /etc/nginx/sites-enabled/default

# Restart Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx
```

---

## SSL/TLS Setup

### 1. Point Domain to Droplet

Before obtaining SSL certificate:

1. Go to your domain registrar (e.g., GoDaddy, Namecheap)
2. Add A records:
   ```
   Type  Name   Value              TTL
   A     @      your_droplet_ip    3600
   A     www    your_droplet_ip    3600
   ```
3. Wait for DNS propagation (5-60 minutes)

**Verify DNS:**
```bash
nslookup yourdomain.com
dig yourdomain.com +short
```

### 2. Obtain SSL Certificate

```bash
# Get certificate for both domain and www subdomain
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Follow prompts:
# - Enter email address
# - Agree to terms
# - Choose to redirect HTTP to HTTPS (option 2)
```

**Test SSL:**
```bash
# Test certificate
sudo certbot certificates

# Test renewal
sudo certbot renew --dry-run
```

**Auto-renewal is automatic via systemd timer:**
```bash
sudo systemctl status certbot.timer
```

### 3. Configure SSL in Nginx

Certbot automatically updates your Nginx config. Verify:
```bash
sudo nano /etc/nginx/sites-available/nutriediet
```

Should include:
```nginx
ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
```

---

## Security Hardening

### 1. Configure Firewall (UFW)

```bash
# Enable UFW
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH (change 22 if using custom port)
sudo ufw allow 22/tcp

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status verbose
```

### 2. Secure MySQL

```bash
# Edit MySQL config
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
```

Ensure these settings:
```ini
# Bind to localhost only (not accessible from outside)
bind-address = 127.0.0.1
mysqlx-bind-address = 127.0.0.1

# Disable remote root login
skip-name-resolve
```

Restart MySQL:
```bash
sudo systemctl restart mysql
```

### 3. Install Fail2Ban (Brute Force Protection)

```bash
# Install
sudo apt install fail2ban -y

# Configure for Nginx
sudo nano /etc/fail2ban/jail.local
```

```ini
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5
destemail = your-email@example.com
sendername = Fail2Ban
action = %(action_mwl)s

[sshd]
enabled = true
port = ssh
logpath = /var/log/auth.log

[nginx-http-auth]
enabled = true
port = http,https
logpath = /var/log/nginx/error.log

[nginx-limit-req]
enabled = true
port = http,https
logpath = /var/log/nginx/error.log
maxretry = 3
findtime = 300
bantime = 7200
```

```bash
# Start Fail2Ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Check status
sudo fail2ban-client status
```

### 4. Automatic Security Updates

```bash
sudo apt install unattended-upgrades -y
sudo dpkg-reconfigure --priority=low unattended-upgrades
```

### 5. Disable Root SSH Login

```bash
sudo nano /etc/ssh/sshd_config
```

Change/add:
```
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
```

```bash
sudo systemctl restart sshd
```

---

## Monitoring & Logging

### 1. Set Up Log Rotation

```bash
sudo nano /etc/logrotate.d/nutriediet
```

```
/opt/nutriediet/logs/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    missingok
    create 0640 nutriediet nutriediet
    sharedscripts
    postrotate
        systemctl reload nutriediet > /dev/null 2>&1 || true
    endscript
}
```

### 2. Monitor System Resources

```bash
# Install monitoring tools
sudo apt install htop iotop nethogs -y

# Check system resources
htop

# Monitor disk usage
df -h

# Monitor MySQL
mysqladmin -u root -p status
mysqladmin -u root -p processlist
```

### 3. Application Monitoring Script

Create monitoring script:
```bash
nano /opt/nutriediet/scripts/health_check.sh
```

```bash
#!/bin/bash

# Health check script
HEALTH_URL="http://localhost:8080/health"
LOG_FILE="/opt/nutriediet/logs/health_check.log"

response=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_URL)

if [ $response -eq 200 ]; then
    echo "$(date): Application is healthy" >> $LOG_FILE
else
    echo "$(date): Application health check failed (HTTP $response)" >> $LOG_FILE
    # Restart service
    sudo systemctl restart nutriediet
    echo "$(date): Service restarted" >> $LOG_FILE
fi
```

```bash
chmod +x /opt/nutriediet/scripts/health_check.sh

# Add to crontab (check every 5 minutes)
crontab -e
```

Add:
```
*/5 * * * * /opt/nutriediet/scripts/health_check.sh
```

### 4. Set Up Alerts with Email

```bash
# Install mail utilities
sudo apt install mailutils -y

# Test email
echo "Test email from NutrieDiet server" | mail -s "Test" your-email@example.com
```

---

## Backup Strategy

### 1. Database Backup Script

```bash
sudo mkdir -p /opt/backups
sudo chown nutriediet:nutriediet /opt/backups

nano /opt/nutriediet/scripts/backup_db.sh
```

```bash
#!/bin/bash

# Database backup script
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="nutriediet_production"
DB_USER="nutriediet_app"
DB_PASS="your_db_password"

# Create backup
mysqldump -u $DB_USER -p$DB_PASS $DB_NAME | gzip > $BACKUP_DIR/db_backup_$DATE.sql.gz

# Keep only last 7 days of backups
find $BACKUP_DIR -name "db_backup_*.sql.gz" -mtime +7 -delete

echo "$(date): Database backup completed - db_backup_$DATE.sql.gz" >> /opt/nutriediet/logs/backup.log
```

```bash
chmod +x /opt/nutriediet/scripts/backup_db.sh

# Add to crontab (daily at 2 AM)
crontab -e
```

Add:
```
0 2 * * * /opt/nutriediet/scripts/backup_db.sh
```

### 2. Digital Ocean Automated Backups

Enable in Digital Ocean dashboard:
- **Droplet** → **Backups** → Enable ($2.40/month)
- Weekly automatic snapshots
- Can restore entire droplet

### 3. Backup Images Folder

```bash
nano /opt/nutriediet/scripts/backup_images.sh
```

```bash
#!/bin/bash

BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d)
SOURCE_DIR="/opt/nutriediet/images"

tar -czf $BACKUP_DIR/images_backup_$DATE.tar.gz -C /opt/nutriediet images

# Keep only last 30 days
find $BACKUP_DIR -name "images_backup_*.tar.gz" -mtime +30 -delete
```

```bash
chmod +x /opt/nutriediet/scripts/backup_images.sh

# Add to crontab (weekly on Sunday at 3 AM)
crontab -e
```

Add:
```
0 3 * * 0 /opt/nutriediet/scripts/backup_images.sh
```

---

## Deployment Workflow

### 1. Initial Deployment Checklist

- [ ] Server created and secured
- [ ] Domain DNS configured
- [ ] MySQL installed and secured
- [ ] Database created with strong password
- [ ] Application built and tested locally
- [ ] Environment variables configured
- [ ] systemd service created and running
- [ ] Nginx configured as reverse proxy
- [ ] SSL certificate obtained
- [ ] Firewall configured
- [ ] Fail2Ban installed
- [ ] Backups configured
- [ ] Monitoring set up
- [ ] Health checks working

### 2. Update/Redeploy Application

```bash
# SSH into server
ssh nutriediet@your_droplet_ip

# Navigate to app directory
cd /opt/nutriediet

# Pull latest changes
git pull origin main

# Rebuild
go build -a -installsuffix cgo -ldflags="-w -s" -o nutriediet-go .

# Restart service
sudo systemctl restart nutriediet

# Check status
sudo systemctl status nutriediet

# Monitor logs
sudo journalctl -u nutriediet -f
```

### 3. Zero-Downtime Deployment (Advanced)

Create deployment script:
```bash
nano /opt/nutriediet/scripts/deploy.sh
```

```bash
#!/bin/bash

echo "Starting deployment..."

# Pull latest code
cd /opt/nutriediet
git pull origin main

# Build new binary
go build -a -installsuffix cgo -ldflags="-w -s" -o nutriediet-go-new .

# Test new binary
./nutriediet-go-new &
NEW_PID=$!
sleep 5

# Check if new binary is running
if ps -p $NEW_PID > /dev/null; then
    echo "New binary tested successfully"
    kill $NEW_PID
    
    # Replace old binary
    mv nutriediet-go nutriediet-go-old
    mv nutriediet-go-new nutriediet-go
    
    # Restart service
    sudo systemctl restart nutriediet
    
    echo "Deployment completed successfully"
else
    echo "New binary failed to start"
    rm nutriediet-go-new
    exit 1
fi
```

---

## Performance Optimization

### 1. Optimize Go Application

Update `main.go`:
```go
func main() {
    // Set GOMAXPROCS to match CPU cores
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    // ... rest of your code
}
```

### 2. Enable Nginx Caching

Add to Nginx config:
```nginx
# Cache zone
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=nutriediet_cache:10m max_size=100m inactive=60m use_temp_path=off;

# In server block for static content
location /images/ {
    alias /opt/nutriediet/images/;
    expires 30d;
    add_header Cache-Control "public, immutable";
    access_log off;
}
```

### 3. Enable Gzip Compression

Add to Nginx config:
```nginx
gzip on;
gzip_vary on;
gzip_proxied any;
gzip_comp_level 6;
gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss application/rss+xml font/truetype font/opentype application/vnd.ms-fontobject image/svg+xml;
```

---

## Troubleshooting

### Application Won't Start

```bash
# Check service status
sudo systemctl status nutriediet

# Check logs
sudo journalctl -u nutriediet -n 100

# Check if port is in use
sudo netstat -tulpn | grep 8080

# Test application manually
cd /opt/nutriediet
./nutriediet-go
```

### Database Connection Issues

```bash
# Check MySQL is running
sudo systemctl status mysql

# Test connection
mysql -u nutriediet_app -p nutriediet_production

# Check MySQL logs
sudo tail -f /var/log/mysql/error.log
```

### Nginx Issues

```bash
# Test configuration
sudo nginx -t

# Check error log
sudo tail -f /var/log/nginx/error.log

# Restart Nginx
sudo systemctl restart nginx
```

### SSL Certificate Issues

```bash
# Check certificate
sudo certbot certificates

# Renew manually
sudo certbot renew

# Check Nginx SSL config
sudo nginx -t
```

---

## Cost Estimate

**Monthly Costs:**
- Droplet (2GB RAM): $12/month
- Automated Backups: $2.40/month
- Domain (if new): $10-15/year
- **Total: ~$14.40/month + domain**

**Optional:**
- Larger droplet for scaling: $24-48/month
- Load balancer (for high traffic): $12/month
- CDN for images: $5-20/month

---

## Quick Reference Commands

```bash
# Application
sudo systemctl status nutriediet
sudo systemctl restart nutriediet
sudo journalctl -u nutriediet -f

# Database
sudo systemctl status mysql
mysql -u nutriediet_app -p
mysqladmin -u root -p processlist

# Nginx
sudo systemctl status nginx
sudo nginx -t
sudo systemctl reload nginx

# SSL
sudo certbot renew --dry-run
sudo certbot certificates

# Firewall
sudo ufw status
sudo ufw allow 22/tcp

# Logs
tail -f /opt/nutriediet/logs/app.log
tail -f /var/log/nginx/error.log
tail -f /var/log/mysql/error.log

# Backups
/opt/nutriediet/scripts/backup_db.sh
ls -lh /opt/backups/
```

---

## Security Checklist

- [ ] SSH key authentication only (no passwords)
- [ ] Root login disabled
- [ ] Firewall (UFW) enabled with minimal open ports
- [ ] Fail2Ban installed and configured
- [ ] MySQL bound to localhost only
- [ ] Strong database passwords
- [ ] SSL/TLS certificates installed
- [ ] Security headers in Nginx
- [ ] Rate limiting configured
- [ ] Application running as non-root user
- [ ] File permissions properly set (600 for .env)
- [ ] Automatic security updates enabled
- [ ] Regular backups configured
- [ ] Monitoring and alerts set up

---

**Document Version:** 1.0  
**Last Updated:** 2025-10-05  
**Tested On:** Ubuntu 22.04 LTS


# Staging Deployment Guide - staging.nutriediet.com

## Overview
This guide covers deploying the NutrieDiet application to a staging environment on DigitalOcean. The staging environment mirrors production but uses separate resources for testing changes before production deployment.

**Staging URL:** `https://staging.nutriediet.com`

---

## ⚠️ Important Configuration Notes

Before starting, be aware of these key configuration details:

| Item | Value | Notes |
|------|-------|-------|
| **Username** | `nutriediet-staging-user` | Not `nutriediet-staging` or `nutriediet` |
| **Service Name** | `nutriediet` | Not `nutriediet-staging` |
| **Backend Directory** | `/opt/nutriediet` | Not `/opt/nutriediet-staging` |
| **Frontend Directory** | `/var/www/nutriediet` | Not `/var/www/nutriediet-staging` |
| **Backend Port** | `8080` | Internal only, not exposed |
| **Database Name** | `nutriediet_staging` | Separate from production |
| **DB User** | `nutriediet_staging_user` | Separate credentials |

**Common Pitfalls to Avoid:**
1. ❌ Using work laptop's `package-lock.json` (contains private registry URLs)
2. ❌ Forgetting to add user to `www-data` group
3. ❌ Not cleaning up Nginx config after certbot
4. ❌ Using wrong username in `chown` commands

---

## Table of Contents
1. [Droplet Configuration](#droplet-configuration)
2. [Initial Server Setup](#initial-server-setup)
3. [Software Installation](#software-installation)
4. [Database Configuration](#database-configuration)
5. [Backend Deployment](#backend-deployment)
6. [Frontend Deployment](#frontend-deployment)
7. [Nginx Configuration](#nginx-configuration)
8. [SSL Setup](#ssl-setup)
9. [Testing & Verification](#testing--verification)
10. [Deployment Workflow](#deployment-workflow)
11. [Staging-Specific Notes](#staging-specific-notes)

---

## Droplet Configuration

### Recommended Specs
- **Size:** Basic - Rs869.15/month (2GB RAM, 1 CPU, 40GB SSD)
- **OS:** Ubuntu 24.04 LTS x64
- **Datacenter:** Choose closest to your users (e.g., Bangalore for India)
- **Hostname:** nutriediet-staging
- **Add-ons:** 
  - ✅ Monitoring (free)
  - ⚠️ Backups (optional for staging - $2.40/month)

### Cost Comparison
| Environment | Monthly Cost | Backups |
|------------|-------------|----------|
| Staging | Rs869.15 (~$10-12) | Optional |
| Production | $12-14 | Recommended |

---

## Initial Server Setup

### Step 1: Create Droplet
1. Login to DigitalOcean dashboard
2. Click **Create** → **Droplets**
3. Configure as per specs above
4. Add your SSH key for authentication
5. Click **Create Droplet**
6. Note the IP address: `your_staging_droplet_ip`

### Step 2: Configure DNS at Namecheap

1. Login to Namecheap account
2. Go to **Domain List** → Select your domain → **Manage**
3. Go to **Advanced DNS** tab
4. Add/Update DNS records:

```
Type    Host        Value                      TTL
A       staging     your_staging_droplet_ip    Automatic
```

**Wait 5-30 minutes for DNS propagation**

Verify DNS:
```bash
# On your local machine
nslookup staging.nutriediet.com
# Should return your staging droplet IP
```

### Step 3: Initial SSH & System Update

```bash
# SSH into droplet as root
ssh root@your_staging_droplet_ip

# Update system packages
apt update && apt upgrade -y

# Set timezone (adjust for your location)
timedatectl set-timezone Asia/Kolkata

# Install basic utilities
apt install -y curl wget vim htop net-tools
```

### Step 4: Create Application User

```bash
# Create user for running the application
adduser nutriediet-staging
# Set a strong password when prompted

# Add to sudo group
usermod -aG sudo nutriediet-staging

# Set up SSH key for new user
mkdir -p /home/nutriediet-staging/.ssh
cp ~/.ssh/authorized_keys /home/nutriediet-staging/.ssh/
chown -R nutriediet-staging-user:nutriediet-staging-user /home/nutriediet-staging/.ssh
chmod 700 /home/nutriediet-staging/.ssh
chmod 600 /home/nutriediet-staging/.ssh/authorized_keys

# Switch to new user
su - nutriediet-staging
```

---

## Software Installation

### Install Go 1.21.5

```bash
# Download Go
cd /tmp
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

# Install
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
# Expected output: go version go1.21.5 linux/amd64
```

### Install MySQL

```bash
# Install MySQL server
sudo apt install -y mysql-server

# Check status
sudo systemctl status mysql
sudo systemctl enable mysql
```

### Install Nginx

```bash
# Install Nginx
sudo apt install -y nginx

# Check status
sudo systemctl status nginx
sudo systemctl enable nginx
```

### Install Certbot (SSL)

```bash
# Install Certbot for Let's Encrypt SSL
sudo apt install -y certbot python3-certbot-nginx
```

### Install Git & Build Tools

```bash
# Install development tools
sudo apt install -y git build-essential
```

---

## Database Configuration

### Step 1: Secure MySQL Installation

```bash
sudo mysql_secure_installation
```

**Recommended answers:**
- Setup VALIDATE PASSWORD plugin? **Y**
- Password validation policy: **1** (MEDIUM)
- Set root password: **YES** (use strong password)
- Remove anonymous users: **YES**
- Disallow root login remotely: **YES**
- Remove test database: **YES**
- Reload privilege tables: **YES**

### Step 2: Create Staging Database

```bash
# Login to MySQL as root
sudo mysql -u root -p
```

**Execute these SQL commands:**

```sql
-- Create staging database
CREATE DATABASE nutriediet_staging CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create staging database user
CREATE USER 'nutriediet_staging'@'localhost' IDENTIFIED BY 'STRONG_PASSWORD_HERE';

-- Grant privileges
GRANT ALL PRIVILEGES ON nutriediet_staging.* TO 'nutriediet_staging'@'localhost';

-- Apply changes
FLUSH PRIVILEGES;

-- Verify
SHOW DATABASES;
SELECT user, host FROM mysql.user WHERE user = 'nutriediet_staging';

EXIT;
```

### Step 3: Test Database Connection

```bash
# Test connection with new credentials
mysql -u nutriediet_staging -p nutriediet_staging
# Enter password, you should connect successfully

# Exit
EXIT;
```

### Step 4: Basic MySQL Configuration for Staging

```bash
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
```

**Add/modify these settings (lighter than production):**

```ini
[mysqld]
# Basic Settings
max_connections = 100
connect_timeout = 10
wait_timeout = 600
max_allowed_packet = 32M

# InnoDB Settings (lighter for staging)
innodb_buffer_pool_size = 256M
innodb_log_file_size = 64M
innodb_file_per_table = 1

# Bind to localhost only
bind-address = 127.0.0.1
mysqlx-bind-address = 127.0.0.1

# Character Set
character_set_server = utf8mb4
collation_server = utf8mb4_unicode_ci
```

**Restart MySQL:**

```bash
sudo systemctl restart mysql
sudo systemctl status mysql
```

---

## Backend Deployment

### Step 1: Create Directory Structure

```bash
# Create application directories
sudo mkdir -p /opt/nutriediet
sudo chown nutriediet-staging-user:nutriediet-staging-user /opt/nutriediet
cd /opt/nutriediet

# Create subdirectories
mkdir -p images logs scripts
```

### Step 2: Clone Repository

```bash
cd /opt/nutriediet

# Clone backend repository (use staging branch if available)
git clone https://github.com/cd-Ishita/nutriediet-go.git .

# If you have a staging branch:
# git checkout staging
# Otherwise stay on main/master
```

### Step 3: Create Staging Environment File

```bash
nano /opt/nutriediet/.env
```

**Staging .env configuration:**

```bash
# Environment
ENVIRONMENT=staging
PORT=8080
GIN_MODE=release

# Database - STAGING
DB_USER=nutriediet_staging
DB_PASSWORD=your_staging_db_password_here
DB_HOST=localhost
DB_PORT=3306
DB_NAME=nutriediet_staging

# JWT Secret - Generate with: openssl rand -base64 64
# Use DIFFERENT secret than production
JWT_SECRET_KEY=your_staging_jwt_secret_minimum_64_characters_here
JWT_EXPIRY=24h

# SMTP Configuration - Can use same as production or separate staging email
SMTP_EMAIL=nutriediet.staging@gmail.com
SMTP_PASSWORD=your_16_character_gmail_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# CORS - Staging Domain
ALLOWED_ORIGINS=https://staging.nutriediet.com,http://localhost:3000

# Rate Limiting (can be more lenient for testing)
RATE_LIMIT_LOGIN=10
RATE_LIMIT_WINDOW=1m

# Logging
LOG_LEVEL=debug
LOG_FILE=/opt/nutriediet/logs/app.log

# Upload Configuration
UPLOAD_DIR=/opt/nutriediet/images
MAX_UPLOAD_SIZE=10485760

# App URLs
APP_ENV=staging
APP_URL=https://staging.nutriediet.com
API_URL=https://staging.nutriediet.com/api
FRONTEND_URL=https://staging.nutriediet.com

# Staging-specific flags
DEBUG_MODE=true
ENABLE_PROFILING=true
```

**Important Notes:**
- Use **different JWT secret** than production
- Use **different database** than production
- Set `LOG_LEVEL=debug` for more detailed logs
- Include `localhost:3000` in CORS for local frontend development

**Secure the file:**

```bash
chmod 600 /opt/nutriediet/.env
```

### Step 4: Build Backend Application

```bash
cd /opt/nutriediet

# Download dependencies
go mod download
go mod verify

# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -a -installsuffix cgo \
  -ldflags="-w -s -X main.Version=staging-$(date +%Y%m%d)" \
  -o nutriediet-go .

# Make executable
chmod +x nutriediet-go

# Verify binary
ls -lh nutriediet-go
file nutriediet-go
```

### Step 5: Run Database Migrations

```bash
cd /opt/nutriediet/migrate

# Run migrations
go run migrate.go

# Verify tables were created
mysql -u nutriediet_staging -p nutriediet_staging -e "SHOW TABLES;"
```

### Step 6: Test Backend Locally

```bash
cd /opt/nutriediet

# Test run the application
./nutriediet-go

# Should start on port 8080
# Press Ctrl+C to stop after verifying it starts successfully
```

### Step 7: Create systemd Service

```bash
sudo nano /etc/systemd/system/nutriediet.service
```

**Service configuration:**

```ini
[Unit]
Description=NutrieDiet Staging Go API Service
After=network.target mysql.service
Requires=mysql.service

[Service]
Type=simple
User=nutriediet-staging
Group=nutriediet-staging
WorkingDirectory=/opt/nutriediet
ExecStart=/opt/nutriediet/nutriediet-go
Restart=always
RestartSec=5
StandardOutput=append:/opt/nutriediet/logs/app.log
StandardError=append:/opt/nutriediet/logs/error.log

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/nutriediet/images /opt/nutriediet/logs

# Environment
Environment="GIN_MODE=release"
EnvironmentFile=/opt/nutriediet/.env

# Resource Limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

**Start and enable service:**

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable nutriediet

# Start service
sudo systemctl start nutriediet

# Check status
sudo systemctl status nutriediet

# Should show "active (running)"
```

**Monitor logs:**

```bash
# Follow live logs
sudo journalctl -u nutriediet -f

# View last 50 lines
sudo journalctl -u nutriediet -n 50

# Check application log file
tail -f /opt/nutriediet/logs/app.log
```

---

## Frontend Deployment

### Step 1: Create Frontend Directory

```bash
sudo mkdir -p /var/www/nutriediet
sudo chown nutriediet-staging-user:nutriediet-staging-user /var/www/nutriediet
cd /var/www/nutriediet
```

### Step 2: Clone Frontend Repository

```bash
# Option 1: Clone from GitHub
git clone https://github.com/YOUR_USERNAME/frontend.git .

# Option 2: Copy from local machine
# On your local machine:
# scp -r /Users/ishitagupta/Documents/Personal/frontend/* nutriediet-staging-user@staging_ip:/var/www/nutriediet/
```

### Step 3: Configure Frontend for Staging

**Update `package.json`:**

```bash
nano package.json
```

Add homepage field (if not already present):

```json
{
  "name": "nutriediet-frontend",
  "version": "0.1.0",
  "homepage": "/",
  ...
}
```

**Create staging environment file:**

```bash
nano .env.production
```

```bash
# Staging API configuration
REACT_APP_API_URL=https://staging.nutriediet.com/api
REACT_APP_ENV=staging
GENERATE_SOURCEMAP=false
PUBLIC_URL=https://staging.nutriediet.com

# Optional: Add staging identifier
REACT_APP_ENVIRONMENT_NAME=Staging
```

**Update `axiosInstance.js` to use environment variable:**

```bash
nano src/api/axiosInstance.js
```

Ensure it uses:

```javascript
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
```

### Step 4: Build Frontend

```bash
cd /var/www/nutriediet

# Install Node.js and npm if not already installed
# Use nvm for version management
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
source ~/.bashrc
nvm install 20
nvm use 20

# IMPORTANT: Fix package-lock.json if it has private registry references
# This happens if you created the lock file on a work laptop
# Check for private registries:
grep -i "unpm\|artifactory" package-lock.json

# If found, delete and regenerate:
rm package-lock.json

# Install dependencies (will create clean package-lock.json)
npm install

# Build for production
GENERATE_SOURCEMAP=false npm run build

# Verify build directory
ls -la build/
```

**Troubleshooting npm install:**

If you get authentication errors like `E401 Incorrect or missing password`:

```bash
# Your package-lock.json contains private registry URLs
# Delete it and regenerate from public npm
rm package-lock.json

# Clear npm cache
npm cache clean --force

# Install fresh
npm install

# Then build
npm run build
```

**Set proper permissions:**

```bash
# Add user to www-data group (required for Nginx to serve files)
sudo usermod -aG www-data nutriediet-staging-user

# Apply group change (logout/login or use newgrp)
newgrp www-data

# Set ownership and permissions
sudo chown -R nutriediet-staging-user:www-data /var/www/nutriediet
sudo chmod -R 755 /var/www/nutriediet
```

---

## Nginx Configuration

### Step 1: Create Nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/staging.nutriediet.com
```

**Nginx configuration for staging:**

```nginx
# Rate limiting zones (lighter than production)
limit_req_zone $binary_remote_addr zone=staging_auth_limit:10m rate=10r/m;
limit_req_zone $binary_remote_addr zone=staging_api_limit:10m rate=200r/m;
limit_req_zone $binary_remote_addr zone=staging_general_limit:10m rate=500r/m;

# Upstream Go backend
upstream staging_nutriediet_backend {
    server 127.0.0.1:8080 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

# HTTP server - redirect to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name staging.nutriediet.com;

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
    server_name staging.nutriediet.com;

    # SSL certificates (will be added by certbot)
    # ssl_certificate /etc/letsencrypt/live/staging.nutriediet.com/fullchain.pem;
    # ssl_certificate_key /etc/letsencrypt/live/staging.nutriediet.com/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    # Staging environment identifier
    add_header X-Environment "Staging" always;

    # Logging
    access_log /var/log/nginx/staging_nutriediet_access.log;
    error_log /var/log/nginx/staging_nutriediet_error.log;

    # Max body size for file uploads
    client_max_body_size 10M;
    client_body_buffer_size 128k;

    # Timeouts
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;

    # Backend image uploads
    location /images/ {
        alias /opt/nutriediet/images/;
        expires 7d;
        add_header Cache-Control "public, no-transform";
        access_log off;
    }

    # Health check endpoint (no rate limit)
    location /health {
        proxy_pass http://staging_nutriediet_backend;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        access_log off;
    }

    # API endpoints
    location /api/ {
        limit_req zone=staging_api_limit burst=50 nodelay;
        limit_req_status 429;

        proxy_pass http://staging_nutriediet_backend/;
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

    # Auth endpoints - moderate rate limiting for testing
    location ~ ^/(signup|login|auth/forgot-password|auth/reset-password) {
        limit_req zone=staging_auth_limit burst=5 nodelay;
        limit_req_status 429;

        proxy_pass http://staging_nutriediet_backend;
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

    # React static files with caching
    location /static/ {
        alias /var/www/nutriediet/build/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
        access_log off;
    }

    # React app - serve index.html for SPA routing
    location / {
        root /var/www/nutriediet/build;
        try_files $uri $uri/ /index.html;
        
        # Don't cache index.html
        location = /index.html {
            add_header Cache-Control "no-cache, no-store, must-revalidate";
            add_header Pragma "no-cache";
            add_header Expires "0";
        }
    }
}
```

### Step 2: Enable Site and Test Configuration

```bash
# Test nginx configuration
sudo nginx -t

# Create symlink to enable site
sudo ln -s /etc/nginx/sites-available/staging.nutriediet.com /etc/nginx/sites-enabled/

# Remove default site if present
sudo rm -f /etc/nginx/sites-enabled/default

# Test again
sudo nginx -t

# If test passes, reload nginx
sudo systemctl reload nginx
```

---

## SSL Setup

### Step 1: Obtain SSL Certificate

**Important:** Ensure DNS is already pointing to your staging server before running certbot.

```bash
# Verify DNS first
nslookup staging.nutriediet.com
# Should return your staging droplet IP

# Obtain SSL certificate
sudo certbot --nginx -d staging.nutriediet.com

# Follow the prompts:
# 1. Enter email address (for renewal notifications)
# 2. Agree to terms of service (Y)
# 3. Share email with EFF (optional, Y or N)
# Certbot will automatically configure Nginx
```

### Step 2: Verify SSL Configuration

```bash
# Check certificates
sudo certbot certificates

# Test SSL renewal
sudo certbot renew --dry-run

# Should see "Congratulations, all simulated renewals succeeded"
```

### Step 3: Verify Nginx Configuration After SSL

```bash
# Check that certbot updated the config correctly
sudo nano /etc/nginx/sites-available/staging.nutriediet.com

# Look for these lines (added by certbot):
# ssl_certificate /etc/letsencrypt/live/staging.nutriediet.com/fullchain.pem;
# ssl_certificate_key /etc/letsencrypt/live/staging.nutriediet.com/privkey.pem;

# Test and reload
sudo nginx -t
sudo systemctl reload nginx
```

### Step 4: Enable Auto-Renewal

```bash
# Certbot installs a systemd timer for auto-renewal
# Verify it's active
sudo systemctl status certbot.timer

# Should show "active (waiting)"
```

---

## Testing & Verification

### Backend API Tests

```bash
# Test from server
curl http://localhost:8080/health
# Expected: {"status":"ok"} or similar

# Test externally
curl https://staging.nutriediet.com/health
```

### Frontend Tests

```bash
# Test in browser
# Visit: https://staging.nutriediet.com
# Should load React app

# Check for mixed content warnings (should be none)
```

### Full Integration Tests

**From your local machine:**

```bash
# Health check
curl https://staging.nutriediet.com/api/health

# Test signup
curl -X POST https://staging.nutriediet.com/api/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "staging-test@example.com",
    "password": "TestPassword123!",
    "first_name": "Staging",
    "last_name": "Test",
    "user_type": "CLIENT"
  }'

# Test login
curl -X POST https://staging.nutriediet.com/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "staging-test@example.com",
    "password": "TestPassword123!",
    "user_type": "CLIENT"
  }'
```

### Verify All Services

```bash
# On the staging server
sudo systemctl status nutriediet
sudo systemctl status mysql
sudo systemctl status nginx

# Check logs
sudo journalctl -u nutriediet -n 50
tail -f /opt/nutriediet/logs/app.log
tail -f /var/log/nginx/staging_nutriediet_error.log
```

---

## Deployment Workflow

### Initial Deployment ✅
You've just completed this by following the guide above!

### Deploying Backend Updates

```bash
# SSH into staging server
ssh nutriediet-staging-user@staging.nutriediet.com

# Navigate to backend directory
cd /opt/nutriediet

# Pull latest changes
git pull origin main  # or staging branch

# Rebuild binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -a -installsuffix cgo \
  -ldflags="-w -s -X main.Version=staging-$(date +%Y%m%d)" \
  -o nutriediet-go .

# Run any new migrations
cd migrate
go run migrate.go
cd ..

# Restart service
sudo systemctl restart nutriediet

# Monitor logs
sudo journalctl -u nutriediet -f
```

### Deploying Frontend Updates

```bash
# SSH into staging server
ssh nutriediet-staging-user@staging.nutriediet.com

# Navigate to frontend directory
cd /var/www/nutriediet

# Pull latest changes
git pull origin main  # or staging branch

# Install any new dependencies
npm ci

# Rebuild
GENERATE_SOURCEMAP=false npm run build

# No service restart needed - files are served directly by Nginx
# Just verify the new build is live
curl -I https://staging.nutriediet.com
```

### Rolling Back

**Backend rollback:**

```bash
cd /opt/nutriediet

# View git history
git log --oneline -n 10

# Rollback to previous commit
git checkout <commit-hash>

# Rebuild
CGO_ENABLED=0 GOOS=linux go build -o nutriediet-go .

# Restart
sudo systemctl restart nutriediet
```

**Frontend rollback:**

```bash
cd /var/www/nutriediet

# Rollback to previous commit
git checkout <commit-hash>

# Rebuild
npm run build
```

---

## Staging-Specific Notes

### Differences from Production

| Feature | Staging | Production |
|---------|---------|------------|
| Domain | staging.nutriediet.com | nutriediet.com / www.nutriediet.com |
| RAM | 2GB | 2GB |
| Storage | 40GB | 50GB+ |
| Rate Limits | Lenient (for testing) | Strict |
| Logging | DEBUG level | INFO level |
| Backups | Optional | Required |
| JWT Secret | Different | Different |
| Database | Separate | Separate |
| CORS | Includes localhost | Production domains only |

### Staging Best Practices

1. **Always test on staging first**
   - Deploy all changes to staging before production
   - Run full integration tests
   - Verify database migrations

2. **Use staging for client demos**
   - Safe environment to showcase features
   - Can reset data easily if needed

3. **Test with production-like data**
   - Anonymize and copy production data to staging periodically
   - Test edge cases and data migrations

4. **Monitor resource usage**
   - Watch disk space (40GB limit)
   - Monitor memory usage
   - Check database size

5. **Keep staging in sync**
   - Regularly update staging to match production configuration
   - Keep dependencies up to date

### Resetting Staging Database

If you need to reset staging data:

```bash
# Backup first (just in case)
mysqldump -u nutriediet_staging -p nutriediet_staging > /tmp/staging_backup.sql

# Drop and recreate database
mysql -u root -p
```

```sql
DROP DATABASE nutriediet_staging;
CREATE DATABASE nutriediet_staging CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

```bash
# Run migrations again
cd /opt/nutriediet/migrate
go run migrate.go
```

### Copying Production Data to Staging

```bash
# On production server: Export database
mysqldump -u nutriediet_app -p nutriediet_production > prod_export.sql

# Transfer to staging
scp prod_export.sql nutriediet-staging-user@staging.nutriediet.com:/tmp/

# On staging: Import
mysql -u nutriediet_staging -p nutriediet_staging < /tmp/prod_export.sql

# Anonymize sensitive data (emails, passwords, etc.)
mysql -u nutriediet_staging -p nutriediet_staging
```

```sql
-- Example: Anonymize user emails
UPDATE userauth SET email = CONCAT('user', user_id, '@staging.test');
-- Add more anonymization as needed
```

---

## Security Configuration

### Firewall Setup

```bash
# Set up UFW firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status verbose
```

### Disable Root SSH Login

```bash
sudo nano /etc/ssh/sshd_config
```

Add/modify:

```
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
```

```bash
# Restart SSH
sudo systemctl restart sshd
```

### Install Fail2Ban (Optional for Staging)

```bash
# Install
sudo apt install fail2ban -y

# Basic configuration
sudo nano /etc/fail2ban/jail.local
```

```ini
[DEFAULT]
bantime = 1800
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = ssh
logpath = /var/log/auth.log
```

```bash
# Start Fail2Ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Check status
sudo fail2ban-client status
```

---

## Monitoring & Maintenance

### Log Files

```bash
# Application logs
tail -f /opt/nutriediet/logs/app.log
tail -f /opt/nutriediet/logs/error.log

# System logs
sudo journalctl -u nutriediet -f

# Nginx logs
tail -f /var/log/nginx/staging_nutriediet_access.log
tail -f /var/log/nginx/staging_nutriediet_error.log

# MySQL logs
sudo tail -f /var/log/mysql/error.log
```

### System Monitoring

```bash
# Install monitoring tools
sudo apt install htop iotop

# Check system resources
htop

# Check disk usage
df -h

# Check disk I/O
sudo iotop

# Check memory
free -h

# Check active connections
sudo netstat -tulpn
```

### Database Backups (Optional)

If you want automated backups on staging:

```bash
# Create backup script
nano /opt/nutriediet/scripts/backup_db.sh
```

```bash
#!/bin/bash

BACKUP_DIR="/opt/nutriediet/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="nutriediet_staging"
DB_USER="nutriediet_staging"
DB_PASS="your_staging_db_password"

mkdir -p $BACKUP_DIR

# Create backup
mysqldump -u $DB_USER -p$DB_PASS $DB_NAME | gzip > $BACKUP_DIR/staging_backup_$DATE.sql.gz

# Keep only last 7 days of backups
find $BACKUP_DIR -name "staging_backup_*.sql.gz" -mtime +7 -delete

echo "$(date): Backup completed - staging_backup_$DATE.sql.gz" >> /opt/nutriediet/logs/backup.log
```

```bash
# Make executable
chmod +x /opt/nutriediet/scripts/backup_db.sh

# Create backups directory
mkdir -p /opt/nutriediet/backups

# Add to crontab (weekly on Sundays at 3 AM)
crontab -e
```

Add:

```
0 3 * * 0 /opt/nutriediet/scripts/backup_db.sh
```

---

## Troubleshooting

### npm Install Fails with E401 Authentication Error

**Problem:** `npm error code E401` or `Incorrect or missing password`

**Cause:** Your `package-lock.json` contains references to a private/work registry (e.g., Uber's internal npm).

**Solution:**

```bash
cd /var/www/nutriediet

# Check for private registries
grep -i "unpm\|artifactory\|registry.*internal" package-lock.json

# If found, delete and regenerate
rm package-lock.json

# Clear npm cache
npm cache clean --force

# Remove any .npmrc files with authentication
rm -f ~/.npmrc
rm -f .npmrc

# Install fresh from public registry
npm install

# Then build
npm run build
```

### Backend Won't Start

```bash
# Check service status
sudo systemctl status nutriediet

# Common issue: Service is stopped
# Solution: Start it
sudo systemctl start nutriediet

# Check logs
sudo journalctl -u nutriediet -n 100
tail -f /opt/nutriediet/logs/error.log

# Check if port 8080 is in use
sudo netstat -tulpn | grep 8080

# Check if .env file exists and has correct values
ls -la /opt/nutriediet/.env
cat /opt/nutriediet/.env | grep -E "DB_|PORT"

# Try running manually to see errors
cd /opt/nutriediet
./nutriediet-go
```

### Frontend Shows 404 or Blank Page

```bash
# Check if build directory exists
ls -la /var/www/nutriediet/build/

# Check file permissions
ls -la /var/www/nutriediet/build/index.html

# Fix permissions if needed
sudo chown -R nutriediet-staging-user:www-data /var/www/nutriediet
sudo chmod -R 755 /var/www/nutriediet

# Check Nginx configuration
sudo nginx -t

# Check Nginx error logs
tail -f /var/log/nginx/staging_nutriediet_error.log

# Rebuild frontend
cd /var/www/nutriediet
npm run build
```

### Database Connection Issues

```bash
# Check MySQL is running
sudo systemctl status mysql

# Test connection
mysql -u nutriediet_staging -p nutriediet_staging

# Check MySQL logs
sudo tail -f /var/log/mysql/error.log

# Verify .env file has correct credentials
cat /opt/nutriediet/.env | grep DB_
```

### SSL Certificate Issues

```bash
# Check certificates
sudo certbot certificates

# Renew manually
sudo certbot renew

# Check Nginx SSL configuration
sudo nano /etc/nginx/sites-available/staging.nutriediet.com

# Test and reload
sudo nginx -t
sudo systemctl reload nginx
```

### Nginx Shows 301 Redirect Loop or Config Error

**Problem:** After running certbot, Nginx has duplicate or conflicting server blocks.

**Solution:**

```bash
# Edit Nginx config
sudo nano /etc/nginx/sites-available/staging.nutriediet.com

# Ensure you have ONLY TWO server blocks:
# 1. HTTPS server (listen 443 ssl) - with all your location blocks
# 2. HTTP redirect (listen 80) - that redirects to HTTPS

# Remove any server block that:
# - Mixes listen 443 with return 301
# - Has SSL config but no locations
# - Is duplicated

# Your config should look like:
# 
# upstream staging_nutriediet_backend { ... }
#
# server {
#     listen 443 ssl http2;
#     ssl_certificate ...;
#     location / { ... }
#     location /api/ { ... }
# }
#
# server {
#     listen 80;
#     return 301 https://$host$request_uri;
# }

# Test
sudo nginx -t

# Reload
sudo systemctl reload nginx
```

### Permission Denied Errors

**Problem:** `chown: invalid user: 'nutriediet-staging:www-data'`

**Solution:**

```bash
# Your username is nutriediet-staging-user (not nutriediet-staging)
whoami
# Should show: nutriediet-staging-user

# Use correct username
sudo chown -R nutriediet-staging-user:www-data /var/www/nutriediet

# Add yourself to www-data group
sudo usermod -aG www-data nutriediet-staging-user

# Apply group change
newgrp www-data
```

### API Returns 502 Bad Gateway

```bash
# Check if backend is running
sudo systemctl status nutriediet

# Check if backend responds locally
curl http://localhost:8080/health

# Check Nginx upstream configuration
sudo nginx -t

# Check Nginx error log
tail -f /var/log/nginx/staging_nutriediet_error.log
```

---

## Common Deployment Issues We Fixed

### 1. Package Lock File with Private Registry

If `npm install` fails with authentication error, your `package-lock.json` has private registry references. Delete it and regenerate.

### 2. Wrong Username

The username is `nutriediet-staging-user` (not `nutriediet-staging` or `nutriediet`).

### 3. Service Name

The systemd service is named `nutriediet` (not `nutriediet-staging`).

### 4. Directory Paths

- Backend: `/opt/nutriediet` (not `/opt/nutriediet-staging`)
- Frontend: `/var/www/nutriediet` (not `/var/www/nutriediet-staging`)

### 5. Certbot Creates Duplicate Nginx Blocks

After running certbot, manually clean up the Nginx config to have only 2 server blocks (one for HTTPS with locations, one for HTTP redirect).

### 6. www-data Group

Don't forget to add your user to www-data group: `sudo usermod -aG www-data nutriediet-staging-user`

---

## Quick Reference Commands

### Service Management

```bash
# Backend
sudo systemctl status nutriediet
sudo systemctl restart nutriediet
sudo systemctl stop nutriediet
sudo systemctl start nutriediet
sudo journalctl -u nutriediet -f

# MySQL
sudo systemctl status mysql
sudo systemctl restart mysql

# Nginx
sudo systemctl status nginx
sudo nginx -t
sudo systemctl reload nginx
sudo systemctl restart nginx

# SSL
sudo certbot certificates
sudo certbot renew
sudo certbot renew --dry-run
```

### Log Viewing

```bash
# Application logs
tail -f /opt/nutriediet/logs/app.log
tail -f /opt/nutriediet/logs/error.log

# System logs
sudo journalctl -u nutriediet -n 100 --no-pager
sudo journalctl -u nutriediet -f

# Nginx logs
tail -f /var/log/nginx/staging_nutriediet_access.log
tail -f /var/log/nginx/staging_nutriediet_error.log

# MySQL logs
sudo tail -f /var/log/mysql/error.log
```

### Deployment Commands

```bash
# Backend update
cd /opt/nutriediet
git pull
go build -o nutriediet-go .
sudo systemctl restart nutriediet

# Frontend update
cd /var/www/nutriediet
git pull
npm run build

# Database migration
cd /opt/nutriediet/migrate
go run migrate.go
```

---

## Checklist

### Initial Deployment Checklist

- [ ] Droplet created on DigitalOcean
- [ ] DNS configured at Namecheap (staging → droplet IP)
- [ ] SSH access configured
- [ ] System packages updated
- [ ] Non-root user created (nutriediet-staging)
- [ ] Go 1.21.5 installed
- [ ] MySQL installed and secured
- [ ] Nginx installed
- [ ] Certbot installed
- [ ] Database created (nutriediet_staging)
- [ ] Database user created with privileges
- [ ] Backend repository cloned
- [ ] Backend .env file configured
- [ ] Backend built successfully
- [ ] Database migrations run
- [ ] systemd service created and running
- [ ] Frontend repository cloned
- [ ] Frontend built successfully
- [ ] Nginx configured
- [ ] SSL certificate obtained
- [ ] Firewall (UFW) configured
- [ ] All services running
- [ ] Health endpoints responding
- [ ] Frontend accessible via HTTPS
- [ ] API accessible via HTTPS
- [ ] Login/signup tested

### Pre-Production Promotion Checklist

Test these on staging before promoting to production:

- [ ] User signup works
- [ ] User login works
- [ ] Password reset flow works
- [ ] Image upload works
- [ ] All CRUD operations work
- [ ] Rate limiting works
- [ ] SSL certificate valid
- [ ] CORS configured correctly
- [ ] Database migrations successful
- [ ] No console errors in browser
- [ ] API returns proper error messages
- [ ] Mobile responsiveness checked
- [ ] Load testing passed
- [ ] Security scan passed

---

## Cost Summary

**Monthly Costs:**
- VPS Droplet: Rs869.15 (~$10-12 USD)
- SSL Certificate: FREE (Let's Encrypt)
- Backups (optional): $2.40/month
- **Total: Rs869.15/month (without backups)**

**Additional one-time/annual:**
- Domain registration (if separate): $10-15/year
- Namecheap DNS: Included with domain

---

## Support & Additional Resources

### Related Documentation
- **Production Deployment:** See `DIGITAL_OCEAN_DEPLOYMENT.md`
- **Quick Start:** See `DO_QUICK_START.md`
- **Security Improvements:** See `PRODUCTION_IMPROVEMENTS.md` or `SECURITY_QUICK_FIXES.md`
- **Frontend Configuration:** See `frontend-axios-update.md`

### Useful Links
- [DigitalOcean Documentation](https://docs.digitalocean.com/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Go Documentation](https://go.dev/doc/)

---

## Next Steps After Staging Deployment

1. **Test thoroughly on staging**
   - All user workflows
   - All API endpoints
   - Database operations
   - File uploads
   - Email notifications

2. **Load testing**
   - Use tools like Apache Bench or k6
   - Simulate concurrent users
   - Identify bottlenecks

3. **Security audit**
   - Run security scans
   - Test authentication flows
   - Verify rate limiting

4. **Documentation**
   - Document any staging-specific configurations
   - Update API documentation
   - Document known issues

5. **Plan production deployment**
   - Review differences between staging and production
   - Plan migration strategy
   - Schedule maintenance window

---

**Document Version:** 1.0  
**Last Updated:** 2025-11-07  
**Tested On:** Ubuntu 24.04 LTS  
**Target Environment:** Staging (staging.nutriediet.com)

---

## Questions or Issues?

If you encounter any issues during deployment:

1. Check the troubleshooting section above
2. Review service logs (systemd, nginx, mysql)
3. Verify all environment variables are set correctly
4. Ensure DNS has propagated
5. Check firewall rules

**Happy Staging! 🚀**


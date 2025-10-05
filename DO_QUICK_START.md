# Digital Ocean Quick Start Guide

## Prerequisites
- Digital Ocean account
- Domain name (e.g., yourdomain.com)
- GitHub repository access
- SSH key pair

---

## Part 1: Create & Configure Droplet (20 minutes)

### Step 1: Create Droplet
1. Login to Digital Ocean
2. Click **Create** → **Droplets**
3. Select:
   - **Image:** Ubuntu 22.04 LTS x64
   - **Plan:** Basic - $12/month (2GB RAM)
   - **Datacenter:** Choose closest to your users
   - **Authentication:** SSH keys (add yours)
   - **Hostname:** nutriediet-production
   - **Enable:** Backups ($2.40/month)
4. Click **Create Droplet**
5. Note the IP address: `your_droplet_ip`

### Step 2: Point Domain to Droplet
1. Go to your domain registrar (GoDaddy, Namecheap, etc.)
2. Add DNS records:
   ```
   Type  Name   Value              TTL
   A     @      your_droplet_ip    3600
   A     www    your_droplet_ip    3600
   ```
3. Wait 5-60 minutes for DNS propagation
4. Verify: `nslookup yourdomain.com`

### Step 3: Initial Server Setup
```bash
# SSH into droplet
ssh root@your_droplet_ip

# Update system
apt update && apt upgrade -y

# Create app user
adduser nutriediet
usermod -aG sudo nutriediet

# Set up SSH for new user
mkdir -p /home/nutriediet/.ssh
cp ~/.ssh/authorized_keys /home/nutriediet/.ssh/
chown -R nutriediet:nutriediet /home/nutriediet/.ssh
chmod 700 /home/nutriediet/.ssh
chmod 600 /home/nutriediet/.ssh/authorized_keys

# Switch to new user
su - nutriediet
```

---

## Part 2: Install Software (15 minutes)

```bash
# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version

# Install MySQL
sudo apt install mysql-server -y

# Install Nginx
sudo apt install nginx -y

# Install Certbot (SSL)
sudo apt install certbot python3-certbot-nginx -y

# Install Git
sudo apt install git build-essential -y
```

---

## Part 3: Configure MySQL (10 minutes)

```bash
# Secure MySQL
sudo mysql_secure_installation
# Answer: YES to all prompts, set strong root password

# Create database
sudo mysql -u root -p
```

In MySQL:
```sql
CREATE DATABASE nutriediet_production CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'nutriediet_app'@'localhost' IDENTIFIED BY 'your_strong_password_here';

GRANT ALL PRIVILEGES ON nutriediet_production.* TO 'nutriediet_app'@'localhost';

FLUSH PRIVILEGES;

EXIT;
```

Test connection:
```bash
mysql -u nutriediet_app -p nutriediet_production
# Enter password, then: EXIT;
```

---

## Part 4: Deploy Application (20 minutes)

```bash
# Create app directory
sudo mkdir -p /opt/nutriediet
sudo chown nutriediet:nutriediet /opt/nutriediet
cd /opt/nutriediet

# Clone repository (replace with your repo URL)
git clone https://github.com/cd-Ishita/nutriediet-go.git .

# Create directories
mkdir -p images logs

# Create .env file
nano .env
```

**Paste this into `.env`** (update the values):
```bash
# Application
ENVIRONMENT=production
PORT=8080
GIN_MODE=release

# Database (LOCAL)
DB_USER=nutriediet_app
DB_PASSWORD=your_strong_password_here
DB_HOST=localhost
DB_PORT=3306
DB_NAME=nutriediet_production

# JWT Secret (generate with: openssl rand -base64 64)
JWT_SECRET_KEY=paste_your_generated_key_here

# SMTP
SMTP_EMAIL=nutriediet.help@gmail.com
SMTP_PASSWORD=your_16_character_gmail_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# CORS
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

Save and exit (Ctrl+X, Y, Enter)

```bash
# Secure .env file
chmod 600 .env

# Install Go dependencies
go mod download

# Build application
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o nutriediet-go .

# Make executable
chmod +x nutriediet-go

# Run migrations
cd migrate
go run migrate.go
cd ..

# Test application
./nutriediet-go
# Press Ctrl+C after verifying it starts
```

---

## Part 5: Set Up systemd Service (5 minutes)

```bash
sudo nano /etc/systemd/system/nutriediet.service
```

**Paste this:**
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

NoNewPrivileges=true
PrivateTmp=true

Environment="GIN_MODE=release"
EnvironmentFile=/opt/nutriediet/.env

[Install]
WantedBy=multi-user.target
```

Save and exit, then:
```bash
# Start service
sudo systemctl daemon-reload
sudo systemctl enable nutriediet
sudo systemctl start nutriediet

# Check status
sudo systemctl status nutriediet

# Should show "active (running)"
```

---

## Part 6: Configure Nginx (15 minutes)

```bash
sudo nano /etc/nginx/sites-available/nutriediet
```

**Paste this** (replace `yourdomain.com` with your actual domain):
```nginx
# Rate limiting
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=5r/m;
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/m;

upstream nutriediet_backend {
    server 127.0.0.1:8080;
    keepalive 32;
}

# HTTP - redirect to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name yourdomain.com www.yourdomain.com;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL certificates (added by certbot)
    # ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    # ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Logs
    access_log /var/log/nginx/nutriediet_access.log;
    error_log /var/log/nginx/nutriediet_error.log;

    client_max_body_size 10M;

    # Serve images
    location /images/ {
        alias /opt/nutriediet/images/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # Health check
    location /health {
        proxy_pass http://nutriediet_backend;
        access_log off;
    }

    # Auth endpoints - strict rate limit
    location ~ ^/(signup|login|auth/) {
        limit_req zone=auth_limit burst=2 nodelay;
        
        proxy_pass http://nutriediet_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # All other endpoints
    location / {
        limit_req zone=api_limit burst=20 nodelay;
        
        proxy_pass http://nutriediet_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Save and exit, then:
```bash
# Test Nginx config
sudo nginx -t

# Enable site
sudo ln -s /etc/nginx/sites-available/nutriediet /etc/nginx/sites-enabled/

# Remove default
sudo rm /etc/nginx/sites-enabled/default

# Restart Nginx
sudo systemctl restart nginx
```

---

## Part 7: Get SSL Certificate (5 minutes)

```bash
# Get certificate (replace yourdomain.com with your domain)
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Follow prompts:
# 1. Enter email
# 2. Agree to terms (Y)
# 3. Share email (N)
# Certbot will automatically configure Nginx

# Test renewal
sudo certbot renew --dry-run
```

---

## Part 8: Configure Firewall (5 minutes)

```bash
# Set up UFW firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

---

## Part 9: Set Up Backups (10 minutes)

```bash
# Create backup script
nano /opt/nutriediet/scripts/backup_db.sh
```

**Paste this:**
```bash
#!/bin/bash
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="nutriediet_production"
DB_USER="nutriediet_app"
DB_PASS="your_db_password"

mkdir -p $BACKUP_DIR
mysqldump -u $DB_USER -p$DB_PASS $DB_NAME | gzip > $BACKUP_DIR/db_backup_$DATE.sql.gz
find $BACKUP_DIR -name "db_backup_*.sql.gz" -mtime +7 -delete
echo "$(date): Backup completed" >> /opt/nutriediet/logs/backup.log
```

Save and exit, then:
```bash
# Make executable
chmod +x /opt/nutriediet/scripts/backup_db.sh

# Add to crontab (daily at 2 AM)
crontab -e
```

Add this line:
```
0 2 * * * /opt/nutriediet/scripts/backup_db.sh
```

---

## Part 10: Verify Everything Works (5 minutes)

```bash
# Check all services
sudo systemctl status nutriediet
sudo systemctl status mysql
sudo systemctl status nginx

# Test health endpoint
curl http://localhost:8080/health

# Test from outside (replace with your domain)
curl https://yourdomain.com/health

# Check logs
tail -f /opt/nutriediet/logs/app.log

# Check SSL certificate
sudo certbot certificates
```

---

## Testing Your API

From your local machine:
```bash
# Test health
curl https://yourdomain.com/health

# Test signup (should get rate limited after 5 requests)
curl -X POST https://yourdomain.com/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"TestPass123!","first_name":"Test","last_name":"User","user_type":"CLIENT"}'

# Test login
curl -X POST https://yourdomain.com/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"TestPass123!","user_type":"CLIENT"}'
```

---

## Common Commands

```bash
# Application logs
sudo journalctl -u nutriediet -f

# Restart application
sudo systemctl restart nutriediet

# Restart Nginx
sudo systemctl restart nginx

# Check disk space
df -h

# Check memory
free -h

# Monitor processes
htop

# Database backup (manual)
/opt/nutriediet/scripts/backup_db.sh

# List backups
ls -lh /opt/backups/
```

---

## Updating Your Application

When you push changes to GitHub:
```bash
# SSH into server
ssh nutriediet@your_droplet_ip

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
tail -f logs/app.log
```

---

## Total Time: ~2 hours

✅ Your application is now live at: `https://yourdomain.com`

## What You Have:
- ✅ Secure HTTPS with auto-renewal
- ✅ Rate limiting on all endpoints
- ✅ Nginx reverse proxy
- ✅ MySQL database (local)
- ✅ Automatic daily backups
- ✅ Firewall protection
- ✅ systemd service management
- ✅ Production-ready environment

## Next Steps:
1. Read `PRODUCTION_IMPROVEMENTS.md` for additional security enhancements
2. Implement remaining critical fixes (JWT secret, rate limiting in Go code)
3. Set up monitoring and alerts
4. Configure log rotation
5. Test all API endpoints

---

## Troubleshooting

**Application won't start:**
```bash
sudo systemctl status nutriediet
sudo journalctl -u nutriediet -n 50
```

**Database connection issues:**
```bash
mysql -u nutriediet_app -p nutriediet_production
sudo systemctl status mysql
```

**SSL certificate issues:**
```bash
sudo certbot certificates
sudo certbot renew
sudo nginx -t
```

**Can't access via domain:**
```bash
# Check DNS
nslookup yourdomain.com

# Check Nginx
sudo systemctl status nginx
sudo tail -f /var/log/nginx/error.log
```

---

## Support
- Full deployment guide: `DIGITAL_OCEAN_DEPLOYMENT.md`
- Security improvements: `PRODUCTION_IMPROVEMENTS.md`
- Quick security fixes: `SECURITY_QUICK_FIXES.md`

**You're all set! 🚀**


# Nutriediet New App Deployment Guide

## Overview
This guide covers deploying the new Go backend + React frontend to `www.nutriediet.com/new` on an existing Digital Ocean droplet while maintaining the existing Node.js application.

## Architecture
```
www.nutriediet.com/
├── /              → Existing Node.js app (port 2299) [UNCHANGED]
├── /libs/         → Existing static files [UNCHANGED]
├── /uploads/      → Existing static files [UNCHANGED]
└── /new/          → New application
    ├── /          → React frontend (static files served by Nginx)
    ├── /api/      → Go backend (port 8080)
    ├── /images/   → Go backend image uploads
    └── /static/   → React build assets
```

## Pre-Deployment Checklist

### 1. Current Setup Verification
- [x] Existing Node.js app running on port 2299
- [x] PM2 managing existing app (name: "app")
- [x] Nginx configured with SSL (Let's Encrypt)
- [x] MySQL database running locally
- [x] User: sk (non-root)

### 2. Prerequisites
- GitHub repositories accessible
- MySQL root password
- SSH access to droplet as user 'sk'
- Sudo privileges for Nginx and system packages

### 3. Ports Required
- 8080: Go backend (new)
- 2299: Node.js backend (existing, unchanged)
- React frontend served as static files (no port needed)

## Deployment Options

### Option A: Automated Deployment (Recommended)

1. **Copy deployment files to droplet:**
```bash
# On your local machine
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
scp -r deployment/ sk@your-droplet-ip:/home/sk/
```

2. **SSH into droplet:**
```bash
ssh sk@your-droplet-ip
```

3. **Run deployment script:**
```bash
cd /home/sk/deployment
chmod +x deploy.sh
./deploy.sh
```

The script will:
- Install Go 1.21.5 if not present
- Upgrade Node.js to v20 if needed
- Create MySQL database and user
- Clone repositories
- Build backend and frontend
- Configure PM2
- Guide you through Nginx configuration

### Option B: Manual Deployment

Follow the steps in the "Manual Deployment Steps" section below.

## Manual Deployment Steps

### Step 1: Install Go

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
source ~/.bashrc

# Verify
go version
```

### Step 2: Upgrade Node.js (Optional but Recommended)

```bash
# Install/update nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
source ~/.bashrc

# Install Node 20
nvm install 20
nvm use 20
nvm alias default 20

# Verify
node -v  # Should show v20.x.x
npm -v
```

### Step 3: Create Directory Structure

```bash
mkdir -p /home/sk/mys/nutriediet-new/{backend,frontend,logs}
mkdir -p /home/sk/mys/nutriediet-new/backend/images
```

### Step 4: Setup MySQL Database

```bash
mysql -u root -p
```

```sql
CREATE DATABASE nutriediet_new_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'nutriediet_new_user'@'localhost' IDENTIFIED BY 'STRONG_PASSWORD_HERE';
GRANT ALL PRIVILEGES ON nutriediet_new_db.* TO 'nutriediet_new_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### Step 5: Clone and Build Backend

```bash
# Clone repository
cd /home/sk/mys/nutriediet-new
git clone https://github.com/cd-Ishita/nutriediet-go.git backend
cd backend

# Create .env file
cat > .env <<'EOF'
PORT=8080
GIN_MODE=release

DB_HOST=localhost
DB_PORT=3306
DB_USER=nutriediet_new_user
DB_PASSWORD=YOUR_DB_PASSWORD_HERE
DB_NAME=nutriediet_new_db

JWT_SECRET=CHANGE_TO_RANDOM_64_CHAR_STRING
JWT_EXPIRY=24h

ALLOWED_ORIGINS=https://nutriediet.com,https://www.nutriediet.com

UPLOAD_DIR=/home/sk/mys/nutriediet-new/backend/images
MAX_UPLOAD_SIZE=10485760

APP_ENV=production
APP_URL=https://nutriediet.com/new
API_URL=https://nutriediet.com/new/api
EOF

# Set secure permissions
chmod 600 .env

# Download dependencies and build
go mod download
go mod verify

# Build the binary (note: the dot '.' at the end is required, not a dash '-')
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o nutriediet-go -ldflags="-s -w" .
chmod +x nutriediet-go

# Run migrations (if applicable)
go run migrate/migrate.go
```

### Step 6: Build React Frontend

```bash
# Clone/copy frontend
cd /home/sk/mys/nutriediet-new
# Option 1: Clone from GitHub
git clone YOUR_FRONTEND_REPO frontend

# Option 2: Copy from local
# scp -r /Users/ishitagupta/Documents/Personal/frontend/* sk@droplet-ip:/home/sk/mys/nutriediet-new/frontend/

cd frontend

# Update package.json - add homepage field
nano package.json
# Add: "homepage": "/new",

# Create production environment file
cat > .env.production <<'EOF'
REACT_APP_API_URL=/new/api
PUBLIC_URL=/new
NODE_ENV=production
EOF

# Install dependencies and build
npm ci
GENERATE_SOURCEMAP=false npm run build
```

### Step 7: Configure PM2

```bash
# Install PM2 globally if not present
npm install -g pm2

# Create ecosystem config
cat > /home/sk/mys/nutriediet-new/ecosystem.config.js <<'EOF'
module.exports = {
  apps: [
    {
      name: 'nutriediet-go-api',
      script: './nutriediet-go',
      cwd: '/home/sk/mys/nutriediet-new/backend',
      instances: 1,
      exec_mode: 'fork',
      autorestart: true,
      watch: false,
      max_memory_restart: '500M',
      env: {
        PORT: '8080',
        GIN_MODE: 'release'
      },
      error_file: '/home/sk/mys/nutriediet-new/logs/go-api-error.log',
      out_file: '/home/sk/mys/nutriediet-new/logs/go-api-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss Z',
      merge_logs: true
    }
  ]
};
EOF

# Start the application
cd /home/sk/mys/nutriediet-new
pm2 start ecosystem.config.js

# Save PM2 configuration
pm2 save

# Enable PM2 startup (if not already done)
pm2 startup systemd -u sk --hp /home/sk

# Verify
pm2 list
```

### Step 8: Update Nginx Configuration

```bash
# Backup current config
sudo cp /etc/nginx/sites-available/nutriediet.com /etc/nginx/sites-available/nutriediet.com.backup

# Edit configuration
sudo nano /etc/nginx/sites-available/nutriediet.com
```

Replace with the configuration from `deployment/nginx-config-new.conf` or manually add these location blocks BEFORE the default `location /` block:

```nginx
# API endpoints for new Go backend
location /new/api/ {
    proxy_pass http://localhost:8080/;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_cache_bypass $http_upgrade;
}

# Static images for new Go backend
location /new/images/ {
    alias /home/sk/mys/nutriediet-new/backend/images/;
    expires 7d;
    add_header Cache-Control "public, no-cache";
}

# React build static files
location /new/static/ {
    alias /home/sk/mys/nutriediet-new/frontend/build/static/;
    expires 1y;
    add_header Cache-Control "public, immutable";
}

# React app root
location /new/ {
    alias /home/sk/mys/nutriediet-new/frontend/build/;
    try_files $uri $uri/ /new/index.html;
    
    location = /new/index.html {
        alias /home/sk/mys/nutriediet-new/frontend/build/index.html;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }
}
```

**Test and reload:**
```bash
# Test configuration
sudo nginx -t

# If test passes, reload (not restart - zero downtime)
sudo systemctl reload nginx

# If there are errors, restore backup
# sudo cp /etc/nginx/sites-available/nutriediet.com.backup /etc/nginx/sites-available/nutriediet.com
```

### Step 9: Verify Deployment

```bash
# Check if Go API is running
curl http://localhost:8080

# Check if existing app still works
curl http://localhost:2299

# Check PM2 status
pm2 list
pm2 logs nutriediet-go-api --lines 50

# Check nginx logs
sudo tail -f /var/log/nginx/error.log
```

**Browser tests:**
1. Visit `https://nutriediet.com` - should show existing app
2. Visit `https://nutriediet.com/new` - should show new React app
3. Test API: `https://nutriediet.com/new/api/health` or similar endpoint

## Important Notes

### React App Configuration
Your current frontend is Create React App (CRA), not Next.js. The deployment serves it as static files, which is simpler and more efficient than Next.js for this use case.

**Key files to update before building:**
1. `package.json` - add `"homepage": "/new"`
2. `.env.production` - set `PUBLIC_URL=/new` and `REACT_APP_API_URL=/new/api`
3. Update API calls in your React code to use relative paths or the environment variable

### API Configuration
Update your Go backend's CORS configuration to allow requests from your domain:

```go
// In main.go
config := cors.Config{
    AllowOrigins: []string{
        "https://nutriediet.com",
        "https://www.nutriediet.com",
    },
    // ... other settings
}
```

### Database Migrations
If you have database migrations in `migrate/migrate.go`, run them:
```bash
cd /home/sk/mys/nutriediet-new/backend
go run migrate/migrate.go
```

## Troubleshooting

### Go API won't start
```bash
# Check logs
pm2 logs nutriediet-go-api

# Check if port 8080 is in use
sudo netstat -tlnp | grep 8080

# Test binary directly
cd /home/sk/mys/nutriediet-new/backend
./nutriediet-go
```

### React app shows 404
- Verify `homepage` field in `package.json`
- Check nginx configuration for `/new/` location
- Verify build files exist: `ls /home/sk/mys/nutriediet-new/frontend/build/`

### API calls fail
- Check CORS configuration in Go backend
- Verify nginx proxy_pass for `/new/api/`
- Check if Go API is running: `pm2 list`

### Existing site broken
```bash
# Restore nginx backup
sudo cp /etc/nginx/sites-available/nutriediet.com.backup /etc/nginx/sites-available/nutriediet.com
sudo nginx -t
sudo systemctl reload nginx

# Check if PM2 app is running
pm2 list
pm2 restart app
```

## Useful Commands

### PM2 Management
```bash
# List all apps
pm2 list

# View logs
pm2 logs nutriediet-go-api
pm2 logs nutriediet-go-api --lines 100
pm2 logs nutriediet-go-api --err

# Restart
pm2 restart nutriediet-go-api

# Stop
pm2 stop nutriediet-go-api

# Delete
pm2 delete nutriediet-go-api

# Monitor
pm2 monit
```

### Nginx Management
```bash
# Test configuration
sudo nginx -t

# Reload (zero downtime)
sudo systemctl reload nginx

# Restart
sudo systemctl restart nginx

# View logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Database Management
```bash
# Connect to database
mysql -u nutriediet_new_user -p nutriediet_new_db

# Backup database
mysqldump -u nutriediet_new_user -p nutriediet_new_db > backup.sql

# Restore database
mysql -u nutriediet_new_user -p nutriediet_new_db < backup.sql
```

## Future Updates

### Backend Updates
```bash
cd /home/sk/mys/nutriediet-new/backend
git pull
# Build the binary (note: the dot '.' at the end is required, not a dash '-')
go build -o nutriediet-go -ldflags="-s -w" .
pm2 restart nutriediet-go-api
```

### Frontend Updates
```bash
cd /home/sk/mys/nutriediet-new/frontend
git pull
npm ci
GENERATE_SOURCEMAP=false npm run build
# No PM2 restart needed - files are served directly by Nginx
```

## Security Checklist
- [x] Database user has limited privileges
- [x] `.env` file has secure permissions (600)
- [x] SSL/TLS enabled via Let's Encrypt
- [x] Security headers configured in Nginx
- [x] CORS properly configured
- [x] JWT secrets are randomly generated
- [x] File upload directory properly secured

## Rollback Plan
If anything goes wrong:
1. Restore Nginx config backup
2. Stop new PM2 app: `pm2 delete nutriediet-go-api`
3. Reload Nginx: `sudo systemctl reload nginx`
4. Verify existing app: `curl http://localhost:2299`

The existing application on port 2299 remains untouched throughout the deployment.

## Support
If you encounter issues:
1. Check PM2 logs: `pm2 logs nutriediet-go-api`
2. Check Nginx logs: `sudo tail -f /var/log/nginx/error.log`
3. Verify ports: `sudo netstat -tlnp`
4. Check PM2 status: `pm2 list`


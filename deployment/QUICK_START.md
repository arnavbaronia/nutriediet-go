# Quick Start Deployment Guide

## 🚀 Fast Track Deployment (15 minutes)

### Prerequisites
- SSH access to Digital Ocean droplet as user `sk`
- MySQL root password
- GitHub repositories accessible

### Step-by-Step Commands

#### 1. Copy Files to Server
```bash
# On your local machine
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
scp -r deployment/ sk@YOUR_DROPLET_IP:/home/sk/nutriediet-deployment/
```

#### 2. SSH into Server
```bash
ssh sk@YOUR_DROPLET_IP
```

#### 3. Run Automated Deployment
```bash
cd /home/sk/nutriediet-deployment
chmod +x deploy.sh
./deploy.sh
```

The script will prompt you for:
- MySQL root password
- New database user password
- Confirmation to update Nginx

#### 4. Manual Nginx Update
After the script completes:

```bash
# Backup current config
sudo cp /etc/nginx/sites-available/nutriediet.com /etc/nginx/sites-available/nutriediet.com.backup

# Edit config
sudo nano /etc/nginx/sites-available/nutriediet.com
```

Add these location blocks BEFORE the `location /` block:

```nginx
# New Go + React app at /new
location /new/api/ {
    proxy_pass http://localhost:8080/;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

location /new/images/ {
    alias /home/sk/mys/nutriediet-new/backend/images/;
    expires 7d;
}

location /new/static/ {
    alias /home/sk/mys/nutriediet-new/frontend/build/static/;
    expires 1y;
    add_header Cache-Control "public, immutable";
}

location /new/ {
    alias /home/sk/mys/nutriediet-new/frontend/build/;
    try_files $uri $uri/ /new/index.html;
}
```

Test and reload:
```bash
sudo nginx -t
sudo systemctl reload nginx
```

#### 5. Verify Deployment
```bash
# Check PM2 status
pm2 list

# Test Go API
curl http://localhost:8080

# Test existing app (should still work)
curl http://localhost:2299
```

Browser tests:
- Old app: https://nutriediet.com
- New app: https://nutriediet.com/new
- New API: https://nutriediet.com/new/api

## ✅ Success Checklist
- [ ] Go API running on port 8080 (check with `pm2 list`)
- [ ] Existing Node app still on port 2299 (check with `pm2 list`)
- [ ] Nginx reloaded successfully
- [ ] New app accessible at /new
- [ ] Old app still works at /
- [ ] No errors in logs: `pm2 logs nutriediet-go-api`

## 🔧 Troubleshooting

### Go API not starting
```bash
pm2 logs nutriediet-go-api
cd /home/sk/mys/nutriediet-new/backend
./nutriediet-go  # Test binary directly
```

### React app 404
```bash
ls /home/sk/mys/nutriediet-new/frontend/build/
# Should show index.html and static/ folder
```

### Existing site broken
```bash
# Rollback nginx
sudo cp /etc/nginx/sites-available/nutriediet.com.backup /etc/nginx/sites-available/nutriediet.com
sudo systemctl reload nginx
```

## 📝 Configuration Files

All deployment files are in `/Users/ishitagupta/Documents/Personal/nutriediet-go/deployment/`:

- `deploy.sh` - Automated deployment script
- `nginx-config-new.conf` - Complete Nginx configuration
- `ecosystem.config.js` - PM2 configuration
- `.env.production.template` - Backend environment template
- `frontend-env-production` - Frontend environment variables
- `DEPLOYMENT_GUIDE.md` - Comprehensive guide

## 🎯 What Gets Deployed

```
Digital Ocean Droplet
├── /home/sk/mys/nutriediet-new/
│   ├── backend/
│   │   ├── nutriediet-go (binary)
│   │   ├── .env (config)
│   │   └── images/ (uploads)
│   ├── frontend/
│   │   └── build/ (static files)
│   ├── logs/
│   └── ecosystem.config.js
│
├── MySQL
│   └── nutriediet_new_db (new database)
│
├── PM2
│   ├── app (existing - unchanged)
│   └── nutriediet-go-api (new - port 8080)
│
└── Nginx
    ├── / → localhost:2299 (existing)
    ├── /new/ → static files
    └── /new/api/ → localhost:8080 (new)
```

## 🔄 Future Updates

### Backend Update
```bash
cd /home/sk/mys/nutriediet-new/backend
git pull
go build -o nutriediet-go .
pm2 restart nutriediet-go-api
```

### Frontend Update
```bash
cd /home/sk/mys/nutriediet-new/frontend
git pull
npm ci
npm run build
# Nginx serves static files - no restart needed
```

## 📞 Need Help?

See the full `DEPLOYMENT_GUIDE.md` for:
- Manual deployment steps
- Detailed troubleshooting
- Security checklist
- Rollback procedures


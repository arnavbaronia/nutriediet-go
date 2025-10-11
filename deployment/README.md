# Nutriediet New App - Deployment Files

This directory contains all files needed to deploy the new Go backend + React frontend to www.nutriediet.com/new on your Digital Ocean droplet.

## 📁 Files Overview

| File | Purpose |
|------|---------|
| `deploy.sh` | Automated deployment script (recommended) |
| `QUICK_START.md` | Fast-track deployment guide (15 min) |
| `DEPLOYMENT_GUIDE.md` | Comprehensive deployment documentation |
| `nginx-config-new.conf` | Complete Nginx configuration file |
| `ecosystem.config.js` | PM2 process manager configuration |
| `.env.production.template` | Backend environment variables template |
| `frontend-env-production` | Frontend environment variables |
| `package-json-update.txt` | Instructions for updating package.json |
| `test-deployment.sh` | Post-deployment verification script |

## 🚀 Quick Start (Choose One Method)

### Method 1: Automated Deployment (Recommended)
```bash
# Copy files to server
scp -r deployment/ sk@YOUR_DROPLET_IP:/home/sk/nutriediet-deployment/

# SSH and run
ssh sk@YOUR_DROPLET_IP
cd /home/sk/nutriediet-deployment
chmod +x deploy.sh
./deploy.sh
```

See `QUICK_START.md` for detailed steps.

### Method 2: Manual Deployment
Follow step-by-step instructions in `DEPLOYMENT_GUIDE.md`.

## ⚠️ Important Notes

### Current Frontend is React (CRA), Not Next.js
Your frontend uses Create React App, which is simpler to deploy than Next.js:
- No server-side rendering needed
- Builds to static files served by Nginx
- No separate frontend server process
- Better performance for this use case

### Zero Downtime Deployment
The deployment:
- ✅ Does NOT touch existing Node.js app on port 2299
- ✅ Does NOT modify existing PM2 "app" process
- ✅ Uses separate MySQL database
- ✅ Uses Nginx reload (not restart) for zero downtime

### Before Deployment

1. **Update Frontend package.json:**
   - Add `"homepage": "/new"` field
   - See `package-json-update.txt`

2. **Update Frontend API calls:**
   - Use environment variable: `process.env.REACT_APP_API_URL`
   - Or use relative path: `/new/api`

3. **Update Go CORS settings:**
   ```go
   AllowOrigins: []string{
       "https://nutriediet.com",
       "https://www.nutriediet.com",
   }
   ```

## 📋 Deployment Checklist

- [ ] SSH access to droplet as user 'sk'
- [ ] MySQL root password available
- [ ] GitHub repository accessible
- [ ] Verified existing app is running (port 2299)
- [ ] Updated frontend package.json with homepage
- [ ] Updated API base URL in frontend
- [ ] Updated CORS in Go backend

## 🧪 After Deployment

Run verification tests:
```bash
cd /home/sk/nutriediet-deployment
chmod +x test-deployment.sh
./test-deployment.sh
```

Manual verification:
- Visit https://nutriediet.com (existing app)
- Visit https://nutriediet.com/new (new app)
- Check PM2: `pm2 list`
- Check logs: `pm2 logs nutriediet-go-api`

## 🔧 Common Issues

### Issue: React app shows blank page
**Solution:** Check browser console for 404 errors on static assets
```bash
# Verify homepage in package.json
cd /home/sk/mys/nutriediet-new/frontend
cat package.json | grep homepage

# Should show: "homepage": "/new"
```

### Issue: API calls fail with CORS error
**Solution:** Update CORS in Go backend
```bash
cd /home/sk/mys/nutriediet-new/backend
# Edit main.go to add your domain to AllowOrigins
nano main.go
go build -o nutriediet-go .
pm2 restart nutriediet-go-api
```

### Issue: Go API won't start
**Solution:** Check logs and environment
```bash
pm2 logs nutriediet-go-api
cd /home/sk/mys/nutriediet-new/backend
cat .env  # Verify configuration
./nutriediet-go  # Test binary directly
```

### Issue: Existing site is down
**Solution:** Rollback Nginx configuration
```bash
sudo cp /etc/nginx/sites-available/nutriediet.com.backup /etc/nginx/sites-available/nutriediet.com
sudo nginx -t
sudo systemctl reload nginx
```

## 📖 Documentation Structure

1. **Start here:** `README.md` (this file)
2. **Fast deployment:** `QUICK_START.md`
3. **Need details:** `DEPLOYMENT_GUIDE.md`
4. **After deployment:** `test-deployment.sh`

## 🏗️ Architecture Overview

```
www.nutriediet.com
├── /                  → Existing Node.js (port 2299) ✅ UNCHANGED
├── /libs/            → Existing static files ✅ UNCHANGED
├── /uploads/         → Existing static files ✅ UNCHANGED
└── /new/             → NEW APPLICATION
    ├── /             → React static files (Nginx)
    ├── /api/         → Go backend (port 8080)
    ├── /images/      → Go backend uploads
    └── /static/      → React assets (CSS, JS)
```

## 🔐 Security Notes

- Database user has limited privileges (only nutriediet_new_db)
- `.env` file has secure permissions (600)
- SSL/TLS via Let's Encrypt (existing)
- Security headers configured in Nginx
- JWT secrets randomly generated
- CORS properly configured

## 📞 Support

If you encounter issues:
1. Check `DEPLOYMENT_GUIDE.md` troubleshooting section
2. Run `test-deployment.sh` for diagnostics
3. Check logs: `pm2 logs nutriediet-go-api`
4. Check Nginx logs: `sudo tail -f /var/log/nginx/error.log`

## 🔄 Updating After Deployment

### Backend Update
```bash
cd /home/sk/mys/nutriediet-new/backend
git pull
go build -o nutriediet-go -ldflags="-s -w" .
pm2 restart nutriediet-go-api
```

### Frontend Update
```bash
cd /home/sk/mys/nutriediet-new/frontend
git pull
npm ci
GENERATE_SOURCEMAP=false npm run build
# No restart needed - Nginx serves static files
```

## 📊 What Gets Installed/Changed

### New Installations
- Go 1.21.5 (if not present)
- Node.js v20 (upgrade from v14)
- PM2 (if not present)

### New Directories
- `/home/sk/mys/nutriediet-new/` - All new app files
- `/home/sk/mys/nutriediet-new/backend/` - Go backend
- `/home/sk/mys/nutriediet-new/frontend/` - React build
- `/home/sk/mys/nutriediet-new/logs/` - Application logs

### New Processes
- PM2: `nutriediet-go-api` (port 8080)

### Modified Files
- `/etc/nginx/sites-available/nutriediet.com` - Add location blocks for /new

### New Database
- `nutriediet_new_db` - Separate from existing database
- `nutriediet_new_user` - Limited privileges

## ✅ Success Criteria

Deployment is successful when:
- ✅ Existing site works at https://nutriediet.com
- ✅ New app works at https://nutriediet.com/new
- ✅ New API works at https://nutriediet.com/new/api
- ✅ PM2 shows both apps running
- ✅ No errors in logs
- ✅ All tests pass in `test-deployment.sh`

---

**Ready to deploy?** Start with `QUICK_START.md` for fast deployment or `DEPLOYMENT_GUIDE.md` for detailed instructions.


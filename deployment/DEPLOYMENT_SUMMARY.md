# Deployment Summary

## 📦 What You Have

A complete deployment package for deploying your Go + React application to `www.nutriediet.com/new` on your Digital Ocean droplet.

## 📂 Deployment Package Contents

```
deployment/
├── README.md                          # Start here - Overview and quick links
├── QUICK_START.md                     # 15-minute fast-track deployment guide
├── DEPLOYMENT_GUIDE.md                # Comprehensive step-by-step guide
├── PRE_DEPLOYMENT_CHECKLIST.md        # Must complete before deploying
├── EXACT_CODE_CHANGES.md              # Line-by-line code changes needed
│
├── deploy.sh                          # Automated deployment script ⭐
├── test-deployment.sh                 # Post-deployment verification
│
├── nginx-config-new.conf              # Complete Nginx configuration
├── ecosystem.config.js                # PM2 process manager config
├── .env.production.template           # Backend environment template
├── frontend-env-production            # Frontend environment variables
│
├── cors-update.md                     # Go backend CORS guide
├── frontend-axios-update.md           # Frontend API config guide
├── frontend-constants-update.md       # Frontend constants guide
└── package-json-update.txt            # package.json changes
```

## 🚀 Three Deployment Options

### Option 1: Automated (Recommended) ⭐
**Time: ~15 minutes**

1. Make code changes from `EXACT_CODE_CHANGES.md`
2. Complete `PRE_DEPLOYMENT_CHECKLIST.md`
3. Run deployment script:
   ```bash
   scp -r deployment/ sk@YOUR_DROPLET_IP:/home/sk/nutriediet-deployment/
   ssh sk@YOUR_DROPLET_IP
   cd /home/sk/nutriediet-deployment
   ./deploy.sh
   ```

**Best for:** Most users, production deployments

### Option 2: Quick Manual Deployment
**Time: ~20 minutes**

Follow `QUICK_START.md` for streamlined manual steps.

**Best for:** Those who want control but quick deployment

### Option 3: Comprehensive Manual
**Time: ~45 minutes**

Follow `DEPLOYMENT_GUIDE.md` for detailed step-by-step instructions.

**Best for:** Learning, customization, troubleshooting

## ⚡ Quick Start (Right Now)

### Step 1: Make Code Changes (5 minutes)

Open `EXACT_CODE_CHANGES.md` and make these 5 changes:

1. ✏️ `frontend/src/utils/constants.js` (2 changes)
2. ✏️ `frontend/package.json` (1 addition)
3. ✏️ `frontend/src/App.js` (1 change)
4. ✏️ `backend/main.go` (1 change)
5. ➕ Create `frontend/.env.production`

### Step 2: Test Locally (5 minutes)

```bash
# Test backend
cd nutriediet-go
go run main.go

# Test frontend build
cd frontend
npm run build
npx serve -s build -l 3001
```

### Step 3: Commit & Push (2 minutes)

```bash
# Backend
cd nutriediet-go
git add .
git commit -m "Configure for /new subpath deployment"
git push

# Frontend (update repo URL in deploy.sh first!)
cd frontend
git add .
git commit -m "Configure for /new subpath deployment"
git push
```

### Step 4: Deploy (15 minutes)

```bash
# Copy deployment files
cd nutriediet-go
scp -r deployment/ sk@YOUR_DROPLET_IP:/home/sk/nutriediet-deployment/

# SSH and run
ssh sk@YOUR_DROPLET_IP
cd /home/sk/nutriediet-deployment
chmod +x *.sh
./deploy.sh
```

### Step 5: Verify (2 minutes)

```bash
./test-deployment.sh
```

**Total time:** ~30 minutes

## 🎯 What Gets Deployed

### New Infrastructure
```
/home/sk/mys/nutriediet-new/
├── backend/
│   ├── nutriediet-go (Go binary)
│   ├── .env (configuration)
│   └── images/ (uploads)
├── frontend/
│   └── build/ (React static files)
├── logs/
│   ├── go-api-error.log
│   └── go-api-out.log
└── ecosystem.config.js
```

### New Processes
- PM2 process: `nutriediet-go-api` (port 8080)

### New Database
- Database: `nutriediet_new_db`
- User: `nutriediet_new_user`

### Updated Config
- Nginx: Add location blocks for `/new/*`

### Existing (UNCHANGED)
- ✅ Node.js app (port 2299)
- ✅ PM2 app "app"
- ✅ Existing database
- ✅ SSL certificates
- ✅ /libs/ and /uploads/ paths

## 🔒 Safety Features

### Zero Downtime
- Existing site remains running during deployment
- Nginx reload (not restart) = no connection drops
- Separate database = no data conflicts

### Isolated Deployment
- New app in separate directory
- New database with limited user
- New PM2 process with different name
- New ports (8080)

### Easy Rollback
If anything goes wrong:
```bash
pm2 delete nutriediet-go-api
sudo cp /etc/nginx/sites-available/nutriediet.com.backup /etc/nginx/sites-available/nutriediet.com
sudo systemctl reload nginx
```
Existing site remains unaffected!

## 📋 Prerequisites

Before starting, you need:

- [ ] SSH access to droplet as user `sk`
- [ ] MySQL root password
- [ ] Sudo password (for Nginx changes)
- [ ] GitHub repositories accessible
- [ ] Frontend repository URL (update in `deploy.sh`)

## ⚠️ Important Notes

### 1. Frontend is React (not Next.js)
Your frontend is Create React App, which is actually simpler:
- Builds to static files
- Served directly by Nginx
- No separate server process needed
- Better performance for your use case

### 2. Update Frontend Repo URL
**Before deploying**, edit `deployment/deploy.sh`:

```bash
GITHUB_FRONTEND_REPO="YOUR_FRONTEND_REPO_URL"
```

Replace with your actual frontend repository URL.

### 3. CORS Configuration
The Go backend needs to allow requests from nutriediet.com. The code changes in `EXACT_CODE_CHANGES.md` handle this.

## 🗺️ Document Navigation

**Start Here:**
1. Read `README.md` (this file)
2. Complete `PRE_DEPLOYMENT_CHECKLIST.md`
3. Make changes from `EXACT_CODE_CHANGES.md`

**For Deployment:**
- Fast: `QUICK_START.md`
- Detailed: `DEPLOYMENT_GUIDE.md`
- Automated: Run `deploy.sh`

**For Configuration:**
- Backend: `cors-update.md`
- Frontend: `frontend-constants-update.md`
- Code: `EXACT_CODE_CHANGES.md`

**For Testing:**
- After deployment: `test-deployment.sh`
- Troubleshooting: See `DEPLOYMENT_GUIDE.md`

## 🎓 Architecture Understanding

### Request Flow
```
User types: https://nutriediet.com/new
                    ↓
            Nginx (Port 443)
            ↙              ↘
      Static Files      /new/api/*
    (React build)     Proxy to Go API
         ↓              (Port 8080)
    Served by Nginx         ↓
                     Go Backend
                          ↓
                   MySQL Database
```

### File Serving
```
/new/              → /home/sk/mys/nutriediet-new/frontend/build/index.html
/new/static/       → /home/sk/mys/nutriediet-new/frontend/build/static/*
/new/api/*         → http://localhost:8080/* (Go backend)
/new/images/*      → /home/sk/mys/nutriediet-new/backend/images/*
```

### Process Management
```
PM2
├── app (existing)
│   └── Port 2299 → Node.js monolith
│
└── nutriediet-go-api (new)
    └── Port 8080 → Go API
```

## 📞 Getting Help

### During Code Changes
- See `EXACT_CODE_CHANGES.md` for exact changes
- See specific update guides for details

### During Deployment
- Check `DEPLOYMENT_GUIDE.md` troubleshooting section
- Run `test-deployment.sh` for diagnostics

### After Deployment
- View logs: `pm2 logs nutriediet-go-api`
- Check status: `pm2 list`
- Nginx logs: `sudo tail -f /var/log/nginx/error.log`

### Common Issues
- **404 on assets**: Check `homepage` in package.json
- **API CORS errors**: Check Go CORS config
- **Go won't start**: Check `.env` file and logs
- **Existing site down**: Rollback Nginx config

## ✅ Success Criteria

Deployment is successful when:

1. **Existing site works**: https://nutriediet.com
2. **New site works**: https://nutriediet.com/new
3. **New API works**: https://nutriediet.com/new/api
4. **Both PM2 apps running**: `pm2 list` shows 2 apps
5. **No errors in logs**: `pm2 logs nutriediet-go-api`
6. **Test script passes**: `./test-deployment.sh` all green

## 🎉 After Successful Deployment

### Test Everything
- [ ] Login/logout on new site
- [ ] Navigate through all pages
- [ ] Test API calls
- [ ] Upload images (if applicable)
- [ ] Check mobile responsiveness

### Monitor
```bash
# Watch logs
pm2 logs nutriediet-go-api

# Monitor resources
pm2 monit

# Check status
pm2 list
```

### Update Later
**Backend:**
```bash
cd /home/sk/mys/nutriediet-new/backend
git pull
go build -o nutriediet-go .
pm2 restart nutriediet-go-api
```

**Frontend:**
```bash
cd /home/sk/mys/nutriediet-new/frontend
git pull
npm ci
npm run build
# No restart needed - Nginx serves static files
```

## 📈 Next Steps

1. **Complete PRE_DEPLOYMENT_CHECKLIST.md** ✅
2. **Make code changes from EXACT_CODE_CHANGES.md** ✅
3. **Test locally** ✅
4. **Choose deployment method** (Automated recommended)
5. **Deploy** 🚀
6. **Verify with test-deployment.sh** ✅
7. **Celebrate** 🎉

## 🚨 Emergency Rollback

If something goes seriously wrong:

```bash
# Stop new app
pm2 delete nutriediet-go-api

# Restore Nginx
sudo cp /etc/nginx/sites-available/nutriediet.com.backup /etc/nginx/sites-available/nutriediet.com
sudo nginx -t
sudo systemctl reload nginx

# Verify existing app
curl http://localhost:2299
pm2 list
```

The existing app remains untouched throughout!

---

**Ready to begin?** Start with `PRE_DEPLOYMENT_CHECKLIST.md` and `EXACT_CODE_CHANGES.md`.

**Questions?** See `DEPLOYMENT_GUIDE.md` for comprehensive documentation.

**Let's deploy!** 🚀


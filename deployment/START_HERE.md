# 🚀 Start Here - Deployment Quick Navigation

Welcome! This guide will get you deploying in the fastest way possible.

## 🎯 Your Goal

Deploy your Go backend + React frontend to:
- **URL:** https://nutriediet.com/new
- **Server:** Existing Digital Ocean droplet
- **Impact:** Zero downtime for existing site

## ⏱️ Time Estimate

- **Code changes:** 5 minutes
- **Local testing:** 5 minutes
- **Deployment:** 15 minutes
- **Total:** ~25 minutes

## 📍 You Are Here

```
┌─────────────────────────────────────────────────────────────┐
│  START HERE                                                 │
│  ↓                                                          │
│  1. Make code changes (5 min)                              │
│  2. Test locally (5 min)                                   │
│  3. Deploy (15 min)                                        │
│  4. Verify (2 min)                                         │
│  ↓                                                          │
│  DONE! 🎉                                                   │
└─────────────────────────────────────────────────────────────┘
```

## 🎬 Step 1: Make Code Changes (5 minutes)

### Option A: Quick Summary (Experienced developers)

Make these 5 changes:
1. `frontend/src/utils/constants.js` - Update API_BASE_URL and ROUTES
2. `frontend/package.json` - Add `"homepage": "/new"`
3. `frontend/src/App.js` - Add basename to BrowserRouter
4. `backend/main.go` - Update CORS configuration
5. Create `frontend/.env.production`

**Full details:** Open `EXACT_CODE_CHANGES.md`

### Option B: I Need Help

Open these files in order:
1. `EXACT_CODE_CHANGES.md` ← Line-by-line instructions
2. `PRE_DEPLOYMENT_CHECKLIST.md` ← Verify nothing is missed

## 🧪 Step 2: Test Locally (5 minutes)

```bash
# Terminal 1: Test backend
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
go run main.go
# Should start on port 8080 without errors

# Terminal 2: Test frontend build
cd /Users/ishitagupta/Documents/Personal/frontend
npm run build
npx serve -s build -l 3001
# Visit http://localhost:3001 - should work
```

✅ **Pass criteria:**
- Backend starts without errors
- Frontend builds successfully
- Can navigate the site at localhost:3001
- No console errors in browser

## 📤 Step 3: Commit & Push (2 minutes)

```bash
# Backend
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
git add .
git commit -m "Configure for /new subpath deployment"
git push

# Frontend
cd /Users/ishitagupta/Documents/Personal/frontend
git add .
git commit -m "Configure for /new subpath deployment"
git push
```

## 🚀 Step 4: Deploy (15 minutes)

### Before You Deploy

**Required information:**
- [ ] Droplet IP address
- [ ] SSH access as user 'sk'
- [ ] MySQL root password
- [ ] Frontend GitHub repo URL

**Update deploy script:**
```bash
# Edit deployment/deploy.sh
# Line ~15: Update GITHUB_FRONTEND_REPO with your frontend repo URL
```

### Choose Your Deployment Method

#### Method 1: Automated (Recommended) ⭐

```bash
# From your local machine
cd /Users/ishitagupta/Documents/Personal/nutriediet-go

# Copy deployment files to server
scp -r deployment/ sk@YOUR_DROPLET_IP:/home/sk/nutriediet-deployment/

# SSH into server
ssh sk@YOUR_DROPLET_IP

# Run deployment
cd /home/sk/nutriediet-deployment
./deploy.sh
```

**Follow prompts:**
- Enter MySQL root password when asked
- Enter new database user password
- Review Nginx configuration when prompted
- Wait for completion

#### Method 2: Manual Quick Start

If you prefer manual control:
1. Open `QUICK_START.md`
2. Follow the step-by-step commands

#### Method 3: Comprehensive Manual

For detailed understanding:
1. Open `DEPLOYMENT_GUIDE.md`
2. Follow the complete guide

## ✅ Step 5: Verify (2 minutes)

### On the Server

```bash
# Still on the server
cd /home/sk/nutriediet-deployment
./test-deployment.sh
```

All tests should pass ✅

### In Your Browser

1. **Existing site:** https://nutriediet.com
   - Should work exactly as before
   
2. **New site:** https://nutriediet.com/new
   - Should show your React app
   
3. **Test new site:**
   - Try logging in
   - Navigate pages
   - Check browser console (no errors)

## 🎉 Success!

If all tests pass, congratulations! Your deployment is complete.

### Monitor Your App

```bash
# View logs
pm2 logs nutriediet-go-api

# Check status
pm2 list

# Monitor resources
pm2 monit
```

## 🆘 Something Went Wrong?

### Quick Fixes

**Go API won't start:**
```bash
pm2 logs nutriediet-go-api
cd /home/sk/mys/nutriediet-new/backend
cat .env  # Check configuration
```

**React app shows 404:**
- Check `homepage` in package.json
- Verify build files exist
- Check Nginx configuration

**Existing site is down (Urgent!):**
```bash
sudo cp /etc/nginx/sites-available/nutriediet.com.backup /etc/nginx/sites-available/nutriediet.com
sudo systemctl reload nginx
```

### Get Detailed Help

- **Troubleshooting:** See `DEPLOYMENT_GUIDE.md` Section: Troubleshooting
- **Code issues:** Review `EXACT_CODE_CHANGES.md`
- **Pre-flight:** Run `PRE_DEPLOYMENT_CHECKLIST.md`

## 📚 Full Documentation Map

```
START_HERE.md (you are here)
├── Quick Path
│   ├── EXACT_CODE_CHANGES.md ← Make these changes
│   ├── QUICK_START.md ← Fast deployment
│   └── test-deployment.sh ← Verify it works
│
├── Automated Path
│   ├── PRE_DEPLOYMENT_CHECKLIST.md ← Verify readiness
│   ├── deploy.sh ← Run this script
│   └── test-deployment.sh ← Verify it works
│
├── Manual Path
│   ├── PRE_DEPLOYMENT_CHECKLIST.md
│   ├── DEPLOYMENT_GUIDE.md ← Step-by-step
│   └── test-deployment.sh
│
├── Reference Guides
│   ├── cors-update.md ← Backend CORS
│   ├── frontend-constants-update.md ← Frontend config
│   ├── frontend-axios-update.md ← API calls
│   └── package-json-update.txt ← package.json
│
├── Configuration Files
│   ├── nginx-config-new.conf ← Nginx
│   ├── ecosystem.config.js ← PM2
│   ├── .env.production.template ← Backend env
│   └── frontend-env-production ← Frontend env
│
└── Overview
    ├── README.md ← Package overview
    └── DEPLOYMENT_SUMMARY.md ← High-level summary
```

## 🎓 Understanding Your Deployment

### What's Being Created

```
New Application at /new
├── React Frontend (static files via Nginx)
├── Go API (port 8080, managed by PM2)
├── MySQL Database (nutriediet_new_db)
└── Images/uploads directory

Existing Application (UNTOUCHED)
├── Node.js App (port 2299, managed by PM2)
├── Existing database
└── /libs/ and /uploads/ paths
```

### How It Works

```
User requests: https://nutriediet.com/new
        ↓
    Nginx (SSL)
    ↙        ↘
/new/     /new/api/*
Static     Go API
Files    (port 8080)
```

## 🔐 Safety Notes

✅ **Safe:**
- Existing site remains running throughout
- Separate database (no data mixing)
- Easy rollback if needed
- Zero downtime deployment

✅ **Isolated:**
- New directory: `/home/sk/mys/nutriediet-new/`
- New PM2 process: `nutriediet-go-api`
- New database: `nutriediet_new_db`
- New port: 8080

## ⚡ Common Questions

### Q: Will this affect my existing site?
**A:** No! The existing site on port 2299 remains completely untouched.

### Q: Can I test before going live?
**A:** Yes! You can build locally and test with `npx serve` before deploying.

### Q: What if something breaks?
**A:** Easy rollback - just stop the new PM2 app and restore Nginx. See DEPLOYMENT_GUIDE.md for steps.

### Q: Do I need to upgrade Node.js?
**A:** The deploy script can do this automatically from v14 to v20 for better React support.

### Q: How do I update after deployment?
**A:** See DEPLOYMENT_GUIDE.md "Future Updates" section.

### Q: My frontend is in a different repo?
**A:** Update `GITHUB_FRONTEND_REPO` in deploy.sh before running.

## 📝 Pre-Flight Checklist

Before you start, make sure:
- [ ] I have SSH access to the droplet
- [ ] I know the MySQL root password
- [ ] My code changes are ready
- [ ] I've tested locally
- [ ] I've committed and pushed to GitHub
- [ ] I've updated deploy.sh with frontend repo URL
- [ ] I have 15 minutes uninterrupted time

## 🎯 Next Action

Choose your path:

**Option 1: I want the fastest deployment (Recommended)**
→ Open `EXACT_CODE_CHANGES.md`, make changes, then run `deploy.sh`

**Option 2: I want to understand everything first**
→ Open `DEPLOYMENT_GUIDE.md` and read through

**Option 3: I need a checklist**
→ Open `PRE_DEPLOYMENT_CHECKLIST.md`

**Option 4: Show me a quick summary**
→ Open `DEPLOYMENT_SUMMARY.md`

---

## 🚀 Ready to Deploy?

1. **Make code changes:** Open `EXACT_CODE_CHANGES.md`
2. **Test locally:** Run commands above
3. **Deploy:** Run `deploy.sh` or follow `QUICK_START.md`
4. **Verify:** Run `test-deployment.sh`

**You've got this!** 💪

---

**Need help?** Every document has detailed troubleshooting sections.

**First time?** Start with `PRE_DEPLOYMENT_CHECKLIST.md` to ensure nothing is missed.

**Experienced?** Jump straight to `EXACT_CODE_CHANGES.md` then `deploy.sh`.


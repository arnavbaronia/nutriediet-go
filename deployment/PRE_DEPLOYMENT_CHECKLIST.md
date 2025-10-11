# Pre-Deployment Checklist

Complete all items in this checklist BEFORE running the deployment script.

## ✅ Local Development - Backend (Go)

### 1. Update CORS Configuration
- [ ] Open `main.go`
- [ ] Update CORS AllowOrigins to use environment variable
- [ ] See `cors-update.md` for detailed instructions
- [ ] Test locally: `go run main.go`

### 2. Verify Environment Variables
- [ ] Check `env.example` is up to date
- [ ] Ensure all required variables are documented
- [ ] Test with local `.env` file

### 3. Test Build
```bash
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
go mod download
go mod verify
go build -o nutriediet-go .
./nutriediet-go  # Should start without errors
```
- [ ] Build completes successfully
- [ ] Binary runs without errors
- [ ] Can connect to database
- [ ] API endpoints respond correctly

### 4. Database Migrations
- [ ] Verify migrations are ready in `migrate/migrate.go`
- [ ] Test migrations on local database
- [ ] Document any manual migration steps

## ✅ Local Development - Frontend (React)

### 1. Update package.json
```bash
cd /Users/ishitagupta/Documents/Personal/frontend
```
- [ ] Add `"homepage": "/new"` to package.json
- [ ] See `package-json-update.txt` for exact placement

### 2. Update Constants
- [ ] Open `src/utils/constants.js`
- [ ] Update `API_BASE_URL` to use environment variable
- [ ] Update all routes to use `BASE_PATH`
- [ ] See `frontend-constants-update.md` for details

### 3. Create Environment Files

Create `.env.production`:
```bash
cat > .env.production <<EOF
REACT_APP_API_URL=/new/api
PUBLIC_URL=/new
NODE_ENV=production
EOF
```
- [ ] `.env.production` created

Create `.env.development` (if not exists):
```bash
cat > .env.development <<EOF
REACT_APP_API_URL=http://localhost:8080
PUBLIC_URL=/
EOF
```
- [ ] `.env.development` created

### 4. Update React Router
- [ ] Open `src/App.js`
- [ ] Add basename to BrowserRouter: `<BrowserRouter basename={process.env.PUBLIC_URL}>`

### 5. Test Development Build
```bash
npm start
```
- [ ] Runs on http://localhost:3000
- [ ] No console errors
- [ ] Can connect to backend API
- [ ] All routes work
- [ ] Login/logout works

### 6. Test Production Build
```bash
GENERATE_SOURCEMAP=false npm run build
npx serve -s build -l 3001
```
- [ ] Build completes successfully
- [ ] Navigate to http://localhost:3001
- [ ] All assets load correctly
- [ ] No 404 errors in console
- [ ] Routing works
- [ ] API calls work (if backend is running)

## ✅ GitHub Repository

### 1. Backend Repository
- [ ] All changes committed and pushed
- [ ] `.env` file is in `.gitignore`
- [ ] `env.example` is included and up to date
- [ ] README has deployment instructions
- [ ] Repository URL: `https://github.com/cd-Ishita/nutriediet-go.git`

### 2. Frontend Repository
- [ ] All changes committed and pushed
- [ ] `.env.*` files are in `.gitignore`
- [ ] `homepage` field in package.json
- [ ] `build/` directory in `.gitignore`
- [ ] Repository URL: _________________ (fill in your frontend repo)

## ✅ Server Prerequisites

### 1. Access
- [ ] Have SSH credentials for user `sk`
- [ ] Can connect: `ssh sk@YOUR_DROPLET_IP`
- [ ] Have sudo password (for Nginx changes)

### 2. Existing Site Verification
```bash
ssh sk@YOUR_DROPLET_IP
curl http://localhost:2299
pm2 list  # Should show "app" running
```
- [ ] Existing site responds on port 2299
- [ ] PM2 shows "app" running
- [ ] Nginx is running: `sudo systemctl status nginx`

### 3. Database Access
- [ ] Have MySQL root password
- [ ] Can connect: `mysql -u root -p`
- [ ] Have chosen strong password for new database user

### 4. Ports Availability
```bash
sudo netstat -tlnp | grep -E '8080|3001'
```
- [ ] Port 8080 is available (not in use)
- [ ] Port 3001 can be skipped (React is static files)

## ✅ Deployment Files

### 1. Copy to Local Staging
```bash
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
ls deployment/
```

Should see:
- [ ] `deploy.sh`
- [ ] `nginx-config-new.conf`
- [ ] `ecosystem.config.js`
- [ ] `.env.production.template`
- [ ] `DEPLOYMENT_GUIDE.md`
- [ ] `QUICK_START.md`
- [ ] `test-deployment.sh`

### 2. Update Configuration
- [ ] Open `deployment/deploy.sh`
- [ ] Update `GITHUB_FRONTEND_REPO` variable with your frontend repository URL
- [ ] Review all configuration variables at top of script

### 3. Make Scripts Executable
```bash
chmod +x deployment/*.sh
```
- [ ] Scripts are executable

## ✅ Final Verification

### 1. Review Architecture
```
www.nutriediet.com
├── /              → Existing Node.js (port 2299) [NO CHANGES]
├── /libs/         → Existing static [NO CHANGES]
├── /uploads/      → Existing static [NO CHANGES]
└── /new/          → NEW APPLICATION
    ├── /          → React static files
    ├── /api/      → Go backend (port 8080)
    ├── /images/   → Go uploads
    └── /static/   → React assets
```
- [ ] Architecture understood
- [ ] Clear on which parts are changing

### 2. Backup Plan
- [ ] Know how to rollback Nginx config
- [ ] Have PM2 rollback plan
- [ ] Existing site will remain untouched

### 3. Downtime Expectations
- [ ] Understand there will be ZERO downtime for existing site
- [ ] New site deployment will take ~15 minutes
- [ ] Nginx reload (not restart) ensures no connection drops

## ✅ Knowledge Check

Answer these questions before proceeding:

1. **What port will the Go backend run on?**
   - [ ] 8080

2. **What port will the existing Node.js app remain on?**
   - [ ] 2299

3. **Where will users access the new application?**
   - [ ] https://nutriediet.com/new

4. **Will the new app have its own database?**
   - [ ] Yes (nutriediet_new_db)

5. **What happens if something goes wrong?**
   - [ ] I can rollback Nginx and remove PM2 app
   - [ ] Existing site remains unaffected

6. **Do I need to restart the existing PM2 app?**
   - [ ] No, it remains untouched

7. **What tool manages the Go backend process?**
   - [ ] PM2 (app name: nutriediet-go-api)

## 🚦 Ready to Deploy?

### All Green? Proceed!
If ALL items above are checked, you're ready to deploy:

```bash
# Copy files to server
scp -r deployment/ sk@YOUR_DROPLET_IP:/home/sk/nutriediet-deployment/

# SSH and deploy
ssh sk@YOUR_DROPLET_IP
cd /home/sk/nutriediet-deployment
./deploy.sh
```

### Any Red? Stop and Fix!
If ANY item is unchecked:
1. Complete that item
2. Test thoroughly
3. Return to this checklist
4. Only proceed when everything is ✅

## 📞 Need Help?

- **Backend CORS issues**: See `cors-update.md`
- **Frontend config issues**: See `frontend-constants-update.md`
- **Deployment questions**: See `DEPLOYMENT_GUIDE.md`
- **Quick reference**: See `QUICK_START.md`

## 📝 Post-Deployment Checklist

After deployment completes, verify:

- [ ] Run `test-deployment.sh`
- [ ] Access https://nutriediet.com (existing site works)
- [ ] Access https://nutriediet.com/new (new site works)
- [ ] Test login/logout on new site
- [ ] Test API calls on new site
- [ ] Check PM2: `pm2 list` shows both apps
- [ ] Check logs: `pm2 logs nutriediet-go-api`
- [ ] No errors in Nginx logs: `sudo tail -f /var/log/nginx/error.log`

---

**Remember:** The deployment script will NOT touch your existing application. Everything is isolated to `/home/sk/mys/nutriediet-new/` and a new database.

**Safety First:** If anything seems wrong during deployment, you can always:
1. Stop the script (Ctrl+C)
2. Remove the new PM2 app: `pm2 delete nutriediet-go-api`
3. Restore Nginx config from backup
4. Your existing site remains unaffected


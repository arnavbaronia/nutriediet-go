# Architecture Diagram - Deployment at www.nutriediet.com/new

## 🏗️ System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    www.nutriediet.com                        │
│                  (Digital Ocean Droplet)                     │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │                    NGINX (Port 443)                 │    │
│  │                     SSL/TLS Enabled                 │    │
│  └───┬──────────────────────────────────────────┬─────┘    │
│      │                                           │          │
│      │                                           │          │
│  ┌───▼───────────────────┐          ┌───────────▼───────┐  │
│  │  EXISTING APP         │          │   NEW APP (/new)   │  │
│  │  (UNCHANGED)          │          │   (DEPLOYING)      │  │
│  │                       │          │                    │  │
│  │  Port: 2299          │          │  ┌──────────────┐  │  │
│  │  PM2: "app"          │          │  │ React App    │  │  │
│  │  Node.js v14         │          │  │ (Static)     │  │  │
│  │                       │          │  └──────────────┘  │  │
│  │  /                   │          │                    │  │
│  │  /libs/              │          │  ┌──────────────┐  │  │
│  │  /uploads/           │          │  │ Go API       │  │  │
│  │                       │          │  │ Port: 8080   │  │  │
│  │  ┌─────────────┐    │          │  │ PM2: "api"   │  │  │
│  │  │   MySQL     │    │          │  └──────────────┘  │  │
│  │  │ (existing)  │    │          │                    │  │
│  │  └─────────────┘    │          │  ┌──────────────┐  │  │
│  │                       │          │  │   MySQL      │  │  │
│  └───────────────────────┘          │  │   (new db)   │  │  │
│                                      │  └──────────────┘  │  │
│                                      └────────────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Directory Structure

```
/home/sk/mys/
├── nutribackend/                    [EXISTING - UNCHANGED]
│   ├── app.js
│   ├── libs/
│   ├── uploads/
│   └── ... (existing files)
│
└── nutriediet-new/                  [NEW - DEPLOYING]
    ├── backend/
    │   ├── nutriediet-go           (Go binary)
    │   ├── .env                    (configuration)
    │   ├── images/                 (uploads)
    │   └── ... (source code)
    │
    ├── frontend/
    │   └── build/
    │       ├── index.html
    │       ├── static/
    │       │   ├── css/
    │       │   ├── js/
    │       │   └── media/
    │       ├── favicon.ico
    │       └── ... (build files)
    │
    ├── logs/
    │   ├── go-api-error.log
    │   └── go-api-out.log
    │
    └── ecosystem.config.js
```

## 🌐 URL Routing

```
User Browser → https://nutriediet.com
                     │
                     ▼
        ┌────────────────────────────┐
        │    NGINX (Port 443/SSL)    │
        └─┬──────────────────────────┘
          │
          ├─ / ──────────────────────────► localhost:2299 (Existing Node.js)
          │
          ├─ /libs/* ────────────────────► /home/sk/mys/nutribackend/libs/
          │
          ├─ /uploads/* ─────────────────► /home/sk/mys/nutribackend/uploads/
          │
          ├─ /new/ ──────────────────────► Static files: .../frontend/build/
          │                                (index.html)
          │
          ├─ /new/static/* ──────────────► Static files: .../frontend/build/static/
          │                                (JS, CSS, images)
          │
          ├─ /new/api/* ─────────────────► localhost:8080 (New Go API)
          │                                Proxy: /new/api/health → http://localhost:8080/health
          │
          └─ /new/images/* ──────────────► /home/sk/mys/nutriediet-new/backend/images/
```

## 🔄 Request Flow Examples

### Example 1: User visits homepage (new app)
```
1. User types: https://nutriediet.com/new
2. Browser sends: GET /new/
3. Nginx receives: /new/
4. Nginx serves: /home/sk/mys/nutriediet-new/frontend/build/index.html
5. Browser loads: HTML file
6. Browser requests: /new/static/js/main.abc123.js
7. Nginx serves: /home/sk/mys/nutriediet-new/frontend/build/static/js/main.abc123.js
```

### Example 2: API call from new app
```
1. React app calls: axios.get('/new/api/clients')
2. Browser sends: GET https://nutriediet.com/new/api/clients
3. Nginx receives: /new/api/clients
4. Nginx proxies to: http://localhost:8080/clients
5. Go API receives: /clients
6. Go API responds: JSON data
7. Nginx forwards: Response to browser
8. React app receives: Data
```

### Example 3: Existing app (unchanged)
```
1. User types: https://nutriediet.com
2. Browser sends: GET /
3. Nginx receives: /
4. Nginx proxies to: http://localhost:2299/
5. Node.js app responds: HTML
6. Browser displays: Existing site
```

## 💾 Database Architecture

```
MySQL Server (localhost:3306)
│
├── Existing Database           [UNCHANGED]
│   ├── Tables: users, posts, etc.
│   └── User: existing_user
│
└── New Database: nutriediet_new_db    [NEW]
    ├── Tables: (from migrations)
    │   ├── userauth
    │   ├── clients
    │   ├── recipes
    │   ├── exercises
    │   └── ... (other tables)
    │
    └── User: nutriediet_new_user
        └── Privileges: Only on nutriediet_new_db
```

## 🔐 Process Management

```
PM2 Process Manager
│
├── app                         [EXISTING - UNCHANGED]
│   ├── Script: app.js
│   ├── Port: 2299
│   ├── Status: online
│   └── Restarts: auto
│
└── nutriediet-go-api          [NEW]
    ├── Script: ./nutriediet-go
    ├── Port: 8080
    ├── Status: online
    ├── Restarts: auto
    ├── Max memory: 500M
    └── Logs:
        ├── out: /home/sk/mys/nutriediet-new/logs/go-api-out.log
        └── err: /home/sk/mys/nutriediet-new/logs/go-api-error.log
```

## 🔌 Port Allocation

```
┌────────────────────────────────────────────────────────┐
│  Port    Service              Status      App          │
├────────────────────────────────────────────────────────┤
│  80      HTTP                 Redirect    (to 443)     │
│  443     HTTPS (Nginx)        Active      Main entry   │
│  2299    Node.js Backend      Active      Existing app │
│  8080    Go API               New         New backend  │
│  3306    MySQL                Active      Both apps    │
└────────────────────────────────────────────────────────┘
```

## 🔒 Security Layers

```
┌──────────────────────────────────────────────┐
│         User's Browser (HTTPS)               │
└───────────────┬──────────────────────────────┘
                │
                ▼
┌──────────────────────────────────────────────┐
│  Layer 1: SSL/TLS (Let's Encrypt)           │
│  - Certificate: /etc/letsencrypt/...         │
│  - Protocols: TLSv1.2, TLSv1.3              │
└───────────────┬──────────────────────────────┘
                │
                ▼
┌──────────────────────────────────────────────┐
│  Layer 2: Nginx Security Headers            │
│  - X-Frame-Options                           │
│  - X-Content-Type-Options                    │
│  - X-XSS-Protection                          │
└───────────────┬──────────────────────────────┘
                │
                ▼
┌──────────────────────────────────────────────┐
│  Layer 3: Application CORS                   │
│  - Go backend validates origins              │
│  - Only allows nutriediet.com                │
└───────────────┬──────────────────────────────┘
                │
                ▼
┌──────────────────────────────────────────────┐
│  Layer 4: JWT Authentication                 │
│  - Token validation                          │
│  - User authorization                        │
└───────────────┬──────────────────────────────┘
                │
                ▼
┌──────────────────────────────────────────────┐
│  Layer 5: Database Access Control            │
│  - Limited user privileges                   │
│  - Separate databases                        │
└──────────────────────────────────────────────┘
```

## 📊 Data Flow

### Authentication Flow
```
User Login Request
    │
    ├─→ Browser: POST /new/api/login
    │
    ├─→ Nginx: Receives /new/api/login
    │
    ├─→ Nginx: Proxies to localhost:8080/login
    │
    ├─→ Go API: Validates credentials
    │
    ├─→ MySQL: SELECT * FROM userauth WHERE email = ?
    │
    ├─→ Go API: Generates JWT token
    │
    ├─→ Nginx: Forwards response
    │
    └─→ Browser: Stores token in localStorage
```

### Data Fetch Flow
```
Fetch User Data
    │
    ├─→ Browser: GET /new/api/clients/123/profile
    │           Header: Authorization: Bearer {token}
    │
    ├─→ Nginx: Receives /new/api/clients/123/profile
    │
    ├─→ Nginx: Proxies to localhost:8080/clients/123/profile
    │
    ├─→ Go API: Validates JWT token
    │
    ├─→ Go API: Checks authorization
    │
    ├─→ MySQL: SELECT * FROM clients WHERE id = 123
    │
    ├─→ Go API: Returns JSON
    │
    ├─→ Nginx: Forwards response
    │
    └─→ Browser: Updates UI
```

## 🔄 Deployment Flow

```
Local Development
    │
    ├─ 1. Make code changes
    │   ├─ Update constants.js
    │   ├─ Update package.json
    │   └─ Update main.go
    │
    ├─ 2. Test locally
    │   ├─ npm run build
    │   └─ go run main.go
    │
    ├─ 3. Commit & push
    │   └─ git push
    │
    └─ 4. Deploy to server
        │
        ▼
Server Deployment (via deploy.sh)
    │
    ├─ 5. Install dependencies
    │   ├─ Go 1.21.5
    │   └─ Node.js v20
    │
    ├─ 6. Setup directories
    │   └─ /home/sk/mys/nutriediet-new/
    │
    ├─ 7. Setup database
    │   ├─ CREATE DATABASE nutriediet_new_db
    │   └─ CREATE USER nutriediet_new_user
    │
    ├─ 8. Clone repositories
    │   ├─ git clone backend
    │   └─ git clone frontend
    │
    ├─ 9. Build backend
    │   ├─ go mod download
    │   ├─ go build
    │   └─ Create .env
    │
    ├─ 10. Build frontend
    │   ├─ npm ci
    │   └─ npm run build
    │
    ├─ 11. Run migrations
    │   └─ go run migrate/migrate.go
    │
    ├─ 12. Setup PM2
    │   ├─ pm2 start ecosystem.config.js
    │   └─ pm2 save
    │
    ├─ 13. Update Nginx
    │   ├─ Backup current config
    │   ├─ Add new location blocks
    │   ├─ Test: nginx -t
    │   └─ Reload: systemctl reload nginx
    │
    └─ 14. Verify
        ├─ Test Go API: curl localhost:8080
        ├─ Test existing: curl localhost:2299
        └─ Browser test: https://nutriediet.com/new
```

## 🎯 Component Interaction

```
┌─────────────────────────────────────────────────────────┐
│                    User's Browser                        │
│  ┌──────────────┐      ┌──────────────┐                 │
│  │  React App   │      │ Local Storage│                 │
│  │  (Frontend)  │◄────►│ - JWT Token  │                 │
│  │  at /new     │      │ - User Data  │                 │
│  └──────┬───────┘      └──────────────┘                 │
└─────────┼────────────────────────────────────────────────┘
          │ HTTPS
          │ /new/api/*
          ▼
┌─────────────────────────────────────────────────────────┐
│               Digital Ocean Droplet                      │
│                                                          │
│  ┌───────────────────────────────────────────────┐     │
│  │              Nginx (Reverse Proxy)            │     │
│  │  - SSL Termination                            │     │
│  │  - Static File Serving                        │     │
│  │  - API Proxying                               │     │
│  └──────────┬────────────────────────────────────┘     │
│             │                                           │
│             ├─────► /new/        → Serve React build   │
│             ├─────► /new/static/ → Serve CSS/JS        │
│             └─────► /new/api/*   → Proxy to Go API     │
│                                        │                │
│                                        ▼                │
│  ┌──────────────────────────────────────────────┐     │
│  │           Go API (Port 8080)                 │     │
│  │  ┌────────────────────────────────────┐     │     │
│  │  │ Gin Framework                      │     │     │
│  │  │  - Routing                         │     │     │
│  │  │  - Middleware (Auth, CORS, etc)    │     │     │
│  │  └────────────────────────────────────┘     │     │
│  │  ┌────────────────────────────────────┐     │     │
│  │  │ Controllers                        │     │     │
│  │  │  - Admin Controller                │     │     │
│  │  │  - Client Controller               │     │     │
│  │  │  - Auth Controller                 │     │     │
│  │  └────────────────────────────────────┘     │     │
│  │  ┌────────────────────────────────────┐     │     │
│  │  │ Models (GORM)                      │     │     │
│  │  └────────────────────────────────────┘     │     │
│  └──────────────────┬───────────────────────────┘     │
│                     │                                   │
│                     ▼                                   │
│  ┌──────────────────────────────────────────────┐     │
│  │       MySQL Database Server                  │     │
│  │  ┌────────────────────────────────────┐     │     │
│  │  │ nutriediet_new_db                  │     │     │
│  │  │  - userauth                        │     │     │
│  │  │  - clients                         │     │     │
│  │  │  - recipes                         │     │     │
│  │  │  - exercises                       │     │     │
│  │  │  - diet_plans                      │     │     │
│  │  └────────────────────────────────────┘     │     │
│  └──────────────────────────────────────────────┘     │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

## 📦 Technology Stack

```
Frontend Layer
├── React 18.3.1
├── React Router 6
├── Axios (HTTP client)
├── Material-UI
├── Bootstrap
└── Chart.js

Backend Layer
├── Go 1.20+
├── Gin Framework
├── GORM (ORM)
├── JWT (Authentication)
└── bcrypt (Password hashing)

Infrastructure Layer
├── Nginx (Web server / Reverse proxy)
├── PM2 (Process manager)
├── MySQL 8+ (Database)
├── Let's Encrypt (SSL/TLS)
└── Ubuntu Linux (OS)

Deployment Layer
├── Git (Version control)
├── GitHub (Repository hosting)
└── Digital Ocean (Cloud hosting)
```

## 🚀 Scalability Considerations

### Current Setup (Single Server)
```
All components on one droplet:
- Nginx
- Go API (1 instance)
- MySQL
- Static files
```

### Future Scaling Options
```
1. Horizontal Scaling
   ├── Multiple Go API instances
   │   └── PM2 cluster mode
   │
   ├── Load balancer
   │   └── Nginx upstream
   │
   └── Separate database server
       └── MySQL on separate droplet

2. Vertical Scaling
   ├── Increase droplet size
   ├── More RAM for Go API
   └── Faster storage for database

3. CDN Integration
   └── Serve static files from CDN
       └── CloudFlare, AWS CloudFront
```

## 🔍 Monitoring Points

```
Application Level
├── PM2 logs: pm2 logs nutriediet-go-api
├── Nginx logs: /var/log/nginx/error.log
└── MySQL logs: /var/log/mysql/error.log

System Level
├── CPU usage: htop
├── Memory: free -h
├── Disk: df -h
└── Network: netstat -tlnp

Application Metrics
├── API response times
├── Error rates
├── Database query performance
└── User sessions
```

---

**This architecture ensures:**
- ✅ Zero downtime for existing application
- ✅ Isolated new application
- ✅ Scalable design
- ✅ Secure communication
- ✅ Easy maintenance and updates


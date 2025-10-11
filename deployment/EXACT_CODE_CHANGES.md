# Exact Code Changes Required for Deployment

This document shows the EXACT changes needed in your codebase for /new subpath deployment.

## 1. Frontend: `src/utils/constants.js`

### Change 1: Update API_BASE_URL
**Find this line (line 7):**
```javascript
export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';
```

**Replace with:**
```javascript
export const API_BASE_URL = process.env.REACT_APP_API_URL || '/new/api';
```

### Change 2: Update ROUTES with BASE_PATH
**Find this section (lines 68-76):**
```javascript
// Routes
export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  SIGNUP: '/signup',
  ADMIN_LOGIN: '/admin/login',
  ADMIN_DASHBOARD: '/admin/dashboard',
  CLIENT_DASHBOARD: (clientId) => `/clients/${clientId}/diet`,
  ACCOUNT_ACTIVATION: '/account-activation',
};
```

**Replace with:**
```javascript
// Routes with support for subpath deployment
const BASE_PATH = process.env.PUBLIC_URL || '';

export const ROUTES = {
  HOME: `${BASE_PATH}/`,
  LOGIN: `${BASE_PATH}/login`,
  SIGNUP: `${BASE_PATH}/signup`,
  ADMIN_LOGIN: `${BASE_PATH}/admin/login`,
  ADMIN_DASHBOARD: `${BASE_PATH}/admin/dashboard`,
  CLIENT_DASHBOARD: (clientId) => `${BASE_PATH}/clients/${clientId}/diet`,
  ACCOUNT_ACTIVATION: `${BASE_PATH}/account-activation`,
};
```

**Complete updated file:**
```javascript
/**
 * Application Constants
 * Centralized configuration and constant values
 */

// API Configuration
export const API_BASE_URL = process.env.REACT_APP_API_URL || '/new/api';

// Environment
export const IS_PRODUCTION = process.env.REACT_APP_ENV === 'production';
export const IS_DEVELOPMENT = process.env.REACT_APP_ENV === 'development';
export const IS_DEBUG = process.env.REACT_APP_DEBUG === 'true';

// API Endpoints
export const API_ENDPOINTS = {
  // Auth
  LOGIN: '/login',
  SIGNUP: '/signup',
  FORGOT_PASSWORD: '/auth/forgot-password',
  RESET_PASSWORD: '/auth/reset-password',
  REFRESH_TOKEN: '/auth/refresh',
  
  // Client
  CLIENT_PROFILE: (clientId) => `/clients/${clientId}/my_profile`,
  CLIENT_DIET: (clientId) => `/clients/${clientId}/diet`,
  CLIENT_EXERCISE: (clientId) => `/clients/${clientId}/exercise`,
  CLIENT_RECIPE: (clientId) => `/clients/${clientId}/recipe`,
  CLIENT_MOTIVATION: (clientId) => `/clients/${clientId}/motivation`,
  CLIENT_WEIGHT: (clientId) => `/clients/${clientId}/weight_update`,
  
  // Admin
  ADMIN_CLIENTS: '/admin/clients',
  ADMIN_CLIENT_DETAILS: (clientId) => `/admin/clients/${clientId}`,
  ADMIN_DIET_TEMPLATES: '/admin/diet_templates',
  ADMIN_RECIPES: '/admin/recipes',
  ADMIN_EXERCISES: '/admin/exercises',
  ADMIN_MOTIVATIONS: '/admin/motivations',
};

// User Types
export const USER_TYPES = {
  CLIENT: 'CLIENT',
  ADMIN: 'ADMIN',
};

// Diet Types
export const DIET_TYPES = {
  REGULAR: '1',
  DETOX: '2',
  DETOX_WATER: '3',
};

// Storage Keys
export const STORAGE_KEYS = {
  TOKEN: 'token',
  REFRESH_TOKEN: 'refreshToken',
  USER_TYPE: 'user_type',
  CLIENT_ID: 'client_id',
  EMAIL: 'email',
  IS_ACTIVE: 'is_active',
  USER: 'user',
  USER_ID: 'userId',
  FIRST_NAME: 'firstName',
  LAST_NAME: 'lastName',
};

// Routes with support for subpath deployment
const BASE_PATH = process.env.PUBLIC_URL || '';

export const ROUTES = {
  HOME: `${BASE_PATH}/`,
  LOGIN: `${BASE_PATH}/login`,
  SIGNUP: `${BASE_PATH}/signup`,
  ADMIN_LOGIN: `${BASE_PATH}/admin/login`,
  ADMIN_DASHBOARD: `${BASE_PATH}/admin/dashboard`,
  CLIENT_DASHBOARD: (clientId) => `${BASE_PATH}/clients/${clientId}/diet`,
  ACCOUNT_ACTIVATION: `${BASE_PATH}/account-activation`,
};

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection.',
  UNAUTHORIZED: 'Your session has expired. Please login again.',
  SERVER_ERROR: 'Server error. Please try again later.',
  NOT_FOUND: 'Resource not found.',
};

// Success Messages
export const SUCCESS_MESSAGES = {
  PROFILE_UPDATED: 'Profile updated successfully!',
  PASSWORD_RESET: 'Password reset successfully!',
  DATA_SAVED: 'Data saved successfully!',
};

// Validation
export const VALIDATION = {
  MIN_PASSWORD_LENGTH: 12,
  MAX_PASSWORD_LENGTH: 128,
  EMAIL_REGEX: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  OTP_LENGTH: 6,
  PASSWORD_REQUIREMENTS: {
    UPPERCASE: /[A-Z]/,
    LOWERCASE: /[a-z]/,
    NUMBER: /[0-9]/,
    SPECIAL_CHAR: /[!@#$%^&*(),.?":{}|<>_\-+=\[\]\\\/;'~`]/,
  },
  WEAK_PASSWORDS: [
    'password123!',
    'admin123456!',
    'welcome12345!',
    'qwerty123456!',
    '123456789abc!',
    'letmein12345!',
  ],
};
```

---

## 2. Frontend: `package.json`

### Change: Add homepage field

**Find (near the top):**
```json
{
  "name": "frontend",
  "version": "0.1.0",
  "private": true,
```

**Add after "version":**
```json
{
  "name": "frontend",
  "version": "0.1.0",
  "homepage": "/new",
  "private": true,
```

---

## 3. Frontend: Create `.env.production`

**Create new file:** `frontend/.env.production`

**Contents:**
```bash
REACT_APP_API_URL=/new/api
PUBLIC_URL=/new
NODE_ENV=production
REACT_APP_ENV=production
```

---

## 4. Frontend: `src/App.js`

### Change: Add basename to BrowserRouter

**Find the BrowserRouter import:**
```javascript
import { BrowserRouter } from 'react-router-dom';
```

**Find where BrowserRouter is used (likely in return statement):**
```javascript
<BrowserRouter>
  {/* routes */}
</BrowserRouter>
```

**Update to:**
```javascript
<BrowserRouter basename={process.env.PUBLIC_URL || ''}>
  {/* routes */}
</BrowserRouter>
```

**Complete example:**
```javascript
import { BrowserRouter } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter basename={process.env.PUBLIC_URL || ''}>
      <Routes>
        {/* Your existing routes */}
      </Routes>
    </BrowserRouter>
  );
}

export default App;
```

---

## 5. Backend: `main.go`

### Change: Update CORS configuration

**Find this section (lines 48-56):**
```go
config := cors.Config{
    AllowOrigins:     []string{"https://nutriediet.netlify.app", "http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Client-Email", "Request-Client-ID"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
router.Use(cors.New(config))
```

**Replace with:**
```go
// Get allowed origins from environment or use defaults
allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
var origins []string

if allowedOrigins != "" {
    origins = strings.Split(allowedOrigins, ",")
} else {
    // Default origins for development and production
    origins = []string{
        "https://nutriediet.com",
        "https://www.nutriediet.com",
        "https://nutriediet.netlify.app",
        "http://localhost:3000",
    }
}

config := cors.Config{
    AllowOrigins:     origins,
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
    AllowHeaders:     []string{
        "Origin",
        "Content-Type",
        "Authorization",
        "Client-Email",
        "Request-Client-ID",
        "X-Requested-With",
    },
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
router.Use(cors.New(config))
```

**Make sure imports include strings:**
```go
import (
    "fmt"
    "log"
    "os"
    "strings"  // Add this if not present
    "time"
    // ... other imports
)
```

---

## 6. Backend: Create `.env.production.template`

This is already created in `deployment/.env.production.template`. You'll copy and customize it during deployment.

---

## Quick Command Reference

### Test Changes Locally

**Backend:**
```bash
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
go run main.go
# Should start on port 8080
```

**Frontend Development:**
```bash
cd /Users/ishitagupta/Documents/Personal/frontend
npm start
# Should work on http://localhost:3000
```

**Frontend Production Build:**
```bash
cd /Users/ishitagupta/Documents/Personal/frontend
GENERATE_SOURCEMAP=false npm run build
npx serve -s build -l 3001
# Test on http://localhost:3001
```

---

## Verification Checklist

After making changes:

### Frontend
- [ ] `constants.js` updated with BASE_PATH
- [ ] `package.json` has `"homepage": "/new"`
- [ ] `.env.production` created
- [ ] `App.js` BrowserRouter has basename
- [ ] `npm start` works for development
- [ ] `npm run build` completes without errors
- [ ] Production build serves correctly

### Backend
- [ ] `main.go` CORS updated
- [ ] Imports include `strings`
- [ ] `go run main.go` starts successfully
- [ ] API endpoints respond

---

## Common Mistakes to Avoid

### ❌ Wrong: Hardcoding /new in constants
```javascript
export const API_BASE_URL = '/new/api';  // Won't work in development
```

### ✅ Right: Use environment variable with fallback
```javascript
export const API_BASE_URL = process.env.REACT_APP_API_URL || '/new/api';
```

### ❌ Wrong: Forgetting basename in Router
```javascript
<BrowserRouter>
```

### ✅ Right: Include basename
```javascript
<BrowserRouter basename={process.env.PUBLIC_URL || ''}>
```

### ❌ Wrong: Not testing production build
```bash
npm run build  # But never testing the build
```

### ✅ Right: Test production build locally
```bash
npm run build
npx serve -s build
# Visit and test
```

---

## Files Summary

**Files to modify:**
1. `frontend/src/utils/constants.js` - 2 changes
2. `frontend/package.json` - 1 addition
3. `frontend/src/App.js` - 1 change
4. `backend/main.go` - 1 change

**Files to create:**
1. `frontend/.env.production` - New file

**Total changes:** 5 files, ~15 minutes of work

---

## Ready to Deploy?

Once all changes are made and tested:

1. **Commit and push:**
```bash
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
git add .
git commit -m "Configure for /new subpath deployment"
git push

cd /Users/ishitagupta/Documents/Personal/frontend
git add .
git commit -m "Configure for /new subpath deployment"
git push
```

2. **Follow QUICK_START.md** for deployment steps

3. **Run PRE_DEPLOYMENT_CHECKLIST.md** to verify everything


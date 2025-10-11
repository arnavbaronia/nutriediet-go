# Frontend Constants Configuration Update

## File: `src/utils/constants.js`

Your frontend uses a centralized constants file for configuration. This needs to be updated for the /new subpath deployment.

## Current Configuration

Your constants file currently has hardcoded URLs. For deployment at `/new`, you need to use environment variables and relative paths.

## Updated Configuration

Replace the API_BASE_URL and ROUTES configuration in `src/utils/constants.js`:

```javascript
// API Configuration
export const API_BASE_URL = process.env.REACT_APP_API_URL || '/new/api';

// Storage Keys (no changes needed)
export const STORAGE_KEYS = {
  TOKEN: 'nutriediet_token',
  REFRESH_TOKEN: 'nutriediet_refresh_token',
  USER: 'nutriediet_user',
  USER_TYPE: 'nutriediet_user_type',
  // ... rest of your storage keys
};

// Routes - Add basePath for /new deployment
const BASE_PATH = process.env.PUBLIC_URL || '/new';

export const ROUTES = {
  HOME: `${BASE_PATH}/`,
  LOGIN: `${BASE_PATH}/login`,
  REGISTER: `${BASE_PATH}/register`,
  DASHBOARD: `${BASE_PATH}/dashboard`,
  ADMIN_DASHBOARD: `${BASE_PATH}/admin/dashboard`,
  CLIENT_DASHBOARD: `${BASE_PATH}/client/dashboard`,
  PROFILE: `${BASE_PATH}/profile`,
  // ... update all your routes with BASE_PATH prefix
};

// Or if you prefer simpler approach for production:
export const ROUTES = {
  HOME: '/new/',
  LOGIN: '/new/login',
  REGISTER: '/new/register',
  DASHBOARD: '/new/dashboard',
  // ... etc
};

// Status codes (no changes needed)
export const STATUS_CODES = {
  SUCCESS: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  SERVER_ERROR: 500,
};

// User types (no changes needed)
export const USER_TYPES = {
  ADMIN: 'admin',
  CLIENT: 'client',
  DIETICIAN: 'dietician',
};

// ... rest of your constants
```

## Environment Variables

### Development: `.env.development`
```bash
REACT_APP_API_URL=http://localhost:8080
PUBLIC_URL=/
```

### Production: `.env.production`
```bash
REACT_APP_API_URL=/new/api
PUBLIC_URL=/new
NODE_ENV=production
```

## Why Relative Paths?

Using relative paths like `/new/api` instead of absolute URLs like `https://nutriediet.com/new/api` has several benefits:

1. **Works across environments** - localhost, staging, production
2. **Automatic HTTPS** - inherits protocol from current page
3. **No CORS issues** - same-origin requests
4. **Easier maintenance** - no hardcoded domains

## Complete Example

Here's a complete example of an updated constants file:

```javascript
// ========================================
// API Configuration
// ========================================
export const API_BASE_URL = process.env.REACT_APP_API_URL || '/new/api';

// ========================================
// Storage Keys
// ========================================
export const STORAGE_KEYS = {
  TOKEN: 'nutriediet_token',
  REFRESH_TOKEN: 'nutriediet_refresh_token',
  USER: 'nutriediet_user',
  USER_TYPE: 'nutriediet_user_type',
  THEME: 'nutriediet_theme',
  LANGUAGE: 'nutriediet_language',
};

// ========================================
// Application Routes
// ========================================
const BASE_PATH = process.env.PUBLIC_URL || '';

export const ROUTES = {
  // Auth routes
  HOME: `${BASE_PATH}/`,
  LOGIN: `${BASE_PATH}/login`,
  REGISTER: `${BASE_PATH}/register`,
  FORGOT_PASSWORD: `${BASE_PATH}/forgot-password`,
  RESET_PASSWORD: `${BASE_PATH}/reset-password`,
  
  // Dashboard routes
  DASHBOARD: `${BASE_PATH}/dashboard`,
  ADMIN_DASHBOARD: `${BASE_PATH}/admin/dashboard`,
  CLIENT_DASHBOARD: `${BASE_PATH}/client/dashboard`,
  DIETICIAN_DASHBOARD: `${BASE_PATH}/dietician/dashboard`,
  
  // User routes
  PROFILE: `${BASE_PATH}/profile`,
  SETTINGS: `${BASE_PATH}/settings`,
  
  // Feature routes
  MEALS: `${BASE_PATH}/meals`,
  RECIPES: `${BASE_PATH}/recipes`,
  DIET_PLANS: `${BASE_PATH}/diet-plans`,
  EXERCISES: `${BASE_PATH}/exercises`,
  CLIENTS: `${BASE_PATH}/clients`,
};

// ========================================
// API Endpoints
// ========================================
export const API_ENDPOINTS = {
  // Auth
  LOGIN: '/auth/login',
  REGISTER: '/auth/register',
  LOGOUT: '/auth/logout',
  REFRESH_TOKEN: '/auth/refresh',
  FORGOT_PASSWORD: '/auth/forgot-password',
  RESET_PASSWORD: '/auth/reset-password',
  
  // User
  GET_PROFILE: '/user/profile',
  UPDATE_PROFILE: '/user/profile',
  CHANGE_PASSWORD: '/user/change-password',
  
  // Meals
  GET_MEALS: '/meals',
  GET_MEAL: (id) => `/meals/${id}`,
  CREATE_MEAL: '/meals',
  UPDATE_MEAL: (id) => `/meals/${id}`,
  DELETE_MEAL: (id) => `/meals/${id}`,
  
  // Add your other endpoints...
};

// ========================================
// Status Codes
// ========================================
export const STATUS_CODES = {
  SUCCESS: 200,
  CREATED: 201,
  ACCEPTED: 202,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  SERVER_ERROR: 500,
  SERVICE_UNAVAILABLE: 503,
};

// ========================================
// User Types
// ========================================
export const USER_TYPES = {
  ADMIN: 'admin',
  CLIENT: 'client',
  DIETICIAN: 'dietician',
};

// ========================================
// Validation Rules
// ========================================
export const VALIDATION = {
  PASSWORD_MIN_LENGTH: 8,
  USERNAME_MIN_LENGTH: 3,
  USERNAME_MAX_LENGTH: 50,
  EMAIL_REGEX: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PHONE_REGEX: /^[0-9]{10}$/,
};

// ========================================
// Application Settings
// ========================================
export const APP_CONFIG = {
  APP_NAME: 'Nutriediet',
  VERSION: '2.0.0',
  API_TIMEOUT: 30000, // 30 seconds
  MAX_FILE_SIZE: 5 * 1024 * 1024, // 5MB
  SUPPORTED_IMAGE_FORMATS: ['image/jpeg', 'image/png', 'image/jpg', 'image/webp'],
};

// ========================================
// Date/Time Formats
// ========================================
export const DATE_FORMATS = {
  DISPLAY: 'MMM DD, YYYY',
  DISPLAY_WITH_TIME: 'MMM DD, YYYY hh:mm A',
  API: 'YYYY-MM-DD',
  API_WITH_TIME: 'YYYY-MM-DD HH:mm:ss',
};
```

## Usage in Components

With these updates, your components should work without changes:

```javascript
import api from '../api/axiosInstance';
import { API_ENDPOINTS, ROUTES, STORAGE_KEYS } from '../utils/constants';

// API calls work automatically
const fetchMeals = async () => {
  const response = await api.get(API_ENDPOINTS.GET_MEALS);
  return response.data;
};

// Navigation works
const handleLogin = () => {
  navigate(ROUTES.LOGIN);
};

// Storage works
const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
```

## React Router Configuration

Update your App.js to use basename:

```javascript
import { BrowserRouter } from 'react-router-dom';

const BASE_PATH = process.env.PUBLIC_URL || '';

function App() {
  return (
    <BrowserRouter basename={BASE_PATH}>
      {/* Your routes */}
    </BrowserRouter>
  );
}
```

## Testing the Changes

### Test 1: Development Build
```bash
cd frontend
npm start
# Should work on http://localhost:3000 without /new prefix
```

### Test 2: Production Build
```bash
cd frontend
REACT_APP_API_URL=/new/api PUBLIC_URL=/new npm run build
npx serve -s build -l 3001
# Navigate to http://localhost:3001
# All assets should load correctly
```

### Test 3: Verify Environment Variables
```javascript
// In any component, console log to verify:
console.log('API URL:', process.env.REACT_APP_API_URL);
console.log('Public URL:', process.env.PUBLIC_URL);
```

## Common Mistakes to Avoid

### ❌ Don't use absolute URLs
```javascript
// Bad
const API_BASE_URL = 'https://nutriediet.com/new/api';
```

### ✅ Use environment variables or relative paths
```javascript
// Good
const API_BASE_URL = process.env.REACT_APP_API_URL || '/new/api';
```

### ❌ Don't hardcode route paths in JSX
```javascript
// Bad
<Link to="/dashboard">Dashboard</Link>
```

### ✅ Use constants
```javascript
// Good
import { ROUTES } from '../utils/constants';
<Link to={ROUTES.DASHBOARD}>Dashboard</Link>
```

### ❌ Don't forget basename in Router
```javascript
// Bad
<BrowserRouter>
```

### ✅ Include basename
```javascript
// Good
<BrowserRouter basename={process.env.PUBLIC_URL}>
```

## Rollout Checklist

- [ ] Update `src/utils/constants.js` with environment variables
- [ ] Create `.env.production` file
- [ ] Update `package.json` with `"homepage": "/new"`
- [ ] Update `App.js` BrowserRouter with basename
- [ ] Test build locally
- [ ] Deploy to server
- [ ] Verify all routes work
- [ ] Verify API calls work
- [ ] Verify image loading works

## After Deployment

If you need to check what environment variables are being used:

```bash
# On server, check the build
cd /home/sk/mys/nutriediet-new/frontend/build
grep -r "REACT_APP" static/js/*.js
# Should show your production values
```


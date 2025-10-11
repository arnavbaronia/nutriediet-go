# Frontend API Configuration Update

## Current Setup
Your frontend uses `axiosInstance.js` to configure API calls.

## Required Changes for /new Subpath Deployment

### Option 1: Use Environment Variable (Recommended)

Update `src/api/axiosInstance.js`:

```javascript
import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_URL || '/new/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor for auth token
axiosInstance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized - clear token and redirect to login
      localStorage.removeItem('token');
      window.location.href = '/new/login';
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;
```

### Option 2: Use Relative Path

If you prefer simpler configuration:

```javascript
const axiosInstance = axios.create({
  baseURL: '/new/api',  // Fixed path for production
  // ... rest of config
});
```

## Environment File

Make sure `.env.production` exists in frontend root:

```
REACT_APP_API_URL=/new/api
PUBLIC_URL=/new
NODE_ENV=production
```

## Development vs Production

For local development, use `.env.development`:

```
REACT_APP_API_URL=http://localhost:8080
PUBLIC_URL=/
```

## Router Configuration

If using React Router, update your BrowserRouter in `src/App.js`:

```javascript
import { BrowserRouter } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter basename="/new">
      {/* Your routes */}
    </BrowserRouter>
  );
}
```

## Image/Asset Paths

For images or assets served from Go backend:

```javascript
// Good - relative to domain
<img src="/new/images/profile.jpg" />

// Good - from environment
<img src={`${process.env.PUBLIC_URL}/images/profile.jpg`} />

// Bad - absolute path will break
<img src="https://nutriediet.com/images/profile.jpg" />
```

## Testing Changes

After making these changes:

1. **Build locally to test:**
```bash
cd frontend
npm run build
npx serve -s build -l 3001
# Visit http://localhost:3001 - should work without /new locally
```

2. **Test production build:**
- Ensure homepage is set in package.json
- Build and deploy to server
- All asset paths should automatically use /new prefix

## Common Issues

### Issue: Assets 404
**Cause:** `homepage` not set in package.json  
**Fix:** Add `"homepage": "/new"` to package.json

### Issue: API calls to wrong URL
**Cause:** baseURL not using environment variable  
**Fix:** Use `process.env.REACT_APP_API_URL`

### Issue: Routing breaks on refresh
**Cause:** Nginx not configured for SPA  
**Fix:** Already handled in nginx-config-new.conf with `try_files`

### Issue: CORS errors
**Cause:** Go backend not allowing your domain  
**Fix:** Update CORS in Go main.go (see cors-update.md)


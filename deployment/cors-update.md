# Go Backend CORS Configuration Update

## Required Changes for Production

Update your `main.go` CORS configuration to work with the /new subpath deployment.

## Current Configuration
```go
config := cors.Config{
    AllowOrigins:     []string{"https://nutriediet.netlify.app", "http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Client-Email", "Request-Client-ID"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

## Updated Production Configuration

Replace the CORS config in `main.go` with:

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
        "http://localhost:3001",
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
    ExposeHeaders:    []string{"Content-Length", "Content-Type"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
router.Use(cors.New(config))
```

## Import Statement

Make sure you have the strings import:

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

## Environment Variable

In your `.env` file on the server:

```bash
ALLOWED_ORIGINS=https://nutriediet.com,https://www.nutriediet.com
```

## Why This Works

1. **Same Origin:** When React is served from `/new` and API is at `/new/api`, they're on the same domain (nutriediet.com), so CORS isn't strictly needed for same-origin requests.

2. **Nginx Proxy:** Nginx proxies requests from `/new/api/*` to `localhost:8080/*`, and the browser sees everything as coming from nutriediet.com.

3. **Explicit CORS:** We still configure CORS explicitly for:
   - Development environments (localhost)
   - Any external frontends (like Netlify)
   - Future flexibility

## Testing CORS

### Test 1: Check CORS Headers
```bash
curl -I -X OPTIONS https://nutriediet.com/new/api/health \
  -H "Origin: https://nutriediet.com" \
  -H "Access-Control-Request-Method: GET"
```

Expected response should include:
```
Access-Control-Allow-Origin: https://nutriediet.com
Access-Control-Allow-Credentials: true
```

### Test 2: Test from Browser Console
Open https://nutriediet.com/new and run in console:
```javascript
fetch('/new/api/health', {
  method: 'GET',
  credentials: 'include'
})
.then(r => r.json())
.then(console.log)
.catch(console.error);
```

Should work without CORS errors.

## Common CORS Issues

### Issue: "CORS policy: No 'Access-Control-Allow-Origin' header"
**Fix:** Check that CORS middleware is loaded before routes in main.go

### Issue: "Credentials flag is true, but Access-Control-Allow-Credentials is not"
**Fix:** Ensure `AllowCredentials: true` is set

### Issue: Preflight OPTIONS requests fail
**Fix:** Ensure OPTIONS is in AllowMethods and MaxAge is set

## Additional Security

For production, you can tighten CORS further:

```go
if os.Getenv("GIN_MODE") == "release" {
    // Production: Only allow production domains
    origins = []string{
        "https://nutriediet.com",
        "https://www.nutriediet.com",
    }
} else {
    // Development: Allow localhost
    origins = []string{
        "https://nutriediet.com",
        "https://www.nutriediet.com",
        "http://localhost:3000",
        "http://localhost:3001",
    }
}
```

## After Making Changes

1. **Build and deploy:**
```bash
go build -o nutriediet-go -ldflags="-s -w" .
```

2. **Restart PM2:**
```bash
pm2 restart nutriediet-go-api
```

3. **Test in browser:**
- Open https://nutriediet.com/new
- Open browser console
- Should see no CORS errors

## Path Handling

Since Nginx proxies `/new/api/` to `localhost:8080/`, your Go routes should NOT include `/new/api` prefix:

```go
// ✅ Correct - routes without /new/api prefix
router.GET("/health", healthHandler)
router.POST("/auth/login", loginHandler)
router.GET("/users", getUsersHandler)

// ❌ Wrong - don't add /new/api prefix
router.GET("/new/api/health", healthHandler)  // This will result in /new/api/new/api/health
```

The Nginx proxy_pass strips `/new/api` and forwards the rest to Go.

Example:
- Browser requests: `https://nutriediet.com/new/api/health`
- Nginx receives: `/new/api/health`
- Nginx forwards to Go: `http://localhost:8080/health`
- Go route handles: `/health`


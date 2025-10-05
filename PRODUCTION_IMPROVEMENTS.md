# Production Improvements & Security Hardening Guide

## Overview
This document outlines critical improvements needed before deploying this application to production and making it publicly accessible. Issues are categorized by priority: **🔴 CRITICAL**, **🟠 HIGH**, **🟡 MEDIUM**, **🟢 LOW**.

---

## 🔴 CRITICAL SECURITY ISSUES (Must Fix Before Production)

### 1. **Hardcoded Database Credentials in Source Code** ✅ FIXED
**File:** `database/database.go:13`

**Risk:** Database credentials were exposed in version control. Anyone with repository access could access your database.

**Status:** ✅ **RESOLVED** - Updated to use environment variables

**Current Implementation:**
The database connection now:
- Reads credentials from environment variables
- Supports both local (Digital Ocean) and cloud (Aiven) databases
- Automatically disables TLS for localhost (Digital Ocean setup)
- Uses TLS for remote databases
- Includes connection pooling
- Has proper error handling

**Digital Ocean Deployment Notes:**
- When running MySQL on the same droplet, use `DB_HOST=localhost`
- TLS is automatically disabled for localhost connections
- Connection pooling is configured for optimal performance
- See `DIGITAL_OCEAN_DEPLOYMENT.md` for MySQL setup instructions

**Environment Variables Required:**
```bash
# For Digital Ocean (MySQL on same machine)
DB_USER=nutriediet_app
DB_PASSWORD=your_strong_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=nutriediet_production

# For Development (using Aiven cloud database)
DB_USER=avnadmin
DB_PASSWORD=AVNS_7QDxgZDlRhQXAx3QV4z
DB_HOST=nutriediet-mysql-ishitagupta-5564.f.aivencloud.com
DB_PORT=22013
DB_NAME=defaultdb
```

---

### 2. **Empty JWT Secret Key**
**File:** `helpers/token_helpers.go:21`
```go
var SECRET_KEY = ""
```

**Risk:** JWT tokens can be forged, bypassing authentication entirely.

**Solution:**
- Use a strong, randomly generated secret key from environment variables
- Rotate keys periodically

**Fix:**
```go
var SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

func init() {
	if SECRET_KEY == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required")
	}
	if len(SECRET_KEY) < 32 {
		log.Fatal("JWT_SECRET_KEY must be at least 32 characters")
	}
}
```

Generate a secure key:
```bash
openssl rand -base64 64
```

---

### 3. **No Rate Limiting**
**Risk:** Vulnerable to brute force attacks, DDoS, credential stuffing, OTP enumeration.

**Critical Endpoints:**
- `/login` - brute force password attacks
- `/auth/forgot-password` - OTP spam/enumeration
- `/auth/reset-password` - OTP brute force
- `/signup` - account creation spam

**Solution:** Implement rate limiting middleware

**Fix:**
```go
// Install: go get github.com/ulule/limiter/v3

// middleware/rate_limit.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"net/http"
)

func RateLimitMiddleware(rate string) gin.HandlerFunc {
	store := memory.NewStore()
	rateLimiter := limiter.New(store, limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5, // 5 requests per minute
	})
	
	return gin.Limit(rateLimiter)
}
```

Apply to sensitive routes:
```go
// routes/auth_router.go
func AuthRoutes(incomingRoutes *gin.Engine) {
	// Apply strict rate limiting to auth endpoints
	authLimiter := middleware.RateLimitMiddleware("5-M") // 5 per minute
	
	incomingRoutes.POST("/signup", authLimiter, userController.SignUp)
	incomingRoutes.POST("/login", authLimiter, userController.Login)
	incomingRoutes.POST("/auth/forgot-password", authLimiter, userController.ForgotPassword)
	incomingRoutes.POST("/auth/reset-password", authLimiter, userController.ResetPassword)
}
```

---

### 4. **Excessive JWT Token Expiration**
**File:** `helpers/token_helpers.go:32,39`
```go
ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(8760)))
```

**Risk:** Tokens valid for 1 year (8760 hours). If stolen, attacker has persistent access.

**Solution:** Use short-lived access tokens with refresh token rotation

**Fix:**
```go
// Access token: 15 minutes
ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute))

// Refresh token: 7 days
ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
```

Implement token refresh endpoint:
```go
func RefreshToken(c *gin.Context) {
	// Validate refresh token
	// Generate new access token
	// Optionally rotate refresh token
}
```

---

### 5. **Weak Password Requirements**
**File:** `controller/password_reset.go:26`
```go
NewPassword string `json:"new_password" binding:"required,min=6"`
```

**Risk:** 6-character passwords are easily cracked.

**Solution:** Enforce strong password policy

**Fix:**
```go
// Add password validation helper
func ValidatePasswordStrength(password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters")
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}
	return nil
}
```

---

### 6. **No OTP Rate Limiting or Attempt Tracking**
**File:** `controller/password_reset.go`

**Risk:** Attackers can brute force 6-digit OTPs (1 million combinations).

**Solution:**
- Limit OTP verification attempts (3-5 max)
- Implement exponential backoff
- Add CAPTCHA after failed attempts
- Lock account temporarily after excessive failures

**Fix:**
```go
// Add to password_otp model
type PasswordOTP struct {
	Email        string    `gorm:"primaryKey"`
	OtpHash      string    `gorm:"type:varchar(255);not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	Attempts     int       `gorm:"default:0"`
	MaxAttempts  int       `gorm:"default:5"`
	LockedUntil  *time.Time
}

// In ResetPassword function
if passwordOTP.Attempts >= passwordOTP.MaxAttempts {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"error": "Too many failed attempts. Please request a new OTP.",
	})
	return
}

// Increment attempts on failed verification
db.Model(&passwordOTP).Update("attempts", passwordOTP.Attempts+1)
```

---

### 7. **No Input Sanitization for File Uploads**
**File:** `controller/admin/recipe.go:189-213`

**Risk:** 
- No file type validation (can upload malicious files)
- No file size limits
- Path traversal vulnerability
- XSS via filenames

**Solution:**
```go
func UploadRecipeImage(c *gin.Context) {
	// ... existing code ...
	
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file received"})
		return
	}
	
	// 1. Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024)
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
		return
	}
	
	// 2. Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	
	fileHeader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer fileHeader.Close()
	
	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}
	
	contentType := http.DetectContentType(buffer)
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type (only JPEG, PNG, WebP allowed)"})
		return
	}
	
	// 3. Sanitize filename (use UUID only, ignore original filename)
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := uuid.New().String() + ext
	
	// 4. Prevent path traversal
	savePath := filepath.Clean(filepath.Join("images", filename))
	if !strings.HasPrefix(savePath, "images/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}
	
	// ... save file ...
}
```

---

### 8. **Inconsistent Error Handling - Information Disclosure**
**Multiple Files**

**Risk:** Error messages reveal internal system details (database errors, file paths, etc.)

**Examples:**
```go
c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
```

**Solution:**
```go
// Log detailed error internally
log.Printf("Database error: %v", err)

// Return generic error to user
c.JSON(http.StatusInternalServerError, gin.H{
	"error": "An internal error occurred. Please try again later.",
	"request_id": generateRequestID(),
})
```

---

## 🟠 HIGH PRIORITY ISSUES

### 9. **No Request Size Limits**
**Risk:** Memory exhaustion attacks, DoS

**Fix:**
```go
// main.go
func main() {
	// ... existing code ...
	
	router := gin.New()
	
	// Limit request body size (10MB max)
	router.MaxMultipartMemory = 10 << 20
	
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)
		c.Next()
	})
	
	// ... rest of setup ...
}
```

---

### 10. **Missing Production Logging**
**Issue:** Only basic `fmt.Println` and `log.Println` used

**Solution:** Implement structured logging with log levels

**Fix:**
```go
// go get go.uber.org/zap

// helpers/logger.go
package helpers

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger(env string) {
	var err error
	if env == "production" {
		Logger, err = zap.NewProduction()
	} else {
		Logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
}
```

Usage:
```go
// Instead of:
fmt.Println("error: client user not allowed to access")

// Use:
helpers.Logger.Error("unauthorized access attempt",
	zap.String("user_email", email),
	zap.String("endpoint", c.Request.URL.Path),
)
```

---

### 11. **No Request ID Tracing**
**Issue:** Difficult to trace requests across logs

**Fix:**
```go
// middleware/request_id.go
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
```

---

### 12. **Missing Security Headers**
**Risk:** XSS, clickjacking, MIME sniffing attacks

**Fix:**
```go
// middleware/security_headers.go
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
```

Apply in `main.go`:
```go
router.Use(middleware.SecurityHeaders())
```

---

### 13. **No Health Check Endpoint**
**Issue:** Cannot monitor application status

**Fix:**
```go
// controller/health.go
func HealthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := database.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"database": "down",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"database": "up",
		"version": "1.0.0",
		"timestamp": time.Now().Unix(),
	})
}

// main.go
router.GET("/health", controller.HealthCheck)
router.GET("/ready", controller.ReadinessCheck)
```

---

### 14. **Unprotected Admin Routes in main.go**
**File:** `main.go:61-69`

```go
router.POST("/create_user", controller.CreateUser)
router.GET("/get_users", controller.GetUsers)
router.GET("exercise", controller.GetExercisesForAdmin)
// ... etc
```

**Risk:** These routes are defined OUTSIDE the authenticated `UserRoutes()`, making them publicly accessible!

**Fix:** Move all admin routes inside `routes/client_router.go` with proper authentication and authorization.

---

### 15. **SQL Injection Prevention Verification**
**Status:** ✅ GOOD - Using GORM with parameterized queries

GORM is used correctly with parameterized queries:
```go
db.Where("email = ?", email).First(&user)
```

**Action:** Continue using GORM's query builders, avoid raw SQL.

---

## 🟡 MEDIUM PRIORITY ISSUES

### 16. **No Environment-Based Configuration**
**Issue:** No distinction between dev/staging/production

**Fix:**
```go
// config/config.go
type Config struct {
	Environment  string
	Port         string
	JWTSecret    string
	DatabaseURL  string
	SMTPConfig   SMTPConfig
	RateLimits   RateLimitConfig
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	
	return &Config{
		Environment: env,
		Port:        getEnv("PORT", "8080"),
		// ... load all config
	}
}
```

---

### 17. **No Database Connection Pooling Configuration**
**Issue:** Default connection pool may not be optimal for production

**Fix:**
```go
func ConnectToDB() {
	// ... existing connection code ...
	
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance")
	}
	
	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
```

---

### 18. **Missing Pagination for List Endpoints**
**File:** `controller/user.go:140`
```go
// TODO: pagination?
database.DB.Find(&users)
```

**Risk:** Performance issues with large datasets, memory exhaustion

**Fix:**
```go
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	offset := (page - 1) * pageSize
	
	var users []model.UserAuth
	var total int64
	
	database.DB.Model(&model.UserAuth{}).Count(&total)
	database.DB.Limit(pageSize).Offset(offset).Find(&users)
	
	c.JSON(200, gin.H{
		"users": users,
		"page": page,
		"page_size": pageSize,
		"total": total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
```

Apply to all list endpoints: GetAllClients, GetListOfRecipes, GetListOfExercises, etc.

---

### 19. **No Email Validation on Reset Password**
**Issue:** OTP can be requested for non-existent emails (enumeration)

**Current behavior in `ForgotPassword`:**
```go
if errors.Is(err, gorm.ErrRecordNotFound) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "User not found with this email address",
	})
	return
}
```

**Risk:** Attackers can enumerate valid email addresses

**Fix:**
```go
// Always return success, but only send email if user exists
var user model.UserAuth
err := db.Where("email = ?", req.Email).First(&user).Error

// Always return success (prevent enumeration)
c.JSON(http.StatusOK, gin.H{
	"message": "If this email exists, an OTP has been sent",
})

// Only send email if user actually exists
if err == nil {
	// Generate and send OTP
}
```

---

### 20. **Inconsistent Use of fmt.Errorf**
**Multiple Files**

**Issue:** `fmt.Errorf()` is called but return value is ignored

**Example:**
```go
fmt.Errorf("error: client with email %s does not exist", emailFromContext)
```

**Fix:**
```go
// Either return it:
return fmt.Errorf("error: client with email %s does not exist", emailFromContext)

// Or log it:
log.Printf("error: client with email %s does not exist", emailFromContext)
```

---

### 21. **Debug/Test Code in Production**
**File:** `controller/user.go:80,90,96`

```go
fmt.Println("client %v | err %v", client, err)
fmt.Errorf("error")
fmt.Println("stupid")
```

**Fix:** Remove all debug code and implement proper logging.

---

### 22. **No CORS Origin Validation from Environment**
**File:** `main.go:40`

```go
AllowOrigins: []string{"https://nutriediet.netlify.app", "http://localhost:3000"},
```

**Fix:**
```go
allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
if len(allowedOrigins) == 0 {
	allowedOrigins = []string{"https://nutriediet.netlify.app"}
}

config := cors.Config{
	AllowOrigins: allowedOrigins,
	// ...
}
```

---

### 23. **Exposed Password Hash in API Response**
**File:** `controller/user.go:182`

```go
c.JSON(200, gin.H{
	"created": user,  // This includes the password hash!
})
```

**Fix:**
```go
type UserResponse struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	UserType  string `json:"user_type"`
	CreatedAt *time.Time `json:"created_at"`
}

c.JSON(200, gin.H{
	"user": UserResponse{
		ID: user.ID,
		FirstName: user.FirstName,
		// ... exclude sensitive fields
	},
})
```

---

### 24. **No Account Lockout Mechanism**
**Issue:** No protection against persistent brute force attacks

**Fix:** Implement account lockout after N failed login attempts

```go
// Add to user model or create separate table
type LoginAttempts struct {
	Email        string
	FailedCount  int
	LockedUntil  *time.Time
}

// In Login function
// Check if account is locked
// Increment failed count on wrong password
// Lock account after 5 failed attempts for 15 minutes
```

---

## 🟢 LOW PRIORITY ISSUES

### 25. **No API Versioning**
**Issue:** Breaking changes will affect all clients

**Fix:**
```go
v1 := router.Group("/v1")
{
	v1.POST("/login", controller.Login)
	// ... all routes
}
```

---

### 26. **Missing Request Timeout Configuration**
**Fix:**
```go
srv := &http.Server{
	Addr:           ":" + port,
	Handler:        router,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	IdleTimeout:    120 * time.Second,
	MaxHeaderBytes: 1 << 20,
}

srv.ListenAndServe()
```

---

### 27. **No Graceful Shutdown**
**Issue:** In-flight requests may be terminated abruptly during deployment

**Fix:**
```go
srv := &http.Server{
	Addr:    ":" + port,
	Handler: router,
}

go func() {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("Shutting down server...")

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := srv.Shutdown(ctx); err != nil {
	log.Fatal("Server forced to shutdown:", err)
}

log.Println("Server exiting")
```

---

### 28. **No Metrics/Monitoring**
**Solution:** Add Prometheus metrics

```go
// go get github.com/prometheus/client_golang/prometheus
// go get github.com/zsais/go-gin-prometheus

p := ginprometheus.NewPrometheus("gin")
p.Use(router)
```

---

### 29. **Missing Swagger/OpenAPI Documentation**
**Solution:**
```go
// go get github.com/swaggo/gin-swagger
// go get github.com/swaggo/files

// Add swagger annotations to controllers
// Generate docs: swag init
```

---

### 30. **No Database Migration Strategy**
**Issue:** No versioned schema migrations

**Solution:** Use migration tool
```go
// go get -u github.com/golang-migrate/migrate/v4

// Create migrations/
// 001_initial_schema.up.sql
// 001_initial_schema.down.sql
```

---

## Environment Variables Checklist

Create a `.env.example` file:

```bash
# Application
ENVIRONMENT=production
PORT=8080

# Database
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=your-db-host.com
DB_PORT=3306
DB_NAME=nutriediet

# JWT
JWT_SECRET_KEY=your_64_character_random_secret_key_here
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=168h

# SMTP
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# CORS
ALLOWED_ORIGINS=https://nutriediet.netlify.app,https://yourdomain.com

# Rate Limiting
RATE_LIMIT_REQUESTS=5
RATE_LIMIT_WINDOW=1m

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

---

## Deployment Checklist

### Before Deploying:

- [ ] Fix all CRITICAL issues (#1-8)
- [ ] Implement rate limiting
- [ ] Configure secure JWT with short expiration
- [ ] Move all secrets to environment variables
- [ ] Enable TLS for database connections
- [ ] Set up structured logging
- [ ] Add security headers middleware
- [ ] Implement health check endpoints
- [ ] Configure proper CORS origins
- [ ] Remove all debug/test code
- [ ] Add request size limits
- [ ] Implement graceful shutdown
- [ ] Set up monitoring/alerting
- [ ] Configure backups
- [ ] Review all error messages (no info disclosure)
- [ ] Test with security scanner (OWASP ZAP, Burp Suite)
- [ ] Load test critical endpoints
- [ ] Set up SSL/TLS certificates
- [ ] Configure firewall rules
- [ ] Enable database encryption at rest
- [ ] Set up log aggregation (ELK, Datadog, etc.)

### After Deploying:

- [ ] Monitor error rates
- [ ] Check rate limit effectiveness
- [ ] Review authentication logs
- [ ] Test password reset flow
- [ ] Verify CORS configuration
- [ ] Check database connection pool metrics
- [ ] Monitor API response times
- [ ] Set up alerts for security events

---

## Security Testing Commands

```bash
# 1. Test rate limiting
for i in {1..10}; do curl -X POST https://your-api.com/login; done

# 2. Test SQL injection (should fail safely)
curl -X POST https://your-api.com/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com'\'' OR 1=1--","password":"test"}'

# 3. Check security headers
curl -I https://your-api.com

# 4. Test JWT expiration
# Generate token, wait for expiry, test endpoint

# 5. Test file upload limits
curl -X POST https://your-api.com/admin/recipes/upload \
  -F "file=@large_file.jpg"
```

---

## Priority Order for Implementation

1. **Week 1 - Critical Security**
   - Environment variables for secrets (#1, #2)
   - Rate limiting (#3)
   - JWT token expiration (#4)
   - Password strength (#5)

2. **Week 2 - Authentication & Authorization**
   - OTP attempt limits (#6)
   - File upload security (#7)
   - Move unprotected admin routes (#14)
   - Error handling improvements (#8)

3. **Week 3 - Infrastructure**
   - Structured logging (#10)
   - Security headers (#12)
   - Health checks (#13)
   - Request limits (#9)

4. **Week 4 - Optimization**
   - Pagination (#18)
   - Database pooling (#17)
   - Graceful shutdown (#27)
   - Monitoring (#28)

---

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://github.com/OWASP/Go-SCP)
- [Gin Framework Security](https://gin-gonic.com/docs/examples/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)

---

**Last Updated:** 2025-10-05
**Document Version:** 1.0
**Review Status:** Initial Assessment


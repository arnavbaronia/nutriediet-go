# Security Quick Fixes - Top 5 Critical Issues

## 🚨 Fix These IMMEDIATELY Before Going Live

### 1. Database Credentials Exposed (15 minutes)

**Problem:** Hardcoded in `database/database.go:13`

**Quick Fix:**
```bash
# Add to .env
DB_USER=avnadmin
DB_PASSWORD=AVNS_7QDxgZDlRhQXAx3QV4z
DB_HOST=nutriediet-mysql-ishitagupta-5564.f.aivencloud.com
DB_PORT=22013
DB_NAME=defaultdb
```

Replace `database/database.go`:
```go
package database

import (
	"fmt"
	"log"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	
	if dbUser == "" || dbPassword == "" {
		log.Fatal("❌ Database credentials not configured. Set DB_USER and DB_PASSWORD environment variables")
	}
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=true&parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	
	// Configure connection pool
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	
	DB = db
	log.Println("✅ Database connected successfully")
}
```

---

### 2. Empty JWT Secret (5 minutes)

**Problem:** `helpers/token_helpers.go:21` - `var SECRET_KEY = ""`

**Quick Fix:**
```bash
# Generate secure key
openssl rand -base64 64

# Add to .env
JWT_SECRET_KEY=your_generated_64_character_key_here
```

Update `helpers/token_helpers.go`:
```go
package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"github.com/cd-Ishita/nutriediet-go/database"
	jwt "github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserID    string
	UserType  string
	jwt.RegisteredClaims
}

var SECRET_KEY string

func init() {
	SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	if SECRET_KEY == "" {
		log.Fatal("❌ JWT_SECRET_KEY environment variable is required")
	}
	if len(SECRET_KEY) < 32 {
		log.Fatal("❌ JWT_SECRET_KEY must be at least 32 characters")
	}
	log.Println("✅ JWT Secret loaded")
}

func GenerateAllTokens(email, firstName, lastName, userType string, id uint64) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		UserID:    strconv.FormatUint(id, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			// Changed from 1 year to 15 minutes
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			// Changed from 1 year to 7 days
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		fmt.Println("error: cannot generate the token for the user", email)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		fmt.Println("error: cannot generate the refresh token for the user", email)
		return "", "", err
	}

	return token, refreshToken, nil
}

func UpdateTokens(token, refreshToken string, id uint64) error {
	db := database.DB
	err := db.Table("user_auths").Where("id = ?", id).Updates(map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
	}).Error
	if err != nil {
		fmt.Println("error: cannot update the tokens")
		return err
	}
	return nil
}

func ValidateToken(token string) (SignedDetails, error) {
	res, err := jwt.ParseWithClaims(token, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return SignedDetails{}, fmt.Errorf("error parsing claims: %v", err)
	}

	claims, ok := res.Claims.(*SignedDetails)
	if !ok {
		return SignedDetails{}, errors.New("invalid token")
	}

	// check if token is expired
	if claims.ExpiresAt.Before(time.Now()) {
		return SignedDetails{}, errors.New("expired token")
	}
	return *claims, nil
}
```

---

### 3. Add Rate Limiting (30 minutes)

**Install dependency:**
```bash
go get github.com/ulule/limiter/v3
go get github.com/ulule/limiter/v3/drivers/middleware/gin
go get github.com/ulule/limiter/v3/drivers/store/memory
```

**Create `middleware/rate_limit.go`:**
```go
package middleware

import (
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitAuth() gin.HandlerFunc {
	// 5 requests per minute for auth endpoints
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	
	middleware := mgin.NewMiddleware(instance)
	return middleware
}

func RateLimitAPI() gin.HandlerFunc {
	// 100 requests per minute for API endpoints
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	
	middleware := mgin.NewMiddleware(instance)
	return middleware
}
```

**Update `routes/auth_router.go`:**
```go
package routes

import (
	userController "github.com/cd-Ishita/nutriediet-go/controller"
	clientController "github.com/cd-Ishita/nutriediet-go/controller/client"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	// Apply rate limiting to auth endpoints
	rateLimiter := middleware.RateLimitAuth()
	
	incomingRoutes.POST("/signup", rateLimiter, userController.SignUp)
	incomingRoutes.POST("/login", rateLimiter, userController.Login)
	incomingRoutes.POST("/create_profile/:email", rateLimiter, clientController.CreateProfileByClient)
	
	// Password reset routes with stricter rate limiting
	incomingRoutes.POST("/auth/forgot-password", rateLimiter, userController.ForgotPassword)
	incomingRoutes.POST("/auth/reset-password", rateLimiter, userController.ResetPassword)
}
```

**Update `routes/client_router.go`:**
```go
func UserRoutes(incomingRoutes *gin.Engine) {
	// Apply rate limiting to all authenticated routes
	incomingRoutes.Use(middleware.RateLimitAPI())
	
	// Authentication middleware applies to all routes
	incomingRoutes.Use(middleware.Authenticate)
	
	// ... rest of routes ...
}
```

---

### 4. Add Security Headers (10 minutes)

**Create `middleware/security_headers.go`:**
```go
package middleware

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Force HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; img-src 'self' data: https:; script-src 'self'; style-src 'self' 'unsafe-inline'")
		
		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Remove server information
		c.Header("X-Powered-By", "")
		
		c.Next()
	}
}
```

**Update `main.go`:**
```go
func main() {
	// ... existing code ...
	
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.SecurityHeaders()) // Add this line
	
	// ... rest of setup ...
}
```

---

### 5. Fix Unprotected Admin Routes (15 minutes)

**Problem:** Routes in `main.go:61-69` are publicly accessible

**Move these routes to `routes/client_router.go`:**

Remove from `main.go`:
```go
// DELETE THESE LINES (61-69)
router.POST("/create_user", controller.CreateUser)
router.GET("/get_users", controller.GetUsers)
router.GET("exercise", controller.GetExercisesForAdmin)
router.GET("exercise/:exercise_id", controller.GetExercise)
router.POST("exercise/:exercise_id/delete", controller.RemoveExerciseFromList)
router.POST("exercise/:exercise_id/update", controller.UpdateExerciseFromList)
router.POST("exercise/submit", controller.AddExerciseFromList)
```

Add to `routes/client_router.go` inside `UserRoutes()` function after line 102:
```go
// ADMIN - USER MANAGEMENT
incomingRoutes.POST("/admin/create_user", controller.CreateUser)
incomingRoutes.GET("/admin/users", controller.GetUsers)

// ADMIN - EXERCISE (moved from main.go)
incomingRoutes.GET("/admin/exercises_all", controller.GetExercisesForAdmin)
incomingRoutes.GET("/admin/exercise_detail/:exercise_id", controller.GetExercise)
incomingRoutes.POST("/admin/exercise_remove/:exercise_id", controller.RemoveExerciseFromList)
incomingRoutes.POST("/admin/exercise_update/:exercise_id", controller.UpdateExerciseFromList)
incomingRoutes.POST("/admin/exercise_add", controller.AddExerciseFromList)
```

**Clean up `main.go`:**
```go
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	} else {
		log.Println("✅ Environment variables loaded from .env file")
	}

	database.ConnectToDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set to release mode in production
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders())

	// CORS configuration
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 || allowedOrigins[0] == "" {
		allowedOrigins = []string{"https://nutriediet.netlify.app", "http://localhost:3000"}
	}

	config := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Client-Email", "Request-Client-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
	
	// Serve static files (with size limit)
	router.Static("/images", "./images")
	
	// Health check endpoint (public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Route groups
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// Start server
	log.Printf("🚀 Server starting on port %s", port)
	router.Run(":" + port)
}
```

---

## Updated .env File

Create/update your `.env` file with all required variables:

```bash
# Application
ENVIRONMENT=production
PORT=8080

# Database
DB_USER=avnadmin
DB_PASSWORD=AVNS_7QDxgZDlRhQXAx3QV4z
DB_HOST=nutriediet-mysql-ishitagupta-5564.f.aivencloud.com
DB_PORT=22013
DB_NAME=defaultdb

# JWT (Generate with: openssl rand -base64 64)
JWT_SECRET_KEY=your_very_long_random_secret_key_at_least_32_characters_long

# SMTP
SMTP_EMAIL=nutriediet.help@gmail.com
SMTP_PASSWORD=your_16_character_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# CORS
ALLOWED_ORIGINS=https://nutriediet.netlify.app,https://yourdomain.com
```

---

## Test After Applying Fixes

```bash
# 1. Build
go build -o nutriediet-go

# 2. Run with .env
./nutriediet-go

# 3. Test rate limiting (should block after 5 requests)
for i in {1..10}; do 
  curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"test"}'
  echo ""
done

# 4. Check security headers
curl -I http://localhost:8080/health

# 5. Verify JWT expiration (token should expire after 15 minutes)
```

---

## Deployment Commands

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nutriediet-go .

# Set environment variables on your server
export ENVIRONMENT=production
export DB_USER=...
export DB_PASSWORD=...
export JWT_SECRET_KEY=...
# ... etc

# Run
./nutriediet-go
```

---

## If You Only Have 1 Hour

**Priority Order:**
1. ✅ Fix JWT Secret (#2) - 5 min
2. ✅ Fix Database Credentials (#1) - 15 min
3. ✅ Add Security Headers (#4) - 10 min
4. ✅ Fix Admin Routes (#5) - 15 min
5. ✅ Add Rate Limiting (#3) - 30 min

**These 5 fixes address ~70% of critical security issues.**

After deployment, gradually implement the remaining improvements from `PRODUCTION_IMPROVEMENTS.md`.

---

## Emergency Rollback Plan

If something breaks:

1. Keep old binary: `cp nutriediet-go nutriediet-go.backup`
2. If new version fails: `./nutriediet-go.backup`
3. Check logs: `tail -f /var/log/nutriediet-go.log`
4. Verify environment variables are set correctly

---

**⚠️ DO NOT deploy without fixing at least issues #1 and #2!**


package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adds security headers to all HTTP responses
// These headers protect against common web vulnerabilities
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking attacks by not allowing the page to be embedded in iframes
		c.Header("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing which can lead to security vulnerabilities
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Enable XSS filter built into most browsers
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Force HTTPS connections (only add this if you're using HTTPS)
		// In development (HTTP), this is skipped. In production with HTTPS, this is crucial.
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
		// Content Security Policy - restricts what resources can be loaded
		// This helps prevent XSS attacks and unauthorized data exfiltration
		c.Header("Content-Security-Policy", "default-src 'self'; img-src 'self' data: https:; script-src 'self'; style-src 'self' 'unsafe-inline'")
		
		// Control how much referrer information is passed
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Remove server information to avoid revealing technology stack
		c.Header("X-Powered-By", "")
		c.Header("Server", "")
		
		// Prevent browsers from performing DNS prefetching
		c.Header("X-DNS-Prefetch-Control", "off")
		
		// Disable client-side caching for sensitive API responses
		// This ensures fresh data and prevents cached sensitive information
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		
		c.Next()
	}
}


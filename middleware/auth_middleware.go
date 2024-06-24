package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Println("no authorization header received")
		c.JSON(http.StatusBadRequest, gin.H{"err": "no token received"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		fmt.Println("invalid authorization header format")
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid token format"})
		c.Abort()
		return
	}

	clientToken := parts[1]
	claims, err := helpers.ValidateToken(clientToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		c.Abort()
		return
	}

	// setting context with this information
	c.Set("email", claims.Email)
	c.Set("first_name", claims.FirstName)
	c.Set("last_name", claims.LastName)
	c.Set("user_type", claims.UserType)
	c.Set("user_id", claims.UserID)
	c.Next()
}

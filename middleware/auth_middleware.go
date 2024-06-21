package middleware

import (
	"fmt"
	"net/http"

	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		fmt.Println("no client token received")
		c.JSON(http.StatusBadRequest, gin.H{"err": "no token received"})
		c.Abort()
		return
	}

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

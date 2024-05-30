package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
}

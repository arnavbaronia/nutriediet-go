package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func MatchUserTypeToUid(c *gin.Context, userID string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")

	if userType == "USER" && uid != userID {
		return errors.New("error: Not Authorised")
	} else if userType == "ADMIN" || (userType == "USER" && uid == userID) {
		return nil
	}
	return errors.New("error: Not Authorised")
}

func CheckUserType(c *gin.Context, userType string) bool {
	if c.GetString("user_type") == userType {
		return true
	}
	return false
}

package controller

import (
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(c *gin.Context) {
	CreateUser(c)
}

func Login(c *gin.Context) {
	user := model.User{}
	if err := c.BindJSON(&user); err != nil {
		fmt.Errorf("error: request cannot be parsed #{err}")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	// find the record with this email id
	dbRecord := model.User{}
	err := db.Where("email = ?", user.Email).First(&dbRecord).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: cannot find a record for email: %s", user.Email)
		c.JSON(http.StatusNotFound, gin.H{"err": "Record Not Found"})
		return
	} else if err != nil {
		fmt.Errorf("error: cannot extract record for email: %s", user.Email)
		c.JSON(http.StatusNotFound, gin.H{"err": "Can't extract record"})
		return
	}

	valid, err := VerifyPassword(user.Password, dbRecord.Password)
	if err != nil {
		fmt.Printf("error: password does not match for user: %s\n", user.Email)
		c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
		// TODO: verify the err makes sense to send back to client
		return
	}

	if !valid {
		fmt.Println("error: invalid password")
		c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
		// TODO: verify the err makes sense to send back to client
		return
	}

	token, refreshToken, err := helpers.GenerateAllTokens(dbRecord.Email, dbRecord.FirstName, dbRecord.LastName, dbRecord.UserType, dbRecord.ID)
	if err != nil {
		fmt.Printf("error: cannot generate tokens for user: %s\n", dbRecord.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	err = helpers.UpdateTokens(token, refreshToken, dbRecord.ID)
	if err != nil {
		fmt.Println("error: cannot update tokens for user: %s", dbRecord.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true"})
	return
}

// only admin should have access or get own user's data?
func GetUser(c *gin.Context) {
	clientId := c.Param("user_id")

	err := helpers.MatchUserTypeToUid(c, clientId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	db := database.DB
	user := model.User{}
	err = db.Where("id = ?", clientId).Find(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
	return
}

func GetUsers(c *gin.Context) {
	var users []model.User

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Println("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "authorized access by client"})
		return
	}
	// TODO: pagination?
	database.DB.Find(&users)
	c.JSON(200, gin.H{
		"message": users,
	})
}

func CreateUser(c *gin.Context) {
	// get data from req
	user := model.User{}
	if err := c.BindJSON(&user); err != nil {
		fmt.Errorf("error: request cannot be parsed #{err}")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: add a struct validation before inserting in DB
	// what if user already exists?
	token, refreshToken, err := helpers.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.UserType, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	user.Token = token
	user.RefreshToken = refreshToken

	hashedPassword, err := HashPassword(user.Password)
	user.Password = hashedPassword
	// store in DB
	err = database.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	// return
	c.JSON(200, gin.H{
		"created": user,
	})
}

func HashPassword(password string) (string, error) {
	newPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("error: generating hash from password")
		return "", err
	}
	return string(newPassword), nil
}

func VerifyPassword(providedPassword string, savedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(providedPassword))
	if err != nil {
		fmt.Println("error: email or password incorrect")
		return false, err
	}
	return true, nil
}

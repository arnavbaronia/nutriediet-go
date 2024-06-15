package helpers

import (
	"errors"
	"fmt"
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

var SECRET_KEY = ""

func GenerateAllTokens(email, firstName, lastName, userType string, id uint64) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		UserID:    strconv.FormatUint(id, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			// TODO: finalise the expires at values
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			// TODO: finalise the expires at values
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(168))),
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
		return SignedDetails{}, err
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

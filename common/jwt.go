package common

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type SignedDetails struct {
	ID    string
	Email string

	jwt.StandardClaims
}

func GenerateToken(id uint64, email string) (signed string, err error) {
	secretKey := []byte(viper.GetString("SECRET_KEY"))
	fmt.Println(secretKey)
	jwtLifeHour, err := strconv.Atoi(viper.GetString("JWT_LIFE"))
	if err != nil {
		log.Panic(err)
	}
	jwtLifeHour = jwtLifeHour * time.Now().Hour()

	claims := &SignedDetails{
		ID:    strconv.FormatUint(id, 10),
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(jwtLifeHour)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	return token, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	secretKey := os.Getenv("SECRET_KEY")

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "The token is expired"
		return
	}
	return claims, msg
}

func Auth(signedToken string) *SignedDetails {
	claims, msg := ValidateToken(signedToken)
	if msg != "" {
		log.Fatalf(msg)
	}
	return claims
}

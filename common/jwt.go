package common

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type SignedDetails struct {
	ID          string
	Email       string
	PhoneNumber string
	jwt.StandardClaims
}

func GenerateToken(id string, email string, phone_number string) (signed string, err error) {
	secretKey := os.Getenv("SECRET_KEY")
	jwtLifeHour, err := strconv.Atoi(os.Getenv("JWT_LIFE_HOUR"))
	if err != nil {
		log.Panic(err)
	}
	jwtLifeHour = jwtLifeHour * time.Now().Hour()

	claims := &SignedDetails{
		ID:          id,
		Email:       email,
		PhoneNumber: phone_number,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(jwtLifeHour)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(secretKey))
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

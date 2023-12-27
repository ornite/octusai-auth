package utils

import (
	models "auth/src/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret  = []byte("your_secret_key") // Replace with your secret key

// CustomClaims includes the standard JWT claims and additional data
type CustomClaims struct {
    jwt.StandardClaims
    UserID string `json:"userId"`
}

func GenerateToken(user models.User) (string, error) {
    // Create the Claims
    claims := CustomClaims{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
        UserID: user.ID,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(hmacSampleSecret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
    claims := &CustomClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return hmacSampleSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil , errors.New("Invalid token")
    }

    return claims, nil
}

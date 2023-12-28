package utils

import (
	models "auth/src/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret = []byte("your_secret_key") // Replace with your secret key

// CustomClaims includes the standard JWT claims and additional data
type CustomClaims struct {
    jwt.StandardClaims
    UserID string `json:"userId"`
}

func GenerateToken(user models.User, isExp bool, expTime *float64) (string, error) {
    // Set expiration time if isExp is true
    var expiresAt int64
    if isExp && expTime != nil {
        expiresAt = time.Now().Add(time.Duration(*expTime) * time.Second).Unix()
    } else {
        // Default to 24 hours if isExp is not specified or expTime is nil
        expiresAt = time.Now().Add(24 * time.Hour).Unix()
    }

    // Create the Claims
    claims := CustomClaims{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expiresAt,
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
    // Initialize a new instance of `CustomClaims`
    claims := &CustomClaims{}

    // Parse the token
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        return hmacSampleSecret, nil
    })

    if err != nil {
        // Handle specific token parsing error (e.g., expired)
        if ve, ok := err.(*jwt.ValidationError); ok {
            if ve.Errors&jwt.ValidationErrorExpired != 0 {
                // Token is expired
                return nil, errors.New("Token is expired")
            } else {
                // Other validation error
                return nil, errors.New("Token is invalid")
            }
        }
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("Token is invalid")
    }

    return claims, nil
}




package helpers

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// Secret key for JWT token signing
	secretKey = []byte(os.Getenv("JWT_SECRET"))
)

var (
    ErrMissingAuthorizationHeader = errors.New("authorization header is missing")
    ErrInvalidTokenFormat         = errors.New("invalid token format")
    ErrInvalidTokenSigningMethod  = errors.New("invalid token signing method")
    ErrInvalidToken               = errors.New("invalid token")
    ErrInvalidTokenClaims         = errors.New("invalid token claims")
)

// ErrExpiredToken represents an expired token error
var ErrExpiredToken = errors.New("token has expired")

// GenerateToken generates a JWT token with the provided user ID and email
func GenerateToken(id int) (string, error) {
    claims := jwt.MapClaims{
        "id":  id,
        "exp": time.Now().Add(time.Hour * 8).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return signedToken, nil
}

// VerifyToken verifies the JWT token provided in the request header
func VerifyToken(ctx *gin.Context) (jwt.MapClaims, error) {
    headerToken := ctx.Request.Header.Get("Authorization")
    if headerToken == "" {
        return nil, ErrMissingAuthorizationHeader
    }

    parts := strings.Split(headerToken, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return nil, ErrInvalidTokenFormat
    }

    stringToken := parts[1]
    token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
        if t.Method != jwt.SigningMethodHS256 {
            return nil, ErrInvalidTokenSigningMethod
        }
        return secretKey, nil
    })
    if err != nil || !token.Valid {
        return nil, ErrInvalidToken
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, ErrInvalidTokenClaims
    }

    exp, ok := claims["exp"].(float64)
    if !ok {
        return nil, errors.New("expiration time is missing or invalid")
    }

    if int64(exp) < time.Now().Unix() {
        return nil, ErrExpiredToken
    }

    return claims, nil
}
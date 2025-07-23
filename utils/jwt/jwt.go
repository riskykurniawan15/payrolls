package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	JWTConfig struct {
		SecretKey string
		Expired   int
	}

	JWTClaims struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		jwt.RegisteredClaims
	}

	UserContext struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}
)

// GenerateToken generates a new JWT token for the given user data
func GenerateToken(config JWTConfig, userID uint, username, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(config.Expired) * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, jwt.ErrInvalidType
	}

	return claims, nil
}

// ExtractUserContext extracts user context from JWT claims
func ExtractUserContext(claims *JWTClaims) UserContext {
	return UserContext{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}
}

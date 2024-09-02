package token_maker

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

type ContextKey string

const ClaimsKey ContextKey = "auth.claims"

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userId, role, secret string) (string, error) {
	claims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userId,
			Issuer:  os.Getenv("ISSUER"),

			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(3 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken(r *http.Request) (*CustomClaims, error) {
	claims := &CustomClaims{}

	tokenString, err := getBearerToken(r)
	if err != nil {
		return nil, errors.New("invalid token_maker")
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(key *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, errors.New("invalid token_maker")
	}

	if issuer, err := token.Claims.GetIssuer(); err != nil || issuer != os.Getenv("ISSUER") {
		return nil, errors.New("invalid token_maker")
	}

	if !token.Valid {
		return nil, errors.New("invalid token_maker")
	}

	return claims, nil
}

func getBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no Authorization header found")
	}

	// The Authorization header should start with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization header is not of type Bearer")
	}

	// Extract the token_maker by trimming the "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

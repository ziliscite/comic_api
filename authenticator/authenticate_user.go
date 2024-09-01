package authenticator

import (
	"bookstore/handlers"
	"bookstore/helpers"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strings"
)

// Have to be positioned after the normal auth

func EnsureAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(handlers.ClaimsKey).(*handlers.CustomClaims)
		if !ok {
			helpers.RespondWithError(w, http.StatusUnauthorized, "not authorized")
			return
		}

		if claims.Role != "admin" {
			helpers.RespondWithError(w, http.StatusUnauthorized, "not authorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getBearerToken(r)
		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		accessToken, claims, err := ValidateToken(tokenString)
		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), handlers.ClaimsKey, claims)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), handlers.AccessTokenKey, accessToken)))
	})
}

// ValidateToken Function to verify JWT tokens
func ValidateToken(tokenString string) (*jwt.Token, *handlers.CustomClaims, error) {
	claims := &handlers.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(key *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("Invalid token: %s", err)
		return nil, nil, fmt.Errorf("invalid token")
	}

	if issuer, err := token.Claims.GetIssuer(); err != nil || issuer != os.Getenv("ISSUER") {
		log.Printf("Invalid issuer: %s", err)
		return nil, nil, fmt.Errorf("invalid token")
	}

	if !token.Valid {
		log.Printf("Invalid token")
		return nil, nil, fmt.Errorf("invalid token")
	}

	return token, claims, nil
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

	// Extract the token by trimming the "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

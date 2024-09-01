package handlers

import (
	"bookstore/helpers"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ContextKey string

const (
	UserIDKey      ContextKey = "auth.userId"
	RoleKey        ContextKey = "auth.role"
	AccessTokenKey ContextKey = "auth.accessToken"
	ClaimsKey      ContextKey = "auth.claims"
)

// Hey, we can expand this so that the user have the choice to either login via email or username

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) (int, error) {
	type LoginWithEmail struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ctx := context.Background()

	userReq := LoginWithEmail{}
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		h.Logger.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	user, err := h.Queries.LoginWithEmail(ctx, userReq.Email)
	if err != nil {
		h.Logger.Printf("Error while logging in: %v", err)
		return http.StatusUnauthorized, errors.New("email or password does not match")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password))
	if err != nil {
		h.Logger.Printf("Password does not match: %v", err)
		return http.StatusUnauthorized, errors.New("email or password does not match")
	}

	// Yeah, we gon use this only when we validate the auth // Nope
	ctx = context.WithValue(r.Context(), UserIDKey, strconv.Itoa(int(user.UserID)))
	req := r.WithContext(ctx)

	ctx = context.WithValue(req.Context(), RoleKey, string(user.Role))
	req = r.WithContext(ctx)

	accessToken, err := GenerateJWT(req, h.JWTSecret)
	if err != nil {
		h.Logger.Printf("Error generating access token: %v", err)
		return http.StatusInternalServerError, errors.New("error generating access token")
	}

	w.Header().Set("Authorization", "Bearer "+accessToken)
	/*	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
	})*/

	helpers.RespondWithMessage(w, http.StatusOK, "login successful")
	return http.StatusOK, nil
}

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(r *http.Request, secret string) (string, error) {
	role, ok := r.Context().Value(RoleKey).(string)
	if !ok {
		return "", errors.New("role cannot be converted to a string")
	}

	subject, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return "", errors.New("user cannot be converted to a string")
	}

	claims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subject,
			Issuer:  os.Getenv("ISSUER"),

			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(3 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

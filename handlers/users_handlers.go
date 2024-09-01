package handlers

import (
	"bookstore/database"
	"bookstore/helpers"
	"bookstore/middlewares"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"unicode"
)

func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email")
	}

	return nil
}

func ValidatePassword(password string) error {
	var hasUpper, hasLower, hasNumber bool
	if len(password) < 8 {
		return errors.New("password is too short, it must contain at least 8 characters")
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain uppercase letters")
	}
	if !hasLower {
		return errors.New("password must contain lowercase letters")
	}
	if !hasNumber {
		return errors.New("password must contain numbers")
	}

	return nil
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) (int, error) {
	ctx := context.Background()

	userReq := database.RegisterUserParams{
		Role: "user",
	}

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		h.Logger.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	err = ValidateEmail(userReq.Email)
	if err != nil {
		h.Logger.Printf("Error validating email: %s", err)
		return http.StatusBadRequest, err
	}

	err = ValidatePassword(userReq.Password)
	if err != nil {
		h.Logger.Printf("Error validating password: %s", err)
		return http.StatusBadRequest, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Printf("Error hashing password: %s", err)
		return http.StatusBadRequest, err
	}

	userReq.Password = string(hashedPassword)

	user, err := h.Queries.RegisterUser(ctx, userReq)
	if err != nil {
		return handleUserError(h.Logger, err, userReq.Email, userReq.Username)
	}

	helpers.RespondWithJSON(w, http.StatusCreated, user)
	return http.StatusCreated, nil
}

func handleUserError(logger *middlewares.Logger, err error, email string, username string) (int, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" { // Unique constraint violation code
			switch pgErr.ConstraintName {
			case "users_email_key":
				logger.Printf("Email already registered: %s", email)
				return http.StatusBadRequest, errors.New("email already registered")
			case "users_username_key":
				logger.Printf("Username already registered: %s", username)
				return http.StatusBadRequest, errors.New("username already registered")
			default:
				logger.Printf("Unexpected unique constraint violation: %s", pgErr.ConstraintName)
				return http.StatusInternalServerError, errors.New("unexpected unique constraint violation")
			}
		}
	}

	logger.Printf("Non-Postgres error: %s", err)
	return http.StatusInternalServerError, errors.New("failed to create user")
}

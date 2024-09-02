package handler

import (
	"bookstore/database"
	"bookstore/token_maker"
	"bookstore/utils/helpers"
	"bookstore/utils/middlewares"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strconv"
	"unicode"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) (int, error) {
	type LoginWithEmail struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	userReq := LoginWithEmail{}
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		h.Middlewares.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	user, err := h.Queries.LoginWithEmail(h.Context, userReq.Email)
	if err != nil {
		h.Middlewares.Printf("Error while logging in: %v", err)
		return http.StatusUnauthorized, errors.New("email or password does not match")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password))
	if err != nil {
		h.Middlewares.Printf("Password does not match: %v", err)
		return http.StatusUnauthorized, errors.New("email or password does not match")
	}

	accessToken, err := token_maker.GenerateJWT(strconv.Itoa(int(user.UserID)), string(user.Role), h.JWTSecret)
	if err != nil {
		h.Middlewares.Printf("Error generating access token_maker: %v", err)
		return http.StatusInternalServerError, errors.New("error generating access token_maker")
	}

	w.Header().Set("Authorization", "Bearer "+accessToken)
	/*	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
	})*/

	helpers.RespondWithMessage(w, http.StatusOK, "login successful")
	return http.StatusOK, nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) (int, error) {
	userReq := database.RegisterUserParams{
		Role: "user",
	}

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		h.Middlewares.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	err = validateEmail(userReq.Email)
	if err != nil {
		h.Middlewares.Printf("Error validating email: %s", err)
		return http.StatusBadRequest, err
	}

	err = validatePassword(userReq.Password)
	if err != nil {
		h.Middlewares.Printf("Error validating password: %s", err)
		return http.StatusBadRequest, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Middlewares.Printf("Error hashing password: %s", err)
		return http.StatusBadRequest, err
	}

	userReq.Password = string(hashedPassword)

	user, err := h.Queries.RegisterUser(h.Context, userReq)
	if err != nil {
		return handleUserError(h.Middlewares, err, userReq.Email, userReq.Username)
	}

	helpers.RespondWithJSON(w, http.StatusCreated, user)
	return http.StatusCreated, nil
}

func handleUserError(logger *middlewares.Middleware, err error, email string, username string) (int, error) {
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

func validateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email")
	}

	return nil
}

func validatePassword(password string) error {
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

package handler

import (
	"bookstore/database"
	"bookstore/token_maker"
	"bookstore/utils/helpers"
	"bookstore/utils/middlewares"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strconv"
	"time"
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

	refreshToken, err := token_maker.GenerateRefreshToken()

	sessionParams := database.AddSessionParams{
		UserID:       &user.UserID,
		SessionToken: refreshToken,
		ExpiresAt: pgtype.Timestamp{
			Valid: true,
			Time:  time.Now().UTC().Add(60 * time.Hour),
		},
	}

	session, err := h.Queries.AddSession(h.Context, sessionParams)
	if err != nil {
		h.Middlewares.Printf("Error adding session: %v", err)
		return http.StatusInternalServerError, errors.New("error generating session")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    session.SessionToken,
		HttpOnly: true,
		Path:     "/refresh",
		Expires:  session.ExpiresAt.Time.UTC(),
	})

	w.Header().Set("Authorization", "Bearer "+accessToken)

	type LoginResponse struct {
		SessionId int32  `json:"session_id"`
		Email     string `json:"email"`

		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`

		AccessTokenExpires  pgtype.Timestamp `json:"access_token_expires_at"`
		RefreshTokenExpires pgtype.Timestamp `json:"refresh_token_expires_at"`
	}

	loginResp := LoginResponse{
		SessionId: session.SessionID,
		Email:     user.Email,

		AccessToken: accessToken,
		AccessTokenExpires: pgtype.Timestamp{
			Valid: true,
			Time:  time.Now().UTC().Add(2 * time.Hour),
		},

		RefreshToken:        refreshToken,
		RefreshTokenExpires: session.ExpiresAt,
	}

	helpers.RespondWithJSON(w, http.StatusOK, loginResp)
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

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) (int, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.Middlewares.Printf("Refresh token not present %s", err.Error())
			return http.StatusUnauthorized, errors.New("refresh token not found")
		}

		h.Middlewares.Printf("Error parsing refresh token: %s", err.Error())
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	// Access the cookie value
	tokenString := cookie.Value

	refreshToken, err := h.Queries.GetSessionFromToken(h.Context, tokenString)
	if err != nil {
		h.Middlewares.Printf("Error getting session from token: %s", err.Error())
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	if !refreshToken.ExpiresAt.Time.UTC().After(time.Now().UTC()) {
		h.Middlewares.Printf("Refresh token expired")

		err = h.Queries.RevokeSession(h.Context, refreshToken.SessionID)
		if err != nil {
			h.Middlewares.Printf("Error revoking session: %s", err.Error())
			return http.StatusInternalServerError, errors.New("something went wrong")
		}

		return http.StatusUnauthorized, errors.New("refresh token expired")
	}

	userRole, err := h.Queries.GetUserRole(h.Context, *refreshToken.UserID)
	if err != nil {
		h.Middlewares.Printf("Error getting user role: %s", err.Error())
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	accessToken, err := token_maker.GenerateJWT(strconv.Itoa(int(*refreshToken.UserID)), string(userRole.Role), h.JWTSecret)
	if err != nil {
		h.Middlewares.Printf("Error generating access token_maker: %v", err)
		return http.StatusInternalServerError, errors.New("error generating access token_maker")
	}

	w.Header().Set("Authorization", "Bearer "+accessToken)

	type RefreshResponse struct {
		SessionId          int32            `json:"session_id"`
		AccessToken        string           `json:"access_token"`
		AccessTokenExpires pgtype.Timestamp `json:"access_token_expires_at"`
	}

	refreshResponse := RefreshResponse{
		SessionId:   refreshToken.SessionID,
		AccessToken: accessToken,
		AccessTokenExpires: pgtype.Timestamp{
			Valid: true,
			Time:  time.Now().UTC().Add(2 * time.Hour),
		},
	}

	helpers.RespondWithJSON(w, http.StatusOK, refreshResponse)
	return http.StatusOK, nil
}

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) (int, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.Middlewares.Printf("Refresh token not present %s", err.Error())
			return http.StatusUnauthorized, errors.New("refresh token not found")
		}

		h.Middlewares.Printf("Error parsing refresh token: %s", err.Error())
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	tokenString := cookie.Value

	refreshToken, err := h.Queries.GetSessionFromToken(h.Context, tokenString)
	if err != nil {
		h.Middlewares.Printf("Error getting session from token: %s", err.Error())
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	err = h.Queries.RevokeSession(h.Context, refreshToken.SessionID)
	if err != nil {
		h.Middlewares.Printf("Error revoking session: %s", err.Error())
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	helpers.RespondWithMessage(w, http.StatusOK, "revoked successfully")
	return http.StatusOK, nil
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

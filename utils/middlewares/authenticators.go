package middlewares

import (
	"bookstore/token_maker"
	"bookstore/utils/helpers"
	"context"
	"net/http"
)

func (m *Middleware) EnsureAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(token_maker.ClaimsKey).(*token_maker.CustomClaims)
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

func (m *Middleware) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Maybe redirect to /refresh if token invalid?
		claims, err := token_maker.ValidateToken(r)
		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), token_maker.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

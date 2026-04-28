package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/context_key"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
)

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.NewUnauthorized(w, "unauthorized")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.NewUnauthorized(w, "unauthorized")
			return
		}

		token := parts[1]
		user, err := m.Model.GetUserByToken(r.Context(), token)
		if err != nil {
			response.NewUnauthorized(w, "unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), context_key.UserContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

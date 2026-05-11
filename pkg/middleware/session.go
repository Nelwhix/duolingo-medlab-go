package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/context_key"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
)

func (m *Middleware) SessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("duolingo_medlab_auth")
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}

		value := make(map[string]string)
		if err = m.CookieHandler.Decode("duolingo_medlab_auth", cookie.Value, &value); err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}

		userID, ok := value["user_id"]
		if !ok || userID == "" {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}

		user, err := m.Model.GetAdminUserById(r.Context(), userID)
		if err != nil {
			m.Logger.Error("Failed to retrieve user from database", slog.String("error", err.Error()))
			response.NewUnauthorized(w, "unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), context_key.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package middleware

import (
	"net/http"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/context_key"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/models"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
)

func (m *Middleware) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := context_key.GetUserFromContext(r.Context())
		if !ok {
			response.NewUnauthorized(w, "unauthorized")
			return
		}

		if user.Role != models.Admin {
			response.NewForbidden(w, "forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}

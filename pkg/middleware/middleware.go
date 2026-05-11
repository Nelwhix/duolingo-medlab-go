package middleware

import (
	"log/slog"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/models"
	"github.com/gorilla/securecookie"
)

type Middleware struct {
	Model         *models.Model
	CookieHandler *securecookie.SecureCookie
	Logger        *slog.Logger
}

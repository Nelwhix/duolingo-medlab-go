package handlers

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
)

func (h *Handler) RenderAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("duolingo_medlab_auth")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			parsedTemplate, _ := template.ParseFiles("./templates/admin/login.html")
			err := parsedTemplate.Execute(w, nil)
			if err != nil {
				h.Logger.Error("Failed to render admin login page", slog.String("error", err.Error()))
				h.renderAdminServerError(w)
				return
			}
		}

		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "duolingo_medlab_auth",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

func (h *Handler) RenderAdminDashboard(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, err := template.ParseFiles("./templates/admin/dashboard.html")
	if err != nil {
		h.Logger.Error("Failed to parse admin dashboard", slog.String("error", err.Error()))
		h.renderAdminServerError(w)
		return
	}

	if err := parsedTemplate.Execute(w, nil); err != nil {
		h.Logger.Error("Failed to render admin dashboard", slog.String("error", err.Error()))
		h.renderAdminServerError(w)
		return
	}
}

package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

const adminErrorTemplate = "./templates/admin/500.html"

func (h *Handler) renderAdminServerError(w http.ResponseWriter) {
	parsedTemplate, err := template.ParseFiles(adminErrorTemplate)
	if err != nil {
		h.Logger.Error("Failed to parse admin error page", slog.String("error", err.Error()))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var page bytes.Buffer
	if err := parsedTemplate.Execute(&page, nil); err != nil {
		h.Logger.Error("Failed to render admin error page", slog.String("error", err.Error()))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := page.WriteTo(w); err != nil {
		h.Logger.Error("Failed to write admin error page", slog.String("error", err.Error()))
	}
}

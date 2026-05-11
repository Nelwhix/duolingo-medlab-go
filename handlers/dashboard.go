package handlers

import (
	"html/template"
	"net/http"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
)

func (h *Handler) RenderAdminDashboard(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("./templates/admin/dashboard.html")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		response.NewInternalServerError(w, "Internal error")
		return
	}
}

func (h *Handler) RenderAdminCreateQuestion(w http.ResponseWriter, r *http.Request) {
	departments, err := h.Model.GetDepartments(r.Context())
	if err != nil {
		response.NewInternalServerError(w, "Internal error")
		return
	}

	parsedTemplate, err := template.ParseFiles("./templates/admin/create-question.html")
	if err != nil {
		response.NewInternalServerError(w, "Internal error")
		return
	}

	viewData := map[string]any{
		"Departments": departments,
	}
	renderErr := parsedTemplate.Execute(w, viewData)
	if renderErr != nil {
		response.NewInternalServerError(w, "Internal error")
		return
	}
}

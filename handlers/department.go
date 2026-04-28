package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
)

func (h *Handler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := h.Model.GetDepartments(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get departments", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Failed to get departments")
		return
	}

	response.NewOkResponseWithData(w, departments)
}

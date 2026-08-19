package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/request"
	"github.com/go-playground/validator/v10"
)

const topicValidationFlashCookie = "admin_topic_validation"

type topicValidationFlash struct {
	Errors map[string]string `json:"errors"`
	Old    map[string]string `json:"old"`
}

func (h *Handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.Logger.Error("Failed to parse request", slog.String("error", err.Error()))
		h.renderAdminServerError(w)
		return
	}

	var cRequest request.CreateTopic
	err = h.SchemaDecoder.Decode(&cRequest, r.PostForm)
	if err != nil {
		h.Logger.Error("Failed to decode request into struct", slog.String("error", err.Error()))
		h.renderAdminServerError(w)
		return
	}

	if err = h.Validator.Struct(cRequest); err != nil {
		validationErrors := make(map[string]string)
		for _, fieldError := range err.(validator.ValidationErrors) {
			switch fieldError.Field() {
			case "Name":
				validationErrors["topic"] = "The topic field is required."
			case "DepartmentID":
				validationErrors["department_id"] = "The department field is required."
			}
		}

		flash := topicValidationFlash{
			Errors: validationErrors,
			Old: map[string]string{
				"topic":         cRequest.Name,
				"department_id": cRequest.DepartmentID,
			},
		}
		encodedFlash, encodeErr := h.CookieHandler.Encode(topicValidationFlashCookie, flash)
		if encodeErr != nil {
			h.Logger.Error("Failed to encode topic validation errors", slog.String("error", encodeErr.Error()))
			h.renderAdminServerError(w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     topicValidationFlashCookie,
			Value:    encodedFlash,
			Path:     "/admin/topics/create",
			MaxAge:   60,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		http.Redirect(w, r, "/admin/topics/create", http.StatusSeeOther)
		return
	}

	err = h.Model.InsertTopic(r.Context(), cRequest)
	if err != nil {
		h.Logger.Error("Failed to create topic", slog.String("error", err.Error()))
		h.renderAdminServerError(w)
		return
	}

	http.Redirect(w, r, "/admin/topics", http.StatusSeeOther)
}

func (h *Handler) RenderAdminTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := h.Model.GetTopics(r.Context())
	if err != nil {
		h.Logger.Error("Failed to fetch topics", err.Error())
		h.renderAdminServerError(w)
		return
	}

	parsedTemplate, err := template.New("topics.html").Funcs(
		template.FuncMap{
			"formatDate": func(date time.Time) string {
				return date.Format("Jan 2, 2006 3:04 PM")
			},
		},
	).ParseFiles("./templates/admin/topics.html")
	if err != nil {
		h.Logger.Error("Failed to parse template files", "error", err)
		h.renderAdminServerError(w)
		return
	}

	viewData := map[string]any{
		"Topics": topics,
	}
	err = parsedTemplate.Execute(w, viewData)
	if err != nil {
		h.Logger.Error("Failed to render template files", err.Error())
		h.renderAdminServerError(w)
		return
	}
}

func (h *Handler) RenderAdminCreateTopic(w http.ResponseWriter, r *http.Request) {
	flash := topicValidationFlash{
		Errors: map[string]string{},
		Old:    map[string]string{},
	}
	if cookie, cookieErr := r.Cookie(topicValidationFlashCookie); cookieErr == nil {
		if decodeErr := h.CookieHandler.Decode(topicValidationFlashCookie, cookie.Value, &flash); decodeErr != nil {
			h.Logger.Warn("Failed to decode topic validation errors", slog.String("error", decodeErr.Error()))
		}
		http.SetCookie(w, &http.Cookie{
			Name:     topicValidationFlashCookie,
			Value:    "",
			Path:     "/admin/topics/create",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	departments, err := h.Model.GetDepartments(r.Context())
	if err != nil {
		h.renderAdminServerError(w)
		return
	}

	parsedTemplate, err := template.ParseFiles("./templates/admin/create-topic.html")
	if err != nil {
		h.Logger.Error("Failed to parse template files", "error", err)
		h.renderAdminServerError(w)
		return
	}

	viewData := map[string]any{
		"Departments": departments,
		"Errors":      flash.Errors,
		"Old":         flash.Old,
	}
	renderErr := parsedTemplate.Execute(w, viewData)
	if renderErr != nil {
		h.Logger.Error("Failed to render template data", renderErr.Error())
		h.renderAdminServerError(w)
		return
	}
}

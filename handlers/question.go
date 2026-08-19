package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/request"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
)

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.Logger.Error("Failed to parse request", slog.String("error", err.Error()))
		response.NewBadRequest(w, "Failed to parse request")
		return
	}

	correctOption, _ := strconv.Atoi(r.FormValue("correct_option"))
	requestData := request.CreateQuestion{
		Topic:         r.FormValue("topic"),
		Type:          request.QuestionType(r.FormValue("type")),
		DepartmentID:  r.FormValue("department_id"),
		Question:      r.FormValue("question"),
		Options:       r.Form["options[]"],
		CorrectOption: correctOption,
		Answer:        r.FormValue("answer"),
	}

	err = h.Validator.Struct(requestData)
	if err != nil {
		response.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	err = h.Model.InsertQuestion(r.Context(), requestData)
	if err != nil {
		h.Logger.Error("Failed to insert question", slog.String("error", err.Error()))
		h.renderAdminServerError(w)
		return
	}

	http.Redirect(w, r, "/admin/questions", http.StatusSeeOther)
}

func (h *Handler) RenderAdminQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.Model.GetQuestions(r.Context())
	if err != nil {
		h.Logger.Error("Failed to fetch questions", err.Error())
		h.renderAdminServerError(w)
		return
	}

	parsedTemplate, err := template.New("questions.html").Funcs(
		template.FuncMap{
			"formatDate": func(date time.Time) string {
				return date.Format("Jan 2, 2006 3:04 PM")
			},
		},
	).ParseFiles("./templates/admin/questions.html")
	if err != nil {
		h.Logger.Error("Failed to parse template files", "error", err)
		h.renderAdminServerError(w)
		return
	}

	viewData := map[string]any{
		"Questions": questions,
	}
	err = parsedTemplate.Execute(w, viewData)
	if err != nil {
		h.Logger.Error("Failed to render template files", err.Error())
		h.renderAdminServerError(w)
		return
	}
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	questionId := r.PathValue("id")
	err := h.Model.DeleteQuestion(r.Context(), questionId)
	if err != nil {
		h.Logger.Error("Failed to delete question", slog.String("error", err.Error()))
		h.renderAdminServerError(w)
		return
	}

	http.Redirect(w, r, "/admin/questions", http.StatusSeeOther)
}

func (h *Handler) RenderSingleAdminQuestion(w http.ResponseWriter, r *http.Request) {
	questionId := r.PathValue("id")
	question, err := h.Model.GetQuestionById(r.Context(), questionId)
	if err != nil {
		h.renderAdminServerError(w)
		return
	}

	if question == nil {
		response.NewNotFound(w, "Question not found")
		return
	}

	parsedTemplate, err := template.New("question.html").Funcs(
		template.FuncMap{
			"formatDate": func(date time.Time) string {
				return date.Format("Jan 2, 2006 3:04 PM")
			},
		},
	).ParseFiles("./templates/admin/question.html")
	if err != nil {
		h.renderAdminServerError(w)
		return
	}

	viewData := map[string]any{
		"Question": question,
	}
	err = parsedTemplate.Execute(w, viewData)
	if err != nil {
		h.renderAdminServerError(w)
		return
	}
}

func (h *Handler) RenderAdminCreateQuestion(w http.ResponseWriter, r *http.Request) {
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
	}
	renderErr := parsedTemplate.Execute(w, viewData)
	if renderErr != nil {
		h.Logger.Error("Failed to render template data", renderErr.Error())
		h.renderAdminServerError(w)
		return
	}
}

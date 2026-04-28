package handlers

import (
	"net/http"

	"github.com/Nelwhix/duolingo-medlab-go/pkg"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/request"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/resource"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
)

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		response.NewUnauthorized(w, "unauthorized")
		return
	}

	if user.ID != userId {
		response.NewForbidden(w, "forbidden")
		return
	}

	cRequest, err := pkg.ParseRequestBody[request.UpdateUser](r)
	if err != nil {
		response.NewUnprocessableEntity(w, err.Error())
		return
	}

	err = h.Validator.Struct(cRequest)
	if err != nil {
		response.NewUnprocessableEntity(w, err.Error())
		return
	}

	_, err = h.Model.GetDepartmentById(r.Context(), cRequest.DepartmentID)
	if err != nil {
		response.NewUnprocessableEntity(w, "Invalid department id")
		return
	}

	err = h.Model.UpdateUser(r.Context(), userId, cRequest)
	if err != nil {
		response.NewInternalServerError(w, "Failed to process request")
		return
	}

	user, err = h.Model.GetUserById(r.Context(), userId)
	if err != nil {
		response.NewInternalServerError(w, "Failed to process request")
		return
	}

	res := resource.UserResource{
		ID:                           user.ID,
		Username:                     user.Username,
		Email:                        user.Email,
		LearningGoalPassExams:        user.LearningGoalPassExams,
		LearningGoalRefreshKnowledge: user.LearningGoalRefreshKnowledge,
		LearningGoalPracticeDaily:    user.LearningGoalPracticeDaily,
		Department:                   resource.DepartmentResource{ID: user.Department.ID, Title: user.Department.Title},
	}

	response.NewOkResponseWithData(w, res)
}

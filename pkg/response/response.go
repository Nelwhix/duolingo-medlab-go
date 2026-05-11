package response

import (
	"encoding/json"
	"net/http"
)

type baseResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func NewUnprocessableEntity(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnprocessableEntity, baseResponse{Message: message})
}

func NewOKResponse(w http.ResponseWriter, message string) {
	JSON(w, http.StatusOK, baseResponse{Message: message})
}

func NewInternalServerError(w http.ResponseWriter, message string) {
	JSON(w, http.StatusInternalServerError, baseResponse{Message: message})
}

func NewNotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, baseResponse{Message: message})
}

func NewBadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, baseResponse{Message: message})
}

func NewUnauthorized(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnauthorized, baseResponse{Message: message})
}

func NewForbidden(w http.ResponseWriter, message string) {
	JSON(w, http.StatusForbidden, baseResponse{Message: message})
}

func NewCreatedResponseWithData(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, baseResponse{
		Message: message,
		Data:    data,
	})
}

func NewOkResponseWithData(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, baseResponse{
		Data: data,
	})
}

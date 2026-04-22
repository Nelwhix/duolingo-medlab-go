package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log/slog"
	"net/http"
	"time"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/mailer"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/thanhpk/randstr"
)

type Handler struct {
	Model     *models.Model
	Logger    *slog.Logger
	Validator *validator.Validate
	Mailer    mailer.Mailer
}

type baseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (h *Handler) JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.Logger.Error("Failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *Handler) NewUnprocessableEntity(w http.ResponseWriter, message string) {
	h.JSON(w, http.StatusUnprocessableEntity, baseResponse{Message: message})
}

func (h *Handler) NewOKResponse(w http.ResponseWriter, message string) {
	h.JSON(w, http.StatusOK, baseResponse{Message: message})
}

func (h *Handler) Pong(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("pong, time is: %v", time.Now().Format("2006-01-02 15:04"))

	h.NewOKResponse(w, message)
}

func (h *Handler) NewInternalServerError(w http.ResponseWriter, message string) {
	h.JSON(w, http.StatusInternalServerError, baseResponse{Message: message})
}

func (h *Handler) NewBadRequest(w http.ResponseWriter, message string) {
	h.JSON(w, http.StatusBadRequest, baseResponse{Message: message})
}

func (h *Handler) NewCreatedResponseWithData(w http.ResponseWriter, message string, data interface{}) {
	h.JSON(w, http.StatusCreated, baseResponse{
		Message: message,
		Data:    data,
	})
}

func (h *Handler) NewOkResponseWithData(w http.ResponseWriter, message string, data interface{}) {
	h.JSON(w, http.StatusOK, baseResponse{
		Message: message,
		Data:    data,
	})
}

func (h *Handler) CreateToken(ctx context.Context, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	expires := time.Now().Add(24 * time.Hour * 7)
	tokenString := h.generateTokenString()
	request := models.CreateTokenRequest{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: expires,
	}

	err := h.Model.InsertIntoTokens(ctx, request)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *Handler) generateTokenString() string {
	tokenEntropy := randstr.String(40)
	crc32bHash := crc32.ChecksumIEEE([]byte(tokenEntropy))

	return fmt.Sprintf(
		"%s-%d",
		tokenEntropy,
		crc32bHash,
	)
}

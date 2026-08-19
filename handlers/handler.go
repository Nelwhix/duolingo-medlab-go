package handlers

import (
	"context"
	"fmt"
	"hash/crc32"
	"log/slog"
	"net/http"
	"time"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/mailer"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/models"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/thanhpk/randstr"
)

type Handler struct {
	Model         *models.Model
	Logger        *slog.Logger
	Validator     *validator.Validate
	SchemaDecoder *schema.Decoder
	Mailer        mailer.Mailer
	CookieHandler *securecookie.SecureCookie
}

func (h *Handler) Pong(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("pong, time is: %v", time.Now().Format("2006-01-02 15:04"))

	response.NewOKResponse(w, message)
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

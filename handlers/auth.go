package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Nelwhix/duolingo-medlab-go/pkg"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/requests"
	"golang.org/x/crypto/bcrypt"
)

type UserResource struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	request, err := pkg.ParseRequestBody[requests.SignUp](r)
	if err != nil {
		h.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	err = h.Validator.Struct(request)
	if err != nil {
		h.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	_, err = h.Model.GetUserByEmail(r.Context(), request.Email)
	if err == nil {
		h.NewUnprocessableEntity(w, "Email already taken")
		return
	}

	user, err := h.Model.InsertIntoUsers(r.Context(), request)
	if err != nil {
		h.Logger.Error("Failed to insert user into database", slog.String("error", err.Error()))
		h.NewInternalServerError(w, "Failed to process request")
		return
	}

	response := UserResource{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	h.NewCreatedResponseWithData(w, "User created successfully.", response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	request, err := pkg.ParseRequestBody[requests.Login](r)
	if err != nil {
		h.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	err = h.Validator.Struct(request)
	if err != nil {
		h.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	user, err := h.Model.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.NewBadRequest(w, "Email or Password is incorrect")
			return
		}

		h.NewBadRequest(w, "Failed to process request")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		h.NewBadRequest(w, "Email or Password is incorrect")
		return
	}

	token, err := h.CreateToken(r.Context(), user.ID)
	if err != nil {
		h.Logger.Error("Failed to create token", slog.String("error", err.Error()))
		h.NewInternalServerError(w, "Failed to process request")
		return
	}

	h.NewOkResponseWithData(w, "Login successful", UserResource{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	})
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	request, err := pkg.ParseRequestBody[requests.ForgotPassword](r)
	if err != nil {
		h.NewUnprocessableEntity(w, err.Error())
		return
	}

	err = h.Validator.Struct(request)
	if err != nil {
		h.NewUnprocessableEntity(w, err.Error())
		return
	}

	_, err = h.Model.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.NewOKResponse(w, "Password reset link sent")
			return
		}

		h.Logger.Error("Failed to get user by email", slog.String("error", err.Error()))
		h.NewInternalServerError(w, "Something went wrong")
		return
	}

	token, _ := generateToken()
	tokenHash := hashToken(token)
	expires := time.Now().Add(time.Hour * 1)

	err = h.Model.InsertIntoPasswordResets(r.Context(), request.Email, tokenHash, expires)
	if err != nil {
		h.Logger.Error("Failed to insert password reset into database", slog.String("error", err.Error()))
		h.NewInternalServerError(w, "Failed to process request")
		return
	}

	err = h.Mailer.SendPasswordResetEmail(r.Context(), request.Email, token)
	if err != nil {
		h.Logger.Error("Failed to send password reset email", slog.String("error", err.Error()))
	}

	h.NewOKResponse(w, "Password reset link sent")
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	request, err := pkg.ParseRequestBody[requests.ResetPassword](r)
	if err != nil {
		h.NewUnprocessableEntity(w, err.Error())
		return
	}

	if err = h.Validator.Struct(request); err != nil {
		h.NewUnprocessableEntity(w, err.Error())
		return
	}

	storedHash, expiresAt, err := h.Model.GetPasswordReset(r.Context(), request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.NewBadRequest(w, "Invalid token or email")
			return
		}

		h.Logger.Error("Failed to get password reset from database", slog.String("error", err.Error()))
		h.NewInternalServerError(w, "Internal error")
	}

	if time.Now().After(expiresAt) {
		h.NewBadRequest(w, "Token has expired")
		return
	}

	if hashToken(request.Token) != storedHash {
		h.NewBadRequest(w, "Invalid token")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		h.Logger.Error("Failed to hash password", slog.String("error", err.Error()))
		h.NewInternalServerError(w, "Internal error")
		return
	}

	err = h.Model.UpdateUserPassword(r.Context(), request.Email, string(passwordHash))
	if err != nil {
		h.Logger.Error("Failed to update user password", slog.String("error", err.Error()))
		h.NewInternalServerError(w, "Internal error")
		return
	}

	_ = h.Model.DeletePasswordReset(r.Context(), request.Email)

	h.NewOKResponse(w, "Password reset successful")
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

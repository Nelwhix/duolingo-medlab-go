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

	"github.com/gorilla/schema"

	"github.com/Nelwhix/duolingo-medlab-go/pkg"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/request"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/resource"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	cRequest, err := pkg.ParseRequestBody[request.SignUp](r)
	if err != nil {
		response.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	err = h.Validator.Struct(cRequest)
	if err != nil {
		response.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	_, err = h.Model.GetUserByEmail(r.Context(), cRequest.Email)
	if err == nil {
		response.NewUnprocessableEntity(w, "Email already taken")
		return
	}

	userID, err := h.Model.InsertIntoUsers(r.Context(), cRequest)
	if err != nil {
		h.Logger.Error("Failed to insert user into database", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Failed to process request")
		return
	}

	user, err := h.Model.GetUserById(r.Context(), userID)
	if err != nil {
		h.Logger.Error("Failed to retrieve user from database", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Failed to process request")
		return
	}

	res := resource.UserResource{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	response.NewCreatedResponseWithData(w, "User created successfully.", res)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.NewBadRequest(w, "Failed to parse request")
		return
	}

	var cRequest request.Login
	decoder := schema.NewDecoder()
	err = decoder.Decode(&cRequest, r.PostForm)
	if err != nil {
		response.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	err = h.Validator.Struct(cRequest)
	if err != nil {
		response.NewUnprocessableEntity(w, "Failed to process request")
		return
	}

	user, err := h.Model.GetUserByEmail(r.Context(), cRequest.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NewBadRequest(w, "Email or Password is incorrect")
			return
		}

		response.NewBadRequest(w, "Failed to process request")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cRequest.Password))
	if err != nil {
		response.NewBadRequest(w, "Email or Password is incorrect")
		return
	}

	value := map[string]string{
		"user_id": user.ID,
	}

	encoded, err := h.CookieHandler.Encode("duolingo_medlab_auth", value)
	if err != nil {
		h.Logger.Error("Failed to encode cookie", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Failed to process request")
		return
	}

	cookie := &http.Cookie{
		Name:     "duolingo_medlab_auth",
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	cRequest, err := pkg.ParseRequestBody[request.ForgotPassword](r)
	if err != nil {
		response.NewUnprocessableEntity(w, err.Error())
		return
	}

	err = h.Validator.Struct(cRequest)
	if err != nil {
		response.NewUnprocessableEntity(w, err.Error())
		return
	}

	_, err = h.Model.GetUserByEmail(r.Context(), cRequest.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NewOKResponse(w, "Password reset link sent")
			return
		}

		h.Logger.Error("Failed to get user by email", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Something went wrong")
		return
	}

	token, _ := generateToken()
	tokenHash := hashToken(token)
	expires := time.Now().Add(time.Hour * 1)

	err = h.Model.InsertIntoPasswordResets(r.Context(), cRequest.Email, tokenHash, expires)
	if err != nil {
		h.Logger.Error("Failed to insert password reset into database", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Failed to process request")
		return
	}

	err = h.Mailer.SendPasswordResetEmail(r.Context(), cRequest.Email, token)
	if err != nil {
		h.Logger.Error("Failed to send password reset email", slog.String("error", err.Error()))
	}

	response.NewOKResponse(w, "Password reset link sent")
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	cRequest, err := pkg.ParseRequestBody[request.ResetPassword](r)
	if err != nil {
		response.NewUnprocessableEntity(w, err.Error())
		return
	}

	if err = h.Validator.Struct(cRequest); err != nil {
		response.NewUnprocessableEntity(w, err.Error())
		return
	}

	storedHash, expiresAt, err := h.Model.GetPasswordReset(r.Context(), cRequest.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NewBadRequest(w, "Invalid token or email")
			return
		}

		h.Logger.Error("Failed to get password reset from database", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Internal error")
	}

	if time.Now().After(expiresAt) {
		response.NewBadRequest(w, "Token has expired")
		return
	}

	if hashToken(cRequest.Token) != storedHash {
		response.NewBadRequest(w, "Invalid token")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cRequest.Password), 12)
	if err != nil {
		h.Logger.Error("Failed to hash password", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Internal error")
		return
	}

	err = h.Model.UpdateUserPassword(r.Context(), cRequest.Email, string(passwordHash))
	if err != nil {
		h.Logger.Error("Failed to update user password", slog.String("error", err.Error()))
		response.NewInternalServerError(w, "Internal error")
		return
	}

	_ = h.Model.DeletePasswordReset(r.Context(), cRequest.Email)

	response.NewOKResponse(w, "Password reset successful")
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

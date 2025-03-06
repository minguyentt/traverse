package handlers

import (
	"log/slog"
	"net/http"
	"traverse/api/errors"
	"traverse/api/json"
	"traverse/internal/services"
	"traverse/internal/storage"
	"traverse/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	RegistrationHandler(http.ResponseWriter, *http.Request)
	LoginHandler(http.ResponseWriter, *http.Request)
	ActivationHandler(http.ResponseWriter, *http.Request)
}

type authHandler struct {
	service  services.UserService
	validate *validator.Validate
}

func NewAuthHandler(s services.UserService, v *validator.Validate) *authHandler {
	return &authHandler{
		service:  s,
		validate: v,
	}
}

func (u *authHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	// Use RegistrationPayload json struct
	var userPayload models.RegistrationPayload

	// read the HTTP request
	err := json.Read(w, r, &userPayload)
	if err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	// validate the json struct
	err = u.validate.Struct(userPayload)
	if err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	// call the service for user creation
	user, err := u.service.RegisterUser(r.Context(), &userPayload)
	if err != nil {
		switch err {
		case storage.ErrDuplicateUsername:
			errors.BadRequestResponse(w, r, err)
			return
		default:
			errors.InternalServerErr(w, r, err)
		}
	}

	if err := json.Response(w, http.StatusCreated, user); err != nil {
		errors.InternalServerErr(w, r, err)
	}
}

func (u *authHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.UserLoginPayload
	if err := json.Read(w, r, &payload); err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	if err := u.validate.Struct(payload); err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

    user,err := u.service.LoginUser(r.Context(), &payload)
	if err != nil {
		errors.UnauthorizedErr(w, r, err)
		return
	}

	if err := json.Response(w, http.StatusAccepted, user); err != nil {
		errors.InternalServerErr(w, r, err)
	}
}

func (a *authHandler) ActivationHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	slog.Info("token", "out", token)

	err := a.service.ActivateUser(r.Context(), token)
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			errors.NotFoundRequest(w, r, err)
		default:
			errors.InternalServerErr(w, r, err)
		}
		return
	}

	if err := json.Response(w, http.StatusNoContent, ""); err != nil {
		errors.InternalServerErr(w, r, err)
	}
}

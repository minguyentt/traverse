package handlers

import (
	errs "errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"github.com/minguyentt/traverse/internal/db/redis/cache"
	"github.com/minguyentt/traverse/internal/services"
	"github.com/minguyentt/traverse/internal/storage"
	"github.com/minguyentt/traverse/models"
	"github.com/minguyentt/traverse/pkg/errors"
	"github.com/minguyentt/traverse/pkg/response"
	"github.com/minguyentt/traverse/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type AuthHandler interface {
	Registration(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	ActivateUser(http.ResponseWriter, *http.Request)
}

type authHandler struct {
	service  services.UserService
	validate *validator.Validate
	cache    cache.Redis
	logger   *slog.Logger
}

func NewAuthHandler(s services.UserService, v *validator.Validate, c cache.Redis, logger *slog.Logger) *authHandler {
	return &authHandler{
		service:  s,
		validate: v,
		cache:    c,
		logger: logger,
	}
}

func (u *authHandler) Registration(w http.ResponseWriter, r *http.Request) {
	// Use RegistrationPayload response struct
	var userPayload models.RegistrationPayload

	// read the HTTP request
	err := response.Read(w, r, &userPayload)
	if err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	// validate the response struct
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

	if err := response.JSON(w, http.StatusCreated, user); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

func (u *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload models.UserLoginPayload
	if err := response.Read(w, r, &payload); err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	if err := u.validate.Struct(payload); err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	res, err := u.service.LoginUser(r.Context(), &payload)
	if err != nil {
		errors.UnauthorizedErr(w, r, err)
		return
	}

	ctx := r.Context()

	key := fmt.Sprintf("activation:%s", res.Token)
	userID, err := utils.Marshal(&res.ID)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
	// cache the token as the key, val will be the userID
	err = u.cache.Set(ctx, key, userID, 24*time.Hour)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusAccepted, res.Token); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

func (a *authHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get token from query
	token := r.URL.Query().Get("token")
	if token == "" {
		errors.BadRequestResponse(w, r, fmt.Errorf("missing token parameter"))
		return
	}

	// Look up user ID using token as key
	key, err := a.cache.Get(ctx, fmt.Sprintf("activation:%s", token))
	if err != nil {
		if errs.Is(err, redis.Nil) {
			errors.NotFoundRequest(w, r, fmt.Errorf("activation token not found or expired"))
		} else {
			errors.InternalServerErr(w, r, err)
		}
		return
	}

	err = a.service.ActivateUser(ctx, token)
	if err != nil {
		if errs.Is(err, storage.ErrNotFound) {
			errors.NotFoundRequest(w, r, err)
		} else {
			errors.InternalServerErr(w, r, err)
		}
		return
	}

	if err := a.cache.Delete(ctx, string(key)); err != nil && err != cache.ErrCacheMiss {
		a.logger.Warn("failed to delete cache key", "key", key, "err", err)
	}

	// Send 204 No Content or 200
	if err := response.JSON(w, http.StatusNoContent, nil); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

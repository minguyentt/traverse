package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"traverse/internal/db/redis/cache"
	"traverse/internal/services"
	"traverse/internal/storage"
	"traverse/models"
	"traverse/pkg/errors"
	"traverse/pkg/response"
	"traverse/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	Registration(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	ActivateUser(http.ResponseWriter, *http.Request)
}

type authHandler struct {
	service  services.UserService
	validate *validator.Validate
	cache    cache.Cache
}

func NewAuthHandler(s services.UserService, v *validator.Validate, c cache.Cache) *authHandler {
	return &authHandler{
		service:  s,
		validate: v,
		cache: c,
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

	userToken, err := u.service.LoginUser(r.Context(), &payload)
	if err != nil {
		errors.UnauthorizedErr(w, r, err)
		return
	}

	ctx := r.Context()

	data, err := utils.Marshal(&userToken)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	// upon login request we will cache the user token session for 24 hrs
	userCacheKey := fmt.Sprintf("user-%d:activation", userToken.ID)
	err = u.cache.Set(ctx, userCacheKey, data, 24*time.Hour)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusAccepted, userToken); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

func (a *authHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value("user").(*models.User)
	if !ok || user == nil {
		errors.UnauthorizedErr(w, r, fmt.Errorf("invalid user token"))
		return
	}

	// get the user token from cache
	// handle misses
	var userTokenData *models.UserToken
	userid := strconv.FormatInt(user.ID, 10)
	data, err := a.cache.Get(r.Context(), userid)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	err = utils.Unmarshal(data, &userTokenData)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	err = a.service.ActivateUser(r.Context(), userTokenData.Token)
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			errors.NotFoundRequest(w, r, err)
		default:
			errors.InternalServerErr(w, r, err)
		}
		return
	}

	if err := response.JSON(w, http.StatusNoContent, ""); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

package handlers

import (
	"net/http"
	"traverse/internal/services"
	"traverse/models"
	"traverse/pkg/errors"
	json "traverse/pkg/validator"

	"github.com/go-playground/validator/v10"
)

type ContractHandler interface {
	CreateContract(http.ResponseWriter, *http.Request)
}

type contract struct {
	service  services.ContractService
	validate *validator.Validate
}

func NewContract(s services.ContractService, v *validator.Validate) *contract {
	return &contract{
		service:  s,
		validate: v,
	}
}

func (h *contract) CreateContract(w http.ResponseWriter, r *http.Request) {
	var contractPayload models.ContractPayload

	err := json.Read(w, r, &contractPayload)
	if err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	err = h.validate.Struct(contractPayload)
	if err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	usr := GetUserCtx(r)

	c, err := h.service.CreateContract(ctx, &contractPayload, usr.ID)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	if err := json.Response(w, http.StatusCreated, c); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

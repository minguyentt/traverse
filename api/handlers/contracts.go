package handlers

import (
	"net/http"
	"strconv"
	"traverse/internal/services"
	"traverse/models"
	"traverse/pkg/errors"
	"traverse/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goforj/godump"
)

type ContractHandler interface {
	CreateContract(w http.ResponseWriter, r *http.Request)
	UpdateContract(w http.ResponseWriter, r *http.Request)
	DeleteContract(w http.ResponseWriter, r *http.Request)
}

type contract struct {
	service  services.ContractService
	validate *validator.Validate
}

func NewContract(cs services.ContractService, v *validator.Validate) *contract {
	return &contract{
		service:  cs,
		validate: v,
	}
}

func (h *contract) CreateContract(w http.ResponseWriter, r *http.Request) {
	var contractPayload models.ContractPayload

	err := response.Read(w, r, &contractPayload)
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
	usr := r.Context().Value("user").(*models.User)

	c, err := h.service.CreateContract(ctx, &contractPayload, usr.ID)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusCreated, c); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

func (h *contract) UpdateContract(w http.ResponseWriter, r *http.Request) {
	// get contract from ctx
	con := r.Context().Value("contract_id").(*models.Contract)

	var pl models.UpdateContractPayload

	err := h.validate.Struct(&pl)
	// TODO: abstract this verbose validation error somewhere
	// used for validating nested structs? idk
	// testing out godump here. will remove later
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				godump.Dump("Validation Error: Field=%s, Tag=%s, Value=%v\n", err.Field(), err.Tag(), err.Value())
			}
		} else {
			// fmt.Printf("Error: %v\n", err)
			godump.Dump("Error: %v\n", err)
		}
	} else {
		// fmt.Println("User struct is valid.")
		godump.Dump("User struct is valid.")
	}

	if pl.JobTitle != nil {
		con.JobTitle = *pl.JobTitle
	}

	if pl.City != nil {
		con.City = *pl.City
	}

	if pl.Agency != nil {
		con.Agency = *pl.Agency
	}

	if pl.JobDetail.Profession != nil {
		con.JobDetails.Profession = *pl.JobDetail.Profession
	}

	if pl.JobDetail.AssignmentLength != nil {
		con.JobDetails.AssignmentLength = *pl.JobDetail.AssignmentLength
	}

	if pl.JobDetail.Experience != nil {
		con.JobDetails.Experience = *pl.JobDetail.Experience
	}

	ctx := r.Context()
	if err := h.service.UpdateContract(ctx, con); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, con); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

func (h *contract) DeleteContract(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	ctx := r.Context()
	if err := h.service.DeleteContract(ctx, id); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

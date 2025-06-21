package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"traverse/internal/ctx"
	"traverse/internal/db/redis/cache"
	"traverse/internal/services"
	"traverse/models"
	"traverse/pkg/errors"
	"traverse/pkg/response"
	"traverse/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goforj/godump"
)

type ContractHandler interface {
	Feed(w http.ResponseWriter, r *http.Request)
	ContractByID(w http.ResponseWriter, r *http.Request)
	CreateContract(w http.ResponseWriter, r *http.Request)
	UpdateContract(w http.ResponseWriter, r *http.Request)
	DeleteContract(w http.ResponseWriter, r *http.Request)
}

type contract struct {
	service  services.ContractService
	validate *validator.Validate
	cache    cache.Redis
	logger   *slog.Logger
}

func NewContract(cs services.ContractService, v *validator.Validate, c cache.Redis, logger *slog.Logger) *contract {
	return &contract{
		service:  cs,
		validate: v,
		cache:    c,
		logger:   logger,
	}
}

func (h *contract) ContractByID(w http.ResponseWriter, r *http.Request) {
	usr := ctx.GetUserFromCTX(r)

	id := utils.Int64ToStr(usr.ID)
	key := fmt.Sprintf("user-%s:contract", id)
	if data, err := h.cache.Get(r.Context(), key); err == nil {
		var contract *models.Contract
		if err := utils.Unmarshal(data, &contract); err == nil {

			if err := response.JSON(w, http.StatusOK, contract); err != nil {
				errors.InternalServerErr(w, r, err)
			}
			return
		}
	} else if err != cache.ErrCacheMiss {
		errors.InternalServerErr(w, r, err)
		return
	}

	contract, err := h.service.ContractByID(r.Context(), usr.ID)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	bytes, err := utils.Marshal(contract)
	if err == nil {
		_ = h.cache.Set(r.Context(), key, bytes, 2*time.Minute)
	} else {
		h.logger.Warn("failed to marshal data to bytes for cache", "context", key, "err", err)
	}

	if err := response.JSON(w, http.StatusOK, contract); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

func (h *contract) Feed(w http.ResponseWriter, r *http.Request) {
	usr := ctx.GetUserFromCTX(r)

	id := strconv.FormatInt(usr.ID, 10)
	cacheKeyFeed := fmt.Sprintf("user-%s:feed", id)

	data, err := h.cache.Get(r.Context(), cacheKeyFeed)
	// if Cache miss, call the db
	if err == cache.ErrCacheMiss {
		contracts, err := h.service.GetAllContracts(r.Context(), usr.ID)
		if err != nil {
			errors.InternalServerErr(w, r, err)
			return
		}

		// Cache the result
		bytes, err := utils.Marshal(contracts)
		if err == nil {
			_ = h.cache.Set(r.Context(), cacheKeyFeed, bytes, 30*time.Second)
		} else {
			h.logger.Warn("failed to marshal data to bytes for cache", "context", cacheKeyFeed, "err", err)
		}

		if err := response.JSON(w, http.StatusOK, contracts); err != nil {
			errors.InternalServerErr(w, r, err)
		}

		return
	}

	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	// Cache hit, unmarshal and return
	var contracts []*models.ContractMetaData
	if err := utils.Unmarshal(data, &contracts); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, contracts); err != nil {
		errors.InternalServerErr(w, r, err)
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

	usr := ctx.GetUserFromCTX(r)

	c, err := h.service.CreateContract(r.Context(), &contractPayload, usr.ID)
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
	con := ctx.GetContractFromCTX(r)

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

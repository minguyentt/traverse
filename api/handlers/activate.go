package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"traverse/api/errors"
	"traverse/api/json"
	"traverse/internal/services"
	"traverse/internal/storage"
)

type ActivateHandler interface {
	ActivateUser(w http.ResponseWriter, r *http.Request)
}

type activateUser struct {
	service *services.ActivateService
}

func (h *activateUser) ActivateUser(w http.ResponseWriter, r *http.Request) {
	// get token from url param
	token := chi.URLParam(r, "token")

	if err := h.service.ActivateUser(r.Context(), token); err != nil {
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

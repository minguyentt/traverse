package handlers

import (
	"net/http"
	"traverse/internal/services"
	"traverse/models"
	"traverse/pkg/errors"
	json "traverse/pkg/validator"
)

type ReviewHandler interface {
	ReviewsWithContractID(w http.ResponseWriter, r *http.Request)
}

type reviewHandler struct {
	service services.ReviewService
}

func NewReviewHandler(s services.ReviewService) *reviewHandler {
	return &reviewHandler{s}
}

func (s *reviewHandler) ReviewsWithContractID(w http.ResponseWriter, r *http.Request) {
	contractCtx := r.Context().Value("contract_id").(*models.Contract)

	reviews, err := s.service.GetReviewsByContractID(r.Context(), contractCtx.ID)
	if err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}

	if err := json.Response(w, http.StatusOK, reviews); err != nil {
		errors.InternalServerErr(w, r, err)
		return
	}
}

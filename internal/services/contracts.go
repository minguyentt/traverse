package services

import (
	"context"
	"traverse/internal/storage"
	"traverse/models"
)

type ContractService interface{
	CreateContract(ctx context.Context, cpl *models.ContractPayload, userID int64) (*models.Contract,error)
}

type contractService struct {
	store *storage.Storage
}

func NewContractService(store *storage.Storage) *contractService {
	return &contractService{store}
}

func (s *contractService) CreateContract(
	ctx context.Context,
	cpl *models.ContractPayload,
	userID int64,
) (*models.Contract,error) {
	jobDetail := &models.ContractJobDetails{
		Profession: cpl.JobDetails.Profession,
		AssignmentLength: cpl.JobDetails.AssignmentLength,
		Experience: cpl.JobDetails.Experience,
	}

	contract := &models.Contract{
		UserID:   userID,
		JobTitle: cpl.JobTitle,
		City:     cpl.City,
		Agency:   cpl.Agency,
		JobDetails: jobDetail,
	}

	if err := s.store.Contracts.Create(ctx, contract); err != nil {
		return nil, err
	}

	return contract, nil
}

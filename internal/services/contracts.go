package services

import (
	"context"
	"github.com/minguyentt/traverse/internal/storage"
	"github.com/minguyentt/traverse/models"
)

type ContractService interface {
	CreateContract(ctx context.Context, cpl *models.ContractPayload, userID int64) (*models.Contract, error)
	ContractByID(ctx context.Context, cID int64) (*models.Contract, error)
	ReviewsWithContractID(ctx context.Context, cID int64) ([]models.Review, error)
	GetAllContracts(ctx context.Context, userID int64) ([]models.ContractMetaData, error)
	UpdateContract(ctx context.Context, cpl *models.Contract) error
	DeleteContract(ctx context.Context, cID int64) error
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
) (*models.Contract, error) {
	jobDetail := &models.ContractJobDetails{
		Profession:       cpl.JobDetails.Profession,
		AssignmentLength: cpl.JobDetails.AssignmentLength,
		Experience:       cpl.JobDetails.Experience,
	}

	contract := &models.Contract{
		UserID:     userID,
		JobTitle:   cpl.JobTitle,
		City:       cpl.City,
		Agency:     cpl.Agency,
		JobDetails: jobDetail,
	}

	if err := s.store.Contracts.Create(ctx, contract); err != nil {
		return nil, err
	}

	return contract, nil
}

func (s *contractService) UpdateContract(ctx context.Context, cpl *models.Contract) error {
	if err := s.store.Contracts.Update(ctx, cpl); err != nil {
		return err
	}

	return nil
}

func (s *contractService) DeleteContract(ctx context.Context, cID int64) error {
	if err := s.store.Contracts.Delete(ctx, cID); err != nil {
		return err
	}

	return nil
}

func (s *contractService) GetAllContracts(ctx context.Context, userID int64) ([]models.ContractMetaData, error) {
	contracts, err := s.store.Contracts.All(ctx, userID)
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

func (s *contractService) ContractByID(ctx context.Context, cID int64) (*models.Contract, error) {
	contract, err := s.store.Contracts.ByID(ctx, cID)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

package service

import (
	"ocr-poc/model"
	"ocr-poc/repository"
)

// SupplierService provides business logic for suppliers.
type SupplierService struct {
	repo repository.Repository[model.Supplier]
}

// NewSupplierService creates a new supplier service.
func NewSupplierService(repo repository.Repository[model.Supplier]) *SupplierService {
	return &SupplierService{repo: repo}
}

// GetSupplierByID returns the supplier given its ID.
func (s *SupplierService) GetByID(id int64) (*model.Supplier, error) {
	supplier, err := s.repo.GetById(id)
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

// GetSupplierNameByID returns the name of the supplier given its ID.
func (s *SupplierService) GetSupplierNameByID(id int64) (string, error) {
	supplier, err := s.repo.GetById(id)
	if err != nil {
		return "", err
	}
	return supplier.Name, nil
}

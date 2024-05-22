package repository

import (
	"ocr-poc/model"
)

// SupplierRepository interface for supplier-specific methods.
type SupplierRepository interface {
	Repository[model.Supplier]
}

// InMemorySupplierRepository is an in-memory implementation of SupplierRepository.
type InMemorySupplierRepository struct {
	*BaseRepository[model.Supplier]
}

// NewInMemorySupplierRepository creates a new in-memory supplier repository.
func NewInMemorySupplierRepository() *InMemorySupplierRepository {
	repo := NewBaseRepository[model.Supplier]()
	repo.data[59581883] = &model.Supplier{ID: 1, Name: "Kamer van Koophandel®"}
	repo.data[73228923] = &model.Supplier{ID: 2, Name: "Orchestra"}
	return &InMemorySupplierRepository{repo}
}

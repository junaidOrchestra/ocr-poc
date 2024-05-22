package repository

import "errors"

// BaseRepository provides common methods for data access.
type BaseRepository[T any] struct {
	data map[int64]*T
}

// NewBaseRepository creates a new BaseRepository.
func NewBaseRepository[T any]() *BaseRepository[T] {
	return &BaseRepository[T]{data: make(map[int64]*T)}
}

// GetById fetches an entity by its ID.
func (repo *BaseRepository[T]) GetById(id int64) (*T, error) {
	entity, exists := repo.data[id]
	if !exists {
		return nil, errors.New("entity not found")
	}
	return entity, nil
}

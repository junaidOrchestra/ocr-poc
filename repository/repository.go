package repository

// Repository interface defines common methods for data access.
type Repository[T any] interface {
	GetById(id int64) (*T, error)
}

package user

import "github.com/google/uuid"

type Repository interface {
	Create(User) (User, error)

	GetByID(uuid.UUID) (User, error)

	Update(User) (User, error)

	Delete(uuid.UUID) error
}

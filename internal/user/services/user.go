package services

import (
	"context"
	"fmt"

	"go-campaign.com/internal/user/repository/sqlc"
	"go-campaign.com/pkg/hash"
)

type UserService struct {
	q *sqlc.Queries
}

func NewUserService(q *sqlc.Queries) *UserService {
	return &UserService{
		q: q,
	}
}
func (s *UserService) CreateUser(c context.Context, user CreateUserDTO) (int64, error) {
	hashPassword, err := hash.Password(user.Password)

	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	u, err := s.q.CreateUser(c, sqlc.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashPassword,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return int64(u.ID), nil
}

func (s *UserService) CheckLoginUser(c context.Context, email, password string) (int, error) {
	user, err := s.FindUserByEmail(c, email)
	if err != nil || user == nil {
		return 0, fmt.Errorf("failed to find user by email: %w", err)
	}

	isMatch, err := hash.ComparePassword(user.Password, password)

	if err != nil {
		return 0, fmt.Errorf("failed to compare password: %w", err)
	}

	if !isMatch {
		return 0, fmt.Errorf("password does not match")
	}

	userID := int(user.ID)

	return userID, nil

}

func (s *UserService) FindUserByEmail(c context.Context, email string) (*UserDTO, error) {
	user, err := s.q.GetUserByEmail(c, email)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:       int64(user.ID),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

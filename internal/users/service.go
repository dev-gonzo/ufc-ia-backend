package users

import (
	"ufc-backend/internal/auth"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	input CreateUserInput,
) (*User, error) {

	if !isValidEmail(input.Email) {
		return nil, ErrInvalidEmail
	}

	if !isValidUsername(input.Username) {
		return nil, ErrInvalidUsername
	}

	if !isValidPassword(input.Password) {
		return nil, ErrInvalidPassword
	}

	hash, err := auth.HashPassword(
		input.Password,
	)

	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: hash,
		Role:         auth.RoleUser,
	}

	err = s.repository.Create(user)

	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return user, nil
}

func (s *Service) ChangePassword(
	userID string,
	input ChangePasswordInput,
	user *User,
) error {

	if !isValidPassword(
		input.NewPassword,
	) {
		return ErrInvalidPassword
	}

	valid := auth.CheckPassword(
		input.CurrentPassword,
		user.PasswordHash,
	)

	if !valid {
		return auth.ErrInvalidCredentials
	}

	hash, err := auth.HashPassword(
		input.NewPassword,
	)

	if err != nil {
		return err
	}

	return s.repository.UpdatePassword(
		userID,
		hash,
	)
}

func (s *Service) ChangeRole(
	input ChangeRoleInput,
) error {

	if !isValidRole(input.Role) {
		return ErrInvalidRole
	}

	return s.repository.UpdateRole(
		input.UserID,
		input.Role,
	)
}

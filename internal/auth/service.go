package auth

type UserRepository interface {
	FindByEmail(email string) (*User, error)
}

type Service struct {
	usersRepository UserRepository
}

func NewService(
	usersRepository UserRepository,
) *Service {
	return &Service{
		usersRepository: usersRepository,
	}
}

func (s *Service) Login(
	email string,
	password string,
) (string, error) {

	user, err := s.usersRepository.FindByEmail(email)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	valid := CheckPassword(
		password,
		user.PasswordHash,
	)

	if !valid {
		return "", ErrInvalidCredentials
	}

	token, err := GenerateAccessToken(User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

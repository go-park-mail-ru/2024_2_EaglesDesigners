package service

import (
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/repository"
)

type AuthService struct {
	userRepo     repository.UserRepository
	tokenService TokenService
}

func NewAuthService(userRepo repository.UserRepository, tokenService TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) Authenticate(username, password string) bool {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return false
	}
	return doPasswordsMatch(user.Password, password)
}

func (s *AuthService) Registation(username, name, password string) error {
	hashed := hashPassword(password)
	err := s.userRepo.CreateUser(username, name, hashed)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetUserDataByUsername(username string) (UserData, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return UserData{}, err
	}

	userData := UserData{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}

	return userData, nil
}

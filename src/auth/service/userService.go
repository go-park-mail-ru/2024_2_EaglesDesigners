package service

import (
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/utils"
)

type AuthService struct {
	userRepo     auth.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(userRepo auth.UserRepository, tokenService auth.TokenService) *AuthService {
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
	return utils.DoPasswordsMatch(user.Password, password)
}

func (s *AuthService) Registration(username, name, password string) error {
	hashed := utils.HashPassword(password)
	err := s.userRepo.CreateUser(username, name, hashed)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetUserDataByUsername(username string) (utils.UserData, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return utils.UserData{}, err
	}

	userData := utils.UserData{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}

	return userData, nil
}

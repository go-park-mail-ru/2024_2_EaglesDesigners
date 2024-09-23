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
	user := s.userRepo.GetUserByUsername(username)
	if user.Username == "" {
		return false
	}
	return doPasswordsMatch(user.Password, password, user.Salt)
}

func (s *AuthService) Registation(username, password string) error {
	salt := generateRandomSalt(saltSize)
	hashed := hashPassword(password, salt)
	err := s.userRepo.CreateUser(username, hashed, salt)
	if err != nil {
		return err
	}

	return nil
}

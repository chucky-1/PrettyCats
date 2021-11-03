package service

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/repository"
	"crypto/sha1"
	"fmt"
)

const salt = "ce15ce51ce1ce5ce51c"

type UserAuthService struct {
	repository repository.Auth
}

type Auth interface {
	CreateUserServ(user models.User) (int, error)
}

func NewUserAuthService(r repository.Auth) *UserAuthService {
	return &UserAuthService{repository: r}
}

func (s *UserAuthService) CreateUserServ(user models.User) (int, error) {
	fmt.Println("generic password start...")
	user.Password = generatePassword(user.Password)
	fmt.Println("Generic password successful")
	return s.repository.CreateUser(user)
}

func generatePassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

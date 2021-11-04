package service

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"time"
)

type UserAuthService struct {
	repository repository.Auth
}

type Auth interface {
	CreateUserServ(user models.User) (int, error)
	GenerateToken(username string, password string) (t string, err error)
}

func NewUserAuthService(r repository.Auth) *UserAuthService {
	return &UserAuthService{repository: r}
}

type JwtCustomClaims struct {
	ID 	  int `json:"id"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

func (s *UserAuthService) CreateUserServ(user models.User) (int, error) {
	user.Password = generatePassword(user.Password)
	return s.repository.CreateUser(user)
}

func (s *UserAuthService) GenerateToken(username string, password string) (t string, err error) {
	user, err := s.repository.GetUser(username, generatePassword(password))
	if err != nil {
		return "", errors.New("error in repository")
	}

	claims := &JwtCustomClaims{
		ID: user.ID,
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err = token.SignedString([]byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")))
	if err != nil {
		return "", errors.New("error during generate token")
	}

	return t, nil
}

func generatePassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(viper.GetString("SALT_FOR_GENERATE_PASSWORD"))))
}

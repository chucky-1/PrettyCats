package service

import (
	"CatsCrud/internal/models"
	rep "CatsCrud/internal/repository"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"fmt"
	"time"
)

const tokenAvailableHour = 72

// UserAuth has an interface of Auth of repository
type UserAuth struct {
	repository rep.Auth
}

// Auth has methods for registration and authorization
type Auth interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (t string, err error)
}

// NewUserAuth is a constructor
func NewUserAuth(r rep.Auth) *UserAuth {
	return &UserAuth{repository: r}
}

// JwtCustomClaims expands the jwt.StandardClaims
type JwtCustomClaims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

// CreateUser sends user into repository and return user's id
func (s *UserAuth) CreateUser(user models.User) (int, error) {
	user.Password = generatePassword(user.Password)
	return s.repository.CreateUser(user)
}

// GenerateToken creates token for authorization
func (s *UserAuth) GenerateToken(username, password string) (t string, err error) {
	user, err := s.repository.GetUser(username, generatePassword(password))
	if err != nil {
		log.Error("error in repository")
		return "", err
	}

	claims := &JwtCustomClaims{ID: user.ID, Name: user.Name, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * tokenAvailableHour).Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err = token.SignedString([]byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")))
	if err != nil {
		log.Error("error during generate token")
		return "", err
	}

	return t, nil
}

func generatePassword(password string) string {
	return fmt.Sprintf(password + viper.GetString("SALT_FOR_GENERATE_PASSWORD"))
}

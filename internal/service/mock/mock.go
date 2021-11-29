// Package mock is mock for service
package mock

import (
	"CatsCrud/internal/models"
)

// Service MockService is mock for service
type Service struct{}

// NewService is constructor
func NewService() *Service {
	return &Service{}
}

// GetAll is a method of mock for service
func (m *Service) GetAll() ([]*models.Cat, error) {
	cat := models.Cat{
		ID:   0,
		Name: "",
	}
	allcats := []*models.Cat{&cat}
	return allcats, nil
}

// Create is a method of mock for service
func (m *Service) Create(cats models.Cat) (*models.Cat, error) {
	return &cats, nil
}

// Get is a method of mock for service
func (m *Service) Get(id string) (*models.Cat, error) {
	cat := models.Cat{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

// Update is a method of mock for service
func (m *Service) Update(id string, cats models.Cat) (*models.Cat, error) {
	return &cats, nil
}

// Delete is a method of mock for service
func (m *Service) Delete(id string) (*models.Cat, error) {
	cat := models.Cat{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

// CreateUser is a method of mock for service
func (m *Service) CreateUser(user models.User) (int, error) {
	user.ID = 1
	return user.ID, nil
}

// GenerateToken is a method of mock for service
func (m *Service) GenerateToken(username, password string) (t string, err error) {
	// Рабочий токен
	t = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6IkpvbiBTbm93IiwiZXhwIjoxOTUxNzQ5NjE5fQ.qdAUCQt2nAdKxgqTVVieqn0gF-yiIKtOevOCSHN7DvU"
	return t, nil
}

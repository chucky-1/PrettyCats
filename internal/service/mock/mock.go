// Package mock is mock for service
package mock

import (
	"CatsCrud/internal/models"
)

// CatServ is mock for service
type CatServ struct{}

// NewMockCatServ is constructor
func NewMockCatServ() *CatServ {
	return &CatServ{}
}

// GetAllCatsServ is a method of mock for service
func (m *CatServ) GetAllCatsServ() ([]*models.Cats, error) {
	cat := models.Cats{
		ID:   0,
		Name: "",
	}
	allcats := []*models.Cats{&cat}
	return allcats, nil
}

// CreateCatsServ is a method of mock for service
func (m *CatServ) CreateCatsServ(cats models.Cats) (*models.Cats, error) {
	return &cats, nil
}

// GetCatServ is a method of mock for service
func (m *CatServ) GetCatServ(id string) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

// UpdateCatServ is a method of mock for service
func (m *CatServ) UpdateCatServ(id string, cats models.Cats) (*models.Cats, error) {
	return &cats, nil
}

// DeleteCatServ is a method of mock for service
func (m *CatServ) DeleteCatServ(id string) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

// CreateUserServ is a method of mock for service
func (m *CatServ) CreateUserServ(user models.User) (int, error) {
	user.ID = 1
	return user.ID, nil
}

// GenerateToken is a method of mock for service
func (m *CatServ) GenerateToken(username, password string) (t string, err error) {
	// Рабочий токен
	t = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6IkpvbiBTbm93IiwiZXhwIjoxOTUxNzQ5NjE5fQ.qdAUCQt2nAdKxgqTVVieqn0gF-yiIKtOevOCSHN7DvU"
	return t, nil
}

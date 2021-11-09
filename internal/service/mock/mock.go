package mock

import (
	"CatsCrud/internal/models"
)

type MockCatServ struct {

}

func NewMockCatServ() *MockCatServ {
	return &MockCatServ{}
}

func (m *MockCatServ) GetAllCatsServ() ([]*models.Cats, error) {
	cat := models.Cats{
		ID:   0,
		Name: "",
	}
	allcats := []*models.Cats{&cat}
	return allcats, nil
}

func (m *MockCatServ) CreateCatsServ(cats models.Cats) (*models.Cats, error) {
	return &cats, nil
}

func (m *MockCatServ) GetCatServ(id string) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

func (m *MockCatServ) UpdateCatServ(id string, cats models.Cats) (*models.Cats, error) {
	return &cats, nil
}

func (m *MockCatServ) DeleteCatServ(id string) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

func (m *MockCatServ) CreateUserServ(user models.User) (int, error) {
	user.ID = 1
	return user.ID, nil
}

func (m *MockCatServ) GenerateToken(username string, password string) (t string, err error) {
	// Рабочий токен
	t = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6IkpvbiBTbm93IiwiZXhwIjoxOTUxNzQ5NjE5fQ.qdAUCQt2nAdKxgqTVVieqn0gF-yiIKtOevOCSHN7DvU"
	return t, nil
}

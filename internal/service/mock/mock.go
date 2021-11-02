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
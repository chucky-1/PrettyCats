package mymock

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

func (m *MockCatServ) CreateCatsServ(jsonMap map[string]interface{}) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

func (m *MockCatServ) GetCatServ(id string) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

func (m *MockCatServ) UpdateCatServ(id string, jsonMap map[string]interface{}) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}

func (m *MockCatServ) DeleteCatServ(id string) (*models.Cats, error) {
	cat := models.Cats{
		ID:   1,
		Name: "Jon Snow",
	}
	return &cat, nil
}
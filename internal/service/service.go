package service

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/repository"
)

type CatService struct {
	repository repository.Repository
}

type Service interface {
	GetAllCatsServ() ([]*models.Cats, error)
	CreateCatsServ(cats models.Cats) (*models.Cats, error)
	GetCatServ(id string) (*models.Cats, error)
	UpdateCatServ(id string, cats models.Cats) (*models.Cats, error)
	DeleteCatServ(id string) (*models.Cats, error)
}

func NewCatService(rps repository.Repository) *CatService {
	return &CatService{repository: rps}
}

func (s *CatService) GetAllCatsServ() ([]*models.Cats, error) {
	return s.repository.GetAllCats()
}

func (s *CatService) CreateCatsServ(cats models.Cats) (*models.Cats, error) {
	return s.repository.CreateCats(cats)
}

func (s *CatService) GetCatServ(id string) (*models.Cats, error) {
	return s.repository.GetCat(id)
}

func (s *CatService) UpdateCatServ(id string, cats models.Cats) (*models.Cats, error) {
	return s.repository.UpdateCat(id, cats)
}

func (s *CatService) DeleteCatServ(id string) (*models.Cats, error) {
	return s.repository.DeleteCat(id)
}

// Package service contains all business logic. He gets params from handler and sends it in repository
package service

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/repository"
	log "github.com/sirupsen/logrus"
)

// CatService has an interface of repository
type CatService struct {
	repository repository.Repository
	hash repository.Redis
}

// Service has methods which get params from handler and send it in repository
type Service interface {
	GetAllCatsServ() ([]*models.Cats, error)
	CreateCatsServ(cats models.Cats) (*models.Cats, error)
	GetCatServ(id string) (*models.Cats, error)
	UpdateCatServ(id string, cats models.Cats) (*models.Cats, error)
	DeleteCatServ(id string) (*models.Cats, error)
}

// NewCatService is constructor
func NewCatService(rps repository.Repository, hash repository.Redis) *CatService {
	return &CatService{repository: rps, hash: hash}
}

// GetAllCatsServ is called by handler and calls func in repository
func (s *CatService) GetAllCatsServ() ([]*models.Cats, error) {
	return s.repository.GetAllCats()
}

// CreateCatsServ is called by handler and calls func in repository
func (s *CatService) CreateCatsServ(cats models.Cats) (*models.Cats, error) {
	err := s.hash.CreateCat(cats)
	if err != nil {
		log.Error(err)
	}

	return s.repository.CreateCats(cats)
}

// GetCatServ is called by handler and calls func in repository
func (s *CatService) GetCatServ(id string) (*models.Cats, error) {
	cat, err := s.hash.GetCat(id)

	if err != nil {
		cat, err = s.repository.GetCat(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = s.hash.CreateCat(*cat)
		if err != nil {
			log.Error(err)
		}
		return cat, nil
	}

	return cat, nil
}

// UpdateCatServ is called by handler and calls func in repository
func (s *CatService) UpdateCatServ(id string, cats models.Cats) (*models.Cats, error) {
	err := s.hash.DeleteCat(id)
	if err != nil {
		log.Error(err)
	}
	return s.repository.UpdateCat(id, cats)
}

// DeleteCatServ is called by handler and calls func in repository
func (s *CatService) DeleteCatServ(id string) (*models.Cats, error) {
	err := s.hash.DeleteCat(id)
	if err != nil {
		log.Error(err)
	}
	return s.repository.DeleteCat(id)
}

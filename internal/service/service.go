// Package service contains all business logic. He gets params from handler and sends it in repository
package service

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/repository"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// CatService has an interface of repository
type CatService struct {
	repository repository.Repository
	cache repository.Cache
	stream RedisStream  // it has a Redis client and local map
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
func NewCatService(rps repository.Repository, cache repository.Cache, stream RedisStream) *CatService {
	return &CatService{repository: rps, cache: cache, stream: stream}
}

// GetAllCatsServ is called by handler and calls func in repository
func (s *CatService) GetAllCatsServ() ([]*models.Cats, error) {
	return s.repository.GetAllCats()
}

// CreateCatsServ is called by handler and calls func in repository
func (s *CatService) CreateCatsServ(cats models.Cats) (*models.Cats, error) {
	err := s.stream.WriteDown("INSERT", cats)
	if err != nil {
		return nil, err
	}

	err = s.cache.CreateCat(cats)
	if err != nil {
		log.Error(err)
	}

	return s.repository.CreateCats(cats)
}

func GetCatFromMap(memory map[int32]string ,id string) (*models.Cats, error) {
	idInt, err := strconv.ParseInt(id, 0, 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	name, ok := memory[int32(idInt)]
	if ok == true {
		cat := new(models.Cats)
		cat.ID = int32(idInt)
		cat.Name = name
		if err != nil {
			log.Error(err)
			return nil, err
		}
		log.Println("Get сat from local map")
		return cat, nil
	}
	return nil, fmt.Errorf("cat didn't find")
}

// GetCatServ is called by handler and calls func in repository
func (s *CatService) GetCatServ(id string) (*models.Cats, error) {
	// Get cat from local map
	cat, err := GetCatFromMap(s.stream.memory, id)
	if err == nil {
		return cat, nil
	}

	// Get cat from cache
	cat, err = s.cache.GetCat(id)

	// Get cat from database
	if err != nil {
		cat, err = s.repository.GetCat(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = s.stream.WriteDown("INSERT", *cat)
		if err != nil {
			log.Error(err)
		}
		err = s.cache.CreateCat(*cat)
		if err != nil {
			log.Error(err)
		}
		return cat, nil
	}

	// Return cat from cache but before add in local map
	err = s.stream.WriteDown("INSERT", *cat)
	if err != nil {
		log.Error(err)
	}
	return cat, nil
}

// UpdateCatServ is called by handler and calls func in repository
func (s *CatService) UpdateCatServ(id string, cats models.Cats) (*models.Cats, error) {
	// Delete from local map
	idInt, err := strconv.ParseInt(id, 0, 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s.stream.UpdateMemory("DELETE", int32(idInt), "")

	// Delete from cache
	err = s.cache.DeleteCat(id)
	if err != nil {
		log.Error(err)
	}

	return s.repository.UpdateCat(id, cats)
}

// DeleteCatServ is called by handler and calls func in repository
func (s *CatService) DeleteCatServ(id string) (*models.Cats, error) {
	// Delete from local map
	idInt, err := strconv.ParseInt(id, 0, 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s.stream.UpdateMemory("DELETE", int32(idInt), "")

	// Delete from cache
	err = s.cache.DeleteCat(id)
	if err != nil {
		log.Error(err)
	}
	return s.repository.DeleteCat(id)
}

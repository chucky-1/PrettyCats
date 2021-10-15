package service

import (
	"CatsCrud/models"
	"CatsCrud/repository"
)

type CatService struct {
	repository *repository.Repository
}

func NewCatService(rps *repository.Repository) *CatService {
	return &CatService{repository: rps}
}

func (s *CatService) GetAllCatsServ() ([]*models.Cats, error) {
	return s.repository.GetAllCats()
}

//func (s *CatService) CreateCatsServ(jsonMap map[string]interface{}) (*models.Cats, error) {
//	rep := CatService{}
//	return rep.repository.CreateCats(jsonMap)
//}
//
//func (s *CatService) GetCatServ(id string) (*models.Cats, error) {
//	rep := CatService{}
//	return rep.repository.GetCat(id)
//}
//
//func (s *CatService) UpdateCatServ(id string, jsonMap map[string]interface{}) (*models.Cats, error) {
//	rep := CatService{}
//	return rep.repository.UpdateCat(id, jsonMap)
//}
//
//func (s *CatService) DeleteCatServ(id string) (*models.Cats, error) {
//	rep := CatService{}
//	return rep.repository.DeleteCat(id)
//}

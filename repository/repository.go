package repository

import (
	"CatsCrud/models"
	"context"
	"github.com/labstack/gommon/log"
)

type Repository struct {
	cats *models.Cats
}

func NewRepository(cats models.Cats) *Repository {
	return &Repository{cats: &cats}
}

func (cm *Repository) GetAllCats() ([]*models.Cats, error) {
	conn := DbCall()

	rows, err := conn.Query(context.Background(), "select ID, name from cats")
	if err != nil {
		log.Fatal(err)
	}

	newRep := NewRepository(models.Cats{})
	cat := newRep.cats

	//cat := cm.cats

	var allcats = []*models.Cats{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}
		cat.ID = values[0].(int32)
		cat.Name = values[1].(string)
		allcats = append(allcats, cat)
	}

	return allcats, nil
}


//func GetAllCatsRepppp(c echo.Context) (map[int32]*models.Сats, error) {
//	conn := DbCall()
//
//	rows, err := conn.Query(c.Request().Context(), "select ID, name from cats")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var allcats = map[int32]*models.Сats{}
//
//	for rows.Next() {
//		values, err := rows.Values()
//		if err != nil {
//			log.Fatal(err)
//		}
//		id := values[0].(int32)
//		name := values[1].(string)
//
//		ct := &models.Сats{
//			ID:   id,
//			Name: name,
//		}
//		allcats[ct.ID] = ct
//	}
//
//	return allcats, nil
//}
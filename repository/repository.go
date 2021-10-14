package repository

import (
	"CatsCrud/models"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"strconv"
)

type Repository struct {
	cats *models.Cats
}

//func NewRepository(cats models.Cats) *Repository {
//	return &Repository{cats: &cats}
//}

func (c *Repository) GetAllCats() ([]*models.Cats, error) {
	conn := RequestDB()
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select ID, name from cats")
	if err != nil {
		log.Fatal(err)
	}

	var allcats = []*models.Cats{}

	for rows.Next() {

		cats := models.Cats{
			ID:   0,
			Name: "",
		}
		cat := Repository{cats: &cats}

		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}
		cat.cats.ID = values[0].(int32)
		cat.cats.Name = values[1].(string)
		allcats = append(allcats, cat.cats)
	}

	return allcats, nil
}

func (c *Repository) CreateCats(jsonMap map[string]interface{}) (*models.Cats, error) {
	conn := RequestDB()
	defer conn.Close(context.Background())

	// Инициализация структуры models.Cats
	cat := models.Cats{
		ID:   0,
		Name: "",
	}

	// Достаём id, name. Id преобразуем в int
	var id interface{}
	id = jsonMap["id"]
	idInt, err := strconv.Atoi(id.(string))
	if err != nil {
		return &cat, err
	}

	// Присваиваем значения в структуру models.Cats
	cat.ID = int32(idInt)
	name := jsonMap["name"]
	cat.Name = fmt.Sprintf("%v", name)

	// Добавляем в базу данных
	commandTag, err := conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", idInt, name)
	if err != nil {
		return &cat, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cat, errors.New("Failed to create cat")
	}
	return &cat, nil
}

func (c *Repository) GetCat(id string) (*models.Cats, error) {
	conn := RequestDB()
	defer conn.Close(context.Background())

	// Инициализация структуры models.Cats
	cat := models.Cats{
		ID:   0,
		Name: "",
	}

	// Достаём name
	var name string
	err := conn.QueryRow(context.Background(), "select name from cats where id=$1", id).Scan(&name)
	if err != nil {
		return &cat, err
	}

	//Присваиваем параметры models.Cats
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, nil
	}
	cat.ID = int32(idInt)
	cat.Name = name

	return &cat, nil
}

func (c *Repository) UpdateCat(id string, jsonMap map[string]interface{}) (*models.Cats, error) {
	conn := RequestDB()
	defer conn.Close(context.Background())

	// Инициализация структуры models.Cats
	cat := models.Cats{
		ID:   0,
		Name: "",
	}

	// Преобразуем id в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, err
	}

	// Достаём name
	name := jsonMap["name"]

	// Обновляем models.Cat
	cat.ID = int32(idInt)
	cat.Name = fmt.Sprintf("%v", name)

	// Вносим изменения в базу данных
	commandTag, err := conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", name, idInt)
	if err != nil {
		return &cat, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cat, errors.New("Row isn't update")
	}

	return &cat, nil
}

func (c *Repository) DeleteCat(id string) (*models.Cats, error) {
	conn := RequestDB()
	defer conn.Close(context.Background())

	// Инициализация структуры models.Cats
	cat := models.Cats{
		ID:   0,
		Name: "",
	}

	// Преобразуем id в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, err
	}

	// Обновляес models.Cats
	cat.ID = int32(idInt)

	// Удаление из базы
	commandTag, err := conn.Exec(context.Background(), "delete from cats where id=$1", id)
	if err != nil {
		return &cat, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cat, errors.New("No row found to delete")
	}

	return &cat, nil
}
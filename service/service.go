package service

import (
	"CatsCrud/models"
	"CatsCrud/repository"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"strconv"
)

type CatService struct {
	repository *repository.Repository
}

func (s *CatService) GetAllCatsServ() ([]*models.Cats, error) {
	return s.repository.GetAllCats()
}

func CreateCatsServ(c echo.Context, conn *pgx.Conn) (map[string]interface{}, error) {
	// Получаем параметры JSON в json_map
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return json_map, err
	}

	// Достаём id, name. Id преобразуем в int
	var id interface{}
	id = json_map["id"]
	idInt, err := strconv.Atoi(id.(string))
	if err != nil {
		return json_map, err
	}
	name := json_map["name"]

	// Добавляем в базу данных
	commandTag, err := conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", idInt, name)
	if err != nil {
		return json_map, err
	}
	if commandTag.RowsAffected() != 1 {
		return json_map, errors.New("Failed to create cat")
	}

	return json_map, nil
}

func GetCatServ(c echo.Context, conn *pgx.Conn) (string, string, error) {
	id := c.Param("id")
	var name string

	err := conn.QueryRow(context.Background(), "select name from cats where id=$1", id).Scan(&name)
	if err != nil {
		return "", "", err
	}

	return id, name, nil
}

func UpdateUserServ(c echo.Context, conn *pgx.Conn) (map[string]interface{}, error) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return make(map[string]interface{}) ,err
	}

	// Получаем параметры JSON в json_map
	json_map := make(map[string]interface{})
	err2 := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err2 != nil {
		return json_map, err2
	}

	// Достаём name
	name := json_map["name"]

	// Вносим изменения в базу данных
	commandTag, err := conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", name, idInt)
	if err != nil {
		return json_map, err
	}
	if commandTag.RowsAffected() != 1 {
		return json_map, errors.New("Row isn't update")
	}

	return json_map, nil
}

func DeleteCatServ(c echo.Context, conn *pgx.Conn) (string, error) {
	id := c.Param("id")

	commandTag, err := conn.Exec(context.Background(), "delete from cats where id=$1", id)
	if err != nil {
		return "", err
	}
	if commandTag.RowsAffected() != 1 {
		return "", errors.New("No row found to delete")
	}
	return id, nil
}

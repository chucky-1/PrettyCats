package main

import (
	"CatsCrud/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

func main() {
	e := echo.New()

	//Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})
	e.GET("/cats", getAllCats)
	e.POST("/cats", createCats)
	e.GET("/cats/:id", getCat)
	e.PUT("/cats/:id", updateUser)
	e.DELETE("/cats/:id", deleteCat)

	// Временное решение проблемы. Вызываю RequestDB, что бы получить port из config
	//conn := repository.RequestDB()
	//conn.Close(context.Background())
	e.Logger.Fatal(e.Start(":8000"))
}

func getAllCats(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select ID, name from cats")
	if err != nil {
		log.Fatal(err)
	}

	type cat struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	}
	var allcats = map[int32]*cat{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}
		id := values[0].(int32)
		name := values[1].(string)

		ct := &cat{
			ID:   id,
			Name: name,
		}
		allcats[ct.ID] = ct
	}
	return c.JSON(http.StatusOK, allcats)
}

func createCats(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	// Получаем параметры JSON в json_map
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	}

	// Достаём id, name. Id преобразуем в int
	var id interface{}
	id = json_map["id"]
	idInt, err := strconv.Atoi(id.(string))
	if err != nil {
		return err
	}
	name := json_map["name"]
	fmt.Println(name)

	// Добавляем в базу данных
	commandTag, err := conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", idInt, name)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("Failed to create cat")
	}

	return c.JSON(http.StatusCreated, json_map)
}

func getCat(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	id := c.Param("id")
	var name string

	err := conn.QueryRow(context.Background(), "select name from cats where id=$1", id).Scan(&name)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("cat's name is %s, id: %s", name, id))
}

func updateUser(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	// Получаем параметры JSON в json_map
	json_map := make(map[string]interface{})
	err2 := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err2 != nil {
		return err2
	}

	// Достаём name
	name := json_map["name"]

	// Вносим изменения в базу данных
	commandTag, err := conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", name, idInt)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("Row isn't update")
	}
	return c.JSON(http.StatusOK, json_map)
}

func deleteCat(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	id := c.Param("id")

	commandTag, err := conn.Exec(context.Background(), "delete from cats where id=$1", id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return c.String(http.StatusOK, fmt.Sprintf("Cats %s delete", id))
}

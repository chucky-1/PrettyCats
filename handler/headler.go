package handler

import (
	"CatsCrud/repository"
	"CatsCrud/service"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func GetAllCats(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select ID, name from cats")
	if err != nil {
		log.Fatal(err)
	}

	allcats := service.GetAllCatsServ(rows)

	return c.JSON(http.StatusOK, allcats)
}

func CreateCats(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	json_map, err := service.CreateCatsServ(c, conn)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, json_map)
}

func GetCat(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	id, name, err := service.GetCatServ(c, conn)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("cat's name is %s, id: %s", name, id))
}

func UpdateUser(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	json_map, err := service.UpdateUserServ(c, conn)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, json_map)
}

func DeleteCat(c echo.Context) error {
	conn := repository.RequestDB()
	defer conn.Close(context.Background())

	id, err := service.DeleteCatServ(c, conn)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("Cats %s delete", id))
}


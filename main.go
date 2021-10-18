package main

import (
	"CatsCrud/handler"
	"CatsCrud/repository"
	"CatsCrud/service"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})

	// Соединение с postgres
	conn := repository.RequestDB()
	defer conn.Close()

	// Соединение с mongo
	client := repository.RequestMongo()
	defer client.Disconnect(context.TODO())

	rps := repository.NewRepository(conn, client)
	srv := service.NewCatService(rps)
	hndlr := handler.NewCatHandler(srv)
	e.GET("/cats", hndlr.GetAllCats)
	e.POST("/cats", hndlr.CreateCats)
	e.GET("/cats/:id", hndlr.GetCat)
	e.PUT("/cats/:id", hndlr.UpdateCat)
	e.DELETE("/cats/:id", hndlr.DeleteCat)

	e.Logger.Fatal(e.Start(":8000"))
}
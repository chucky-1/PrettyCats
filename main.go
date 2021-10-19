package main

import (
	"CatsCrud/handler"
	"CatsCrud/repository"
	"CatsCrud/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

// 1 - postgres, 2 - mongo
const flag = 2

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})

	var rps repository.Repository

	if flag == 1 {
		// Соединение с postgres
		conn := repository.RequestDB()
		defer conn.Close()

		rps = repository.NewPostgresRepository(conn)
	} else if flag == 2 {
		// Соединение с mongo
		client, cancel := repository.RequestMongo()
		defer cancel()

		rps = repository.NewMongoRepository(client)
	}

	srv := service.NewCatService(rps)
	hndlr := handler.NewCatHandler(srv)
	e.GET("/cats", hndlr.GetAllCats)
	e.POST("/cats", hndlr.CreateCats)
	e.GET("/cats/:id", hndlr.GetCat)
	e.PUT("/cats/:id", hndlr.UpdateCat)
	e.DELETE("/cats/:id", hndlr.DeleteCat)

	e.Logger.Fatal(e.Start(":8000"))
}
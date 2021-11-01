package main

import (
	"CatsCrud/internal/handler"
	repository2 "CatsCrud/internal/repository"
	"CatsCrud/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

// 1 - postgres, 2 - mongo
const flag = 1

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})

	var rps repository2.Repository
	if flag == 1 {
		// Соединение с postgres
		conn := repository2.RequestDB()
		defer conn.Close()

		rps = repository2.NewPostgresRepository(conn)
	} else if flag == 2 {
		// Соединение с mongo
		client, cancel := repository2.RequestMongo()
		defer cancel()

		rps = repository2.NewMongoRepository(client)
	}

	var srv service.Service
	srv = service.NewCatService(rps)
	hndlr := handler.NewCatHandler(srv)
	e.GET("/cats", hndlr.GetAllCats)
	e.POST("/cats", hndlr.CreateCats)
	e.GET("/cats/:id", hndlr.GetCat)
	e.PUT("/cats/:id", hndlr.UpdateCat)
	e.DELETE("/cats/:id", hndlr.DeleteCat)

	e.Logger.Fatal(e.Start(":8000"))
}
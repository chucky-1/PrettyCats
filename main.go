package main

import (
	"CatsCrud/handler"
	"CatsCrud/repository"
	"CatsCrud/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	//Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})
	conn := repository.RequestDB()
	rps := repository.NewRepository(conn)
	srv := service.NewCatService(rps)
	hndlr := handler.NewCatHandler(srv)
	e.GET("/cats", hndlr.GetAllCats)
	//e.POST("/cats", hndlr.CreateCats)
	//e.GET("/cats/:id", hndlr.GetCat)
	//e.PUT("/cats/:id", hndlr.UpdateCat)
	//e.DELETE("/cats/:id", hndlr.DeleteCat)

	// Временное решение проблемы. Вызываю RequestDB, что бы получить port из config
	//conn := repository.RequestDB()
	//conn.Close(context.Background())
	e.Logger.Fatal(e.Start(":8000"))
}
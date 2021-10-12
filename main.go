package main

import (
	"CatsCrud/handler"

	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	//Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})
	e.GET("/cats", handler.GetAllCats)
	e.POST("/cats", handler.CreateCats)
	e.GET("/cats/:id", handler.GetCat)
	e.PUT("/cats/:id", handler.UpdateUser)
	e.DELETE("/cats/:id", handler.DeleteCat)

	// Временное решение проблемы. Вызываю RequestDB, что бы получить port из config
	//conn := repository.RequestDB()
	//conn.Close(context.Background())
	e.Logger.Fatal(e.Start(":8000"))
}
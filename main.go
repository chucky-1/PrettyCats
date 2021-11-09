package main

import (
	"CatsCrud/internal/handler"
	repository2 "CatsCrud/internal/repository"
	"CatsCrud/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"net/http"
)

// 1 - postgres, 2 - mongo
const flag = 1

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		 c.Set("Dima", "27")
		return c.String(http.StatusOK, "Hello, Cats!")
	})

	var rps repository2.Repository
	var rpsAuth repository2.Auth
	if flag == 1 {
		// Соединение с postgres
		conn := repository2.RequestDB()
		defer conn.Close()

		rps = repository2.NewPostgresRepository(conn)
		rpsAuth = repository2.NewPostgresRepository(conn)
	} else if flag == 2 {
		//Соединение с mongo
		client, cancel := repository2.RequestMongo()
		defer cancel()

		rps = repository2.NewMongoRepository(client)
		rpsAuth = repository2.NewMongoRepository(client)
	}

	var srv service.Service
	srv = service.NewCatService(rps)
	hndlr := handler.NewCatHandler(srv)
	e.GET("/cats", hndlr.GetAllCats)
	e.POST("/cats", hndlr.CreateCats)
	e.GET("/cats/:id", hndlr.GetCat)
	e.PUT("/cats/:id", hndlr.UpdateCat)
	e.DELETE("/cats/:id", hndlr.DeleteCat)

	var srvAuth service.Auth
	srvAuth = service.NewUserAuthService(rpsAuth)
	hndlrAuth := handler.NewUserAuthHandler(srvAuth)
	e.POST("/register", hndlrAuth.SignUp)
	e.POST("/login", hndlrAuth.SignIn)

	r := e.Group("/restrict")
	{
		config := middleware.JWTConfig{
			Claims:     new(service.JwtCustomClaims),
			SigningKey: []byte(viper.GetString("KEY_FOR_SIGNATURE_JWT")),
		}
		r.Use(middleware.JWTWithConfig(config))
		r.GET("", hndlrAuth.Restricted)
	}

	e.Logger.Fatal(e.Start(":8000"))
}
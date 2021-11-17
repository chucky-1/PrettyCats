package main

import (
	"CatsCrud/internal/handler"
	rep "CatsCrud/internal/repository"
	"CatsCrud/internal/request"
	"CatsCrud/internal/service"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"time"

	_ "CatsCrud/docs"
)

// 1 - postgres, 2 - mongo
const flag = 1

// @title Cats CRUD
// @version 1.0
// @description This simple application is written for teaching Go.

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	e := echo.New()
	e.Validator = &request.CustomValidator{Validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})

	var rps rep.Repository
	var rpsAuth rep.Auth
	if flag == 1 {
		// Соединение с postgres
		conn, err := rep.RequestDB()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		rps = rep.NewPostgresRepository(conn)
		rpsAuth = rep.NewPostgresRepository(conn)
	} else if flag == 2 {
		//Соединение с mongo
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := rep.RequestMongo(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()

		rps = rep.NewMongoRepository(client)
		rpsAuth = rep.NewMongoRepository(client)
	}

	var srv service.Service = service.NewCatService(rps)
	hndlr := handler.NewCatHandler(srv)
	e.GET("/cats", hndlr.GetAllCats)
	e.POST("/cats", hndlr.CreateCats)
	e.GET("/cats/:id", hndlr.GetCat)
	e.PUT("/cats/:id", hndlr.UpdateCat)
	e.DELETE("/cats/:id", hndlr.DeleteCat)

	var srvAuth service.Auth = service.NewUserAuthService(rpsAuth)
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

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":8000"))
}
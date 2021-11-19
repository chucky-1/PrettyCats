// Package main starts application
package main

import (
	_ "CatsCrud/docs"
	"CatsCrud/grpc/server"
	"CatsCrud/internal/handler"
	rep "CatsCrud/internal/repository"
	"CatsCrud/internal/request"
	"CatsCrud/internal/service"
	myGrpc "CatsCrud/protocol"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	sw "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"

	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	flag     = true // true - postgres; false - mongo
	portEcho = ":8000"
	portGrpc = "localhost:10000"
)

// NewGrpcServer is constructor
func NewGrpcServer(portGrpc string, srv service.Service) {
	lis, err := net.Listen("tcp", portGrpc)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	myGrpc.RegisterCatsCrudServer(s, server.NewCats(srv))
	fmt.Printf("server listening at %v\n", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
		return
	}
}

// @title Pretty Cats
// @version 1.0
// @description This simple application is written for teaching Go.

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// echo
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
	if flag {
		// Соединение с postgres
		conn, err := rep.RequestDB()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		rps = rep.NewPostgresRepository(conn)
		rpsAuth = rep.NewPostgresRepository(conn)
	} else {
		// Соединение с mongo
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		client, err := rep.RequestMongo(ctx)
		if err != nil {
			log.Panic(err)
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

	go NewGrpcServer(portGrpc, srv)

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

	e.GET("/swagger/*", sw.WrapHandler)
	e.Logger.Fatal(e.Start(portEcho))
}

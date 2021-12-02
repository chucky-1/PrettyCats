// Package main starts application
package main

import (
	_ "CatsCrud/docs"
	"CatsCrud/grpc/server"
	"CatsCrud/internal/configs"
	"CatsCrud/internal/handler"
	rep "CatsCrud/internal/repository"
	"CatsCrud/internal/request"
	"CatsCrud/internal/service"
	myGrpc "CatsCrud/protocol"
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	sw "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"fmt"
	"net"
	"net/http"
)

const (
	flag     = "postgres" // What database do you use? postgres / mongo
	portEcho = ":8000"
	portGrpc = "localhost:10000"
)

// newPgxPool sets a connection with postgres
func newPgxPool(ctx context.Context, cfg *configs.Config) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PsUsername, cfg.PsPassword, cfg.PsHost, cfg.PsPort, cfg.PsDBName)
	conn, err := pgxpool.Connect(ctx, url)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		return nil, fmt.Errorf("we can't connect to database")
	}
	return conn, nil
}

// newMongoClient sets a connection with mongodb
func newMongoClient(ctx context.Context, cfg *configs.Config) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.Mongo))
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("we can't connect to database")
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("we can't connect to database")
	}

	return client, nil
}

// newCache sets a connection with redis cache
func newCache(cfg *configs.Config) (*cache.Cache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			cfg.CacheHost: cfg.CachePort,
		},
	})

	myCache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Hour),
	})

	return myCache, nil
}

// newRedisClient sets a connection with redis
func newRedisClient(cfg *configs.Config) (*redis.Client, error) {
	hostAndPort := cfg.RedisHost + ":" + cfg.RedisPort
	rdb := redis.NewClient(&redis.Options{
		Addr:	  hostAndPort,
		Password: "", // no password set
		DB:		  0,  // use default DB
	})
	return rdb, nil
}

// newGrpcServer is constructor
func newGrpcServer(portGrpc string, srv service.Service) {
	lis, err := net.Listen("tcp", portGrpc)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	myGrpc.RegisterCatsCrudServer(s, server.NewServer(srv))
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

	// Configuration
	cfg := &configs.Config{}
	opts := &env.Options{}
	if err := env.Parse(cfg, *opts); err != nil {
		log.Fatal(err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})

	var rps rep.Repository
	var rpsAuth rep.Auth

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if flag == "postgres" {
		// Соединение с postgres
		conn, err := newPgxPool(ctx, cfg)
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()
		rps = rep.NewPostgresRepository(conn)
		rpsAuth = rep.NewPostgresRepository(conn)
	} else if flag == "mongo" {
		client, err := newMongoClient(ctx, cfg)
		if err != nil {
			log.Panic(err)
		}
		rps = rep.NewMongoRepository(client, cfg)
		rpsAuth = rep.NewMongoRepository(client, cfg)
	}

	cc, err := newCache(cfg)
	if err != nil {
		log.Panic(err)
	}
	myCache := rep.NewCache(cc)

	redisClient, err := newRedisClient(cfg)
	if err != nil {
		log.Panic(err)
	}

	ctx = context.TODO()
	var srv service.Service = service.NewCatService(ctx, rps, *myCache, redisClient)

	hndlr := handler.NewHandler(srv)
	e.GET("/cats", hndlr.GetAll)
	e.POST("/cats", hndlr.Create)
	e.GET("/cats/:id", hndlr.Get)
	e.PUT("/cats/:id", hndlr.Update)
	e.DELETE("/cats/:id", hndlr.Delete)

	var srvAuth service.Auth = service.NewUserAuth(rpsAuth, cfg)
	hndlrAuth := handler.NewAuthHandler(srvAuth)

	go newGrpcServer(portGrpc, srv)

	e.POST("/register", hndlrAuth.SignUp)
	e.POST("/login", hndlrAuth.SignIn)

	r := e.Group("/restrict")
	{
		config := middleware.JWTConfig{
			Claims:     new(service.JwtCustomClaims),
			SigningKey: []byte(cfg.KeyForSignatureJwt),
		}
		r.Use(middleware.JWTWithConfig(config))
		r.GET("", hndlrAuth.Restricted)
	}

	// Download file
	e.GET("/download", func(c echo.Context) error {
		return c.File("public/template/index.html")
	})
	e.GET("/download/file", func(c echo.Context) error {
		return c.File("public/media/echo-logo.svg")
	})

	// Upload file
	e.GET("/upload", func(c echo.Context) error {
		return c.File("public/template/upload.html")
	})
	e.POST("/upload", func(c echo.Context) error {
		name := c.FormValue("name")

		// Source
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer func(src multipart.File) {
			err = src.Close()
			if err != nil {
			}
		}(src)

		// Destination
		dst, err := os.Create(name + " " + file.Filename)
		if err != nil {
			return err
		}
		defer func(dst *os.File) {
			err = dst.Close()
			if err != nil {

			}
		}(dst)

		//Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully", file.Filename))
	})

	e.GET("/swagger/*", sw.WrapHandler)
	e.Logger.Fatal(e.Start(portEcho))
}

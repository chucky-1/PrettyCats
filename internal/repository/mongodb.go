package repository

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func RequestMongo() (*mongo.Client, context.CancelFunc) {
	if err := initConfig(); err != nil {
		log.Fatal("error config files")
	}

	//url := fmt.Sprintf("mongodb://%s:%s/",
	//	viper.GetString("mongodb.host"),
	//	viper.GetString("mongodb.port"))
	//
	//fmt.Println(url)

	url := os.Getenv("MONGODB_CONNSTRING")

	fmt.Println(url)

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Середина")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Всё ок!")
	return client, cancel
}
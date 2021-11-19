package repository

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"fmt"
)

// RequestMongo sets a connection with mongodb
func RequestMongo(ctx context.Context) (*mongo.Client, error) {
	if err := initConfig(); err != nil {
		log.Error("error config files")
		return nil, fmt.Errorf("we can't connect to database")
	}

	if err := godotenv.Load(); err != nil {
		err = godotenv.Load("C:/Users/User/GolandProjects/CatsCrud/.env")
		if err != nil {
			log.Error("error loading env variables")
			return nil, fmt.Errorf("we can't connect to database")
		}
	}

	// Для локальной разработке, закоментить при билдинге
	url := "mongodb://root:example@localhost:27017/"

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
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

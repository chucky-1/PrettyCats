package repository

import (
	"context"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func RequestMongo() (*mongo.Client, context.CancelFunc) {
	if err := initConfig(); err != nil {
		log.Fatal("error config files")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading env variables")
	}

	//url := fmt.Sprintf("mongodb://%s:%s/",
	//	viper.GetString("mongodb.host"),
	//	viper.GetString("mongodb.port"))
	//
	//fmt.Println(url)

	url := os.Getenv("MONGODB_CONNSTRING")

	// Для локальной разработке, закоментить при билдинге
	url = "mongodb://root:example@localhost:27017/"

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, cancel
}
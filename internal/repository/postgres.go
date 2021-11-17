package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func RequestDB() (*pgxpool.Pool, error) {
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

	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		viper.GetString("db.pos"),
		viper.GetString("db.username"),
		os.Getenv("DB_PASSWORD"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.dbase"))

	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Error(os.Stderr)
		log.Errorf("Unable to connect to database: %v\n", err)
		return nil, fmt.Errorf("we can't connect to database")
	}
	return conn, nil
}
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

func RequestDB() *pgxpool.Pool {
	if err := initConfig(); err != nil {
		log.Fatal("error config files")
	}

	if err := godotenv.Load(); err != nil {
		err = godotenv.Load("C:/Users/User/GolandProjects/CatsCrud/.env")
		if err != nil {
			log.Fatal("error loading env variables")
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
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
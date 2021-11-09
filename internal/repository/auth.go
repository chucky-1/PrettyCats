package repository

import (
	"CatsCrud/internal/models"
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

func (c *PostgresRepository) CreateUser(user models.User) (int, error) {
	var id int

	err := c.conn.QueryRow(context.Background(), "INSERT INTO users (ID, Name, Username, Password) " +
		"VALUES (nextval('users_sequence'), $1, $2, $3) RETURNING ID",
		user.Name, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, errors.New("error when adding to the database")
	}

	return id, nil
}

func (c *PostgresRepository) GetUser (username, password string) (models.User, error) {
	var user models.User

	err := c.conn.QueryRow(context.Background(), "SELECT id, name, username, password " +
		"FROM users WHERE username = $1", username).Scan(&user.ID, &user.Name, &user.Username, &user.Password)

	if err != nil {
		return *new(models.User), errors.New("error at working with database")
	}

	if user.Password != password {
		return *new(models.User), errors.New("password isn't corrected")
	}

	return user, nil
}


func (c *MongoRepository) CreateUser(user models.User) (int, error) {
	collection := c.client.Database("users").Collection("users")

	docs := []interface{}{
		bson.D{primitive.E{Key: "name", Value: user.Name}, {Key: "username", Value: user.Username},
			{Key: "password", Value: user.Password}},
	}

	_, insertErr := collection.InsertMany(context.TODO(), docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	id := 1

	return id, nil
}

func (c *MongoRepository) GetUser(username, password string) (models.User, error) {
	return *new(models.User), nil
}

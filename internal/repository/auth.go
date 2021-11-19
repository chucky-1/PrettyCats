package repository

import (
	"CatsCrud/internal/models"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"context"
	"errors"
)

// Auth is responsible for registration and authorization
type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

// CreateUser creates user in database
func (c *PostgresRepository) CreateUser(user models.User) (int, error) {
	var id int
	err := c.conn.QueryRow(context.Background(),
		"INSERT INTO users (ID, Name, Username, Password) VALUES (nextval('users_sequence'), $1, $2, $3) RETURNING ID",
		user.Name, user.Username, user.Password).Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}

// GetUser gets user from database
func (c *PostgresRepository) GetUser(username, password string) (models.User, error) {
	var user models.User
	err := c.conn.QueryRow(context.Background(),
		"SELECT id, name, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Name, &user.Username, &user.Password)

	if err != nil {
		log.Error("error at working with database")
		return models.User{}, err
	}

	if user.Password != password {
		return models.User{}, errors.New("password isn't corrected")
	}

	return user, nil
}

// CreateUser creates user in database
func (c *MongoRepository) CreateUser(user models.User) (int, error) {
	collection := c.client.Database("users").Collection("users")

	docs := []interface{}{
		bson.D{primitive.E{Key: "name", Value: user.Name}, {Key: "username", Value: user.Username},
			{Key: "password", Value: user.Password}},
	}

	_, insertErr := collection.InsertMany(context.TODO(), docs)
	if insertErr != nil {
		log.Error(insertErr)
		return 0, insertErr
	}

	id := 1

	return id, nil
}

// GetUser gets user from database
func (c *MongoRepository) GetUser(username, password string) (models.User, error) {
	return models.User{}, nil
}

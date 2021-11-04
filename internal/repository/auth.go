package repository

import (
	"CatsCrud/internal/models"
	"context"
	"errors"
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
	return 0, nil
}

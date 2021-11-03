package repository

import (
	"CatsCrud/internal/models"
	"context"
	"fmt"
)

type Auth interface {
	CreateUser(user models.User) (int, error)
}

func (c *PostgresRepository) CreateUser(user models.User) (int, error) {
	fmt.Println("Create user start...")
	var id int

	err := c.conn.QueryRow(context.Background(), "INSERT INTO users (ID, Name, Username, Password) " +
		"VALUES (nextval('users_sequence'), $1, $2, $3) RETURNING ID", user.Name, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	fmt.Println("Insert successful")

	return id, nil
}

func (c *MongoRepository) CreateUser(user models.User) (int, error) {
	return 0, nil
}

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

const (
	pos = "postgres"
	username = "postgres"
	password = "220095Qwe"
	//host = "localhost:5432"
	host = "db:5436"
	dbase = "postgres"
)

func main() {
	e := echo.New()

	//Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Cats!")
	})
	e.GET("/cats", getAllCats)
	e.POST("/cats", createCats)
	e.GET("/cats/:id", getCat)
	e.PUT("/cats/:id", updateUser)
	e.DELETE("/cats/:id", deleteCat)

	e.Logger.Fatal(e.Start(":8000"))
}

func requestDB() *pgx.Conn {
	url := fmt.Sprintf("%s://%s:%s@%s/%s", pos, username, password, host, dbase)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

func getAllCats(c echo.Context) error {
	conn := requestDB()
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select ID, name from cats")
	if err != nil {
		log.Fatal(err)
	}

	type cat struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	}
	var allcats = map[int32]*cat{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}
		id := values[0].(int32)
		name := values[1].(string)

		ct := &cat {
			ID: id,
			Name: name,
		}
		allcats[ct.ID] = ct
	}
	return c.JSON(http.StatusOK, allcats)
}

func createCats(c echo.Context) error {
	conn := requestDB()
	defer conn.Close(context.Background())

	//Не смог достучаться до c.Param
	//test := c.QueryParam("name")
	//fmt.Println(test)

	// Получить последний id
	rows, err := conn.Query(context.Background(), "select ID, name from cats")
	if err != nil {
		log.Fatal(err)
	}
	var id int32
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}
		id = values[0].(int32)
	}

	id ++
	name := fmt.Sprintf("cat%d", id)
	commandTag, err := conn.Exec(context.Background(),
		fmt.Sprintf("INSERT INTO cats VALUES (%d, '%s')", id, name))
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}

func getCat(c echo.Context) error {
	conn := requestDB()
	defer conn.Close(context.Background())

	id := c.Param("id")

	rows, err := conn.Query(context.Background(), fmt.Sprintf("select ID, name from cats where id=%s", id))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}
		id := values[0].(int32)
		name := values[1].(string)
		return c.String(http.StatusOK, fmt.Sprintf("id: %d, Имя: %s", id, name))
	}
	return nil
}

func updateUser(c echo.Context) error {
	conn := requestDB()
	defer conn.Close(context.Background())

	id := c.Param("id")
	name := c.QueryParam("name")

	commandTag, err := conn.Exec(context.Background(),
		fmt.Sprintf("UPDATE cats SET name = %s WHERE id = %s", name, id))
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}

func deleteCat(c echo.Context) error {
	conn := requestDB()
	defer conn.Close(context.Background())

	id := c.Param("id")

	commandTag, err := conn.Exec(context.Background(), fmt.Sprintf("delete from cats where id=%s", id))
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}
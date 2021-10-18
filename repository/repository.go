package repository

import (
	"CatsCrud/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

type Repository struct {
	conn *pgxpool.Pool
	client *mongo.Client
}

func NewRepository(conn *pgxpool.Pool, client *mongo.Client) *Repository {
	return &Repository{conn: conn, client: client}
}

func usedb() int {
	usedb := viper.GetInt("usedb")
	return usedb
}

func (c *Repository) GetAllCats() ([]*models.Cats, error) {

	var allcats = []*models.Cats{}

	if usedb() == 1 {
		rows, err := c.conn.Query(context.Background(), "select ID, name from cats")
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {

			cats := models.Cats{
				ID:   0,
				Name: "",
			}

			values, err := rows.Values()
			if err != nil {
				log.Fatal(err)
			}
			cats.ID = values[0].(int32)
			cats.Name = values[1].(string)
			allcats = append(allcats, &cats)
		}

		return allcats, nil
	} else if usedb() == 2 {
		collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
		cur, currErr := collection.Find(context.TODO(), bson.D{})
		if currErr != nil { panic(currErr) }
		defer cur.Close(context.TODO())

		if err := cur.All(context.TODO(), &allcats); err != nil {
			panic(err)
		}
		return allcats, nil
	}
	return allcats, nil
}

func (c *Repository) CreateCats(jsonMap map[string]interface{}) (*models.Cats, error) {

	// Инициализация структуры models.Cats
	cat := models.Cats{
		ID:   0,
		Name: "",
	}

	// Достаём id, name. Id преобразуем в int
	var id interface{}
	id = jsonMap["id"]
	idInt, err := strconv.Atoi(id.(string))
	if err != nil {
		return &cat, err
	}

	// Присваиваем значения в структуру models.Cats
	cat.ID = int32(idInt)
	name := jsonMap["name"]
	cat.Name = fmt.Sprintf("%v", name)

	if usedb() == 1 {
		// Добавляем в базу данных
		commandTag, err := c.conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", idInt, name)
		if err != nil {
			return &cat, err
		}
		if commandTag.RowsAffected() != 1 {
			return &cat, errors.New("Failed to create cat")
		}
		return &cat, nil
	} else if usedb() == 2 {
		collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

		docs := []interface{}{
			bson.D{{"id" , cat.ID},{"name" , cat.Name}},
		}

		_, insertErr := collection.InsertMany(context.TODO(), docs)
		if insertErr != nil {
			log.Fatal(insertErr)
		}
		return &cat, nil
	}
	return &cat, nil
}

func (c *Repository) GetCat(id string) (*models.Cats, error) {

	var cat models.Cats

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, nil
	}

	if usedb() == 1 {
		// Достаём name
		var name string
		err := c.conn.QueryRow(context.Background(), "select name from cats where id=$1", id).Scan(&name)
		if err != nil {
			return &cat, err
		}

		//Присваиваем параметры models.Cats
		cat.ID = int32(idInt)
		cat.Name = name

	} else if usedb() == 2 {
		collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

		err := collection.FindOne(context.TODO(), bson.D{{"id", idInt}}).Decode(&cat)
		if err != nil {
			return &cat, err
		}
	}

	return &cat, nil
}

func (c *Repository) UpdateCat(id string, jsonMap map[string]interface{}) (*models.Cats, error) {

	var cat models.Cats

	// Преобразуем id в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, err
	}

	// Достаём name
	name := jsonMap["name"]

	// Обновляем models.Cat
	cat.ID = int32(idInt)
	cat.Name = fmt.Sprintf("%v", name)

	if usedb() == 1 {
		// Вносим изменения в базу данных
		commandTag, err := c.conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", name, idInt)
		if err != nil {
			return &cat, err
		}
		if commandTag.RowsAffected() != 1 {
			return &cat, errors.New("Row isn't update")
		}
	} else if usedb() == 2 {
		collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

		filter := bson.D{{"id", idInt}}
		update := bson.D{{"$set", bson.D{{"name", name}}}}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &cat, nil
}

func (c *Repository) DeleteCat(id string) (*models.Cats, error) {

	var cat models.Cats

	// Преобразуем id в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, err
	}

	// Обновляес models.Cats
	cat.ID = int32(idInt)

	if usedb() == 1 {
		// Удаление из базы
		commandTag, err := c.conn.Exec(context.Background(), "delete from cats where id=$1", id)
		if err != nil {
			return &cat, err
		}
		if commandTag.RowsAffected() != 1 {
			return &cat, errors.New("No row found to delete")
		}
	} else if usedb() == 2 {
		collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

		_, err :=collection.DeleteOne(context.TODO(), bson.D{{"id", idInt}})
		if err != nil {
			log.Fatal(err)
		}
	}

	return &cat, nil
}
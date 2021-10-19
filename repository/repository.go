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

type PostgresRepository struct {
	conn *pgxpool.Pool
}

type MongoRepository struct {
	client *mongo.Client
}

type Repository interface {
	GetAllCats() ([]*models.Cats, error)
	CreateCats(jsonMap map[string]interface{}) (*models.Cats, error)
	GetCat(id string) (*models.Cats, error)
	UpdateCat(id string, jsonMap map[string]interface{}) (*models.Cats, error)
	DeleteCat(id string) (*models.Cats, error)
}

func NewPostgresRepository(conn *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{client: client}
}

func (c *PostgresRepository) GetAllCats() ([]*models.Cats, error) {

	var allcats = []*models.Cats{}

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
}

func (c *PostgresRepository) CreateCats(jsonMap map[string]interface{}) (*models.Cats, error) {

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

	// Добавляем в базу данных
	commandTag, err := c.conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", idInt, name)
	if err != nil {
		return &cat, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cat, errors.New("Failed to create cat")
	}

	return &cat, nil
}

func (c *PostgresRepository) GetCat(id string) (*models.Cats, error) {

	var cat models.Cats

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, nil
	}

	// Достаём name
	var name string
	err = c.conn.QueryRow(context.Background(), "select name from cats where id=$1", id).Scan(&name)
	if err != nil {
		return &cat, err
	}

	//Присваиваем параметры models.Cats
	cat.ID = int32(idInt)
	cat.Name = name

	return &cat, nil
}

func (c *PostgresRepository) UpdateCat(id string, jsonMap map[string]interface{}) (*models.Cats, error) {

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

	// Вносим изменения в базу данных
	commandTag, err := c.conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", name, idInt)
	if err != nil {
		return &cat, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cat, errors.New("Row isn't update")
	}

	return &cat, nil
}

func (c *PostgresRepository) DeleteCat(id string) (*models.Cats, error) {

	var cat models.Cats

	// Преобразуем id в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, err
	}

	// Обновляес models.Cats
	cat.ID = int32(idInt)

	// Удаление из базы
	commandTag, err := c.conn.Exec(context.Background(), "delete from cats where id=$1", id)
	if err != nil {
		return &cat, err
	}
	if commandTag.RowsAffected() != 1 {
		return &cat, errors.New("No row found to delete")
	}

	return &cat, nil
}

func (c *MongoRepository) GetAllCats() ([]*models.Cats, error) {

	var allcats = []*models.Cats{}

	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	cur, currErr := collection.Find(context.TODO(), bson.D{})
	if currErr != nil { panic(currErr) }
	defer cur.Close(context.TODO())

	if err := cur.All(context.TODO(), &allcats); err != nil {
		panic(err)
	}
	return allcats, nil
}

func (c *MongoRepository) CreateCats(jsonMap map[string]interface{}) (*models.Cats, error) {

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

func (c *MongoRepository) GetCat(id string) (*models.Cats, error) {
	var cat models.Cats

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, nil
	}

	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

	err = collection.FindOne(context.TODO(), bson.D{{"id", idInt}}).Decode(&cat)
	if err != nil {
		return &cat, err
	}
	return &cat, nil
}

func (c *MongoRepository) UpdateCat(id string, jsonMap map[string]interface{}) (*models.Cats, error) {
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

	// Вносим изменения в базу данных
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	filter := bson.D{{"id", idInt}}
	update := bson.D{{"$set", bson.D{{"name", name}}}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return &cat, nil
}

func (c *MongoRepository) DeleteCat(id string) (*models.Cats, error) {
	var cat models.Cats

	// Преобразуем id в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &cat, err
	}

	// Обновляес models.Cats
	cat.ID = int32(idInt)

	//	Удаляем из базы
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	_, err =collection.DeleteOne(context.TODO(), bson.D{{"id", idInt}})
	if err != nil {
		log.Fatal(err)
	}

	return &cat, nil
}
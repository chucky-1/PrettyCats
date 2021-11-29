// Package repository connects with database and send given files in service
package repository

import (
	"CatsCrud/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
	"fmt"
	"strconv"
)

const (
	bitSize = 32
)

// PostgresRepository provides a connection with postgres
type PostgresRepository struct {
	conn *pgxpool.Pool
}

// MongoRepository provides a connection with postgres
type MongoRepository struct {
	client *mongo.Client
}

// Repository has methods which work with database
type Repository interface {
	GetAll() ([]*models.Cat, error)
	Create(cats models.Cat) (*models.Cat, error)
	Get(id string) (*models.Cat, error)
	Update(id string, cats models.Cat) (*models.Cat, error)
	Delete(id string) (*models.Cat, error)
}

// NewPostgresRepository is constructor
func NewPostgresRepository(conn *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

// NewMongoRepository is constructor
func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{client: client}
}

// GetAll gets all cats from database
func (c *PostgresRepository) GetAll() ([]*models.Cat, error) {
	var allcats []*models.Cat

	rows, err := c.conn.Query(context.Background(), "select ID, name from cats")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for rows.Next() {
		cats := models.Cat{
			ID:   0,
			Name: "",
		}

		values, err := rows.Values()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		cats.ID = values[0].(int32)
		cats.Name = values[1].(string)
		allcats = append(allcats, &cats)
	}

	return allcats, nil
}

// Create creates cats in database
func (c *PostgresRepository) Create(cats models.Cat) (*models.Cat, error) {
	commandTag, err := c.conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", cats.ID, cats.Name)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if commandTag.RowsAffected() != 1 {
		log.Error("Failed to create cat")
		return nil, err
	}

	return &cats, nil
}

// Get gets one cats from database by ID
func (c *PostgresRepository) Get(id string) (*models.Cat, error) {
	var cat models.Cat

	idInt, err := strconv.ParseInt(id, 0, bitSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Достаём name
	var name string
	err = c.conn.QueryRow(context.Background(), "select name from cats where id=$1", id).Scan(&name)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Присваиваем параметры models.Cats
	cat.ID = int32(idInt)
	cat.Name = name

	return &cat, nil
}

// Update updates cat in database by ID
func (c *PostgresRepository) Update(id string, cats models.Cat) (*models.Cat, error) {
	// Преобразуем id в int
	idInt, err := strconv.ParseInt(id, 0, bitSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Обновляем models.Cat
	cats.ID = int32(idInt)

	// Вносим изменения в базу данных
	commandTag, err := c.conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", cats.Name, cats.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if commandTag.RowsAffected() != 1 {
		log.Fatal("Row isn't update")
		return nil, fmt.Errorf("cat doesn't update")
	}

	return &cats, nil
}

// Delete deletes cat from database by ID
func (c *PostgresRepository) Delete(id string) (*models.Cat, error) {
	var cat models.Cat

	// Преобразуем id в int
	idInt, err := strconv.ParseInt(id, 0, bitSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Обновляем models.Cats
	cat.ID = int32(idInt)
	// Достаём name
	var name string
	err = c.conn.QueryRow(context.Background(), "select name from cats where id=$1", id).Scan(&name)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	cat.Name = name

	// Удаление из базы
	commandTag, err := c.conn.Exec(context.Background(), "delete from cats where id=$1", id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if commandTag.RowsAffected() != 1 {
		log.Error("No row found to delete")
		return nil, fmt.Errorf("cat can't delete")
	}

	return &cat, nil
}

// GetAll gets all cats from database
func (c *MongoRepository) GetAll() ([]*models.Cat, error) {
	var allcats = []*models.Cat{}

	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	cur, currErr := collection.Find(context.TODO(), bson.D{})
	if currErr != nil {
		log.Error(currErr)
		return nil, currErr
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error(err)
			return
		}
	}(cur, context.TODO())

	if err := cur.All(context.TODO(), &allcats); err != nil {
		log.Error(err)
		return nil, err
	}
	return allcats, nil
}

// Create creates cats in database
func (c *MongoRepository) Create(cats models.Cat) (*models.Cat, error) {
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

	docs := []interface{}{
		bson.D{primitive.E{Key: "id", Value: cats.ID}, {Key: "name", Value: cats.Name}},
	}

	_, insertErr := collection.InsertMany(context.TODO(), docs)
	if insertErr != nil {
		log.Error(insertErr)
		return nil, insertErr
	}
	return &cats, nil
}

// Get gets one cats from database by ID
func (c *MongoRepository) Get(id string) (*models.Cat, error) {
	var cat models.Cat

	idInt, err := strconv.ParseInt(id, 0, bitSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))

	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "id", Value: idInt}}).Decode(&cat)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &cat, nil
}

// Update updates cat in database by ID
func (c *MongoRepository) Update(id string, cats models.Cat) (*models.Cat, error) {
	// Преобразуем id в int
	idInt, err := strconv.ParseInt(id, 0, bitSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Обновляем models.Cat
	cats.ID = int32(idInt)

	// Вносим изменения в базу данных
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	filter := bson.D{primitive.E{Key: "id", Value: cats.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "name", Value: cats.Name}}}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &cats, nil
}

// Delete deletes cat from database by ID
func (c *MongoRepository) Delete(id string) (*models.Cat, error) {
	var cat models.Cat

	// Преобразуем id в int
	idInt, err := strconv.ParseInt(id, 0, bitSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Обновляес models.Cats
	cat.ID = int32(idInt)

	//	Удаляем из базы
	collection := c.client.Database(viper.GetString("mongodb.dbase")).Collection(viper.GetString("mongodb.collection"))
	_, err = collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "id", Value: idInt}})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &cat, nil
}

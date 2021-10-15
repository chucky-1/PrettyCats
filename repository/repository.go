package repository

import (
	"CatsCrud/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
)

type Repository struct {
	cats *models.Cats
	conn *pgxpool.Pool
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func NewRepository(pool pgxpool.Pool) *Repository {

	// Инициализация models.Cats
	cats := models.Cats{
		ID:   0,
		Name: "",
	}

	// Соединение с БД
	if err := initConfig(); err != nil {
		log.Fatal("error config files")
	}
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading env variables")
	}
	url := fmt.Sprintf("%s://%s:%s@%s/%s",
		viper.GetString("db.pos"),
		viper.GetString("db.username"),
		os.Getenv("DB_PASSWORD"),
		viper.GetString("db.host"),
		viper.GetString("db.dbase"))

	conn, err := pgxpool.Connect(context.Background(), url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &Repository{cats: &cats, conn: conn}
}

func (c *Repository) GetAllCats() ([]*models.Cats, error) {
	//conn := RequestDB()
	defer c.conn.Close()

	rows, err := c.conn.Query(context.Background(), "select ID, name from cats")
	if err != nil {
		log.Fatal(err)
	}

	var allcats = []*models.Cats{}

	for rows.Next() {

		////cats := models.Cats{
		////	ID:   0,
		////	Name: "",
		////}
		//cat := NewRepository()

		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}
		c.cats.ID = values[0].(int32)
		c.cats.Name = values[1].(string)
		allcats = append(allcats, c.cats)
	}

	return allcats, nil
}

//func (c *Repository) CreateCats(jsonMap map[string]interface{}) (*models.Cats, error) {
//	conn := RequestDB()
//	defer conn.Close()
//
//	// Инициализация структуры models.Cats
//	cat := models.Cats{
//		ID:   0,
//		Name: "",
//	}
//
//	// Достаём id, name. Id преобразуем в int
//	var id interface{}
//	id = jsonMap["id"]
//	idInt, err := strconv.Atoi(id.(string))
//	if err != nil {
//		return &cat, err
//	}
//
//	// Присваиваем значения в структуру models.Cats
//	cat.ID = int32(idInt)
//	name := jsonMap["name"]
//	cat.Name = fmt.Sprintf("%v", name)
//
//	// Добавляем в базу данных
//	commandTag, err := conn.Exec(context.Background(), "INSERT INTO cats VALUES ($1, $2)", idInt, name)
//	if err != nil {
//		return &cat, err
//	}
//	if commandTag.RowsAffected() != 1 {
//		return &cat, errors.New("Failed to create cat")
//	}
//	return &cat, nil
//}
//
//func (c *Repository) GetCat(id string) (*models.Cats, error) {
//	conn := RequestDB()
//	defer conn.Close()
//
//	// Инициализация структуры models.Cats
//	cat := models.Cats{
//		ID:   0,
//		Name: "",
//	}
//
//	// Достаём name
//	var name string
//	err := conn.QueryRow(context.Background(), "select name from cats where id=$1", id).Scan(&name)
//	if err != nil {
//		return &cat, err
//	}
//
//	//Присваиваем параметры models.Cats
//	idInt, err := strconv.Atoi(id)
//	if err != nil {
//		return &cat, nil
//	}
//	cat.ID = int32(idInt)
//	cat.Name = name
//
//	return &cat, nil
//}
//
//func (c *Repository) UpdateCat(id string, jsonMap map[string]interface{}) (*models.Cats, error) {
//	conn := RequestDB()
//	defer conn.Close()
//
//	// Инициализация структуры models.Cats
//	cat := models.Cats{
//		ID:   0,
//		Name: "",
//	}
//
//	// Преобразуем id в int
//	idInt, err := strconv.Atoi(id)
//	if err != nil {
//		return &cat, err
//	}
//
//	// Достаём name
//	name := jsonMap["name"]
//
//	// Обновляем models.Cat
//	cat.ID = int32(idInt)
//	cat.Name = fmt.Sprintf("%v", name)
//
//	// Вносим изменения в базу данных
//	commandTag, err := conn.Exec(context.Background(), "UPDATE cats SET name = $1 WHERE id = $2", name, idInt)
//	if err != nil {
//		return &cat, err
//	}
//	if commandTag.RowsAffected() != 1 {
//		return &cat, errors.New("Row isn't update")
//	}
//
//	return &cat, nil
//}
//
//func (c *Repository) DeleteCat(id string) (*models.Cats, error) {
//	conn := RequestDB()
//	defer conn.Close()
//
//	// Инициализация структуры models.Cats
//	cat := models.Cats{
//		ID:   0,
//		Name: "",
//	}
//
//	// Преобразуем id в int
//	idInt, err := strconv.Atoi(id)
//	if err != nil {
//		return &cat, err
//	}
//
//	// Обновляес models.Cats
//	cat.ID = int32(idInt)
//
//	// Удаление из базы
//	commandTag, err := conn.Exec(context.Background(), "delete from cats where id=$1", id)
//	if err != nil {
//		return &cat, err
//	}
//	if commandTag.RowsAffected() != 1 {
//		return &cat, errors.New("No row found to delete")
//	}
//
//	return &cat, nil
//}
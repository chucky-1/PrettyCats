package repository

import (
	"CatsCrud/internal/models"
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 180 * time.Second
	if err = pool.Retry(func() error {
		db, err = pgxpool.Connect(context.Background(), databaseUrl)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// Временное решешие по миграциям, что бы написать тесты
func MyMigrate(pool *pgxpool.Pool) {
	exec, err := pool.Exec(context.Background(), "CREATE TABLE cats (ID integer, Name varchar(120));")
	log.Println(exec)
	if err != nil {
		log.Fatal("Fatal 2", err)
	}
	commandTag, err := pool.Exec(context.Background(), "INSERT INTO cats VALUES(1, 'cat'),(2, 'dog');")
	if err != nil {
		log.Fatal("Fatal 3", err)
	}
	//if commandTag.RowsAffected() != 1 {
	//	log.Fatal("Failed to create cat")
	//}
	log.Println(commandTag.String())
}

// Количество элементов в таблице
const countOfCats = 2

// Инициализация репозитория
var rps Repository
func NewPostgresRepositoryTest(db *pgxpool.Pool) {
	rps = NewPostgresRepository(db)
}

func TestInit(t *testing.T)  {
	// Делаем миграции и инициализируем репозиторий
	MyMigrate(db)
	NewPostgresRepositoryTest(db)
}

func TestPostgresRepository_GetAllCats(t *testing.T) {
	// Тест
	allcats, err := rps.GetAllCats()

	// Проверка сопоставления типов
	typeAllcats := fmt.Sprintf("%T", allcats)
	var tr []*models.Cats
	TypeTrue := fmt.Sprintf("%T", tr)
	assert.Equal(t, typeAllcats, TypeTrue)

	// Проверка количества элементов в базе
	count := len(allcats)
	assert.Equal(t, count, countOfCats)

	// Проверка ошибок
	assert.Nil(t, err)
}

func TestPostgresRepository_CreateCats(t *testing.T) {
	// Входящие значения
	jsonMap := make(map[string]interface{})
	jsonMap["id"] = "3"
	jsonMap["name"] = "cat3"
	catTrue := models.Cats{
		ID:   3,
		Name: "cat3",
	}

	// Тест
	cat, err := rps.CreateCats(jsonMap)

	// Проверка возвращаемых значений
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)

	// Проверка добавления в базу
	allcats, err := rps.GetAllCats()
	if err != nil {
		log.Fatal(err)
	}
	count := len(allcats)
	assert.Equal(t, count, countOfCats + 1)
}

func TestPostgresRepository_GetCat(t *testing.T) {
	// Входящие значения
	catTrue := models.Cats{
		ID:   3,
		Name: "cat3",
	}
	id := "3"

	// Тест
	cat, err := rps.GetCat(id)

	// Проверка возвращаемых значений
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)
}

func TestPostgresRepository_UpdateCat(t *testing.T) {
	// Входящие значения
	id := "1"
	jsonMap := make(map[string]interface{})
	jsonMap["name"] = "dog"
	catTrue := models.Cats{
		ID:   1,
		Name: "dog",
	}

	//Тест
	cat, err := rps.UpdateCat(id, jsonMap)

	// Проверка возвращаемых значений
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)

	// Проверка внесения изменений в базу
	cat, err = rps.GetCat(id)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, cat, &catTrue)
}

func TestPostgresRepository_DeleteCat(t *testing.T) {
	// Входящие значения
	catTrue := models.Cats{
		ID:   3,
		Name: "cat3",
	}
	id := "3"

	// Тест
	cat, err := rps.DeleteCat(id)

	// Проверка возвращаемых значений
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)
}
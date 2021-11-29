package repository

import (
	"CatsCrud/internal/models"
	"errors"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

var db *pgxpool.Pool
var hostAndPort string

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Errorf("Could not connect to docker: %s", err)
		return
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
		log.Errorf("Could not start resource: %s", err)
	}

	hostAndPort = resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseURL)

	err = resource.Expire(120)
	if err != nil {
		log.Error(err)
		return
	} // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 180 * time.Second
	if err = pool.Retry(func() error {
		db, err = pgxpool.Connect(context.Background(), databaseURL)
		if err != nil {
			log.Error(err)
		}
		return nil
	}); err != nil {
		log.Errorf("Could not connect to docker: %s", err)
		return
	}

	// Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Errorf("Could not purge resource: %s", err)
		return
	}

	os.Exit(code)
}

// Миграции
func MyMigrate() {
	cmd := exec.Command("flyway", "-url=jdbc:postgresql://", hostAndPort, "/dbname", "-user=user_name", "-password=secret", "migrate")
	cmd.Dir = "C:/Program Files/flyway-8.0.3"
	err := cmd.Run()
	if err != nil {
		log.Error(err)
		return
	}
}

// Количество элементов в таблице
const countOfCats = 2

// Инициализация репозитория
var rps Repository
var rpsAuth Auth

func NewPostgresRepositoryTest(db *pgxpool.Pool) {
	rps = NewPostgresRepository(db)
	rpsAuth = NewPostgresRepository(db)
}

func TestInit(t *testing.T) {
	// Делаем миграции и инициализируем репозиторий
	MyMigrate()
	NewPostgresRepositoryTest(db)
}

func TestPostgresRepository_GetAllCats(t *testing.T) {
	// Тест
	allcats, err := rps.GetAll()

	// Проверка сопоставления типов
	typeAllcats := fmt.Sprintf("%T", allcats)
	var tr []*models.Cat
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
	catsTrue := models.Cat{
		ID:   3,
		Name: "cat3",
	}

	// Тест
	cat, err := rps.Create(catsTrue)

	// Проверка возвращаемых значений
	assert.Equal(t, cat, &catsTrue)
	assert.Nil(t, err)

	// Проверка добавления в базу
	allcats, err := rps.GetAll()
	if err != nil {
		log.Error(err)
		return
	}
	count := len(allcats)
	assert.Equal(t, count, countOfCats+1)
}

func TestPostgresRepository_GetCat(t *testing.T) {
	// Входящие значения
	catTrue := models.Cat{
		ID:   3,
		Name: "cat3",
	}
	id := "3"

	// Тест
	cat, err := rps.Get(id)

	// Проверка возвращаемых значений
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)
}

func TestPostgresRepository_UpdateCat(t *testing.T) {
	// Входящие значения
	id := "1"
	catTrue := models.Cat{
		ID:   1,
		Name: "dogs",
	}

	// Тест
	cat, err := rps.Update(id, catTrue)

	// Проверка возвращаемых значений
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)

	// Проверка внесения изменений в базу
	cat, err = rps.Get(id)
	if err != nil {
		log.Error(err)
		return
	}
	assert.Equal(t, cat, &catTrue)
}

func TestPostgresRepository_DeleteCat(t *testing.T) {
	// Входящие значения
	catTrue := models.Cat{
		ID:   3,
		Name: "cat3",
	}
	id := "3"

	// Тест
	cat, err := rps.Delete(id)

	// Проверка возвращаемых значений
	assert.Equal(t, cat, &catTrue)
	assert.Nil(t, err)
}

func TestPostgresRepository_CreateUser(t *testing.T) {
	TestTable := []struct {
		name        string
		user        models.User
		exceptError error
	}{
		{
			name: "OK",
			user: models.User{
				ID: 1, Name: "Jon Snow", Username: "Jonny", Password: "Jon25@ew5",
			},
			exceptError: nil,
		},
	}

	for _, TestCase := range TestTable {
		TestCase := TestCase // pin!
		t.Run(TestCase.name, func(t *testing.T) {
			id, err := rpsAuth.CreateUser(TestCase.user)

			// Проверка на возвращаемые значения
			assert.Equal(t, TestCase.user.ID, id)
			assert.Equal(t, TestCase.exceptError, err)

			// Проверка на добавление в базу
			var newUser models.User
			err = db.QueryRow(context.Background(), "SELECT id, name, username, password FROM users WHERE username = $1",
				TestCase.user.Username).Scan(&newUser.ID, &newUser.Name, &newUser.Username, &newUser.Password)
			if err != nil {
				log.Error(err)
				return
			}

			assert.Equal(t, TestCase.user, newUser, "User didn't corrected to save")
		})
	}
}

func TestMongoRepository_GetUser(t *testing.T) {
	TestTable := []struct {
		name          string
		inputUsername string
		inputPassword string
		expectUser    models.User
		exceptError   error
	}{
		{
			name: "OK", inputUsername: "Jonny", inputPassword: "Jon25@ew5",
			expectUser: models.User{
				ID: 1, Name: "Jon Snow", Username: "Jonny", Password: "Jon25@ew5",
			}, exceptError: nil,
		},
		{
			name: "user isn't in database", inputUsername: "Bob", inputPassword: "random", expectUser: models.User{}, exceptError: errors.New("error at working with database"),
		},
		{
			name: "password isn't corrected", inputUsername: "Jonny", inputPassword: "random", expectUser: models.User{}, exceptError: errors.New("password isn't corrected"),
		},
	}

	for _, TestCase := range TestTable {
		TestCase := TestCase // pin!
		t.Run(TestCase.name, func(t *testing.T) {
			user, err := rpsAuth.GetUser(TestCase.inputUsername, TestCase.inputPassword)

			assert.Equal(t, TestCase.expectUser, user)
			assert.Equal(t, TestCase.exceptError, err)
		})
	}
}

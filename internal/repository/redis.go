package repository

import (
	"CatsCrud/internal/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

// RedisRepository provides a connection with redis
type RedisRepository struct {
	rdb *redis.Client
}

// NewRedisRepository is constructor
func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	if err := initConfig(); err != nil {
		log.Error("error config files")
		return nil, fmt.Errorf("we can't connect to database")
	}

	if err := godotenv.Load(); err != nil {
		err = godotenv.Load("C:/Users/User/GolandProjects/CatsCrud/.env")
		if err != nil {
			log.Error("error loading env variables")
			return nil, fmt.Errorf("we can't connect to database")
		}
	}

	hostAndPort := viper.GetString("redis.host") + ":" + viper.GetString("redis.port")

	rdb := redis.NewClient(&redis.Options{
		Addr:     hostAndPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb, nil
}

func (c *RedisRepository) CreateCat(cats models.Cats) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := c.rdb.Set(ctx, strconv.Itoa(int(cats.ID)), cats.Name, 0).Err()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (c *RedisRepository) GetCat(id string) (*models.Cats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	val, err := c.rdb.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	cat := new(models.Cats)
	idInt, err := strconv.ParseInt(id, 0, bitSize)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	cat.ID = int32(idInt)
	cat.Name = val

	return cat, nil
}

func (c *RedisRepository) DeleteCat(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c.rdb.Del(ctx, id)
}

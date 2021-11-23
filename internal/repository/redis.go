package repository

import (
	"CatsCrud/internal/models"
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// RedisRepository provides a connection with redis
type RedisRepository struct {
	rdb *cache.Cache
}

// NewRedisRepository is constructor
func NewRedisRepository(rdb *cache.Cache) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

func NewRedisClient() (*cache.Cache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
			"server2": ":6380",
		},
	})

	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Hour),
	})

	return mycache, nil
}

func (c *RedisRepository) CreateCat(cats models.Cats) error {
	ctx := context.TODO()

	if err := c.rdb.Set(&cache.Item{
		Ctx:   ctx,
		Key:   strconv.Itoa(int(cats.ID)),
		Value: cats,
	}); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (c *RedisRepository) GetCat(id string) (*models.Cats, error) {
	ctx := context.TODO()

	var cat models.Cats
	err := c.rdb.Get(ctx, id, &cat)

	if err != nil {
		if err.Error() == "cache: key is missing" {
			return nil, err
		} else {
			log.Error(err)
			return nil, err
		}
	}

	return &cat, nil
}

func (c *RedisRepository) DeleteCat(id string) error {
	ctx := context.TODO()

	err := c.rdb.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

package repository

import (
	"CatsCrud/internal/models"
	"context"
	"github.com/go-redis/cache/v8"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// Cache provides a connection with redis
type Cache struct {
	rdb *cache.Cache
}

// NewCache is constructor
func NewCache(rdb *cache.Cache) *Cache {
	return &Cache{rdb: rdb}
}

func (c *Cache) CreateCat(cats models.Cat) error {
	ctx := context.TODO()

	if err := c.rdb.Set(&cache.Item{
		Ctx:   ctx,
		Key:   strconv.Itoa(int(cats.ID)),
		Value: cats,
	}); err != nil {
		log.Error(err)
		return err
	}

	log.Println("Create cat in cache")
	return nil
}

func (c *Cache) GetCat(id string) (*models.Cat, error) {
	ctx := context.TODO()

	var cat models.Cat
	err := c.rdb.Get(ctx, id, &cat)

	if err != nil {
		if err.Error() == "cache: key is missing" {
			return nil, err
		} else {
			log.Error(err)
			return nil, err
		}
	}

	log.Println("Get cat in cache")
	return &cat, nil
}

func (c *Cache) DeleteCat(id string) error {
	ctx := context.TODO()

	err := c.rdb.Delete(ctx, id)
	if err != nil {
		return err
	}

	log.Println("Delete cat from cache")
	return nil
}

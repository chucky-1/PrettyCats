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

type Redis interface {
	CreateCat(cats models.Cats) error
	GetCat(id string) (*models.Cats, error)
	DeleteCat(id string) error
}

// Cache provides a connection with redis
type Cache struct {
	rdb *cache.Cache
}

// Stream provides a connection with redis
type Stream struct {
	rdb *redis.Client
}

// NewCache is constructor
func NewCache(rdb *cache.Cache) *Cache {
	return &Cache{rdb: rdb}
}

// NewStream is constructor
func NewStream(rdb *redis.Client) *Stream {
	return &Stream{rdb: rdb}
}

func CacheConnect() (*cache.Cache, error) {
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

func StreamConnect(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
		Password: "", // no password set
		DB:		  0,  // use default DB
	})

	rdb.XGroupCreateMkStream(ctx, "streamCats", "receiver", "0")

	return rdb, nil
}

func (c *Cache) CreateCat(cats models.Cats) error {
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

func (c *Cache) GetCat(id string) (*models.Cats, error) {
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

func (c *Cache) DeleteCat(id string) error {
	ctx := context.TODO()

	err := c.rdb.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s Stream) CreateCat(cats models.Cats) error {
	ctx := context.TODO()

	key := strconv.Itoa(int(cats.ID))
	val := cats.Name

	err := s.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "streamCats",
		Values: []interface{}{key, val},
		ID: key,
	}).Err()

	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s Stream) GetCat(id string) (*models.Cats, error) {
	ctx := context.TODO()

	entries, err := s.rdb.XRangeN(ctx, "streamCats", id, "+", 1).Result()
	if err != nil {
		return nil, err
	}
	//s.rdb.XAck(ctx, "streamCats", "receiver", id)

	cat := new(models.Cats)
	idInt, err := strconv.ParseInt(entries[0].ID[:len(id)], 0, 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	cat.ID = int32(idInt)
	cat.Name = entries[0].Values[id].(string)

	return cat, nil
}

func (s Stream) DeleteCat(id string) error {
	panic("implement me")
}

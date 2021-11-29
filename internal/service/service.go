// Package service contains all business logic. He gets params from handler and sends it in repository
package service

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/repository"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"sync"
)

// Service has methods which get params from handler and send it in repository
type Service interface {
	GetAll() ([]*models.Cat, error)
	Create(cats models.Cat) (*models.Cat, error)
	Get(id string) (*models.Cat, error)
	Update(id string, cats models.Cat) (*models.Cat, error)
	Delete(id string) (*models.Cat, error)
}

// CatService has an interface of repository
type CatService struct {
	repository repository.Repository
	cache repository.Cache
	stream *redis.Client
	memory map[int32]string
	mu sync.Mutex
}

// preNewCatService goes before NewCatService
func preNewCatService(rps repository.Repository, cache repository.Cache, stream *redis.Client) *CatService {
	return &CatService{repository: rps, cache: cache, stream: stream, memory: make(map[int32]string)}
}

// NewCatService is constructor
func NewCatService(ctx context.Context, rps repository.Repository, cache repository.Cache, stream *redis.Client) *CatService {
	srv := preNewCatService(rps, cache, stream)
	go srv.listenStream(ctx)
	return srv
}

func RedisConnect() (*redis.Client, error) {
	hostAndPort := viper.GetString("redis.host") + ":" + viper.GetString("redis.port")
	rdb := redis.NewClient(&redis.Options{
		Addr:	  hostAndPort,
		Password: "", // no password set
		DB:		  0,  // use default DB
	})
	return rdb, nil
}

// GetAll is called by handler and calls func in repository
func (s *CatService) GetAll() ([]*models.Cat, error) {
	return s.repository.GetAll()
}

// Create is called by handler and calls func in repository
func (s *CatService) Create(cats models.Cat) (*models.Cat, error) {
	err := s.write("INSERT", cats)
	if err != nil {
		return nil, err
	}

	err = s.cache.CreateCat(cats)
	if err != nil {
		log.Error(err)
	}

	return s.repository.Create(cats)
}

// Get is called by handler and calls func in repository
func (s *CatService) Get(id string) (*models.Cat, error) {
	// Get cat from local map
	cat, err := s.getFromMemory(id)
	if err == nil {
		return cat, nil
	}

	// Get cat from cache
	cat, err = s.cache.GetCat(id)

	// Get cat from database
	if err != nil {
		cat, err = s.repository.Get(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = s.write("INSERT", *cat)
		if err != nil {
			log.Error(err)
		}
		err = s.cache.CreateCat(*cat)
		if err != nil {
			log.Error(err)
		}
		return cat, nil
	}

	// Return cat from cache but before add in local map
	err = s.write("INSERT", *cat)
	if err != nil {
		log.Error(err)
	}
	return cat, nil
}

// Update is called by handler and calls func in repository
func (s *CatService) Update(id string, cats models.Cat) (*models.Cat, error) {
	// Delete from local map
	idInt, err := strconv.ParseInt(id, 0, 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s.updateMemory("DELETE", int32(idInt), "")

	// Delete from cache
	err = s.cache.DeleteCat(id)
	if err != nil {
		log.Error(err)
	}

	return s.repository.Update(id, cats)
}

// Delete is called by handler and calls func in repository
func (s *CatService) Delete(id string) (*models.Cat, error) {
	// Delete from local map
	idInt, err := strconv.ParseInt(id, 0, 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s.updateMemory("DELETE", int32(idInt), "")

	// Delete from cache
	err = s.cache.DeleteCat(id)
	if err != nil {
		log.Error(err)
	}
	return s.repository.Delete(id)
}

func(s *CatService) listenStream(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			entries, err := s.stream.XRead(ctx, &redis.XReadArgs{
				Streams:  []string{"streamCats", "$"},
				Count:    1,
				Block:    0,
			}).Result()
			if err != nil {
				log.Error(err)
			}

			act, ok := entries[0].Messages[0].Values["act"].(string)
			if ok != true {
				log.Error("Stop listening server")
				return
			}
			id, ok := entries[0].Messages[0].Values["id"].(string)
			if ok != true {
				log.Error("Stop listening server")
				return
			}
			idInt, err := strconv.Atoi(id)
			if err != nil {
				log.Error("Stop listening server")
				return
			}
			name, ok := entries[0].Messages[0].Values["name"].(string)
			if ok != true {
				log.Error("Stop listening server")
				return
			}
			s.updateMemory(act, int32(idInt), name)
		}
	}
}

func (s *CatService) getFromMemory(id string) (*models.Cat, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	idInt, err := strconv.ParseInt(id, 0, 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	name, ok := s.memory[int32(idInt)]
	if ok == true {
		cat := new(models.Cat)
		cat.ID = int32(idInt)
		cat.Name = name
		if err != nil {
			log.Error(err)
			return nil, err
		}
		log.Println("Get сat from local map")
		return cat, nil
	}
	return nil, fmt.Errorf("cat didn't find")
}

func (s *CatService) updateMemory(act string, id int32, name string) {
	log.Println("Update memory")
	s.mu.Lock()
	if act == "INSERT" {
		s.memory[id] = name
	} else if act == "DELETE" {
		delete(s.memory, id)
	}
	s.mu.Unlock()
}

func (s *CatService) write(act string, cats models.Cat) error {
	log.Println("Write down")
	ctx := context.TODO()

	err := s.stream.XAdd(ctx, &redis.XAddArgs{
		Stream: "streamCats",
		Values: []interface{}{"act", act, "id", cats.ID, "name", cats.Name},
	}).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

package service

import (
	"CatsCrud/internal/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type RedisStream struct {
	stream *redis.Client
	memory map[int32]string
}

func NewRedisStream(stream *redis.Client) *RedisStream {
	return &RedisStream{stream: stream, memory: make(map[int32]string)}
}

func RedisConnect() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
		Password: "", // no password set
		DB:		  0,  // use default DB
	})

	return rdb, nil
}

func (s *RedisStream) UpdateMemory(act string, id int32, name string) {
	log.Println("Update memory")
	if act == "INSERT" {
		s.memory[id] = name
	} else if act == "DELETE" {
		delete(s.memory, id)
	}
}

func (s *RedisStream) ListenStream() {
	for {
		ctx := context.TODO()
		entries, err := s.stream.XRead(ctx, &redis.XReadArgs{
			Streams:  []string{"streamCats", "$"},
			Count:    1,
			Block:    0,
		}).Result()
		if err != nil {
			log.Error(err)
		}

		fmt.Println(entries)

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

		s.UpdateMemory(act, int32(idInt), name)
	}
}

func (s *RedisStream) WriteDown(act string, cats models.Cats) error {
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

package main

import (
	"CatsCrud/internal/models"
	rep "CatsCrud/internal/repository"
	"CatsCrud/internal/service"
	"CatsCrud/protocol"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

const (
	flag = 1
	port = ":10000"
)

var srv *service.CatService

type catsCrud struct {
	myGrpc.UnimplementedCatsCrudServer
}

func (s *catsCrud) GetAll(ctx context.Context, r *myGrpc.Request) (*myGrpc.AllCats, error) {
	allCats, err := srv.GetAllCatsServ()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	alc := new(myGrpc.AllCats)
	for _, i := range allCats {
		cts := new(myGrpc.Cats)
		cts.Id = i.ID
		cts.Name = i.Name
		alc.Cats = append(alc.Cats, cts)
	}

	return alc, nil
}

func (s *catsCrud) Create(ctx context.Context, r *myGrpc.RequestCats) (*myGrpc.Cats, error) {
	modCat := new(models.Cats)
	modCat.ID = r.GetId()
	modCat.Name = r.GetName()

	modCat, err := srv.CreateCatsServ(*modCat)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cat := new(myGrpc.Cats)
	cat.Id = modCat.ID
	cat.Name = modCat.Name
	return cat, nil
}

func (s *catsCrud) Get(ctx context.Context, r *myGrpc.Id) (*myGrpc.Cats, error) {
	modCat, err := srv.GetCatServ(r.GetId())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cat := new(myGrpc.Cats)
	cat.Id = modCat.ID
	cat.Name = modCat.Name

	return cat, nil
}

func (s *catsCrud) Update(ctx context.Context, r *myGrpc.RequestCats) (*myGrpc.Cats, error) {
	idStr := strconv.Itoa(int(r.GetId()))
	modCats := new(models.Cats)
	modCats.ID = r.GetId()
	modCats.Name = r.GetName()
	modCats, err := srv.UpdateCatServ(idStr, *modCats)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cat := new(myGrpc.Cats)
	cat.Id = modCats.ID
	cat.Name = modCats.Name

	return cat, nil
}

func (s *catsCrud) Delete(ctx context.Context, r *myGrpc.Id) (*myGrpc.Cats, error) {
	modCat, err := srv.DeleteCatServ(r.GetId())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cat := new(myGrpc.Cats)
	cat.Id = modCat.ID
	cat.Name = modCat.Name

	return cat, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	myGrpc.RegisterCatsCrudServer(s, new(catsCrud))
	fmt.Printf("server listening at %v\n", lis.Addr())

	var rps rep.Repository
	if flag == 1 {
		// Соединение с postgres
		conn, err := rep.RequestDB()
		if err != nil {
			log.Error(err)
			return
		}
		defer conn.Close()

		rps = rep.NewPostgresRepository(conn)
	} else if flag == 2 {
		//Соединение с mongo
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := rep.RequestMongo(ctx)
		if err != nil {
			log.Error(err)
			return
		}
		defer cancel()

		rps = rep.NewMongoRepository(client)
	}
	srv = service.NewCatService(rps)

	if err = s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
		return
	}
}

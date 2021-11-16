package main

import (
	"CatsCrud/internal/models"
	repository2 "CatsCrud/internal/repository"
	"CatsCrud/internal/service"
	"CatsCrud/protocol"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	log "github.com/sirupsen/logrus"
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
		log.Fatal(err)
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
		log.Fatal(err)
	}

	cat := new(myGrpc.Cats)
	cat.Id = modCat.ID
	cat.Name = modCat.Name
	return cat, nil
}

func (s *catsCrud) Get(ctx context.Context, r *myGrpc.Id) (*myGrpc.Cats, error) {
	modCat, err := srv.GetCatServ(r.GetId())
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	cat := new(myGrpc.Cats)
	cat.Id = modCats.ID
	cat.Name = modCats.Name

	return cat, nil
}

func (s *catsCrud) Delete(ctx context.Context, r *myGrpc.Id) (*myGrpc.Cats, error) {
	modCat, err := srv.DeleteCatServ(r.GetId())
	if err != nil {
		log.Fatal(err)
	}

	cat := new(myGrpc.Cats)
	cat.Id = modCat.ID
	cat.Name = modCat.Name

	return cat, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	myGrpc.RegisterCatsCrudServer(s, new(catsCrud))
	log.Printf("server listening at %v", lis.Addr())

	var rps repository2.Repository
	if flag == 1 {
		// Соединение с postgres
		conn := repository2.RequestDB()
		defer conn.Close()

		rps = repository2.NewPostgresRepository(conn)
	} else if flag == 2 {
		//Соединение с mongo
		client, cancel := repository2.RequestMongo()
		defer cancel()

		rps = repository2.NewMongoRepository(client)
	}
	srv = service.NewCatService(rps)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

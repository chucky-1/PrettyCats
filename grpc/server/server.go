// Package server has methods on the server side of grpc
package server

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/service"
	myGrpc "CatsCrud/protocol"

	log "github.com/sirupsen/logrus"

	"context"
	"strconv"
)

// Cats have application's methods in service side
type Cats struct {
	myGrpc.UnimplementedCatsCrudServer
	srv service.Service
}

// NewCats is constructor
func NewCats(srv service.Service) *Cats {
	return &Cats{srv: srv}
}

// GetAll is method of server of grpc
func (s *Cats) GetAll(ctx context.Context, r *myGrpc.Request) (*myGrpc.AllCats, error) {
	allCats, err := s.srv.GetAllCatsServ()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	alc := new(myGrpc.AllCats)
	for _, i := range allCats {
		cts := new(myGrpc.Cat)
		cts.Id = i.ID
		cts.Name = i.Name
		alc.Cat = append(alc.Cat, cts)
	}

	return alc, nil
}

// Create is method of server of grpc
func (s *Cats) Create(ctx context.Context, r *myGrpc.RequestCat) (*myGrpc.Cat, error) {
	modCat := new(models.Cats)
	modCat.ID = r.GetId()
	modCat.Name = r.GetName()

	modCat, err := s.srv.CreateCatsServ(*modCat)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cat := new(myGrpc.Cat)
	cat.Id = modCat.ID
	cat.Name = modCat.Name
	return cat, nil
}

// Get is method of server of grpc
func (s *Cats) Get(ctx context.Context, r *myGrpc.Id) (*myGrpc.Cat, error) {
	modCat, err := s.srv.GetCatServ(r.GetId())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cat := new(myGrpc.Cat)
	cat.Id = modCat.ID
	cat.Name = modCat.Name

	return cat, nil
}

// Update is method of server of grpc
func (s *Cats) Update(ctx context.Context, r *myGrpc.RequestCat) (*myGrpc.Cat, error) {
	idStr := strconv.Itoa(int(r.GetId()))
	modCats := new(models.Cats)
	modCats.ID = r.GetId()
	modCats.Name = r.GetName()
	modCats, err := s.srv.UpdateCatServ(idStr, *modCats)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cat := new(myGrpc.Cat)
	cat.Id = modCats.ID
	cat.Name = modCats.Name

	return cat, nil
}

// Delete is method of server of grpc
func (s *Cats) Delete(ctx context.Context, r *myGrpc.Id) (*myGrpc.Cat, error) {
	modCat, err := s.srv.DeleteCatServ(r.GetId())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cat := new(myGrpc.Cat)
	cat.Id = modCat.ID
	cat.Name = modCat.Name

	return cat, nil
}

package main

import (
	myGrpc "CatsCrud/protocol"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	description string = "Input: 1 - get all cats; 2 - create cats; 3 - get cat\n4 - update cat; 5 - delete cat; exit - to exit"
	hostAndPortGrpc string = "localhost:10000"
)

// getAll gets all cats
func getAll(ctx context.Context, client myGrpc.CatsCrudClient) error {
	r, err := client.GetAll(ctx, &myGrpc.Request{})
	if err != nil {
		log.Error(err)
		return err
	}
	fmt.Println("All cats: ", r.Cat)
	fmt.Println(description)
	return nil
}

// create creates a cat
func create(ctx context.Context, client myGrpc.CatsCrudClient, bitSize int) error {
	fmt.Println("Creating cat...")
	fmt.Print("Input name: ")
	name, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Error(err)
		return err
	}

	fmt.Print("Input ID: ")
	id, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Error(err)
		return err
	}

	idInt, errA := strconv.ParseInt(string(id[:len(id)-2]), 0, bitSize)
	if errA != nil {
		log.Error(errA)
		return err
	}

	cat := new(myGrpc.RequestCat)
	cat.Id = int32(idInt)
	cat.Name = string(name[:len(name)-2])

	r, err := client.Create(ctx, cat)
	if err != nil {
		log.Error(err)
		return err
	}
	fmt.Println("New cat added: ", r)
	fmt.Println(description)

	return nil
}

// get one of the cats
func get(ctx context.Context, client myGrpc.CatsCrudClient) error {
	fmt.Println("Getting cat...")
	fmt.Print("Input id: ")
	id, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Error(err)
		return err
	}

	idStruct := new(myGrpc.Id)
	idStruct.Id = string(id[:len(id)-2])

	r, err := client.Get(ctx, idStruct)
	if err != nil {
		log.Error(err)
		return err
	}

	fmt.Println(r.Id, r.Name)
	fmt.Println(description)
	return nil
}

// update changes one of cats
func update(ctx context.Context, client myGrpc.CatsCrudClient, bitSize int) error {
	fmt.Println("Updating cat...")
	fmt.Print("Input id: ")
	id, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Error(err)
		return err
	}
	idInt, errA := strconv.ParseInt(string(id[:len(id)-2]), 0, bitSize)
	if errA != nil {
		log.Error(errA)
		return err
	}
	fmt.Print("Input name: ")
	name, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Error(err)
		return err
	}

	cat := new(myGrpc.RequestCat)
	cat.Id = int32(idInt)
	cat.Name = string(name[:len(name)-2])

	r, errU := client.Update(ctx, cat)
	if errU != nil {
		log.Error(errU)
		return err
	}

	fmt.Println("Cat updated!", r.Id, r.Name)
	fmt.Println(description)
	return nil
}

// delete deletes a cat by ID
func delete(ctx context.Context, client myGrpc.CatsCrudClient) error {
	fmt.Println("Delete cat...")
	fmt.Print("Input id, please: ")
	id, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Error(err)
		return err
	}
	idStruct := new(myGrpc.Id)
	idStruct.Id = string(id[:len(id)-2])

	r, err := client.Delete(ctx, idStruct)
	if err != nil {
		log.Error(err)
		return err
	}
	fmt.Println(r.Name, " deleted")
	fmt.Println(description)
	return nil
}

// consoleApp is an application that supports some operations with cats
func consoleApp(ctx context.Context, client myGrpc.CatsCrudClient) error {
	input := []byte{1, 2}
	var bitSize = 32
	for string(input[:len(input)-2]) != "exit" {
		input, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			log.Error(err)
			return err
		}
		switch {
		case string(input[:len(input)-2]) == "1":
			err = getAll(ctx, client)
			if err != nil {
				return err
			}
		case string(input[:len(input)-2]) == "2":
			err = create(ctx, client, bitSize)
			if err != nil {
				return nil
			}
		case string(input[:len(input)-2]) == "3":
			err = get(ctx, client)
			if err != nil {
				return err
			}
		case string(input[:len(input)-2]) == "4":
			err = update(ctx, client, bitSize)
			if err != nil {
				return err
			}
		case string(input[:len(input)-2]) == "5":
			err = delete(ctx, client)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	conn, err := grpc.Dial(hostAndPortGrpc, grpc.WithInsecure())
	if err != nil {
		log.Errorf("fail to dial: %v", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(conn)
	client := myGrpc.NewCatsCrudClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	fmt.Println(description)

	err = consoleApp(ctx, client)
	if err != nil {
		return
	}
}

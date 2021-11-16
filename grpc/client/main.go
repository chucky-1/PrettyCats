package main

import (
	myGrpc "CatsCrud/protocol"
	"bufio"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"time"
)

const (
	adress = "localhost:10000"
)

var description string = "Input\n1 - get all cats\n2 - create cats\n3 - get cat\n4 - update cat\n5 - delete cat\nexit - to exit"

func main()  {
	conn, err := grpc.Dial(adress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := myGrpc.NewCatsCrudClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	fmt.Println(description)

	input := []byte{1,2}
	for string(input[:len(input)-2]) != "exit" {
		input, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case string(input[:len(input)-2]) == "1":
			r, err := client.GetAll(ctx, &myGrpc.Request{})
			if err != nil {
				log.Fatal(err, r)
			}
			fmt.Print("All cats: ", r.Cats, "\n")
		case string(input[:len(input)-2]) == "2":
			fmt.Println("Creating cat...")
			fmt.Print("Input name: ")
			name, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("Input ID: ")
			id, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}

			idInt, err := strconv.Atoi(string(id[:len(id)-2]))
			if err != nil {
				log.Fatal(err)
			}

			cat := new(myGrpc.RequestCats)
			cat.Id = int32(idInt)
			cat.Name = string(name[:len(name)-2])

			r, err := client.Create(ctx, cat)
			fmt.Println("New cat added: ", r)
		case string(input[:len(input)-2]) == "3":
			fmt.Println("Getting cat...")
			fmt.Print("Input id: ")
			id, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}

			idStruct := new(myGrpc.Id)
			idStruct.Id = string(id[:len(id)-2])

			r, err := client.Get(ctx, idStruct)
			fmt.Println("Gotten cat is", r.Id, r.Name)
		case string(input[:len(input)-2]) == "4":
			fmt.Println("Updating cat...")
			fmt.Print("Input id: ")
			id, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			idInt, err := strconv.Atoi(string(id[:len(id)-2]))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print("Input name: ")
			name, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}

			cat := new(myGrpc.RequestCats)
			cat.Id = int32(idInt)
			cat.Name = string(name[:len(name)-2])

			r, err := client.Update(ctx, cat)
			fmt.Println("Updated cat is", r.Id, r.Name)
		case string(input[:len(input)-2]) == "5":
			fmt.Println("Delete cat...")
			fmt.Print("Input id: ")
			id, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			idStruct := new(myGrpc.Id)
			idStruct.Id = string(id[:len(id)-2])

			r, err := client.Delete(ctx, idStruct)
			fmt.Println("Deleted cat is", r.Id, r.Name)
		}
	}
}

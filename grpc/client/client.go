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
	description string = "Input\n1 - get all cats\n2 - create cats\n3 - get cat\n4 - update cat\n5 - delete cat\nexit - to exit"
)

func main() {
	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	fmt.Println(description)

	input := []byte{1, 2}
	var bitSize = 32
	for string(input[:len(input)-2]) != "exit" {
		input, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			log.Error(err)
			return
		}
		switch {
		case string(input[:len(input)-2]) == "1":
			r, errGA := client.GetAll(ctx, &myGrpc.Request{})
			if errGA != nil {
				log.Error(errGA)
				return
			}
			fmt.Print("All cats: ", r.Cat, "\n")
		case string(input[:len(input)-2]) == "2":
			fmt.Println("Creating cat...")
			fmt.Print("Input name: ")
			name, errNR := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Error(errNR)
				return
			}

			fmt.Print("Input ID: ")
			id, errNR := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Error(errNR)
				return
			}

			idInt, errA := strconv.ParseInt(string(id[:len(id)-2]), 0, bitSize)
			if errA != nil {
				log.Error(errA)
				return
			}

			cat := new(myGrpc.RequestCat)
			cat.Id = int32(idInt)
			cat.Name = string(name[:len(name)-2])

			r, errC := client.Create(ctx, cat)
			if errC != nil {
				log.Error(errC)
				return
			}
			fmt.Println("New cat added: ", r)
		case string(input[:len(input)-2]) == "3":
			fmt.Println("Getting cat...")
			fmt.Print("Input id: ")
			id, errNR := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if errNR != nil {
				log.Error(errNR)
				return
			}

			idStruct := new(myGrpc.Id)
			idStruct.Id = string(id[:len(id)-2])

			r, errG := client.Get(ctx, idStruct)
			if errG != nil {
				log.Error(errG)
				return
			}
			fmt.Println("Gotten cat is", r.Id, r.Name)
		case string(input[:len(input)-2]) == "4":
			fmt.Println("Updating cat...")
			fmt.Print("Input id: ")
			id, errNR := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if errNR != nil {
				log.Error(errNR)
				return
			}
			idInt, errA := strconv.ParseInt(string(id[:len(id)-2]), 0, bitSize)
			if errA != nil {
				log.Error(errA)
				return
			}
			fmt.Print("Input name: ")
			name, errNR := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Error(errNR)
				return
			}

			cat := new(myGrpc.RequestCat)
			cat.Id = int32(idInt)
			cat.Name = string(name[:len(name)-2])

			r, errU := client.Update(ctx, cat)
			if errU != nil {
				log.Error(errU)
				return
			}
			fmt.Println("Updated cat is", r.Id, r.Name)
		case string(input[:len(input)-2]) == "5":
			fmt.Println("Delete cat...")
			fmt.Print("Input id: ")
			id, errNR := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if errNR != nil {
				log.Error(errNR)
				return
			}
			idStruct := new(myGrpc.Id)
			idStruct.Id = string(id[:len(id)-2])

			r, errD := client.Delete(ctx, idStruct)
			if errD != nil {
				log.Error(errD)
				return
			}
			fmt.Println("Deleted cat is", r.Id, r.Name)
		}
	}
}

package main

import (
	"net/http"
	"basa/config"
	"basa/controller"
	"basa/storage/postgres"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()

	store, err := postgres.New(cfg)
	if err != nil {
		log.Fatalln("error while connecting to db err: ", err.Error())
		return
	}
	defer store.DB.Close()

	con := controller.New(store)


	http.HandleFunc("/car", con.OrderProduct)
    fmt.Println("listening at: ")
	http.ListenAndServe(":8080", nil)



		}


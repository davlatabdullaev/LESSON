package main

import (
	"context"
	"log"
	"test/api"
	"test/config"
	"test/storage/postgres"
)

func main() {
	cfg := config.Load()

	pgStore, err := postgres.New(context.Background(), cfg)
	if err != nil {
		log.Fatalln("error while connecting to db err:", err.Error())
		return
	}
	defer pgStore.Close()

	server := api.New(pgStore)

	if err = server.Run("localhost:8080"); err != nil {
		panic(err)
	}
}

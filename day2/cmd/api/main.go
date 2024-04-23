package main

import (
	"github.com/narekps/go2/day2/internal/app/api"
	"log"
)

func main() {
	log.Println("It's work.")

	config := api.NewConfig()
	server := api.New(config)

	//api server start
	log.Fatal(server.Start())
}

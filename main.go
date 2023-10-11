package main

import (
	"log"
	"storage-api/config"
	"storage-api/server"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting Storage API")
	log.Println("Initializing configuration")
	config := config.InitConfig("storage")
	log.Println("Starting Database")
	dbHandler := server.InitDatabase(config)
	log.Println("Initializing HTTP server")
	httpServer := server.InitHttpServer(config, dbHandler)
	httpServer.Start()
}

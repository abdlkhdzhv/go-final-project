package main

import (
	"go-final-project/pkg/api"
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"log"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("failed to init db: %v", err) // остановка программы
	}

	api.Init()
	server.Run()
}

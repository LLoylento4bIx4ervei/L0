package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LLoylento4bIx4ervei/L0/storage"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error load .env")
	}

	store := storage.NewStorage()

	if err := store.Open(); err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	defer store.Close()

	if err := store.LoadCache(); err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}

	log.Println("Successfully connect to db")

	server := NewServer(store)

	go func() {
		if err := server.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server:%v", err)
		}
	}()

	log.Println("Server starting")
	log.Println("Ctrl+C to stop")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server")

}

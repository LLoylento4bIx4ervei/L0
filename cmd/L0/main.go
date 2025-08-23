package main

import (
	"log"

	"github.com/LLoylento4bIx4ervei/L0/storage"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error load .env")
	}

	store := &storage.Storage{}

	if err := store.Open(); err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	defer store.Close()

	log.Println("Successfully connect to db")
}

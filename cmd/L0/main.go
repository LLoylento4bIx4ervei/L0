package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LLoylento4bIx4ervei/L0/storage"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Note: .env file not found, using environment variables")
	}

	store := storage.NewStorage()

	if err := store.Open(); err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	defer store.Close()

	if err := store.LoadCache(); err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}

	log.Println("Successfully connected to db and loaded cache")

	brokers := []string{"localhost:9092"}
	topic := "orders"
	groupID := "order-service"

	kafkaConsumer := NewKafkaConsumer(brokers, topic, groupID, store)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafkaConsumer.Start(ctx)

	server := NewServer(store)

	go func() {
		if err := server.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started on :8080")
	log.Println("Kafka consumer started")
	log.Println("Kafka UI available at: http://localhost:8081")
	log.Println("Web interface available at: http://localhost:8080")
	log.Println("Press Ctrl+C to stop")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()

	time.Sleep(1 * time.Second)
	if err := kafkaConsumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}

	log.Println("Server stopped")
}

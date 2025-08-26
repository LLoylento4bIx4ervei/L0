package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/LLoylento4bIx4ervei/L0/storage"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	storage *storage.Storage
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, storage *storage.Storage) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  time.Second,
	})

	return &KafkaConsumer{
		reader:  reader,
		storage: storage,
	}
}

func (kc *KafkaConsumer) Start(ctx context.Context) {
	log.Println("Starting Kafka consumer...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer...")
			return
		default:
			msg, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			log.Printf("Received message from Kafka: topic=%s partition=%d offset=%d\n",
				msg.Topic, msg.Partition, msg.Offset)

			var order storage.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Error unmarshaling order: %v", err)
				log.Printf("Raw message: %s", string(msg.Value))
				continue
			}

			if order.OrderUID == "" {
				log.Printf("Invalid order: missing order_uid")
				continue
			}

			if err := kc.storage.Save(&order); err != nil {
				log.Printf("Error saving order: %v", err)
				continue
			}

			log.Printf("Order %s saved from Kafka", order.OrderUID)
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}

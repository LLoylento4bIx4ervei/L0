package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
	})

	defer writer.Close()

	for i := 0; i < 5; i++ {
		order := generateTestOrder()
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Fatal("Error marshaling order:", err)
		}

		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: orderJSON,
			},
		)
		if err != nil {
			log.Fatal("Error writing message:", err)
		}

		log.Printf("Sent order: %s", order)
		time.Sleep(2 * time.Second)
	}
}

func generateTestOrder() map[string]interface{} {
	rand.Seed(time.Now().UnixNano())
	orderUID := generateRandomID(10)

	return map[string]interface{}{
		"order_uid":    orderUID,
		"track_number": "WBILMTESTTRACK",
		"entry":        "WBIL",
		"delivery": map[string]interface{}{
			"name":    "Test Testov",
			"phone":   "+9720000000",
			"zip":     "2639809",
			"city":    "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region":  "Kraiot",
			"email":   "test@gmail.com",
		},
		"payment": map[string]interface{}{
			"transaction":   orderUID,
			"request_id":    "",
			"currency":      "USD",
			"provider":      "wbpay",
			"amount":        1817,
			"payment_dt":    1637907727,
			"bank":          "alpha",
			"delivery_cost": 1500,
			"goods_total":   317,
			"custom_fee":    0,
		},
		"items": []map[string]interface{}{
			{
				"chrt_id":      9934930,
				"track_number": "WBILMTESTTRACK",
				"price":        453,
				"rid":          "ab4219087a764ae0btest",
				"name":         "Mascaras",
				"sale":         30,
				"size":         "0",
				"total_price":  317,
				"nm_id":        2389212,
				"brand":        "Vivienne Sabo",
				"status":       202,
			},
		},
		"locale":             "en",
		"internal_signature": "",
		"customer_id":        "test",
		"delivery_service":   "meest",
		"shardkey":           "9",
		"sm_id":              99,
		"date_created":       time.Now().Format(time.RFC3339),
		"oof_shard":          "1",
	}
}

func generateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

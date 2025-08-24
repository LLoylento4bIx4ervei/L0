package main

import (
	"log"
	"time"

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

	testOrder := &storage.Order{
		OrderUID:    "test123",
		TrackNumber: "TRACK123",
		Entry:       "WBIL",
		Delivery: storage.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: storage.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []storage.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	if err := store.Save(testOrder); err != nil {
		log.Fatalf("Failed to save order: %v", err)
	}
	log.Println("Test order saved successfully")

	retrievedOrder, err := store.GetOrderByUID("test123")
	if err != nil {
		log.Fatalf("Failed to get order: %v", err)
	}
	log.Printf("Retrieved order: %+v", retrievedOrder)

	log.Println("Database operations completed successfully")

}

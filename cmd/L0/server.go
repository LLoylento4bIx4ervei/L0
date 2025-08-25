package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LLoylento4bIx4ervei/L0/storage"
)

type Server struct {
	storage *storage.Storage
}

func NewServer(storage *storage.Storage) *Server {
	return &Server{storage: storage}
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Path[len("/order/"):]
	if orderID == "" {

		return
	}

	log.Printf("Request for orders: %s", orderID)

	order, err := s.storage.GetOrderCache(orderID)
	if err != nil {
		order, err = s.storage.GetOrderByUID(orderID)
		if err != nil {
			return
		}

		s.storage.UpdateCache(order)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (s *Server) handleGetAll(w http.ResponseWriter, r *http.Request) {
	orders, err := s.storage.AllOrders()
	if err != nil {

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/order/", s.handleGet)
	http.HandleFunc("/orders/", s.handleGetAll)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

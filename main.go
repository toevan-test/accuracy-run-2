// Command accuracyrun is a small order service used by the Toevan report accuracy test.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/accuracylab/internal/orders"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	store := orders.NewStore()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /orders", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, store.List())
	})
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var in orders.Order
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&in); err != nil {
			http.Error(w, "invalid order", http.StatusBadRequest)
			return
		}
		created, err := store.Add(in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	})
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}

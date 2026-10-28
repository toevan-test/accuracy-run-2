// Package orders keeps orders in memory.
package orders

import (
	"errors"
	"strings"
	"sync"
)

// Order is one customer order.
type Order struct {
	ID       int    `json:"id"`
	Customer string `json:"customer"`
	Quantity int    `json:"quantity"`
}

// ErrInvalid is returned for an order without a customer or with a quantity outside 1 to 1000.
var ErrInvalid = errors.New("order needs a customer and a quantity from 1 to 1000")

// Store is a concurrency-safe in-memory list of orders.
type Store struct {
	mu     sync.Mutex
	nextID int
	orders []Order
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{nextID: 1}
}

// Add validates the order, gives it the next id and keeps it.
func (s *Store) Add(o Order) (Order, error) {
	o.Customer = strings.TrimSpace(o.Customer)
	if o.Customer == "" || o.Quantity < 1 || o.Quantity > 1000 {
		return Order{}, ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	o.ID = s.nextID
	s.nextID++
	s.orders = append(s.orders, o)
	return o, nil
}

// List returns a copy of every order, oldest first.
func (s *Store) List() []Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Order(nil), s.orders...)
}

package orders

import (
	"errors"
	"testing"
)

func TestAddAssignsIDsInOrder(t *testing.T) {
	s := NewStore()
	a, err := s.Add(Order{Customer: "ada", Quantity: 2})
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.Add(Order{Customer: "bob", Quantity: 1})
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != 1 || b.ID != 2 {
		t.Fatalf("ids = %d, %d; want 1, 2", a.ID, b.ID)
	}
	if got := len(s.List()); got != 2 {
		t.Fatalf("len(List()) = %d; want 2", got)
	}
}

func TestAddRejectsInvalid(t *testing.T) {
	s := NewStore()
	for _, o := range []Order{{Customer: " ", Quantity: 1}, {Customer: "ada", Quantity: 0}, {Customer: "ada", Quantity: 1001}} {
		if _, err := s.Add(o); !errors.Is(err, ErrInvalid) {
			t.Fatalf("Add(%+v) error = %v; want ErrInvalid", o, err)
		}
	}
}

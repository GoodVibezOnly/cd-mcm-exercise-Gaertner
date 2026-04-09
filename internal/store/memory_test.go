package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Widget", Price: 9.99})
	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != created.ID || got.Name != created.Name || got.Price != created.Price {
		t.Errorf("expected %+v, got %+v", created, got)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

// TODO: Add tests for Update, Delete of existing product, and GetByID with invalid ID

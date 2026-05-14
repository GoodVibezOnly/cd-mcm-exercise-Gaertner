package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Widget", Price: 9.99})
	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}
	got, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Widget" {
		t.Errorf("expected name Widget, got %s", got.Name)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestGetAllWithProducts(t *testing.T) {
	s := NewMemoryStore()
	s.Create(model.Product{Name: "Widget", Price: 9.99})
	s.Create(model.Product{Name: "Gadget", Price: 19.99})
	products := s.GetAll()
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	s := NewMemoryStore()
	s.Create(model.Product{Name: "Widget", Price: 9.99})
	updated, err := s.Update(1, model.Product{Name: "Updated", Price: 15.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("expected name Updated, got %s", updated.Name)
	}
	if updated.ID != 1 {
		t.Errorf("expected ID 1, got %d", updated.ID)
	}
}

func TestUpdateNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.Update(999, model.Product{Name: "Updated", Price: 15.0})
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestDelete(t *testing.T) {
	s := NewMemoryStore()
	s.Create(model.Product{Name: "Widget", Price: 9.99})
	if err := s.Delete(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := s.GetByID(1)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound after delete")
	}
}

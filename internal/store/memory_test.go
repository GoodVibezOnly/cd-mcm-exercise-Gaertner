package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	store := NewMemoryStore()
	created := model.Product{ID: 1, Name: "maus", Price: 10.50}
	store.Create(created)
	prod, err := store.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if prod.ID != created.ID || prod.Name != created.Name || prod.Price != created.Price {
		t.Errorf("expected %+v, got %+v", created, prod)
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

func TestUpdateProduct(t *testing.T) {
	store := NewMemoryStore()
	created := model.Product{ID: 1, Name: "maus", Price: 10.50}
	store.Create(created)
	updated := model.Product{Name: "maus updated", Price: 12.00}
	store.Update(1, updated)
	prod, err := store.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if prod.ID != created.ID || prod.Name != updated.Name || prod.Price != updated.Price {
		t.Errorf("expected %+v, got %+v", updated, prod)
	}
}

func TestDeleteProduct(t *testing.T) {
	store := NewMemoryStore()
	created := model.Product{ID: 1, Name: "maus", Price: 10.50}
	store.Create(created)
	store.Delete(1)
	prod, err := store.GetByID(1)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if prod != (model.Product{}) {
		t.Errorf("expected empty product, got %+v", prod)
	}
}
func TestGetByIDNotFound(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{"non-existent ID", 999},
		{"zero ID", 0},
		{"negative ID", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore()
			created := model.Product{ID: 1, Name: "maus", Price: 10.50}
			store.Create(created)
			prod, err := store.GetByID(tt.id)

			if err != ErrNotFound {
				t.Fatalf("expected ErrNotFound, got %v", err)
			}
			if prod != (model.Product{}) {
				t.Errorf("expected empty product, got %+v", prod)
			}
		})
	}
}

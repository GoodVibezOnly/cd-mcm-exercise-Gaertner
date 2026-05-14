package store

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func newMockStore(t *testing.T) (*PostgresStore, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return &PostgresStore{DB: db}, mock
}

func TestPostgresEnsureTable(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS products").WillReturnResult(sqlmock.NewResult(0, 0))
	if err := s.EnsureTable(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPostgresGetAllWithRows(t *testing.T) {
	s, mock := newMockStore(t)
	mockRows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(mockRows)
	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 1 || products[0].Name != "Widget" {
		t.Errorf("unexpected products: %v", products)
	}
}

func TestPostgresGetAllEmpty(t *testing.T) {
	s, mock := newMockStore(t)
	mockRows := sqlmock.NewRows([]string{"id", "name", "price"})
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(mockRows)
	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestPostgresGetByID(t *testing.T) {
	s, mock := newMockStore(t)
	mockRows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").WillReturnRows(mockRows)
	p, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Widget" {
		t.Errorf("expected Widget, got %s", p.Name)
	}
}

func TestPostgresGetByIDNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").WillReturnError(sql.ErrNoRows)
	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresCreate(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectQuery("INSERT INTO products").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	p, err := s.Create(model.Product{Name: "Widget", Price: 9.99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}
}

func TestPostgresUpdate(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(0, 1))
	p, err := s.Update(1, model.Product{Name: "Updated", Price: 15.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}
}

func TestPostgresUpdateNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(0, 0))
	_, err := s.Update(999, model.Product{Name: "Updated", Price: 15.0})
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresDelete(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(0, 1))
	if err := s.Delete(1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(0, 0))
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

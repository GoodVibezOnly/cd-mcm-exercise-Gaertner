package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func setupPostgresRouter(t *testing.T) (*mux.Router, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	s := &store.PostgresStore{DB: db}
	h := NewPostgresHandler(s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r, mock
}

func TestPostgresHealthOK(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectPing()
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresHealthDBError(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectPing().WillReturnError(errors.New("connection refused"))
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rr.Code)
	}
}

func TestPostgresGetProducts(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mockRows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(mockRows)
	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProductsError(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnError(errors.New("db error"))
	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestPostgresGetProduct(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mockRows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").WillReturnRows(mockRows)
	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProductNotFound(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").WillReturnError(sql.ErrNoRows)
	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresCreateProduct(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectQuery("INSERT INTO products").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestPostgresCreateProductInvalidJSON(t *testing.T) {
	r, _ := setupPostgresRouter(t)
	req := httptest.NewRequest("POST", "/products", strings.NewReader("not-json"))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresCreateProductInvalidProduct(t *testing.T) {
	r, _ := setupPostgresRouter(t)
	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"","price":9.99}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresCreateProductDBError(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectQuery("INSERT INTO products").WillReturnError(errors.New("db error"))
	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestPostgresUpdateProduct(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(0, 1))
	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{"name":"Updated","price":15.0}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresUpdateProductInvalidJSON(t *testing.T) {
	r, _ := setupPostgresRouter(t)
	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader("bad-json"))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresUpdateProductNotFound(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(0, 0))
	req := httptest.NewRequest("PUT", "/products/999", strings.NewReader(`{"name":"Updated","price":15.0}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresDeleteProduct(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(0, 1))
	req := httptest.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresDeleteProductNotFound(t *testing.T) {
	r, mock := setupPostgresRouter(t)
	mock.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(0, 0))
	req := httptest.NewRequest("DELETE", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

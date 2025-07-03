package controller

import (
	"RestGoTest/httpserver/service"
	"database/sql"
	"net/http"

	_ "modernc.org/sqlite"
)

func AllProductsController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.AllProducts(w, r, db)
	}
}

func GetProductController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.FetchProduct(w, r, db)
	}
}

func CreateProductController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.CreateProduct(w, r, db)
	}
}

func DeleteProductController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.DeleteProduct(w, r, db)
	}
}

func DeleteAllProductsController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.DeleteAllProducts(w, r, db)
	}
}

func UpdateProductController(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.UpdateProduct(w, r, db)
	}
}

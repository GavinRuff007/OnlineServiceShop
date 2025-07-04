package controller

import (
	"RestGoTest/httpserver/service"
	"net/http"

	_ "modernc.org/sqlite"
)

func AllProductsController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.AllProducts(w, r)
	}
}

func GetProductController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.FetchProduct(w, r)
	}
}

func CreateProductController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.CreateProduct(w, r)
	}
}

func DeleteProductController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.DeleteProduct(w, r)
	}
}

func DeleteAllProductsController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.DeleteAllProducts(w, r)
	}
}

func UpdateProductController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.UpdateProduct(w, r)
	}
}

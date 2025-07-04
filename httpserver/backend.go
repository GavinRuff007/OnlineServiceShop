package httpserver

import (
	"RestGoTest/httpserver/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

type App struct {
	Port   string
	Router *mux.Router
}

func (a *App) Init() {

	a.Router = mux.NewRouter()
	a.initalizeRoutes()
}

func (a *App) initalizeRoutes() {

	//Create Product
	a.Router.HandleFunc("/createProduct", controller.CreateProductController()).Methods("POST")

	//Read Product
	a.Router.HandleFunc("/products", controller.AllProductsController()).Methods("GET")
	a.Router.HandleFunc("/product/{id}", controller.GetProductController()).Methods("GET")

	//Update Product
	a.Router.HandleFunc("/update", controller.UpdateProductController()).Methods("PUT")

	//Delete Product
	a.Router.HandleFunc("/delete/{id}", controller.DeleteProductController()).Methods("DELETE")
	a.Router.HandleFunc("/deleteAll", controller.DeleteAllProductsController()).Methods("DELETE")

}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

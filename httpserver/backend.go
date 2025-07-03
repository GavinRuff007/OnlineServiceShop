package httpserver

import (
	"RestGoTest/httpserver/controller"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

type App struct {
	DB     *sql.DB
	Port   string
	Router *mux.Router
}

func (a *App) Init() {

	DB, err := sql.Open("sqlite", "product.db")
	if err != nil {
		log.Fatal(err)
	}
	a.DB = DB

	a.Router = mux.NewRouter()
	a.initalizeRoutes()
}

func (a *App) initalizeRoutes() {

	//Create Product
	a.Router.HandleFunc("/createProduct", controller.CreateProductController(a.DB)).Methods("POST")

	//Read Product
	a.Router.HandleFunc("/products", controller.AllProductsController(a.DB)).Methods("GET")
	a.Router.HandleFunc("/product/{id}", controller.GetProductController(a.DB)).Methods("GET")

	//Update Product
	a.Router.HandleFunc("/update", controller.UpdateProductController(a.DB)).Methods("PUT")

	//Delete Product
	a.Router.HandleFunc("/delete/{id}", controller.DeleteProductController(a.DB)).Methods("DELETE")
	a.Router.HandleFunc("/deleteAll", controller.DeleteAllProductsController(a.DB)).Methods("DELETE")

}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

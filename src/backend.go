package httpserver

import (
	"RestGoTest/src/controller"
	"RestGoTest/src/middleware"
	"log"
	"net/http"
	"time"

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

	a.Router.Handle("/createProduct", middleware.ContextAbortMiddleware(controller.CreateProductController())).Methods("POST")
	a.Router.Handle("/products", middleware.ContextAbortMiddleware(controller.AllProductsController())).Methods("GET")
	a.Router.Handle("/product/{id}", middleware.ContextAbortMiddleware(controller.GetProductController())).Methods("GET")
	a.Router.Handle("/update", middleware.ContextAbortMiddleware(controller.UpdateProductController())).Methods("PUT")
	a.Router.Handle("/delete/{id}", middleware.ContextDelayAbortMiddleware(controller.DeleteProductController())).Methods("DELETE")
	a.Router.Handle("/deleteAll", middleware.ContextDelayAbortMiddleware(controller.DeleteAllProductsController())).Methods("DELETE")

	a.Router.Use(middleware.TimeoutMiddleware(7 * time.Second))
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

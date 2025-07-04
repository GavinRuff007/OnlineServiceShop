package httpserver

import (
	"RestGoTest/src/controller"
	"RestGoTest/src/middleware"
	"log"
	"net/http"
	"time"

	_ "RestGoTest/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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

	a.Router.Handle("/products", middleware.ContextAbortMiddleware(controller.AllProductsController())).Methods("GET")
	a.Router.Handle("/products/{id}", middleware.ContextAbortMiddleware(controller.GetProductController())).Methods("GET")
	a.Router.Handle("/products", middleware.ContextAbortMiddleware(controller.CreateProductController())).Methods("POST")
	a.Router.Handle("/products", middleware.ContextAbortMiddleware(controller.UpdateProductController())).Methods("PUT")
	a.Router.Handle("/products/{id}", middleware.ContextDelayAbortMiddleware(controller.DeleteProductController())).Methods("DELETE")
	a.Router.Handle("/products", middleware.ContextDelayAbortMiddleware(controller.DeleteAllProductsController())).Methods("DELETE")

	a.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	a.Router.Use(middleware.TimeoutMiddleware(7 * time.Second))
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

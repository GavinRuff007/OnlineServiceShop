package httpserver

import (
	"RestGoTest/src/controller"
	"RestGoTest/src/middleware"
	"log"
	"net/http"
	"time"

	_ "RestGoTest/docs"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @Summary      تست سرویس
// @Description  این یک سرویس تست است
// @Tags         Test
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/v1/test [get]
func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Working!"})
}

type App struct {
	Port   string
	Router *mux.Router
}

func (a *App) Init() {
	a.Router = mux.NewRouter()

	a.InitializeGinService()
	a.InitializeHttpService()
}

func (a *App) InitializeHttpService() {

	a.Router.Use(middleware.EnableCORS)

	a.Router.Handle("/products", middleware.ContextAbortMiddleware(controller.AllProductsController())).Methods("GET")
	a.Router.Handle("/products/{id}", middleware.ContextAbortMiddleware(controller.GetProductController())).Methods("GET")
	a.Router.Handle("/products", middleware.ContextAbortMiddleware(controller.CreateProductController())).Methods("POST")
	a.Router.Handle("/products", middleware.ContextAbortMiddleware(controller.UpdateProductController())).Methods("PUT")
	a.Router.Handle("/products/{id}", middleware.ContextDelayAbortMiddleware(controller.DeleteProductController())).Methods("DELETE")
	a.Router.Handle("/products", middleware.ContextDelayAbortMiddleware(controller.DeleteAllProductsController())).Methods("DELETE")

	a.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	a.Router.Use(middleware.TimeoutMiddleware(7 * time.Second))
}

func (a *App) InitializeGinService() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	v1 := r.Group("/api/v1/")
	{
		v1.GET("/test", TestHandler)
	}
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))
	a.Router.PathPrefix("/api/v1/").Handler(r)
}

func (a *App) Run() {

	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

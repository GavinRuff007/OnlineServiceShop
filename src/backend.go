package httpserver

import (
	"RestGoTest/src/config"
	"RestGoTest/src/router"

	"RestGoTest/src/middleware"

	"log"
	"net/http"

	_ "RestGoTest/docs"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	Port   string
	Router *mux.Router
}

func (a *App) Init(cfg *config.Config) {
	a.Router = mux.NewRouter()

	a.InitializeGinService(cfg)
}

func (a *App) InitializeGinService(cfg *config.Config) {
	r := gin.New()
	r.Use(middleware.DefaultStructuredLogger(cfg))
	r.Use(middleware.LimitByRequest())
	r.Use(gin.Logger(), gin.Recovery())
	v1 := r.Group("/api/v1/")
	{
		health := v1.Group("/health")
		router.Health(health)
	}
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))
	a.Router.PathPrefix("/api/v1/").Handler(r)
}

func (a *App) Run() {

	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

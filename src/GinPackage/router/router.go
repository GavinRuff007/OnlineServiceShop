package router

import (
	"RestGoTest/src/GinPackage/handler"

	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup) {
	hand := handler.NewHealth()
	r.GET("/health", hand.HealthTest)
}

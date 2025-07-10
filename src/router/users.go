package router

import (
	"RestGoTest/src/config"
	"RestGoTest/src/handler"
	"RestGoTest/src/middleware"

	"github.com/gin-gonic/gin"
)

func User(r *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewUserHandler(cfg)
	r.POST("/send-otp", middleware.OtpLimiter(cfg), h.SendOtp)
}

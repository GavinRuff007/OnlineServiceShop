package router

import (
	"RestGoTest/src/config"
	"RestGoTest/src/handler"
	"RestGoTest/src/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewUserHandler(cfg)

	router.POST("/send-otp", middleware.OtpLimiter(cfg), h.SendOtp)
	router.POST("/login-by-username", h.LoginByUsername)
	router.POST("/register-by-username", h.RegisterByUsername)
	router.POST("/login-by-mobile", h.RegisterLoginByMobileNumber)
	router.POST("/refresh-token", h.RefreshToken)
}

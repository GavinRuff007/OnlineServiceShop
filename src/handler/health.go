package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHealth struct{}

func NewHealth() *TestHealth {
	return &TestHealth{}
}

// @Summary      تست سرویس
// @Description  این یک سرویس تست است
// @Tags         Test
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/v1/health/health [get]
func (h *TestHealth) HealthTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Working!",
	})
}

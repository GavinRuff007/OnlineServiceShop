package handler

import (
	"RestGoTest/src/config"
	"RestGoTest/src/dependency"
	"RestGoTest/src/dto"
	"RestGoTest/src/helper"
	"RestGoTest/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	orderUsecase *services.OrderService
	config       *config.Config
}

// ساخت هندلر و وابستگی‌ها
func NewOrdersHandler(cfg *config.Config) *OrdersHandler {
	orderRepo := dependency.GetOrderRepository(cfg)
	orderService := services.NewOrderService(cfg, orderRepo)
	return &OrdersHandler{orderUsecase: orderService, config: cfg}
}

// @Security BearerAuth
// CreateOrder godoc
// @Summary ایجاد سفارش جدید
// @Description ایجاد یک سفارش برای کاربر
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param Request body dto.OrderCreateDTO true "OrderCreateRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Validation Failed"
// @Failure 500 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/orders [post]
func (h *OrdersHandler) CreateOrder(c *gin.Context) {
	req := new(dto.OrderCreateDTO)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	order, err := h.orderUsecase.CreateOrder(c, *req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(order, true, helper.Success))
}

// @Security BearerAuth
// GetOrderByID godoc
// @Summary دریافت سفارش
// @Description دریافت اطلاعات یک سفارش با شناسه
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param Request body dto.OrderGetByIDDTO true "OrderGetRequest"
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Validation Failed"
// @Failure 404 {object} helper.BaseHttpResponse "Not Found"
// @Router /v1/orders/get [post]
func (h *OrdersHandler) GetOrderByID(c *gin.Context) {
	req := new(dto.OrderGetByIDDTO)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	order, err := h.orderUsecase.GetOrderByID(c, *req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.NotFoundError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(order, true, helper.Success))
}

// @Security BearerAuth
// GetOrdersByUser godoc
// @Summary دریافت همه سفارشات کاربر
// @Description دریافت لیست سفارشات یک کاربر
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param Request body dto.OrdersByUserDTO true "OrdersByUserRequest"
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Validation Failed"
// @Router /v1/orders/by-user [post]
func (h *OrdersHandler) GetOrdersByUser(c *gin.Context) {
	req := new(dto.OrdersByUserDTO)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	orders, err := h.orderUsecase.GetOrdersByUserID(c, *req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(orders, true, helper.Success))
}

// @Security BearerAuth
// UpdateOrderStatus godoc
// @Summary بروزرسانی وضعیت سفارش
// @Description تغییر وضعیت سفارش (pending, paid, delivered, failed)
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param Request body dto.OrderStatusUpdateDTO true "UpdateOrderStatusRequest"
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Validation Failed"
// @Failure 500 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/orders/update-status [put]
func (h *OrdersHandler) UpdateOrderStatus(c *gin.Context) {
	req := new(dto.OrderStatusUpdateDTO)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	if err := h.orderUsecase.UpdateOrderStatus(c, *req); err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("وضعیت سفارش با موفقیت بروزرسانی شد", true, helper.Success))
}

// @Security BearerAuth
// DeleteOrder godoc
// @Summary حذف سفارش
// @Description حذف یک سفارش
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param Request body dto.OrderDeleteDTO true "DeleteOrderRequest"
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Validation Failed"
// @Failure 500 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/orders/delete [delete]
func (h *OrdersHandler) DeleteOrder(c *gin.Context) {
	req := new(dto.OrderDeleteDTO)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	if err := h.orderUsecase.DeleteOrder(c, *req); err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("سفارش با موفقیت حذف شد", true, helper.Success))
}

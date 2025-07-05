package controller

import (
	"RestGoTest/src/repository"
	"RestGoTest/src/service"
	"RestGoTest/src/util"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

// AllProductsController godoc
// @Summary      دریافت همه محصولات
// @Description  این API لیست کامل محصولات را برمی‌گرداند.
// @Tags         products
// @Produce      json
// @Success      200  {array}  repository.Product
// @Failure      500  {object}  util.ErrorResponse
// @Router       /products [get]
func AllProductsController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := service.AllProducts(r.Context())
		if err != nil {
			if r.Context().Err() == context.DeadlineExceeded {
				util.ResponseWithError(w, http.StatusRequestTimeout, "درخواست شما Timeout شد")
				return
			}
			util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.ResponseWithJSON(w, products, "فراخوانی با موفقیت انجام شد")
	}
}

// GetProductController godoc
// @Summary      دریافت محصول با شناسه
// @Description  دریافت اطلاعات یک محصول با استفاده از شناسه
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "شناسه محصول"
// @Success      200  {object}  repository.Product
// @Failure      400  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /products/{id} [get]
func GetProductController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "شناسه نامعتبر است")
			return
		}
		product, err := service.FetchProduct(r.Context(), id)
		if err != nil {
			if r.Context().Err() == context.DeadlineExceeded {
				util.ResponseWithError(w, http.StatusRequestTimeout, "درخواست شما Timeout شد")
				return
			}
			util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.ResponseWithJSON(w, product, "فراخوانی با موفقیت انجام شد")
	}
}

// CreateProductController godoc
// @Summary      ایجاد محصول جدید
// @Description  ایجاد یک محصول جدید با استفاده از اطلاعات ارسال شده
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body  repository.Product  true  "اطلاعات محصول جدید"
// @Success      201  {object}  repository.Product
// @Failure      400  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /products [post]
func CreateProductController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "خطا در خواندن بدنه درخواست")
			return
		}
		var p repository.Product
		if err := json.Unmarshal(reqBody, &p); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "فرمت JSON نامعتبر است")
			return
		}

		if err := validate.Struct(p); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "اعتبارسنجی ورودی نامعتبر است: "+err.Error())
			return
		}

		createResponse, err := service.CreateProduct(r.Context(), &p)
		if err != nil {
			if r.Context().Err() == context.DeadlineExceeded {
				util.ResponseWithError(w, http.StatusRequestTimeout, "درخواست شما Timeout شد")
				return
			}
			util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.ResponseWithJSON(w, createResponse, "اطلاعات جدید با موفقیت ذخیره شد")
	}
}

// DeleteProductController godoc
// @Summary      حذف یک محصول
// @Description  حذف یک محصول با استفاده از شناسه
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "شناسه محصول"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /products/{id} [delete]
func DeleteProductController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "شناسه نامعتبر است")
			return
		}
		if err := service.DeleteProduct(r.Context(), id); err != nil {
			if r.Context().Err() == context.DeadlineExceeded {
				util.ResponseWithError(w, http.StatusRequestTimeout, "درخواست شما Timeout شد")
				return
			}
			util.ResponseWithError(w, http.StatusInternalServerError, "خطا در حذف محصول")
			return
		}
		util.ResponseWithJSON(w, nil, "محصول با موفقیت حذف شد")
	}
}

// DeleteAllProductsController godoc
// @Summary      حذف همه محصولات
// @Description  حذف همه محصولات موجود در دیتابیس
// @Tags         products
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  util.ErrorResponse
// @Router       /products [delete]
func DeleteAllProductsController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := service.DeleteAllProducts(r.Context()); err != nil {
			if r.Context().Err() == context.DeadlineExceeded {
				util.ResponseWithError(w, http.StatusRequestTimeout, "درخواست شما Timeout شد")
				return
			}
			util.ResponseWithError(w, http.StatusInternalServerError, "خطا در حذف همه محصولات")
			return
		}
		util.ResponseWithJSON(w, nil, "همه محصولات با موفقیت حذف شدند")
	}
}

// UpdateProductController godoc
// @Summary      بروزرسانی محصول
// @Description  بروزرسانی اطلاعات یک محصول موجود
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body  repository.Product  true  "اطلاعات جدید محصول"
// @Success      200  {object}  repository.Product
// @Failure      400  {object}  util.ErrorResponse
// @Failure      500  {object}  util.ErrorResponse
// @Router       /products [put]
func UpdateProductController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p repository.Product
		body, err := io.ReadAll(r.Body)
		if err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "خطا در خواندن بدنه درخواست")
			return
		}
		if err := json.Unmarshal(body, &p); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "فرمت JSON نامعتبر است")
			return
		}
		if p.ID == 0 {
			util.ResponseWithError(w, http.StatusBadRequest, "شناسه (ID) الزامی است")
			return
		}
		if err := service.UpdateProduct(r.Context(), &p); err != nil {
			if r.Context().Err() == context.DeadlineExceeded {
				util.ResponseWithError(w, http.StatusRequestTimeout, "درخواست شما Timeout شد")
				return
			}
			util.ResponseWithError(w, http.StatusInternalServerError, "خطا در بروزرسانی محصول")
			return
		}
		util.ResponseWithJSON(w, p, "محصول با موفقیت بروزرسانی شد")
	}
}

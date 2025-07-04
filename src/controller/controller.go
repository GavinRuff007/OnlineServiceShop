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

	"github.com/gorilla/mux"
)

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

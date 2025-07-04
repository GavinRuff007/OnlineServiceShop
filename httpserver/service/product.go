package service

import (
	"RestGoTest/httpserver/dto"
	"RestGoTest/httpserver/repository"
	"RestGoTest/httpserver/util"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := repository.GetProducts()
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.ResponseWithJSON(w, products, "فراخوانی با موفقیت انجام شد")
}

func FetchProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var p repository.Product
	p.ID, _ = strconv.Atoi(vars["id"])
	err := p.GetProduct()
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
	}
	util.ResponseWithJSON(w, p, "فراخوانی با موفقیت انجام شد")
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var p repository.Product
	json.Unmarshal(reqBody, &p)
	err := p.CreateProduct()
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, err.Error())
	}

	createResponse := dto.CreateResponse{
		ID:          p.ID,
		ProductCode: p.ProductCode,
		Name:        p.Name,
	}

	util.ResponseWithJSON(w, createResponse, "اطلاعات جدید با موفقیت ذخیره شد")

}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, "شناسه نامعتبر است")
	}

	p := repository.Product{}
	p.ID = id
	err = p.DeleteProduct()
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, "خطا در حذف محصول")
	}

	util.ResponseWithJSON(w, nil, "محصول با موفقیت حذف شد")
}

func DeleteAllProducts(w http.ResponseWriter, r *http.Request) {
	err := repository.DeleteAllProducts()
	if err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, "خطا در حذف همه محصولات")

	}
	util.ResponseWithJSON(w, nil, "همه محصولات با موفقیت حذف شدند")
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var p repository.Product

	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, "خطا در خواندن بدنه درخواست")

	}

	if err := json.Unmarshal(body, &p); err != nil {
		util.ResponseWithError(w, http.StatusBadRequest, "فرمت JSON نامعتبر است")

	}

	if p.ID == 0 {
		util.ResponseWithError(w, http.StatusBadRequest, "شناسه (ID) الزامی است")

	}

	if err := p.UpdateProduct(); err != nil {
		util.ResponseWithError(w, http.StatusInternalServerError, "خطا در بروزرسانی محصول")
	}

	util.ResponseWithJSON(w, p, "محصول با موفقیت بروزرسانی شد")

}

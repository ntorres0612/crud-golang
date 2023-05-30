package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ntorres0612/ionix-crud/api/auth"
	"github.com/ntorres0612/ionix-crud/api/utils/formaterror"
	"github.com/ntorres0612/ionix-crud/models"
	"github.com/ntorres0612/ionix-crud/responses"
)

func (server *Server) CreateProductType(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productType := &models.ProductType{}
	err = json.Unmarshal(body, &productType)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	productType.Prepare()
	err = productType.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	exist := productType.CheckExistCreate(server.DB)
	if exist {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the productType already exist"))
		return
	}

	productTypeCreated, err := productType.SaveProductType(server.DB)
	if err != nil {
		fmt.Println(err.Error())
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, productTypeCreated.ID))
	responses.JSON(w, http.StatusCreated, productTypeCreated)
}

func (server *Server) GetProductTypes(w http.ResponseWriter, r *http.Request) {

	productType := models.ProductType{}

	productTypes, err := productType.FindAllProductType(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, productTypes)
}

func (server *Server) GetProductType(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	productType := models.ProductType{}

	productTypeReceived, err := productType.FindProductTypeByID(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, productTypeReceived)
}

func (server *Server) UpdateProductType(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	productType := models.ProductType{}
	err = server.DB.Debug().Model(models.ProductType{}).Where("id = ?", id).Take(&productType).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("productType not found"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productTypeUpdate := models.ProductType{}
	err = json.Unmarshal(body, &productTypeUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productTypeUpdate.Prepare()
	err = productTypeUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	exist := productTypeUpdate.CheckExistUpdate(server.DB)
	if exist {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the productType already exist"))
		return
	}

	productTypeUpdate.ID = productType.ID

	productTypeUpdated, _ := productTypeUpdate.UpdateProductType(server.DB)

	responses.JSON(w, http.StatusOK, productTypeUpdated)
}

func (server *Server) DeleteProductType(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	productType := models.ProductType{}
	err = server.DB.Debug().Model(models.ProductType{}).Where("id = ?", id).Take(&productType).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = productType.DeleteProductType(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
}

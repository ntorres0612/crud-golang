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

func (server *Server) CreateDelivery(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	delivery := &models.Delivery{}
	err = json.Unmarshal(body, &delivery)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	delivery.Prepare()
	err = delivery.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	deliveryCreated, err := delivery.SaveDelivery(server.DB)
	if err != nil {
		fmt.Println(err.Error())
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, deliveryCreated.ID))
	responses.JSON(w, http.StatusCreated, deliveryCreated)
}

func (server *Server) GetDeliverys(w http.ResponseWriter, r *http.Request) {

	delivery := models.Delivery{}

	deliveries, err := delivery.FindAllDeliverys(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, deliveries)
}

func (server *Server) GetDelivery(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	delivery := models.Delivery{}

	deliveryReceived, err := delivery.FindDeliveryByID(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, deliveryReceived)
}

func (server *Server) UpdateDelivery(w http.ResponseWriter, r *http.Request) {

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

	delivery := models.Delivery{}
	err = server.DB.Debug().Model(models.Delivery{}).Where("id = ?", id).Take(&delivery).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("delivery not found"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	deliveryUpdate := models.Delivery{}
	err = json.Unmarshal(body, &deliveryUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	deliveryUpdate.Prepare()
	err = deliveryUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// exist := deliveryUpdate.CheckExistUpdate(server.DB)
	// if exist {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the delivery name already exist"))
	// 	return
	// }

	deliveryUpdate.ID = delivery.ID

	deliveryUpdated, _ := deliveryUpdate.UpdateDelivery(server.DB)

	responses.JSON(w, http.StatusOK, deliveryUpdated)
}

func (server *Server) DeleteDelivery(w http.ResponseWriter, r *http.Request) {

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

	delivery := models.Delivery{}
	err = server.DB.Debug().Model(models.Delivery{}).Where("id = ?", id).Take(&delivery).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = delivery.DeleteDelivery(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
}

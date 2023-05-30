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

func (server *Server) CreateTruck(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	truck := &models.Truck{}
	err = json.Unmarshal(body, &truck)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	truck.Prepare()
	err = truck.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	exist := truck.CheckExistCreate(server.DB)
	if exist {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the license plate already exist"))
		return
	}

	truckCreated, err := truck.SaveTruck(server.DB)
	if err != nil {
		fmt.Println(err.Error())
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, truckCreated.ID))
	responses.JSON(w, http.StatusCreated, truckCreated)
}

func (server *Server) GetTrucks(w http.ResponseWriter, r *http.Request) {

	truck := models.Truck{}

	trucks, err := truck.FindAllTruck(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, trucks)
}

func (server *Server) GetTruck(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	truck := models.Truck{}

	truckReceived, err := truck.FindTruckByID(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, truckReceived)
}

func (server *Server) UpdateTruck(w http.ResponseWriter, r *http.Request) {

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

	truck := models.Truck{}
	err = server.DB.Debug().Model(models.Truck{}).Where("id = ?", id).Take(&truck).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("truck not found"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	truckUpdate := models.Truck{}
	err = json.Unmarshal(body, &truckUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	truckUpdate.Prepare()
	err = truckUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	exist := truckUpdate.CheckExistUpdate(server.DB)
	if exist {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the license plate already exist"))
		return
	}

	truckUpdate.ID = truck.ID

	truckUpdated, _ := truckUpdate.UpdateTruck(server.DB)

	responses.JSON(w, http.StatusOK, truckUpdated)
}

func (server *Server) DeleteTruck(w http.ResponseWriter, r *http.Request) {

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

	truck := models.Truck{}
	err = server.DB.Debug().Model(models.Truck{}).Where("id = ?", id).Take(&truck).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = truck.DeleteTruck(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
}

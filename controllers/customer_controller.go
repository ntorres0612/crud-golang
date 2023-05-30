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
	"github.com/ntorres0612/ionix-crud/models"
	"github.com/ntorres0612/ionix-crud/responses"
)

func (server *Server) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	customer := &models.Customer{}
	err = json.Unmarshal(body, &customer)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	customer.Prepare()
	err = customer.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	exist := customer.CheckExistCreate(server.DB)
	if exist {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the customer name already exist"))
		return
	}

	customerCreated, _ := customer.SaveCustomer(server.DB)

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, customerCreated.ID))
	responses.JSON(w, http.StatusCreated, customerCreated)
}

func (server *Server) GetCustomers(w http.ResponseWriter, r *http.Request) {

	customer := models.Customer{}

	customers, err := customer.FindAllCustomer(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, customers)
}

func (server *Server) GetCustomer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	customer := models.Customer{}

	customerReceived, err := customer.FindCustomerByID(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, customerReceived)
}

func (server *Server) UpdateCustomer(w http.ResponseWriter, r *http.Request) {

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

	customer := models.Customer{}
	err = server.DB.Debug().Model(models.Customer{}).Where("id = ?", id).Take(&customer).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("customer not found"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	customerUpdate := models.Customer{}
	err = json.Unmarshal(body, &customerUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	customerUpdate.Prepare()
	err = customerUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	exist := customerUpdate.CheckExistUpdate(server.DB)
	if exist {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the customer name already exist"))
		return
	}

	customerUpdate.ID = customer.ID

	customerUpdated, _ := customerUpdate.UpdateCustomer(server.DB)

	responses.JSON(w, http.StatusOK, customerUpdated)
}

func (server *Server) DeleteCustomer(w http.ResponseWriter, r *http.Request) {

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

	customer := models.Customer{}
	err = server.DB.Debug().Model(models.Customer{}).Where("id = ?", id).Take(&customer).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = customer.DeleteCustomer(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
}

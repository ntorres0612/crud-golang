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

func (server *Server) CreateStore(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	store := &models.Store{}
	err = json.Unmarshal(body, &store)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	store.Prepare()
	err = store.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	exist := store.CheckExistCreate(server.DB)
	if exist {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the store name already exist"))
		return
	}

	storeCreated, err := store.SaveStore(server.DB)
	if err != nil {
		fmt.Println(err.Error())
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, storeCreated.ID))
	responses.JSON(w, http.StatusCreated, storeCreated)
}

func (server *Server) GetStores(w http.ResponseWriter, r *http.Request) {

	store := models.Store{}

	stores, err := store.FindAllStore(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, stores)
}

func (server *Server) GetStore(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	store := models.Store{}

	storeReceived, err := store.FindStoreByID(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, storeReceived)
}

func (server *Server) UpdateStore(w http.ResponseWriter, r *http.Request) {

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

	store := models.Store{}
	err = server.DB.Debug().Model(models.Store{}).Where("id = ?", id).Take(&store).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("store not found"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	storeUpdate := models.Store{}
	err = json.Unmarshal(body, &storeUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	storeUpdate.Prepare()
	err = storeUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	exist := storeUpdate.CheckExistUpdate(server.DB)
	if exist {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("the store name already exist"))
		return
	}

	storeUpdate.ID = store.ID

	storeUpdated, _ := storeUpdate.UpdateStore(server.DB)

	responses.JSON(w, http.StatusOK, storeUpdated)
}

func (server *Server) DeleteStore(w http.ResponseWriter, r *http.Request) {

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

	store := models.Store{}
	err = server.DB.Debug().Model(models.Store{}).Where("id = ?", id).Take(&store).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = store.DeleteStore(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
}

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ntorres0612/ionix-crud/models"
	"github.com/ntorres0612/ionix-crud/responses"
)

func (server *Server) GetLogisticTypes(w http.ResponseWriter, r *http.Request) {

	logisticType := models.LogisticType{}

	logisticTypes, err := logisticType.FindAllLogisticType(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, logisticTypes)
}

func (server *Server) GetLogisticType(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	logisticType := models.LogisticType{}

	storeReceived, err := logisticType.FindLogisticTypeByID(server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, storeReceived)
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"user-management/internal/errors"
	"user-management/internal/model"
	"user-management/internal/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}
func (h *UserHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, "unable to retrieve users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
func (h *UserHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}
	user, err := h.service.GetUser(id)
	if err!=nil{
		h.handleServiceError(w,err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func (h *UserHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w,"invalid json ",http.StatusBadRequest)
		return
	}
	err=h.service.CreateUser(&user)
	if err!=nil{
		h.handleServiceError(w,err)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateHandler(w http.ResponseWriter,r *http.Request){
	vars:=mux.Vars(r)
	id,err:=strconv.Atoi(vars["id"])
	fmt.Println("update handler called",id)
	if err!=nil{
		http.Error(w,"invalid id format ",http.StatusBadRequest)
		return
	}
	var user model.User
	err=json.NewDecoder(r.Body).Decode(&user)
	if err!=nil{
		http.Error(w,"invalid json",http.StatusBadRequest)
		return
	}
	fmt.Println("ddata to update:",user)

	err=h.service.UpdateUser(id,user)
	if err!=nil{
		fmt.Println("encountered error while updating",err)
		h.handleServiceError(w,err)
		return
	}
	updatedUser, err := h.service.GetUser(id)
	fmt.Println("updated user:",updatedUser)
    if err != nil {
        h.handleServiceError(w, err)
        return
    }
		w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func (h *UserHandler) DeleteHandler(w http.ResponseWriter,r *http.Request){
	vars:=mux.Vars(r)
	id,err:=strconv.Atoi(vars["id"])
	if err!=nil{
		http.Error(w,"invalid id format ",http.StatusBadRequest)
		return
	}
	
	err=h.service.DeleteUser(id)
	if err!=nil{
		h.handleServiceError(w,err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *UserHandler) handleServiceError(w http.ResponseWriter,err error){
	switch e:=err.(type){
	case *errors.ValidationError:
		http.Error(w,e.Error(),http.StatusBadRequest)
	case *errors.NotFoundError:
		http.Error(w,e.Error(),http.StatusNotFound)
	case *errors.DuplicateError:
		http.Error(w,e.Error(),http.StatusConflict)
	}
}
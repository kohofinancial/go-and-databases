package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kohofinancial/go-and-databases/services"
	"net/http"
)

var NoIDError = fmt.Errorf("no id received")

func (a *app) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("no id sent in url")
		return
	}

	user, err := a.UserService.Get(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(userJson)
	return
}

func (a *app) AddUser(w http.ResponseWriter, r *http.Request) {
	var user services.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ReturnError(w, err)
		return
	}

	err = a.UserService.Create(&user)
	if err != nil {
		ReturnError(w, err)
		return
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		ReturnError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(userJson)
	return
}

func (a *app) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		ReturnError(w, NoIDError)
		return
	}

	err := a.UserService.Delete(vars["id"])
	if err != nil {
		ReturnError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (a *app) DeleteAll(w http.ResponseWriter, _ *http.Request) {
	err := a.UserService.DeleteAll()
	if err != nil {
		ReturnError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (a *app) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user services.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ReturnError(w, err)
		return
	}

	vars := mux.Vars(r)

	if vars["id"] == "" {
		ReturnError(w, NoIDError)
		return
	}

	user.ID = vars["id"]

	err = a.UserService.Update(&user)
	if err != nil {
		ReturnError(w, err)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		ReturnError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(userJson)
	return
}

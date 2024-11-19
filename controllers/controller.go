package controllers

import (
	"github.com/gorilla/mux"
)

func RootingRegister(registerUserController *RegisterUserController)*mux.Router{
	r := mux.NewRouter()

	r.HandleFunc("/user/create",registerUserController.CreateUser).Methods("POST")

	return r
}
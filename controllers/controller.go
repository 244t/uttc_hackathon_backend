package controllers

import (
	"github.com/gorilla/mux"
)

func RootingRegister(registerUserController *RegisterUserController){
	r := mux.NewRouter()

	r.HandleFunc("/user/create",registerUserController.CreateUser).Methods("POST")

}
package controllers

import (
	"github.com/gorilla/mux"
)

func RootingRegister(userController *UserController)*mux.Router{
	r := mux.NewRouter()

	r.HandleFunc("/user/create",userController.CreateUser).Methods("POST")
	r.HandleFunc("/user/{userId}",userController.GetUserProfile).Methods("Get")
	return r
}
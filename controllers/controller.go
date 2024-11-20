package controllers

import (
	"github.com/gorilla/mux"
)

func RootingRegister(userController *UserController)*mux.Router{
	r := mux.NewRouter()

	r.HandleFunc("/user/create",userController.CreateUser).Methods("POST")
	r.HandleFunc("/user/{userId}",userController.GetUserProfile).Methods("Get")
	r.HandleFunc("/user/{userId}/followers",userController.GetFollowers).Methods("Get")
	r.HandleFunc("/user/{userId}/following",userController.GetFollowing).Methods("Get")
	return r
}
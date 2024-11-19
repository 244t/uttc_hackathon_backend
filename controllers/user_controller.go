package controllers

import (
	"net/http"
	"myproject/usecase"
	"encoding/json"
)

type RegisterUserController struct {
	RegisterUserUseCase *usecase.RegisterUserUseCase
}

// NewRegisterUserControllerはRegisterUserControllerのインスタンスを返します。
func NewRegisterUserController(registerUserUseCase *usecase.RegisterUserUseCase) *RegisterUserController {
	return &RegisterUserController{RegisterUserUseCase: registerUserUseCase}
}

// ユーザー登録
func (c *RegisterUserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRegister usecase.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&userRegister); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.RegisterUserUseCase.RegisterUser(userRegister); err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

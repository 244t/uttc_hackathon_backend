// package controllers

// import (
// 	"net/http"
// 	"myproject/usecase"
// 	"encoding/json"
// )

// type RegisterUserController struct {
// 	RegisterUserUseCase *usecase.RegisterUserUseCase
// }

// // NewRegisterUserControllerはRegisterUserControllerのインスタンスを返します。
// func NewRegisterUserController(registerUserUseCase *usecase.RegisterUserUseCase) *RegisterUserController {
// 	return &RegisterUserController{RegisterUserUseCase: registerUserUseCase}
// }

// // ユーザー登録
// func (c *RegisterUserController) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	var userRegister usecase.UserRegister
// 	if err := json.NewDecoder(r.Body).Decode(&userRegister); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if err := c.RegisterUserUseCase.RegisterUser(userRegister); err != nil {
// 		http.Error(w, "Error registering user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }

// // プロフィール取得
// func (c *RegisterUserController) GetUser(w http.ResponseWriter, r *http.Request) {
// 	var userRegister usecase.UserRegister
// 	if err := json.NewDecoder(r.Body).Decode(&userRegister); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if err := c.RegisterUserUseCase.RegisterUser(userRegister); err != nil {
// 		http.Error(w, "Error registering user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }
package controllers

import (
	"net/http"
	"myproject/usecase"
	"myproject/dao"
	"encoding/json"
	"github.com/gorilla/mux"
)

type UserController struct {
	RegisterUserUseCase *usecase.RegisterUserUseCase
	GetProfileUserUseCase *usecase.GetProfileUserUseCase
}

// NewUserControllerはUserControllerのインスタンスを返します。
func NewUserController(db dao.TweetDAOInterface) *UserController {

	// UseCaseのインスタンスを作成
	registerUserUseCase := usecase.NewRegisterUserUseCase(db)
	getProfileUserUseCase := usecase.NewGetProfileUserUseCase(db)

	// UserControllerを作成して返す
	return &UserController{
		RegisterUserUseCase: registerUserUseCase,
		GetProfileUserUseCase: getProfileUserUseCase,
	}
}

// ユーザー登録
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
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

// プロフィール取得
func (c *UserController) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // mux.Varsでパスパラメータを取得
	userID := vars["userId"]

	// GetProfileUserUseCaseを呼び出してユーザーのプロフィールを取得
	profile, err := c.GetProfileUserUseCase.GetUserProfile(userID)
	if err != nil {
		http.Error(w, "Error fetching user profile", http.StatusInternalServerError)
		return
	}

	// プロフィール情報をJSONで返す
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "Error encoding profile", http.StatusInternalServerError)
		return
	}
}

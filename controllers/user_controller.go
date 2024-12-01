package controllers

import (
	"net/http"
	"myproject/usecase"
	"myproject/dao"
	"myproject/model"
	"encoding/json"
	"github.com/gorilla/mux"
)

type UserController struct {
	RegisterUserUseCase *usecase.RegisterUserUseCase
	GetProfileUserUseCase *usecase.GetProfileUserUseCase
	GetFollowingUserUseCase *usecase.GetFollowingUserUseCase
	GetFollowersUserUseCase *usecase.GetFollowersUserUseCase
	FollowUserUseCase *usecase.FollowUserUseCase
	UnFollowUserUseCase *usecase.UnFollowUserUseCase
	UpdateProfileUserUseCase *usecase.UpdateProfileUserUseCase
	GetUserPostsUserUseCase *usecase.GetUserPostsUserUseCase
	SearchUserUseCase *usecase.SearchUserUseCase
}

// NewUserControllerはUserControllerのインスタンスを返します。
func NewUserController(db dao.TweetDAOInterface) *UserController {

	// UseCaseのインスタンスを作成
	registerUserUseCase := usecase.NewRegisterUserUseCase(db)
	getProfileUserUseCase := usecase.NewGetProfileUserUseCase(db)
	getFollowingUserUseCase := usecase.NewGetFollowingUserUseCase(db)
	getFollowersUserUseCase := usecase.NewGetFollowersUserUseCase(db)
	followUserUseCase := usecase.NewFollowUserUseCase(db)
	unfollowUserUseCase := usecase.NewUnFollowUserUseCase(db)
	updateProfileUserUSeCase := usecase.NewUpdateProfileUserUseCase(db)
	getUserPostsUserUseCase := usecase.NewGetUserPostsUserUseCase(db)
	searchUserUseCase := usecase.NewSearchUserUseCase(db)


	// UserControllerを作成して返す
	return &UserController{
		RegisterUserUseCase: registerUserUseCase,
		GetProfileUserUseCase: getProfileUserUseCase,
		GetFollowingUserUseCase: getFollowingUserUseCase,
		GetFollowersUserUseCase: getFollowersUserUseCase,
		FollowUserUseCase: followUserUseCase,
		UnFollowUserUseCase: unfollowUserUseCase,
		UpdateProfileUserUseCase: updateProfileUserUSeCase,
		GetUserPostsUserUseCase: getUserPostsUserUseCase,
		SearchUserUseCase : searchUserUseCase,
	}
}

// ユーザー登録
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRegister model.Profile
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

//userIdが示すuserがフォローしているアカウントを返す
func (c *UserController) GetFollowing(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) // mux.Varsでパスパラメータを取得
	userID := vars["userId"]

	// GetFollowingUserUseCaseを呼び出してフォローしている人のプロフィールを取得
	users, err := c.GetFollowingUserUseCase.GetFollowing(userID)
	if err != nil {
		http.Error(w, "Error fetching following profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

//userIdが示すuserがフォローしているアカウントを返す
func (c *UserController) GetFollowers(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) // mux.Varsでパスパラメータを取得
	userID := vars["userId"]

	// GetFollowersUserUseCaseを呼び出してフォローしている人のプロフィールを取得
	users, err := c.GetFollowersUserUseCase.GetFollowers(userID)
	if err != nil {
		http.Error(w, "Error fetching following profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

//フォロー
func (c* UserController) Follow(w http.ResponseWriter, r *http.Request){
	var followRegister usecase.FollowRegister
	if err := json.NewDecoder(r.Body).Decode(&followRegister); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.FollowUserUseCase.Follow(followRegister); err != nil{
		http.Error(w, "Error registering follow", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//フォロー外す
func (c* UserController) UnFollow(w http.ResponseWriter, r *http.Request){
	var unFollowRegister usecase.UnFollowRegister
	if err := json.NewDecoder(r.Body).Decode(&unFollowRegister); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.UnFollowUserUseCase.UnFollow(unFollowRegister); err != nil{
		http.Error(w, "Error registering unfollow", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ユーザープロフィール更新
func (c *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var ur model.Profile
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.UpdateProfileUserUseCase.UpdateProfile(ur); err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *UserController) GetUserPosts(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	userId := vars["userId"]

	posts, err := c.GetUserPostsUserUseCase.GetUserPosts(userId)
	if err != nil {
		http.Error(w,"Error fetching user posts",http.StatusInternalServerError)
		return 
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w,"Error encoding user posts",http.StatusInternalServerError)
		return 
	}
}

func (c *UserController) SearchUser(w http.ResponseWriter, r *http.Request){
	var searchWord usecase.Search
	if err := json.NewDecoder(r.Body).Decode(&searchWord); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	users, err := c.SearchUserUseCase.SearchUser(searchWord)
	if err != nil {
		http.Error(w, "Error fetching search profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// OPTIONSリクエストに対する処理
func (uc *UserController) CORSOptionsHandler(w http.ResponseWriter, r *http.Request) {
    // 必要なCORSヘッダーを設定
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT,POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK) // 200 OKを返す
}

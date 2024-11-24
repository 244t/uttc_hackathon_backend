package controllers

import (
	"net/http"
	"myproject/usecase"
	"myproject/dao"
	"encoding/json"
	"github.com/gorilla/mux"
)

type PostController struct{
	CreatePostUseCase *usecase.CreatePostUseCase
	GetPostUseCase *usecase.GetPostUseCase
	UpdatePostUseCase *usecase.UpdatePostUseCase
	DeletePostUseCase *usecase.DeletePostUseCase
	ReplyPostUseCase *usecase.ReplyPostUseCase
	LikePostUseCase *usecase.LikePostUseCase
	UnLikePostUseCase *usecase.UnLikePostUseCase
	GetLikesPostUseCase *usecase.GetLikesPostUseCase
	TimelineUseCase *usecase.TimelineUseCase
}

func NewPostController(db dao.PostDAOInterface) *PostController{
	createPostUseCase := usecase.NewCreatePostUseCase(db)
	getPostUseCase := usecase.NewGetPostUseCase(db)
	updatePostUseCase := usecase.NewUpdatePostUSeCase(db)
	deletePostUseCase := usecase.NewDeletePostUseCase(db)
	replyPostUseCase := usecase.NewReplyPostUseCase(db)
	likePostUseCase := usecase.NewLikePostUseCase(db)
	unlikePostUseCase := usecase.NewUnLikePostUseCase(db)
	getLikesPostUseCase := usecase.NewGetLikesPostUseCase(db)
	timelineUseCase := usecase.NewTimelineUseCase(db)

	return &PostController{
		CreatePostUseCase: createPostUseCase,
		GetPostUseCase: getPostUseCase,
		UpdatePostUseCase: updatePostUseCase,
		DeletePostUseCase: deletePostUseCase,
		ReplyPostUseCase: replyPostUseCase,
		LikePostUseCase : likePostUseCase,
		UnLikePostUseCase : unlikePostUseCase,
		GetLikesPostUseCase : getLikesPostUseCase,
		TimelineUseCase : timelineUseCase,
	}
}

//Postを作成
func (c *PostController) CreatePost(w http.ResponseWriter, r* http.Request){
	var postRegister usecase.PostRegister
	if err := json.NewDecoder(r.Body).Decode(&postRegister); err != nil {
		http.Error(w,"Error create post",http.StatusInternalServerError)
		return 
	}

	if err := c.CreatePostUseCase.CreatePost(postRegister); err != nil {
		http.Error(w,"Error creating post",http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusCreated)
}

//Post取得
func (c *PostController) GetPost(w http.ResponseWriter, r* http.Request){
	vars := mux.Vars(r)
	postId := vars["postId"]

	post, err := c.GetPostUseCase.GetPost(postId)
	if err != nil {
		http.Error(w, "Error fetching post", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Error encoding profile", http.StatusInternalServerError)
		return
	}

}

//Postを作成
func (c *PostController) UpdatePost(w http.ResponseWriter, r* http.Request){
	var pu  usecase.PostUpdate
	if err := json.NewDecoder(r.Body).Decode(&pu); err != nil {
		http.Error(w,"Error update post",http.StatusInternalServerError)
		return 
	}

	if err := c.UpdatePostUseCase.UpdatePost(pu); err != nil {
		http.Error(w,"Error updating post",http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *PostController) DeletePost(w http.ResponseWriter, r* http.Request){
	vars := mux.Vars(r)
	postId := vars["postId"]
	err := c.DeletePostUseCase.DeletePost(postId)
	if err != nil {
		http.Error(w, "Error fetching post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

//返信を作成
func (c *PostController) ReplyPost(w http.ResponseWriter, r* http.Request){
	var postRegister usecase.PostRegister
	vars := mux.Vars(r)
	postId := vars["postId"]
	if err := json.NewDecoder(r.Body).Decode(&postRegister); err != nil {
		http.Error(w,"Error create post",http.StatusInternalServerError)
		return 
	}

	if err := c.ReplyPostUseCase.ReplyPost(postId,postRegister); err != nil {
		http.Error(w,"Error reply post",http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusCreated)
}

//いいね
func (c *PostController) LikePost (w http.ResponseWriter, r* http.Request){
	var like usecase.Like

	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w,"Error like",http.StatusInternalServerError)
		return 
	}


	if err := c.LikePostUseCase.LikePost(like); err != nil {
		http.Error(w,"Error like post",http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusCreated)
}

//いいねを削除
func (c *PostController) UnLikePost (w http.ResponseWriter, r* http.Request){
	var unlike usecase.Like

	if err := json.NewDecoder(r.Body).Decode(&unlike); err != nil {
		http.Error(w,"Error like",http.StatusInternalServerError)
		return 
	}


	if err := c.UnLikePostUseCase.UnLikePost(unlike); err != nil {
		http.Error(w,"Error like post",http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusCreated)
}

//指定したpostIdのいいね数を取得
func (c *PostController) GetLikes(w http.ResponseWriter, r* http.Request){
	vars := mux.Vars(r)
	postId := vars["postId"]
	likes, err := c.GetLikesPostUseCase.GetLikesPost(postId)
	if err != nil {
		http.Error(w, "Error fetching likes", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(likes); err != nil {
		http.Error(w, "Error encoding likes", http.StatusInternalServerError)
		return
	}
}

//フォローしているユーザーの投稿を取得
func (c *PostController) Timeline(w http.ResponseWriter, r* http.Request){
	vars := mux.Vars(r)
	userId := vars["userId"]

	posts, err := c.TimelineUseCase.Timeline(userId)
	if err != nil {
		http.Error(w,"Error fetching timeline",http.StatusInternalServerError)
		return 
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w,"Error encoding timeline",http.StatusInternalServerError)
		return 
	}
}

// OPTIONSリクエストに対する処理
func (pc *PostController) CORSOptionsHandler(w http.ResponseWriter, r *http.Request) {
    // 必要なCORSヘッダーを設定
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK) // 200 OKを返す
}

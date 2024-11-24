package controllers

import (
	"github.com/gorilla/mux"
)


func (uc *UserController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/register",uc.CreateUser).Methods("POST")
	r.HandleFunc("/user/register",uc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/user/{userId}",uc.GetUserProfile).Methods("Get")

	r.HandleFunc("/user/{userId}/followers",uc.GetFollowers).Methods("Get")

	r.HandleFunc("/user/{userId}/following",uc.GetFollowing).Methods("Get")

	r.HandleFunc("/user/follow",uc.Follow).Methods("Post")
	r.HandleFunc("/user/follow",uc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/user/unfollow",uc.UnFollow).Methods("Post")
	r.HandleFunc("/user/unfollow",uc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/user/{userId}/update",uc.UpdateProfile).Methods("PUT")
	r.HandleFunc("/user/{userId}/update",uc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/user/{userId}/posts",uc.GetUserPosts).Methods("Get")
}

func (pc *PostController) RegisterRoutes(r *mux.Router){
	r.HandleFunc("/post/create",pc.CreatePost).Methods("Post")
	r.HandleFunc("/post/create",pc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/post/{postId}",pc.GetPost).Methods("Get")

	r.HandleFunc("/post/update",pc.UpdatePost).Methods("Put")
	r.HandleFunc("/post/update",pc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/post/create",pc.CreatePost).Methods("Post")
	r.HandleFunc("/post/create",pc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/post/{postId}",pc.DeletePost).Methods("Delete")

	r.HandleFunc("/post/{postId}/reply",pc.ReplyPost).Methods("POST")
	r.HandleFunc("/post/{postId}/reply",pc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/post/{postId}/like",pc.LikePost).Methods("Post")
	r.HandleFunc("/post/{postId}/like",pc.CORSOptionsHandler).Methods("OPTIONS")  

	r.HandleFunc("/post/{postId}/unlike",pc.UnLikePost).Methods("Post")
	r.HandleFunc("/post/{postId}/unlike",pc.CORSOptionsHandler).Methods("OPTIONS")

	r.HandleFunc("/post/{postId}/likes",pc.GetLikes).Methods("Get")

	r.HandleFunc("/timeline/{userId}",pc.Timeline).Methods("Get")
}
package model

type Profile struct{
	Id string `json:"user_id"`
	Name string `json:"name"`
	Bio string `json:"bio"`
	FireBaseId string `json:"firebase_id"`
	ImgUrl string `json:"image_url"`
}

type Follow struct{
	UserId string `json:"user_id"`
	FollowingId string `json:"following_id"`
}

type UnFollow struct{
	UserId string `json:"user_id"`
	FollowingId string `json:"following_id"`
}
package model

type Profile struct{
	Id string `json:"user_id"`
	Name string `json:"name"`
	Bio string `json:"bio"`
	ImgUrl string `json:"img_url"`
	HeaderUrl string `json:"header_url"`
	Location string `json:"location"`
}

type Follow struct{
	UserId string `json:"user_id"`
	FollowingId string `json:"following_id"`
}

type UnFollow struct{
	UserId string `json:"user_id"`
	FollowingId string `json:"following_id"`
}

type NotificationInfo struct{
	UserId string `json:"user_id"`
	Flag string `json:"flag"`
	UserProfileImg string `json:"profile_img"`
	NotificationId string `json:"notification_id"`
}
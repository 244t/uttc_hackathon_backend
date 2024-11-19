package model

type Profile struct{
	Id string `json:"user_id"`
	Name string `json:"name"`
	Bio string `json:"bio"`
	FireBaseId string `json:"firebase_id"`
}
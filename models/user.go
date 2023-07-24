package models

import "github.com/gorilla/websocket"

type RegisterUser struct {
	FirstName  string `json:"first_name" msgpack:"FirstName"`
	SecondName string `json:"second_name" msgpack:"SecondName"`
	Birthdate  string `json:"birthdate" msgpack:"Birthdate"`
	Biography  string `json:"biography" msgpack:"Biography"`
	City       string `json:"city" msgpack:"City"`
	Password   string `json:"password" msgpack:"Password"`
}

type UserInfo struct {
	Id         string `json:"id" msgpack:"Id"`
	FirstName  string `json:"first_name" msgpack:"FirstName"`
	SecondName string `json:"second_name" msgpack:"SecondName"`
	Age        int    `json:"age" msgpack:"Age"`
	Birthdate  string `json:"birthdate" msgpack:"Birthdate"`
	Biography  string `json:"biography" msgpack:"Biography"`
	City       string `json:"city" msgpack:"City"`
	Password   string `json:"-" msgpack:"Password"`
}

type AuthInfo struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type RegisterResult struct {
	UserId string `json:"user_id"`
}

type LoginResult struct {
	Token string `json:"token"`
}

type UserSearchRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ActiveWsUsers struct {
	User string
	Conn *websocket.Conn
}

type UserSession struct {
	Id        string `json:"id" msgpack:"Id"`
	UserId    string `json:"userId" msgpack:"user_id"`
	Token     string `json:"token" msgpack:"token"`
	CreatedAt uint64 `json:"createdAt" msgpack:"created_at"`
}

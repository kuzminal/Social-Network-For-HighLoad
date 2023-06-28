package models

import "github.com/gorilla/websocket"

type RegisterUser struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Password   string `json:"password"`
}

type UserInfo struct {
	Id         string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Age        int    `json:"age"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Password   string `json:"-"`
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

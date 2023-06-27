package models

type Message struct {
	Id        string `json:"-"`
	Text      string `json:"text"`
	FromUser  string `json:"from"`
	ToUser    string `json:"to"`
	ChatId    string `json:"-"`
	CreatedAt string `json:"-"`
}

type Chat struct {
	Id       string
	FromUser string
	ToUser   string
}

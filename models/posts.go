package models

type Post struct {
	Id           string `json:"id"`
	Text         string `json:"text"`
	AuthorUserId string `json:"author_user_id,omitempty"`
	CreatedAt    string `json:"-"`
}

type PostCreateResponse struct {
	Id string `json:"id"`
}

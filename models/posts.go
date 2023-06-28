package models

type Post struct {
	Id           string `json:"id" msgpack:"id"`
	Text         string `json:"text" msgpack:"text"`
	AuthorUserId string `json:"author_user_id,omitempty" msgpack:"author_user_id"`
	CreatedAt    string `json:"-" msgpack:"created_at"`
}

type PostCreateResponse struct {
	Id string `json:"id"`
}

type Feed struct {
	Posts []Post
}

type Cache struct {
	UserId string
	Posts  []Post
}

type UpdateFeedRequest struct {
	UserId string
	Post   Post
}

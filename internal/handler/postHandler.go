package handler

import (
	"SocialNetHL/internal/store"
	"SocialNetHL/models"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	"strconv"
)

func (i *Instance) HandlePostCreate(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not add post to empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	post.AuthorUserId = userId
	post.Id = uuid.Must(uuid.NewV4()).String()
	postId, err := i.store.AddPost(context.Background(), post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not add new post")
		return
	}
	w.Write([]byte(postId))
}

func (i *Instance) HandlePostDelete(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "id")
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not delete friend from empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := i.store.DeletePost(context.Background(), userId, postId)
	if err != nil {
		log.Println("Could not delete friend from user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (i *Instance) HandlePostUpdate(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	err = i.store.UpdatePost(context.Background(), post)
	if err != nil {
		log.Printf("Could not update post for user, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (i *Instance) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	readStorage := store.GetReadNode(i.readStorages)
	post, _ := readStorage.GetPost(context.Background(), id)
	postDTO, _ := json.Marshal(post)
	w.Header().Add("Content-Type", "application/json")
	w.Write(postDTO)
}

func (i *Instance) HandleFeed(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	readStorage := store.GetReadNode(i.readStorages)
	if len(limit) == 0 {
		limit = "10"
	}
	if len(offset) == 0 {
		offset = "0"
	}
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not delete friend from empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limitNum, _ := strconv.Atoi(limit)
	offsetNum, _ := strconv.Atoi(offset)
	posts, _ := readStorage.FeedPost(context.Background(), offsetNum, limitNum, userId)
	postsDTO, _ := json.Marshal(posts)
	w.Header().Add("Content-Type", "application/json")
	w.Write(postsDTO)
}
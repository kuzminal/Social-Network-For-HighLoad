package handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (i *Instance) HandleFriendAdd(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not add friend to empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	friendId := chi.URLParam(r, "user_id")

	err := i.store.AddFriend(context.Background(), userId, friendId)
	if err != nil {
		log.Println("Could not add friend to user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (i *Instance) HandleFriendDelete(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not delete friend from empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	friendId := chi.URLParam(r, "user_id")

	err := i.store.DeleteFriend(context.Background(), userId, friendId)
	if err != nil {
		log.Println("Could not delete friend from user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

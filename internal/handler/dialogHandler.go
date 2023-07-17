package handler

import (
	"SocialNetHL/models"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	"time"
)

func (i *Instance) SendMessage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not delete friend from empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId, _ := i.dialogueStore.GetChatId(context.Background(), id, userId)

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	msg.Id = uuid.Must(uuid.NewV4()).String()
	msg.FromUser = userId
	msg.ToUser = id
	msg.CreatedAt = time.Now().Format("2006-01-02")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	msg.ChatId = chatId
	err = i.dialogueStore.SaveMessage(context.Background(), msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (i *Instance) GetMessages(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not delete friend from empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg, _ := i.dialogueStore.GetMessages(context.Background(), id, userId)
	msgDTO, _ := json.Marshal(msg)
	w.Header().Add("Content-Type", "application/json")
	w.Write(msgDTO)
}

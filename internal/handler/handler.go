package handler

import (
	"SocialNetHL/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (i *Instance) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	saveUser, err := i.store.SaveUser(context.Background(), &user)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\n  \"user_id\": \"%s\" \n}", saveUser)))
}

func (i *Instance) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, _ := i.store.LoadUser(context.Background(), id)
	userDTO, _ := json.Marshal(user)
	w.Write(userDTO)
	w.WriteHeader(http.StatusOK)
}

func (i *Instance) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

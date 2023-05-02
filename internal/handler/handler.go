package handler

import (
	"SocialNetHL/models"
	"context"
	"crypto/sha256"
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
	var authInfo models.AuthInfo
	err := json.NewDecoder(r.Body).Decode(&authInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	userInfo, _ := i.store.LoadUser(context.Background(), authInfo.Id)
	passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(authInfo.Password)))
	if len(userInfo.Id) > 0 && passHash == userInfo.Password {
		saveUser, err := i.store.CreateSession(context.Background(), &authInfo)
		if err != nil {
			return
		}
		registerRes := models.RegisterResult{UserId: saveUser}
		rr, _ := json.Marshal(registerRes)
		w.WriteHeader(http.StatusOK)
		w.Write(rr)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

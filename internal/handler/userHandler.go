package handler

import (
	"SocialNetHL/internal/store"
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
	saveUser, err := i.userStore.SaveUser(context.Background(), user)
	if err != nil {
		return
	}
	w.Write([]byte(fmt.Sprintf("{\n  \"user_id\": \"%s\" \n}", saveUser)))
}

func (i *Instance) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	readStorage := store.GetReadNode(i.readStorages)
	user, _ := readStorage.LoadUser(context.Background(), id)
	userDTO, _ := json.Marshal(user)
	w.Header().Add("Content-Type", "application/json")
	w.Write(userDTO)
}

func (i *Instance) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var authInfo models.AuthInfo
	err := json.NewDecoder(r.Body).Decode(&authInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}
	userInfo, _ := i.userStore.LoadUser(context.Background(), authInfo.Id)
	passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(authInfo.Password)))
	if len(userInfo.Id) > 0 && passHash == userInfo.Password {
		saveUser, err := i.sessionStore.CreateSession(context.Background(), &authInfo)
		if err != nil {
			return
		}
		loginRes := models.LoginResult{Token: saveUser.Token}
		rr, _ := json.Marshal(loginRes)
		w.Write(rr)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (i *Instance) HandleSearchUser(w http.ResponseWriter, r *http.Request) {
	var userSearchRequest models.UserSearchRequest
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	readStorage := store.GetReadNode(i.readStorages)
	if len(firstName) > 0 || len(lastName) > 0 {
		userSearchRequest.LastName = lastName
		userSearchRequest.FirstName = firstName
		users, _ := readStorage.SearchUser(context.Background(), userSearchRequest)
		if users == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			userDTO, _ := json.Marshal(users)
			w.Header().Add("Content-Type", "application/json")
			w.Write(userDTO)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

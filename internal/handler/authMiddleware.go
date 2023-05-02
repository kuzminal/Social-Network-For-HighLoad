package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (i *Instance) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			token := authHeader[1]
			user, _ := i.store.LoadSession(context.Background(), token)
			if len(user) > 0 {
				log.Printf("User with id: %s logged in", user)
				next.ServeHTTP(w, r)
				return
			} else {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}
	})
}

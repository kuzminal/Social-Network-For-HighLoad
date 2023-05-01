package router

import (
	"SocialNetHL/internal/handler"
	"SocialNetHL/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter(i *handler.Instance) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth)
		r.Post("/user/register", i.HandleRegister)
		r.Get("/user/get/{id}", i.HandleGetUser)
	})

	r.Post("/login", i.HandleLogin)

	return r
}

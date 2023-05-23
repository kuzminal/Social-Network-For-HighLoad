package router

import (
	"SocialNetHL/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(i *handler.Instance) http.Handler {
	r := chi.NewRouter()
	r.Mount("/debug", middleware.Profiler())
	r.Group(func(r chi.Router) {
		r.Use(i.BasicAuth)

		r.Post("/user/register", i.HandleRegister)
		r.Get("/user/get/{id}", i.HandleGetUser)
		r.Get("/user/search", i.HandleSearchUser)
	})

	r.Post("/login", i.HandleLogin)

	return r
}

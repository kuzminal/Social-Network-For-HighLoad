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

		r.Put("/friend/set/{user_id}", i.HandleFriendAdd)
		r.Put("/friend/delete/{user_id}", i.HandleFriendDelete)

		r.Post("/post/create", i.HandlePostCreate)
		r.Put("/post/update", i.HandlePostUpdate)
		r.Delete("/post/delete/{id}", i.HandlePostDelete)
		r.Get("/post/get/{id}", i.HandleGetPost)
		r.Get("/post/feed", i.HandleFeed)
		r.HandleFunc("/post/feed/posted", i.HandlePostedWs)

		r.Get("/dialog/{user_id}/list", i.GetMessages)
		r.Post("/dialog/{user_id}/send", i.SendMessage)
	})

	r.Post("/login", i.HandleLogin)

	return r
}

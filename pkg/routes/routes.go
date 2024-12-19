package routes

import (
	"github.com/Fyefhqdishka/deadlock_v.2/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h handlers.Handlers) {
	UserRoutes(r, h)
	DialogsRoutes(r, h)
	PostRoutes(r, h)
}

func UserRoutes(r *mux.Router, h handlers.Handlers) {
	r.HandleFunc("/api/register", h.Create).Methods("POST")
	r.HandleFunc("/api/login", h.Login)
	r.HandleFunc("/api/users", h.GetUsers).Methods("GET")
}

func DialogsRoutes(r *mux.Router, h handlers.Handlers) {
	r.HandleFunc("/api/dialogs", h.CreateDialog).Methods("POST")
	r.HandleFunc("/api/dialogs", h.GetDialogs).Methods("GET")
}

func PostRoutes(r *mux.Router, h handlers.Handlers) {
	r.HandleFunc("/api/posts", h.CreatePost).Methods("POST")
	r.HandleFunc("/api/posts", h.GetAllPosts).Methods("GET")
	r.HandleFunc("/api/posts/{id}", h.GetPost).Methods("GET")
}

package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/models"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handlers struct {
	Service service.ServiceIface
}

func NewHandlers(service service.ServiceIface) Handlers {
	return Handlers{
		Service: service,
	}
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		h.response(w, SendError(fmt.Sprintf("Can't decode json body: %v", err)), http.StatusBadRequest)
		return
	}

	err := h.Service.Create(r.Context(), &user)
	if err != nil {
		h.response(w, SendError("Can't create user"), http.StatusInternalServerError)
		return
	}

	successMessage := fmt.Sprintf("Welcome %s!", user.Username)
	h.response(w, SendSuccess(successMessage), http.StatusCreated)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неправильный формат запроса", http.StatusBadRequest)
		return
	}

	sessionID, err := h.Service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "неправильный логин или пароль", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authorization", sessionID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving users: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding users to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CreateDialog(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("Authorization")
	if sessionID == "" {
		fmt.Println("Session ID is missing in the headers")
		h.response(w, SendError("Missing session ID in request headers"), http.StatusBadRequest)
		return
	}

	userID, err := h.Service.GetUserByID(r.Context(), sessionID)
	if err != nil {
		h.response(w, SendError(fmt.Sprintf("Error retrieving user for session: %v", err)), http.StatusUnauthorized)
		return
	}

	var dialog struct {
		UserIDOne string `json:"user_id_1"`
		UserIDTwo string `json:"user_id_2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&dialog); err != nil {
		fmt.Println("Error decoding JSON:", err)
		h.response(w, SendError("Can't decode json body"), http.StatusBadRequest)
		return
	}

	dialogID, err := h.Service.CreateDialog(r.Context(), userID, dialog.UserIDTwo)
	if err != nil {
		h.response(w, SendError("Can't create dialog"), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(dialogID), http.StatusCreated)
}

func (h *Handlers) GetDialogs(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("Authorization")
	if sessionID == "" {
		fmt.Println("Session ID is missing in the headers")
		h.response(w, SendError("Missing session ID in request headers"), http.StatusBadRequest)
		return
	}

	userID, err := h.Service.GetUserByID(r.Context(), sessionID)
	if err != nil {
		h.response(w, SendError(fmt.Sprintf("Error retrieving user for session: %v", err)), http.StatusUnauthorized)
		return
	}

	dialogs, err := h.Service.GetUserDialogs(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch dialogs", http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(dialogs), http.StatusOK)
}

func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	sessionID := r.Header.Get("Authorization")
	if sessionID == "" {
		fmt.Println("Session ID is missing in the headers")
		h.response(w, SendError("Missing session ID in request headers"), http.StatusBadRequest)
		return
	}

	userID, err := h.Service.GetUserByID(r.Context(), sessionID)
	if err != nil {
		h.response(w, SendError(fmt.Sprintf("Error retrieving user for session: %v", err)), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		h.response(w, SendError("Invalid input data"), http.StatusBadRequest)
		return
	}

	err = h.Service.CreatePost(r.Context(), &post, userID)
	if err != nil {
		http.Error(w, "Failed to fetch dialogs", http.StatusInternalServerError)
		h.response(w, SendError(fmt.Sprintf("Failed to fetch dialogs, err: %v", err)), http.StatusUnauthorized)
		return
	}

	h.response(w, SendSuccess("dialog created"), http.StatusOK)
}

func (h *Handlers) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Service.GetAllPosts(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving users: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding users to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDStr := vars["id"]
	if postIDStr == "" {
		h.response(w, SendError("Missing post ID in request"), http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	post.ID = postID

	err = h.Service.GetPost(r.Context(), &post)
	if err != nil {
		h.response(w, SendError(fmt.Sprintf("Error retrieving post: %v", err)), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(post), http.StatusOK)
}

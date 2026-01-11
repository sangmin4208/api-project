package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sangmin4208/api-project/internal/models"
)

type UserStorer interface {
	GetUsers() []models.User
	CreateUser(user models.User)
	GetUser(id string) (models.User, error)
	DeleteUser(id string) error
}

type UserHandler struct {
	userStore UserStorer
}

func NewUserHandler(s UserStorer) *UserHandler {
	return &UserHandler{
		userStore: s,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.userStore.GetUsers())
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser models.User
	json.NewDecoder(r.Body).Decode(&newUser)
	h.userStore.CreateUser(newUser)
	json.NewEncoder(w).Encode(newUser)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	user, err := h.userStore.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	err := h.userStore.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// 삭제 성공 시 204 상태 코드 반환
	w.WriteHeader(http.StatusNoContent)
}

package store

import (
	"errors"
	"sync"

	"github.com/sangmin4208/api-project/internal/models"
)

type UserStore struct {
	mu    sync.Mutex
	users []models.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: []models.User{},
	}
}

func (s *UserStore) GetUsers() []models.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.users
}

func (s *UserStore) CreateUser(user models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = append(s.users, user)
}

func (s *UserStore) GetUser(id string) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (s *UserStore) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, user := range s.users {
		if user.ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

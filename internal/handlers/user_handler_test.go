package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sangmin4208/api-project/internal/models"
	"github.com/stretchr/testify/assert"
)

type MockUserStore struct {
	users []models.User
}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		users: []models.User{},
	}
}

func (s *MockUserStore) GetUsers() []models.User {
	return s.users
}

func (s *MockUserStore) CreateUser(user models.User) {
	s.users = append(s.users, user)
}

func (s *MockUserStore) GetUser(id string) (models.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (s *MockUserStore) DeleteUser(id string) error {
	for i, user := range s.users {
		if user.ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func TestGetUsers(t *testing.T) {
	mockStore := &MockUserStore{
		users: []models.User{
			{ID: "1", Email: "test@example.com", Name: "Test User"},
		},
	}
	handler := NewUserHandler(mockStore)
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	response := httptest.NewRecorder()
	handler.GetUsers(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `[{"id":"1","email":"test@example.com","name":"Test User"}]`, response.Body.String())
}

func TestCreateUser(t *testing.T) {
	mockStore := NewMockUserStore()
	handler := NewUserHandler(mockStore)
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"id":"1","email":"test@example.com","name":"Test User"}`))
	response := httptest.NewRecorder()
	handler.CreateUser(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `{"id":"1","email":"test@example.com","name":"Test User"}`, response.Body.String())
}

func TestGetUser(t *testing.T) {
	mockStore := &MockUserStore{
		users: []models.User{
			{ID: "1", Name: "To Be Deleted", Email: "bye@test.com"},
		},
	}
	handler := NewUserHandler(mockStore)
	request, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	request.SetPathValue("id", "1")
	response := httptest.NewRecorder()
	handler.GetUser(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `{"id":"1","email":"bye@test.com","name":"To Be Deleted"}`, response.Body.String())
}

func TestDeleteUser(t *testing.T) {
	mockStore := &MockUserStore{
		users: []models.User{
			{ID: "1", Name: "To Be Deleted", Email: "bye@test.com"},
		},
	}
	handler := NewUserHandler(mockStore)
	request, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
	request.SetPathValue("id", "1")
	response := httptest.NewRecorder()
	handler.DeleteUser(response, request)
	assert.Equal(t, http.StatusNoContent, response.Code)
	_, err := mockStore.GetUser("1")
	assert.Error(t, err)
}

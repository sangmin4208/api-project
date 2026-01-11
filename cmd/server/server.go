package main

import (
	"net/http"

	"github.com/sangmin4208/api-project/internal/handlers"
	"github.com/sangmin4208/api-project/internal/middleware"
	"github.com/sangmin4208/api-project/internal/store"
)

// Server 구조체: 애플리케이션의 모든 의존성(Store, Router)을 가지고 있습니다.
type Server struct {
	listenAddr string
	store      handlers.UserStorer // 인터페이스에 의존!
	handler    http.Handler
}

// NewServer: 서버를 초기화하고 라우팅을 설정합니다.
func NewServer(addr string, dbName string) *Server {
	// 1. 저장소 생성 (나중에 DB로 바꿀 때 여기만 수정하면 됨)
	// userStore := store.NewUserStore()
	userStore := store.NewSQLStore(dbName)

	// 2. 핸들러 생성
	userHandler := handlers.NewUserHandler(userStore)

	// 3. 라우터 설정 (Go 1.22 스타일!)
	mux := http.NewServeMux()

	// [핵심] 이제 메서드별로 경로를 명확하게 지정합니다.
	// "GET /users" -> GetUsers 실행
	mux.HandleFunc("GET /users", userHandler.GetUsers)
	mux.HandleFunc("POST /users", userHandler.CreateUser)

	// {id} 와일드카드 사용 (Go 1.22+)
	// /users/1, /users/abc 등을 잡아냅니다.
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

	wrappedHandler := middleware.Logger(mux)

	return &Server{
		listenAddr: addr,
		store:      userStore,
		handler:    wrappedHandler,
	}
}

// Start: 서버를 시작합니다.
func (s *Server) Start() error {
	// 여기서 로그를 찍거나 추가 설정을 할 수 있습니다.
	return http.ListenAndServe(s.listenAddr, s.handler)
}

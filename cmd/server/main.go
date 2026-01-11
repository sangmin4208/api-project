package main

import (
	"log"
)

func main() {
	// 1. 서버 인스턴스 생성
	server := NewServer(":8080")

	// 2. 서버 실행
	log.Println("Server starting on :8080...")
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

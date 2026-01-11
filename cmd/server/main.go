package main

import (
	"log"
	"os" // 환경 변수를 읽기 위해 필요

	"github.com/joho/godotenv" // 패키지 임포트
)

func main() {
	// 1. .env 파일 로드 (없어도 죽지 않게 처리하는 게 좋지만, 일단은 로그만 찍겠습니다)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// 2. 환경 변수 읽기 (값이 없으면 기본값 설정)
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "app.db"
	}

	// 3. 서버 생성 (환경 변수 주입!)
	server := NewServer(port, dbName)

	log.Printf("Server starting on %s with DB %s...", port, dbName)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

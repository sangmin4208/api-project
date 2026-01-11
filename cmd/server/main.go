package main

import (
	"context"
	"log"
	"net/http"
	"os" // 환경 변수를 읽기 위해 필요
	"os/signal"
	"syscall"
	"time"

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

	// 1. 서버를 별도의 고루틴(Goroutine)에서 실행합니다.
	// 왜냐하면 ListenAndServe는 블로킹(무한 대기) 함수라서,
	// 메인 스레드에서 실행하면 그 밑에 있는 종료 코드가 절대 실행되지 않기 때문입니다.
	go func() {
		log.Printf("Server starting on %s with DB %s...", port, dbName)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			// http.ErrServerClosed는 우리가 의도적으로 종료했을 때 나는 에러이므로 무시합니다.
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 2. 종료 신호(Signal)를 기다리는 채널을 만듭니다.
	// SIGINT(Ctrl+C)나 SIGTERM(배포 환경 종료 신호)을 감지합니다.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 여기서 코드는 멈춰서 신호가 올 때까지 대기합니다. (Block)
	<-quit
	log.Println("\nShutting down server...")

	// 3. 우아한 종료 시작 (Graceful Shutdown)
	// 5초의 시간을 줍니다. "하던 작업 5초 안에 끝내고 문 닫아!"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

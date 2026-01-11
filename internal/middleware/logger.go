package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger: 들어오는 모든 요청의 메서드, 경로, 처리 시간을 기록합니다.
func Logger(next http.Handler) http.Handler {
	// http.HandlerFunc는 함수를 http.Handler 인터페이스로 변환해 줍니다.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // 1. 시작 시간 기록

		// 2. 실제 핸들러(라우터) 실행!
		// 여기서 실제 비즈니스 로직(GetUsers 등)이 돌아갑니다.
		next.ServeHTTP(w, r)

		// 3. 핸들러가 끝나고 돌아오면 걸린 시간 계산
		duration := time.Since(start)

		// 4. 로그 출력
		log.Printf("%s %s | %v", r.Method, r.URL.Path, duration)
	})
}

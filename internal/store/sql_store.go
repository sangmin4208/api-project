package store

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3" // 드라이버 등록 (직접 안 써도 import 해야 함)
	"github.com/sangmin4208/api-project/internal/models"
)

type SQLStore struct {
	db *sql.DB
}

// NewSQLStore: DB 파일을 열고, 테이블이 없으면 만듭니다.
func NewSQLStore(dbName string) *SQLStore {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	// 통신 테스트
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// 테이블 생성 (없으면 생성)
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT,
        name TEXT
    );`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}

	return &SQLStore{db: db}
}

func (s *SQLStore) CreateUser(user models.User) {
	// SQL Injection 방지를 위해 ? (플레이스홀더)를 사용합니다.
	query := `INSERT INTO users (id, email, name) VALUES (?, ?, ?)`

	_, err := s.db.Exec(query, user.ID, user.Email, user.Name)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
}

func (s *SQLStore) GetUsers() []models.User {
	query := `SELECT id, email, name FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return nil
	}
	defer rows.Close()

	var users []models.User

	// 한 줄씩 읽어오기
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, u)
	}
	return users
}

// 인터페이스 만족을 위해 나머지 메서드도 일단 껍데기만 만들어둡니다.
// internal/store/sql_store.go

// 상단 import에 "errors"가 있는지 확인해주세요! 없으면 추가해야 합니다.

func (s *SQLStore) GetUser(id string) (models.User, error) {
	var u models.User
	query := `SELECT id, email, name FROM users WHERE id = ?`

	// QueryRow: 결과가 딱 1개일 때 사용
	err := s.db.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.Name)

	if err == sql.ErrNoRows {
		// DB에 데이터가 없으면 sql.ErrNoRows라는 특별한 에러가 나옵니다.
		// 우리 인터페이스 규칙에 맞춰서 일반 에러로 바꿔줍니다.
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, err
	}

	return u, nil
}
func (s *SQLStore) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	// 몇 개의 행(row)이 삭제되었는지 확인
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// 지워진 게 0개라면? -> 그런 유저는 없었다는 뜻!
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

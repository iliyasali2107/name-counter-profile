package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	InsertURL(int, string) (int, error)
	GetURL(userID int) (string, error)
}

type storage struct {
	DB *pgx.Conn
}

func Init(url string) Storage {
	ctx := context.Background()

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		user_id SMALLINT,
		url VARCHAR(255),
		active BOOLEAN
	);`

	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	return &storage{conn}
}

func (s *storage) InsertURL(userID int, url string) (int, error) {
	query := `INSERT INTO urls(user_id, url) VALUES($1, $2) RETURNING id`
	var id int
	err := s.DB.QueryRow(context.Background(), query, userID, url).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert url %w", err)
	}

	return id, nil
}

func (s *storage) GetURL(userID int) (string, error) {
	query := `SELECT url FROM urls WHERE user_id = $1`
	var url string
	err := s.DB.QueryRow(context.Background(), query, userID).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}

	return url, nil
}

func (s *storage) SetActive(url string) error {
	panic("TODO")
}

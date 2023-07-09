package db

import (
	"context"
	"fmt"
	"log"

	"name-counter-url/pkg/models"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	InsertURL(int64, string) (int64, error)
	GetActiveURL(int64) (models.URL, error)
	SetActive(int64) (int64, error)
	SetNotActive(int64) (int64, error)
	GetURL(int64) (models.URL, error)
	GetUserURLs(userID int64) ([]models.URL, error)
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

func (s *storage) InsertURL(userID int64, url string) (int64, error) {
	query := `INSERT INTO urls(user_id, url) VALUES($1, $2) RETURNING id`
	var id int64
	err := s.DB.QueryRow(context.Background(), query, userID, url).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert url %w", err)
	}

	return id, nil
}

func (s *storage) GetActiveURL(userID int64) (models.URL, error) {
	query := `SELECT id, url FROM urls WHERE user_id = $1 AND active = true`
	var url models.URL
	err := s.DB.QueryRow(context.Background(), query, userID).Scan(&url.ID, &url.UserID, &url.URL)
	if err != nil {
		return models.URL{}, fmt.Errorf("failed to get url: %w", err)
	}

	return url, nil
}

func (s *storage) SetActive(urlID int64) (int64, error) {
	query := `UPDATE urls SET active = true WHERE id = $1;`
	var id int64
	err := s.DB.QueryRow(context.Background(), query, urlID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) GetURL(userID int64) (models.URL, error) {
	query := `SELECT * FROM urls WHERE user_id = $1`
	var url models.URL
	err := s.DB.QueryRow(context.Background(), query, userID).Scan(&url.ID, &url.UserID, &url.URL)
	if err != nil {
		return models.URL{}, fmt.Errorf("failed to get url: %w", err)
	}

	return url, nil
}

func (s *storage) SetNotActive(urlID int64) (int64, error) {
	query := `UPDATE urls SET active = false WHERE id = $1 AND active = false`
	var id int64
	err := s.DB.QueryRow(context.Background(), query, urlID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) GetUserURLs(userID int64) ([]models.URL, error) {
	query := `SELECT * FROM urls where user_id = $1`

	var urls []models.URL

	rows, err := s.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.UserID, &url.URL)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return urls, err
}

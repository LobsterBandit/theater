package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func InitStore(dir string) *Store {
	if err := os.MkdirAll(dir, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s", path.Join(dir, "theater.db")))
	if err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`
CREATE TABLE IF NOT EXISTS plex_webhooks (
	id INTEGER PRIMARY KEY,
	date TEXT NOT NULL,
	type TEXT NOT NULL,
	user TEXT NOT NULL,
	payload BLOB NOT NULL
)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`
CREATE INDEX idx_plex_webhooks_date_type
ON plex_webhooks (id, date, type, user)`); err != nil {
		log.Println(err)
	}

	return &Store{DB: db}
}

func (s *Store) SavePlexWebhook(webhook *WebhookResult) error {
	if _, err := s.DB.Exec(
		"INSERT INTO plex_webhooks (date, type, user, payload) VALUES(?, ?, ?, ?)",
		time.Now().UTC().Format(time.RFC3339),
		webhook.Payload.Event,
		webhook.Payload.Account.Title,
		webhook.RawPayload,
	); err != nil {
		return fmt.Errorf("error saving plex webhook: %w", err)
	}

	return nil
}

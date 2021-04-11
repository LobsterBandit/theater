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
	player TEXT NOT NULL,
	payload BLOB NOT NULL
)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`
CREATE INDEX IF NOT EXISTS idx_plex_webhooks
ON plex_webhooks (id, date, type, user, player)`); err != nil {
		log.Println(err)
	}

	return &Store{DB: db}
}

func (s *Store) Insert(webhook *WebhookResult) error {
	if _, err := s.DB.Exec(
		"INSERT INTO plex_webhooks (date, type, user, player, payload) VALUES(?, ?, ?, ?, ?)",
		time.Now().UTC().Format(time.RFC3339),
		webhook.Payload.Event,
		webhook.Payload.Account.Title,
		webhook.Payload.Player.Title,
		webhook.RawPayload,
	); err != nil {
		return fmt.Errorf("error saving plex webhook: %w", err)
	}

	return nil
}

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/hekmon/plexwebhooks"
	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

type PaginationOptions struct {
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	OrderBy string `json:"orderBy"`
	SortBy  string `json:"sortBy"`
}

type PlexWebhook struct {
	ID      int                  `json:"id"`
	Date    string               `json:"date"`
	Payload plexwebhooks.Payload `json:"payload"`
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

func (s *Store) GetAll() (list []*PlexWebhook, err error) {
	rows, err := s.DB.Query("SELECT id, date, payload FROM plex_webhooks ORDER BY id")
	if err != nil {
		err = fmt.Errorf("error querying plex webhooks: %w", err)

		return
	}

	return mapRowsToPlexWebhooks(rows)
}

func (s *Store) GetAllPaginated(options *PaginationOptions) ([]*PlexWebhook, error) {
	// sqlite doesn't like consecutive space-separated parameters, i.e. ORDER BY ? ?, and throws syntax error on the second ?
	rows, err := s.DB.Query(
		"SELECT id, date, payload FROM plex_webhooks ORDER BY ? LIMIT ? OFFSET ?",
		fmt.Sprintf("%s %s", options.OrderBy, options.SortBy), options.Limit, options.Offset)
	if err != nil {
		err = fmt.Errorf("error querying plex webhooks: %w", err)

		return nil, err
	}

	return mapRowsToPlexWebhooks(rows)
}

func mapRowsToPlexWebhooks(rows *sql.Rows) (list []*PlexWebhook, err error) {
	defer rows.Close()

	for rows.Next() {
		var rawPayload []byte

		plexWebhook := &PlexWebhook{}

		if err = rows.Scan(&plexWebhook.ID, &plexWebhook.Date, &rawPayload); err != nil {
			err = fmt.Errorf("error scanning query results: %w", err)

			return
		}

		if err = json.Unmarshal(rawPayload, &plexWebhook.Payload); err != nil {
			err = fmt.Errorf("error converting raw payload: %w", err)

			return
		}

		list = append(list, plexWebhook)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("error iterating query rows: %w", err)

		return
	}

	return
}

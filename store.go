package main

import (
	"fmt"
	"log"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

type Store struct {
	DB *badger.DB
}

func InitStore(dir string) *Store {
	db, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		log.Fatal(err)
	}

	return &Store{DB: db}
}

func (s *Store) SavePlexWebhook(webhook *WebhookResult) error {
	return s.DB.Update(func(txn *badger.Txn) error { //nolint:wrapcheck
		key := fmt.Sprintf("%s:%s:%s",
			webhook.Payload.Event,
			time.Now().UTC().Format(time.RFC3339),
			webhook.Payload.Account.Title)

		err := txn.Set([]byte(key), webhook.RawPayload)
		if err != nil {
			err = fmt.Errorf("error writing plex webhook to db: %w", err)
		}

		return err
	})
}

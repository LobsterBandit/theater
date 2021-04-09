package main

import (
	"fmt"
	"log"

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
		err := txn.Set([]byte(webhook.Payload.Event), webhook.RawPayload)
		if err != nil {
			err = fmt.Errorf("error writing plex webhook to db: %w", err)
		}

		return err
	})
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	WebhookActionHandler *ActionHandler
	Port                 string
	Router               *chi.Mux
	Store                *Store
}

func (s *Server) Start() {
	s.configureRouter()

	s.configureWebhookActions()

	addr := fmt.Sprintf(":%s", s.Port)
	log.Println("Starting server at", addr)
	log.Fatal(http.ListenAndServe(addr, s.Router))
}

func (s *Server) configureRouter() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", s.ping)
	r.Post("/plex", s.acceptPlexWebhook)
	r.Get("/plex", s.listPlexWebhooks)

	s.Router = r
}

func (s *Server) configureWebhookActions() {
	s.WebhookActionHandler = &ActionHandler{}
	s.WebhookActionHandler.add(&HueAction{})
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("pong"))
}

func (s *Server) acceptPlexWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	multiPartReader, err := r.MultipartReader()
	if err != nil {
		if errors.Is(err, http.ErrNotMultipart) || errors.Is(err, http.ErrMissingBoundary) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		log.Println("unable to create a multipart reader from request:", err)

		return
	}

	result, err := ParseWebhook(multiPartReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Println("unable to parse webhook request:", err)

		return
	}

	log.Printf("received plex webhook: %s\n", result.Payload.Event)

	go s.WebhookActionHandler.processAll(result.Payload)

	go func() {
		if err := s.Store.Insert(result); err != nil {
			log.Println("unable to save webhook:", err)
		}
	}()
}

func (s *Server) listPlexWebhooks(w http.ResponseWriter, r *http.Request) {
	list, err := s.Store.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	listJSON, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(listJSON)
}

func env(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Port string
}

func (s *Server) Start() {
	s.setupRoutes()

	addr := fmt.Sprintf(":%s", s.Port)
	fmt.Println("Starting server at", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s *Server) setupRoutes() {
	http.HandleFunc("/ping", s.handlePing())
	http.HandleFunc("/plex", s.handlePlexWebhook())
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

			return
		}

		fmt.Fprintf(w, "pong")
	}
}

func (s *Server) handlePlexWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

			return
		}

		defer r.Body.Close()

		multiPartReader, err := r.MultipartReader()
		if err != nil {
			if errors.Is(err, http.ErrNotMultipart) || errors.Is(err, http.ErrMissingBoundary) {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			_, wErr := w.Write([]byte(err.Error()))
			if wErr != nil {
				err = fmt.Errorf("request error: %v | write error: %v", err, wErr)
			}

			fmt.Println("can't create a multipart reader from request:", err)

			return
		}

		result := ParsePlexWebhook(multiPartReader)
		if result.err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			_, wErr := w.Write([]byte(result.err.Error()))
			if wErr != nil {
				result.err = fmt.Errorf("request error: %w | write error: %v", result.err, wErr)
			}

			fmt.Println("can't create a multipart reader from request:", result.err)

			return
		}

		if result.Thumbnail != nil {
			fmt.Printf("Name: %s | Size: %d\n", result.Thumbnail.Filename, len(result.Thumbnail.Data))
		}

		fmt.Println(string(result.RawPayload))
		fmt.Println()
	}
}

func env(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}

func main() {
	server := Server{
		Port: env("PORT", "9501"),
	}
	server.Start()
}

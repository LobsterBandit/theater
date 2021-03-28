package main

import (
	"fmt"
	"log"
	"net/http"
)

func createRoute(route string, method string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	return route, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != route {
			http.Error(w, "Not Found", http.StatusNotFound)

			return
		}

		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

			return
		}

		handler(w, r)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func plexWebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Plex webhook receiver")
}

func main() {
	http.HandleFunc(createRoute("/ping", http.MethodGet, pingHandler))

	http.HandleFunc(createRoute("/plex", http.MethodPost, plexWebhookHandler))

	fmt.Printf("Starting server at port 5005\n")

	log.Fatal(http.ListenAndServe(":5005", nil))
}

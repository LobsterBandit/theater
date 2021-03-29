package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hekmon/plexwebhooks"
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

	payload, thumb, err := plexwebhooks.Extract(multiPartReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, wErr := w.Write([]byte(err.Error()))
		if wErr != nil {
			err = fmt.Errorf("request error: %v | write error: %v", err, wErr)
		}

		fmt.Println("can't create a multipart reader from request:", err)

		return
	}

	fmt.Println()
	fmt.Println(time.Now())
	fmt.Printf("%+v\n", *payload)

	if thumb != nil {
		fmt.Printf("Name: %s | Size: %d\n", thumb.Filename, len(thumb.Data))
	}

	fmt.Println()
}

func main() {
	http.HandleFunc(createRoute("/ping", http.MethodGet, pingHandler))

	http.HandleFunc(createRoute("/plex", http.MethodPost, plexWebhookHandler))

	fmt.Printf("Starting server at port 9501\n")

	log.Fatal(http.ListenAndServe(":9501", nil))
}

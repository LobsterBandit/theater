package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/amimof/huego"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	plex "github.com/hekmon/plexwebhooks"
	"github.com/lobsterbandit/theater/internal/actions"
)

type Server struct {
	ActionHandler *actions.Handler
	Port          string
	Router        *chi.Mux
	Store         *Store
}

func (s *Server) Start() {
	s.configureRouter()

	s.configureActions()

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

	// serve SPA
	webDir := http.Dir("./web")

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(webDir))
		fs.ServeHTTP(w, r)
	})

	s.Router = r
}

func (s *Server) configureActions() {
	s.ActionHandler = &actions.Handler{}
	s.ActionHandler.Add(actions.DefaultLogger())

	// add hue action only if ip and user are provided
	bridgeIP, bridgeUser := env("BRIDGE_IP", ""), env("BRIDGE_USER", "")
	if bridgeIP != "" && bridgeUser != "" {
		bridge := huego.New(bridgeIP, bridgeUser)
		plexPlayer, plexUser := env("PLEX_PLAYER", ""), env("PLEX_USER", "")
		hueActions := []actions.Action{
			&actions.Hue{
				Bridge:     bridge,
				Group:      0,
				Scene:      "oSEa7yRNILjLel0",
				PlexEvents: map[plex.EventType]struct{}{plex.EventTypePlay: {}, plex.EventTypeResume: {}},
				PlexPlayer: plexPlayer,
				PlexUser:   plexUser,
			},
			&actions.Hue{
				Bridge:     bridge,
				Group:      1,
				Scene:      "XuQAIFCMsK9raFX",
				PlexEvents: map[plex.EventType]struct{}{plex.EventTypePause: {}, plex.EventTypeStop: {}},
				PlexPlayer: plexPlayer,
				PlexUser:   plexUser,
			},
		}

		s.ActionHandler.Add(hueActions...)
	}
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

	go s.ActionHandler.ProcessAll(result.Payload)

	go func() {
		if err := s.Store.Insert(result); err != nil {
			log.Println("unable to save webhook:", err)
		}
	}()
}

func (s *Server) listPlexWebhooks(w http.ResponseWriter, r *http.Request) {
	var list []*PlexWebhook

	var err error

	if queryParams := r.URL.Query(); len(queryParams) > 0 {
		log.Println(queryParams)

		limit, _ := strconv.Atoi(queryParams.Get("limit"))
		offset, _ := strconv.Atoi(queryParams.Get("offset"))
		options := &PaginationOptions{
			Limit:   limit,
			Offset:  offset,
			OrderBy: queryParams.Get("orderBy"),
			SortBy:  queryParams.Get("sortBy"),
		}

		list, err = s.Store.GetAllPaginated(options)
	} else {
		list, err = s.Store.GetAll()
	}

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	listJSON, err := json.Marshal(list)
	if err != nil {
		log.Println(err.Error())
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

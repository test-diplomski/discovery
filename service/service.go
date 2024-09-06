package service

import (
	"context"
	"fmt"
	"github.com/c12s/discovery/heartbeat"
	"github.com/c12s/discovery/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Service struct {
	w       heartbeat.Heartbeat
	r       *mux.Router
	db      storage.DB
	address string
}

func (ns *Service) sub(ctx context.Context) {
	if ns.w != nil {
		ns.w.Watch(ctx, func(msg string) {
			_, err := ns.db.Store(ctx, msg)
			if err != nil {
				fmt.Println(err)
			}
		})
	}
}

func createBaseRouter(version string) *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	prefix := strings.Join([]string{"/api", version}, "/")
	return r.PathPrefix(prefix).Subrouter()
}

func (s *Service) setupEndpoints() {
	d := s.r.PathPrefix("/discovery").Subrouter()
	d.HandleFunc("/discover", s.discovery()).Methods("GET")
	d.HandleFunc("/heartbeat", s.heartbeat()).Methods("POST")
}

func (s *Service) heartbeat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := read(body)
		if err != nil {
			sendErrorMessage(w, "Could not decode the request body as JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err = s.db.Store(ctx, form(data))
		cancel()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sendJSONResponse(w, map[string]string{"message": "registered"})
	}
}

func (s *Service) discovery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.URL.Query()["service"]; !ok {
			sendErrorMessage(w, "missing service name", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		service := r.URL.Query()["service"][0]
		address, err := s.db.Get(ctx, service)
		if err != nil {
			sendErrorMessage(w, "invalid request try again latter", http.StatusBadRequest)
			return
		}
		defer cancel()

		data, err := resp(service, address)
		if err != nil {
			sendErrorMessage(w, "invalid request try again latter", http.StatusBadRequest)
			return
		}

		sendJSONResponse(w, data)
	}
}

func Run(v, address string, db storage.DB, w heartbeat.Heartbeat) {
	server := &Service{
		db:      db,
		r:       createBaseRouter(v),
		address: address,
		w:       w,
	}
	server.setupEndpoints()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server.sub(ctx)
	db.Watcher(ctx)

	fmt.Println("Discovery Service Started")
	http.ListenAndServe(server.address, handlers.LoggingHandler(os.Stdout, server.r))
}

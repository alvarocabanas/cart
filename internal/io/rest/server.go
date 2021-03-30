package rest

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type (
	Address    string
	PrefixPath string
)

type MetricsHandler http.Handler

func NewServer(address Address, handler CartHandler, metricsHandler MetricsHandler) *http.Server {

	return &http.Server{
		Addr: string(address),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,

		Handler: NewRouter(handler, metricsHandler),
	}
}

func NewRouter(h CartHandler, metricsHandler http.Handler) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/metrics", metricsHandler).Methods(http.MethodGet)
	r.HandleFunc("/", h.AddItem).Methods(http.MethodPost)
	r.HandleFunc("/", h.GetCartStatus).Methods(http.MethodGet)

	return r
}

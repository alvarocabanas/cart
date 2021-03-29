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

func NewServer(address Address, handler CartHandler) *http.Server {

	return &http.Server{
		Addr: string(address),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,

		Handler: NewRouter(handler),
	}
}

func NewRouter(h CartHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.AddItem).Methods(http.MethodPost)
	r.HandleFunc("/", h.GetCartStatus).Methods(http.MethodGet)

	return r
}

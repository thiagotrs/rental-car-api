package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type muxRouter struct {
	*mux.Router
}

func NewMuxRouter(router *mux.Router) *muxRouter {
	return &muxRouter{router}
}

func appJSON(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
}

func (r *muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(uri, appJSON(f)).Methods("GET")
}

func (r *muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(uri, appJSON(f)).Methods("POST")
}

func (r *muxRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(uri, appJSON(f)).Methods("PUT")
}

func (r *muxRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(uri, appJSON(f)).Methods("DELETE")
}

func (r *muxRouter) SERVE(port int) {
	fmt.Printf("API is running on port %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), r)
	if err != nil {
		log.Fatal(err)
	}
}

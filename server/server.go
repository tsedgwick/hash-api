package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tsedgwick/hash-api/api"
)

//Server enables the ability to run an http server
type Server struct {
	client api.Client
	server *http.Server
	mux    *http.ServeMux
}

//New returns a new instance of Server
func New(addr string, client api.Client) *Server {
	s := &Server{
		client: client,
		mux:    http.NewServeMux(),
	}

	s.mux.HandleFunc("/v1/hash", s.v1HashHandler)
	s.mux.HandleFunc("/v2/hash", s.v2HashHandler)
	s.mux.HandleFunc("/v3/hash/", s.v3HashHandler)
	s.mux.HandleFunc("/shutdown", s.shutdownHandler)

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
	return s
}

//Start starts up the http server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

//Shutdown stops the server gracefully
func (s *Server) Shutdown() error {
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nShutdown with timeout: %s\n", timeout)
	//handle shutdown
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown had an error : %v", err)
	}
	return err
}

//ServerHTTP handles the requests
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) v1HashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseForm()
	pw := r.Form.Get("password")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s.client.Tokenize([]byte(pw)))
}

func (s *Server) v2HashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseForm()
	pw := r.Form.Get("password")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s.client.Save([]byte(pw)))
}

func (s *Server) v3HashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, s.client.Token(strings.TrimPrefix(r.URL.Path, "/v3/hash/")))
}

func (s *Server) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go s.Shutdown()
	w.Write([]byte("Shutdown initialized"))
}

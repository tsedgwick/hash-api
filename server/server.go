package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tsedgwick/hash-api/api"
)

//Server enables the ability to run an http server
type Server struct {
	client  api.Client
	server  *http.Server
	mux     *http.ServeMux
	metrics map[string]metricsTag
}

type metricsTag struct {
	Count     int     `json:"count"`
	Average   float64 `json:"average"`
	Method    string  `json:"method"`
	Resource  string  `json:"resource"`
	Code      int     `json:"code"`
	totalTime float64
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{w, http.StatusOK}
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

//New returns a new instance of Server
func New(addr string, client api.Client) *Server {
	s := &Server{
		client:  client,
		mux:     http.NewServeMux(),
		metrics: map[string]metricsTag{},
	}

	handlers := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{pattern: "/v1/hash", handler: s.v1HashHandler},
		{pattern: "/v2/hash", handler: s.v2HashHandler},
		{pattern: "/v3/hash/", handler: s.v3HashHandler},
		{pattern: "/shutdown", handler: s.shutdownHandler},
		{pattern: "/stats", handler: s.stats},
	}

	for _, h := range handlers {
		s.mux.HandleFunc(h.pattern, s.metricsMiddleware(h.handler))
	}

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

func (s *Server) metricsMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		statusW := newResponseWriter(w)

		handler(statusW, r)
		key := r.URL.Path + r.Method + strconv.Itoa(statusW.statusCode)
		tag := s.metrics[key]

		tag.Count++
		tag.totalTime = tag.totalTime + float64(time.Since(startTime))/float64(time.Millisecond)
		tag.Average = tag.totalTime / float64(tag.Count)
		tag.Resource = r.URL.Path
		tag.Method = r.Method
		tag.Code = statusW.statusCode
		s.metrics[key] = tag
	}
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

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metrics := make([]metricsTag, 0, len(s.metrics))
	for _, metric := range s.metrics {
		metrics = append(metrics, metric)
	}
	response, err := json.Marshal(metrics)
	if err != nil {
		//if we cannot marshal our own struct, throw internal server error
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error occurred : %v", err)
		return
	}

	w.Write(response)
}

func (s *Server) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go s.Shutdown()
	w.Write([]byte("Shutdown initialized"))
}

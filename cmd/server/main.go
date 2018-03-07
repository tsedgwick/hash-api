package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/tsedgwick/hash-api/api"
	"github.com/tsedgwick/hash-api/server"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", ":8080", "specify port to use.  defaults to 8080.")
}

func main() {
	flag.Parse()

	s := server.New(port, api.New())
	err := s.Start()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

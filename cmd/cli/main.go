package main

import (
	"flag"
	"log"

	"github.com/tsedgwick/hash-api/api"
)

var password string

func init() {
	flag.StringVar(&password, "password", "", "password to be encoded")
}

func main() {
	flag.Parse()
	c := api.New()
	i := c.Encode([]byte(password))
	log.Printf("result : %v \n", i)
}

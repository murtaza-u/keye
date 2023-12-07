package main

import (
	"log"

	"github.com/murtaza-u/keye/srv"
)

func main() {
	srv, err := srv.New(srv.WithReflection())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Run())
}

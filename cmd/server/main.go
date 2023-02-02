package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/csdev/htecho/internal/handlers"
)

func main() {
	opts := handlers.DefaultOptions()
	flag.BoolVar(&opts.IncludeAuth, "include-auth", opts.IncludeAuth,
		"don't strip authorization headers from the response (this may leak credentials)")

	addr := flag.String("addr", "127.0.0.1:80", "the listen address for the server")
	flag.Parse()

	log.Printf("htecho: listening on %s", *addr)
	err := http.ListenAndServe(*addr, handlers.NewMux(opts))
	if err != nil {
		log.Fatalf("htecho: exited with error: %v", err)
	}
}

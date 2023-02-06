package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/csdev/htecho/internal/handlers"
)

func main() {
	opts := handlers.DefaultOptions()
	flag.BoolVar(&opts.AccessLog, "access-log", opts.AccessLog, "log requests to stdout")
	flag.BoolVar(&opts.IncludeAuth, "include-auth", opts.IncludeAuth,
		"don't strip authorization headers from the response (this may leak credentials)")
	flag.BoolVar(&opts.IncludeIps, "include-ips", opts.IncludeIps,
		"don't strip IP address headers from the response (this may leak network info)")

	includeAll := flag.Bool("A", false, "synonym for --include-auth --include-ips")

	addr := flag.String("addr", "127.0.0.1:80", "the listen address for the server")
	readTimeout := flag.Duration("read-timeout", time.Minute,
		"max time allowed for reading the request")
	writeTimeout := flag.Duration("write-timeout", time.Minute,
		"max time allowed for sending the response")

	flag.Parse()

	if *includeAll {
		opts.IncludeAll()
	}

	server := &http.Server{
		Handler:      handlers.NewMux(opts),
		Addr:         *addr,
		ReadTimeout:  *readTimeout,
		WriteTimeout: *writeTimeout,
	}

	log.Printf("htecho.server: listening on %s", *addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("htecho.server: exited with error: %v", err)
	}
}

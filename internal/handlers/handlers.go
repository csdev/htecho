package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func IsHeaderEqual(a, b string) bool {
	return http.CanonicalHeaderKey(a) == http.CanonicalHeaderKey(b)
}

type EchoHandler struct {
	HandlerOptions *Options
}

func NewEchoHandler(o *Options) http.Handler {
	return &EchoHandler{
		HandlerOptions: o,
	}
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Echo-Method", r.Method)
	w.Header().Add("X-Echo-Query", r.URL.RawQuery)

	for header, vals := range r.Header {
		if IsHeaderEqual(header, "Authorization") || IsHeaderEqual(header, "Proxy-Authorization") {
			if !h.HandlerOptions.IncludeAuth {
				continue // exclude sensitive headers
			}
		}
		for _, val := range vals {
			w.Header().Add(fmt.Sprintf("X-Echo-Header-%s", header), val)
		}
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	_, err := io.Copy(w, r.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

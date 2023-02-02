package handlers

import "net/http"

func NewMux(o *Options) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", NewEchoHandler(o))
	return mux
}

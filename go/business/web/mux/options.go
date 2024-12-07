package mux

import (
	"net/http"
)

// Options represents optional parameters.
type Options struct {
	corsOrigin string
	staticFS   http.Handler
}

// WithCORS provides configuration options for CORS.
func WithCORS(origin string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origin
	}
}

func WithStaticFS(fs http.Handler) func(opts *Options) {

	return func(opts *Options) {
		opts.staticFS = fs
	}
}

// Package web contains a small web framework extension.
package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"syscall"
	"time"

	"github.com/google/uuid"
)

const (
	HTMLMime = "text/html"
	HXMLMime = "application/vnd.hyperview+xml"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	mux      *http.ServeMux
	shutdown chan os.Signal
	mw       []Middleware
	group    string
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	mux := http.NewServeMux()

	return &App{
		mux:      mux,
		shutdown: shutdown,
		mw:       mw,
	}
}

func (a *App) Group() *App {
	return &App{
		mux: a.mux,
		mw:  slices.Clone(a.mw),
	}
}

func (a *App) Mount(subRoute string) *App {
	return &App{
		mux:   a.mux,
		mw:    slices.Clone(a.mw),
		group: subRoute,
	}
}

func (a *App) Use(mw ...Middleware) {
	a.mw = append(a.mw, mw...)
}

func (a *App) Get(path string, fn Handler, mw ...Middleware) {
	a.handle(http.MethodGet, a.group, path, fn, mw...)
}

func (a *App) Post(path string, fn Handler, mw ...Middleware) {
	a.handle(http.MethodPost, a.group, path, fn, mw...)
}

func (a *App) Put(path string, fn Handler, mw ...Middleware) {
	a.handle(http.MethodPut, a.group, path, fn, mw...)
}

func (a *App) Delete(path string, fn Handler, mw ...Middleware) {
	a.handle(http.MethodDelete, a.group, path, fn, mw...)
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) ServeFS(handler http.Handler) {
	a.mux.Handle("/", handler)
}

// EnableCORS enables CORS preflight requests to work in the middleware. It
// prevents the MethodNotAllowedHandler from being called. This must be enabled
// for the CORS middleware to work.
func (a *App) EnableCORS(mw Middleware) {
	a.mw = append(a.mw, mw)

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return RespondJSON(ctx, w, "OK", http.StatusOK)
	}
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)

		handler(ctx, w, r)
	}

	a.mux.HandleFunc("OPTIONS /", h)
}

// HandleNoMiddleware sets a handler function for a given HTTP method and path pair
// to the application server mux. Does not include the application middleware or
// OTEL tracing.
func (a *App) HandleNoMiddleware(method string, group string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			if validateError(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(path, h)
}

// handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) handle(method string, group string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			if validateError(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}

	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, h)
}

// validateError validates the error for special conditions that do not
// warrant an actual shutdown by the system.
func validateError(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}

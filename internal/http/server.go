package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
)

// Server is a basic HTTP http.
type Server struct {
	svr             *http.Server
	shutdownTimeout time.Duration
}

// NewServer returns a new Server.
func NewServer(bind string, shutdownTimeout time.Duration, handler http.Handler) *Server {
	return &Server{
		svr:             &http.Server{Addr: bind, Handler: handler},
		shutdownTimeout: shutdownTimeout,
	}
}

// Serve binds the http to it's bind address and starts serving requests.
func (s *Server) Serve(ctx context.Context) (err error) {
	ch := make(chan error)
	defer close(ch)

	// Start serving.
	go s.serve(ch)

	// Wait on our context to be cancelled.
	<-ctx.Done()

	// Shutdown http within 5 seconds
	ctxShutdown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err = suppressServerClosed(s.svr.Shutdown(ctxShutdown)); err == nil {
		return suppressServerClosed(<-ch)
	}

	if chErr := suppressServerClosed(<-ch); chErr != nil {
		return multierror.Append(err, chErr)
	}

	return err
}

func suppressServerClosed(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) serve(ch chan error) {
	ch <- s.svr.ListenAndServe()
}

package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/app"
)

// Server ...
type Server struct {
	srv *http.Server
}

// NewServer ...
func NewServer(addr string, a app.Application) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: NewHealthcheckHandler(a),
		},
	}
}

// Start starts http server.
func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop stops http server.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

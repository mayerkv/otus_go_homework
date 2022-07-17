package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

type Server struct {
	logger Logger
	app    Application
	host   string
	port   string
	server *http.Server
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	return &Server{
		logger: logger,
		app:    app,
		host:   host,
		port:   port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/health", HealthCheckHandler{})

	s.server = &http.Server{
		Addr:         net.JoinHostPort(s.host, s.port),
		Handler:      loggingMiddleware(mux, s.logger),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.Info("server listen at " + s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

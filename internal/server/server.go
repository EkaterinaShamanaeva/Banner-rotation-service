package server

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/app"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	router *mux.Router
	logg   Logger
	app    *app.App
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

func NewServer(logger Logger, app *app.App) *Server {
	serv := &Server{logg: logger, app: app}
	serv.router = mux.NewRouter()

	serv.router.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, "test")
	}).Methods("GET")

	return serv
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logg.Info("HTTP server " + addr + " starting...")
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	errChan := make(chan error)

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errChan:
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("HTTP server was stopped...")
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

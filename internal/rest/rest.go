package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"url-shortner/internal/datastore"
	"url-shortner/internal/logger"
)

const (
	readTimeout   = 10
	writeTimeout  = 5
	serverTimeout = 5
)

type RESTServer struct {
	httpServer *http.Server
	logger     logger.Logger
	datastore  datastore.UrlStorer
}

func NewRESTServer(
	logger logger.Logger,
	datastore datastore.UrlStorer,
	host string,
	port int,
	handler http.Handler,
) *RESTServer {
	server := &RESTServer{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			ReadTimeout:  readTimeout * time.Second,
			WriteTimeout: writeTimeout * time.Second,
			Handler:      handler,
		},
		logger:    logger,
		datastore: datastore,
	}
	return server
}

func (s *RESTServer) Start() error {
	s.logger.With("address", s.httpServer.Addr).Info("Starting HTTP server")
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *RESTServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), serverTimeout*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	s.logger.Info("Stopping HTTP server")
	return nil
}

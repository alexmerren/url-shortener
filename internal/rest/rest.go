package rest

import (
	"context"
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
	ctx        context.Context
	httpServer *http.Server
	logger     logger.Logger
	datastore  datastore.Datastorer
}

type RESTServerInput struct {
	Logger    logger.Logger
	Datastore datastore.Datastorer
	Address   int
}

func NewRESTServer(input *RESTServerInput) (*RESTServer, error) {
	server := &RESTServer{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", input.Address),
			ReadTimeout:  readTimeout * time.Second,
			WriteTimeout: writeTimeout * time.Second,
		},
		logger:    input.Logger,
		datastore: input.Datastore,
	}
	server.addRoutes()
	return server, nil
}

func (s *RESTServer) Start() error {
	s.logger.With("address", s.httpServer.Addr).Info("Starting HTTP server")
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.logger.Fatal(err)
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

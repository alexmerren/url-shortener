package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *RESTServer) addRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/shortner/add/", s.addURL).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/shortner/get/{id}/", s.getURL).Methods(http.MethodGet)
	router.HandleFunc("/healthz/", s.Health).Methods(http.MethodGet)
	router.Use(s.loggingMiddleware)
	router.Use(s.formatMiddleware)
	s.httpServer.Handler = router
}

package rest

import (
	"net/http"
	"url-shortner/internal/datastore"
	"url-shortner/internal/logger"
	"url-shortner/internal/rest/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(logger logger.Logger, datastore datastore.UrlStorer) http.Handler {

	addUrlHandler := handlers.NewAddUrlHandler(logger, datastore)
	getUrlHandler := handlers.NewGetUrlHandler(logger, datastore)

	router := mux.NewRouter()
	router.Handle("/api/v1/shortener/", addUrlHandler).Methods(http.MethodPost)
	router.Handle("/api/v1/shortener/{id}/", getUrlHandler).Methods(http.MethodGet)
	router.HandleFunc("/health/", handlers.HandleHealth).Methods(http.MethodGet)
	router.Use(formatMiddleware)
	return router
}

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"url-shortner/internal/datastore"
	urlerrors "url-shortner/internal/errors"
	"url-shortner/internal/logger"

	"github.com/gorilla/mux"
)

type GetUrlHandler struct {
	logger    logger.Logger
	datastore datastore.UrlStorer
}

func NewGetUrlHandler(
	logger logger.Logger,
	datastore datastore.UrlStorer,
) *GetUrlHandler {
	return &GetUrlHandler{
		logger:    logger,
		datastore: datastore,
	}
}

func (h *GetUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		HandleError(w, r, urlerrors.NewBadRequestError(errors.New("there was no id provided")))
		return
	}

	url, err := h.datastore.GetUrl(id)
	if err != nil {
		HandleError(w, r, urlerrors.NewInternalServerError(err))
		return
	}

	output := &getURLResponseBody{
		Url: url,
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		HandleError(w, r, urlerrors.NewInternalServerError(err))
		return
	}
}

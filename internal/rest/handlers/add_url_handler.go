package handlers

import (
	"encoding/json"
	"net/http"

	"url-shortner/internal/datastore"
	urlerrors "url-shortner/internal/errors"
	"url-shortner/internal/logger"
)

type AddUrlHandler struct {
	logger    logger.Logger
	datastore datastore.UrlStorer
}

func NewAddUrlHandler(
	logger logger.Logger,
	datastore datastore.UrlStorer,
) *AddUrlHandler {
	return &AddUrlHandler{
		logger:    logger,
		datastore: datastore,
	}
}

func (h *AddUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody := &addURLRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
		HandleError(w, r, urlerrors.NewBadRequestError(err))
		return
	}

	id, err := h.datastore.InsertUrl(requestBody.Url)
	if err != nil {
		HandleError(w, r, urlerrors.NewInternalServerError(err))
		return
	}

	output := &addURLResponseBody{
		ID: id,
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		HandleError(w, r, urlerrors.NewInternalServerError(err))
		return
	}
}

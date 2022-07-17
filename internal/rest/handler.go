package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	urlerror "url-shortner/internal/errors"

	"github.com/gorilla/mux"
)

type addURLRequestBody struct {
	Url string `json:"url"`
}

type addURLResponseBody struct {
	ID string `json:"id"`
}

func (s *RESTServer) addURL(w http.ResponseWriter, r *http.Request) {
	requestBody := &addURLRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
		handleError(w, r, urlerror.NewBadRequestError(err))
		return
	}

	id, err := s.datastore.InsertURL(requestBody.Url)
	if err != nil {
		handleError(w, r, urlerror.NewInternalServerError(err))
		return
	}

	output := &addURLResponseBody{
		ID: id,
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		handleError(w, r, urlerror.NewInternalServerError(err))
		return
	}
}

type getURLResponseBody struct {
	Url string `json:"url"`
}

func (s *RESTServer) getURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		handleError(w, r, urlerror.NewBadRequestError(errors.New("there was no id provided")))
		return
	}

	url, err := s.datastore.GetURL(id)
	if err != nil {
		handleError(w, r, urlerror.NewInternalServerError(err))
		return
	}

	output := &getURLResponseBody{
		Url: url,
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		handleError(w, r, urlerror.NewInternalServerError(err))
		return
	}
}

func (s *RESTServer) Health(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("healthy!")); err != nil {
		handleError(w, r, err)
		return
	}
}

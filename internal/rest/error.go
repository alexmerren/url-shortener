package rest

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	StatusCode  int    `json:"status_code"`
	ErrorString string `json:"error_string"`
}

// nolint:errcheck // No point in trying to check for encoding errors here.
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&ErrorResponse{
		StatusCode:  http.StatusInternalServerError,
		ErrorString: err.Error(),
	})
}

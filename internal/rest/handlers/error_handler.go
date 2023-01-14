package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&errorResponse{
		StatusCode:  http.StatusInternalServerError,
		ErrorString: err.Error(),
	})
}

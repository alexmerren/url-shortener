package handlers

import "net/http"

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("healthy!")); err != nil {
		HandleError(w, r, err)
		return
	}
}

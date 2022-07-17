package rest

import "net/http"

func (s *RESTServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.
			With("path", r.URL.Path).
			Debug(r.Method)
		next.ServeHTTP(w, r)
	})
}

func (s *RESTServer) formatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

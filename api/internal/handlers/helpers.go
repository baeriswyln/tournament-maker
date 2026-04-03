package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"baeriswyl.dev/tournament-maker/internal/apierror"
)

func writeError(w http.ResponseWriter, r *http.Request, e *apierror.APIError, cause ...error) {
	switch e.ErrorCode {
	case "INTERNAL_ERROR":
		if len(cause) > 0 && cause[0] != nil {
			log.Printf("[ERROR] %s %s - %s: %v", r.Method, r.URL.Path, e.ErrorCode, cause[0])
		} else {
			log.Printf("[ERROR] %s %s - %s", r.Method, r.URL.Path, e.ErrorCode)
		}
	case "AUTH_TOKEN_INVALID", "AUTH_INVALID_CREDENTIALS", "AUTH_REFRESH_INVALID", "AUTH_FORBIDDEN":
		log.Printf("[SECURITY] %s %s - %s from %s", r.Method, r.URL.Path, e.ErrorCode, r.RemoteAddr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	_, _ = w.Write(e.JSON())
}

func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
// internal/handler/response_util.go

package handler

import (
	"encoding/json"
	"net/http"
)

// sendJSONResponse sends a JSON-encoded HTTP response.
func sendJSONResponse(w http.ResponseWriter, response interface{}, errorCode int) {
	statusCode := http.StatusOK
	if errorCode > 0 && errorCode < 100 {
		statusCode = http.StatusInternalServerError
	}

	if errorCode > 99 && errorCode < 200 {
		statusCode = http.StatusUnauthorized
	}

	if errorCode > 199 && errorCode < 300 {
		statusCode = http.StatusBadRequest
	}

	if errorCode > 299 && errorCode < 400 {
		statusCode = http.StatusServiceUnavailable
	}

	if errorCode > 399 && errorCode < 500 {
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

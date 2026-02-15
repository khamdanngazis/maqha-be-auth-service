// internal/handler/handler.go

package handler

import (
	"fmt"
	"net/http"
)

type pingHandler struct {
}

type PingHandler interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

func NewPingHandler() PingHandler {
	return &pingHandler{}
}

// PingHandler handles the "/ping" endpoint.
func (c *pingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Pong!")
}

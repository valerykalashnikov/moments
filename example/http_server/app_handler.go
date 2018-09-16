package main

import (
	"fmt"
	"net/http"

	"github.com/valerykalashnikov/moments"
)

// AppHandler tracks requests count for the previous n seconds
type AppHandler struct {
	counter *moments.MomentsCounter
}

// NewAppHandler initializes handler to track requests count for the previous n seconds
func NewAppHandler(counter *moments.MomentsCounter) *AppHandler {
	return &AppHandler{counter}
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.counter.Track()
	fmt.Fprintf(w, fmt.Sprint(h.counter.Count()))
}

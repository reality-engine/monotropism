package opendream

import (
	"net/http"
)

type handler struct {
	provider Provider
}

func NewHandler(provider Provider) *handler {
	return &handler{provider: provider}
}

func (h *handler) ServeCSV(w http.ResponseWriter, r *http.Request) {
	err := h.provider.ServeCSV(w, r)
	if err != nil {
		http.Error(w, "Unable to serve the CSV file", http.StatusInternalServerError)
	}
}

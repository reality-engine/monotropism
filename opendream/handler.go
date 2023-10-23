package opendream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type EEGHandler struct {
	store DataStore
}

func NewEEGHandler(store DataStore) *EEGHandler {
	return &EEGHandler{store: store}
}

func (h *EEGHandler) ServeEEGTextData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rowLimit := r.URL.Query().Get("rows")
	if rowLimit == "" {
		rowLimit = "1000"
	}

	_, err := strconv.Atoi(rowLimit)
	if err != nil {
		http.Error(w, "Invalid row limit provided", http.StatusBadRequest)
		return
	}

	records, err := h.store.QueryEEGData(ctx, "skillful-flow-399108", rowLimit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to query data: %v", err), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(records)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

package delegation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore"
)

// APIHandler handles the API requests.
type APIHandler struct {
	datastore datastore.Datastorer
}

// New creates a new APIHandler.
func New(datastore datastore.Datastorer) *APIHandler {
	return &APIHandler{
		datastore: datastore,
	}
}

// GetDelegationsHandler handles /xtz/delegations endpoint.
func (a *APIHandler) GetDelegationsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the year parameter from the query string
	yearParam := r.URL.Query().Get("year")

	// If yearParam is not empty, filter delegations by year
	year := 0

	var err error

	if yearParam != "" {
		year, err = strconv.Atoi(yearParam)
		if err != nil {
			zap.L().Error("error parsing year parameter", zap.Error(err))
			http.Error(w, "Bad Request", http.StatusBadRequest)

			return
		}
	}

	delegations, err := a.datastore.GetDelegations(r.Context(), year)
	if err != nil {
		zap.L().Error("couldn't get delegations from datastore", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	// Convert the delegations to JSON
	responseJSON, err := json.Marshal(delegations)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON response
	_, err = w.Write(responseJSON)
	if err != nil {
		zap.L().Error("error writing JSON response", zap.Error(err))
	}
}

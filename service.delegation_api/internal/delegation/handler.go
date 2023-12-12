package delegation

import (
	"encoding/json"
	"fmt"
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
//
//nolint:funlen
func (a *APIHandler) GetDelegationsHandler(w http.ResponseWriter, r *http.Request) {
	yearParam := r.URL.Query().Get("year")

	year, err := parseParamInt("year", yearParam)
	if err != nil {
		zap.L().Error("error parsing year parameter", zap.Error(err))
		http.Error(w, fmt.Sprintf("Bad Request: %s", err), http.StatusBadRequest)

		return
	}

	pageParam := r.URL.Query().Get("page")

	pageNumber, err := parseParamInt("page", pageParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request: %s", err), http.StatusBadRequest)

		return
	}

	size := r.URL.Query().Get("size")

	pageSize, err := parseParamInt("size", size)
	if err != nil {
		zap.L().Error("error parsing size parameter", zap.Error(err))
		http.Error(w, fmt.Sprintf("Bad Request: %s", err), http.StatusBadRequest)

		return
	}

	// Ensure default values for pageNumber and pageSize
	if pageNumber <= 0 {
		pageNumber = 1
	}

	if pageSize <= 0 {
		pageSize = 100
	}

	delegations, err := a.datastore.GetDelegations(
		r.Context(),
		pageNumber,
		pageSize,
		year,
	)
	if err != nil {
		zap.L().Error("couldn't get delegations from datastore", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	responseJSON, err := json.Marshal(delegations)
	if err != nil {
		zap.L().Error("error marshalling delegations to JSON", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	// Calculate the maximum number of pages based on the total number of documents and page size
	totalDocuments, err := a.datastore.GetDelegationsCount(r.Context(), year)
	if err != nil {
		zap.L().Error("couldn't get delegations count from datastore", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	maxPages := totalDocuments / pageSize
	if totalDocuments%pageSize != 0 {
		maxPages++
	}

	w.Header().Set("X-Total-Pages", strconv.Itoa(maxPages))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)
	if err != nil {
		zap.L().Error("error writing JSON response", zap.Error(err))
	}
}

func parseParamInt(paramName, paramValue string) (int, error) {
	if paramValue == "" {
		return 0, nil
	}

	parsedValue, err := strconv.Atoi(paramValue)
	if err != nil {
		zap.L().Error(
			"couldn't parse query parameter value",
			zap.String("paramName", paramName),
			zap.String("paramValue", paramValue),
			zap.Error(err),
		)

		return 0, fmt.Errorf(
			"couldn't parse value %s for query parameter %s",
			paramName,
			paramValue,
		)
	}

	return parsedValue, nil
}

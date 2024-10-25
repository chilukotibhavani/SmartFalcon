package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Response struct for API responses
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateAssetRequest struct for creating new assets
type CreateAssetRequest struct {
	DealerID string  `json:"dealerId"`
	MSISDN   string  `json:"msisdn"`
	MPIN     string  `json:"mpin"`
	Balance  float64 `json:"balance"`
}

// UpdateBalanceRequest struct for balance updates
type UpdateBalanceRequest struct {
	MSISDN    string  `json:"msisdn"`
	MPIN      string  `json:"mpin"`
	Amount    float64 `json:"amount"`
	TransType string  `json:"transType"`
	Remarks   string  `json:"remarks"`
}

func main() {
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/api/assets", createAsset).Methods("POST")
	router.HandleFunc("/api/assets/{msisdn}", getAsset).Methods("GET")
	router.HandleFunc("/api/assets/{msisdn}/balance", updateBalance).Methods("PUT")
	router.HandleFunc("/api/assets/{msisdn}/history", getAssetHistory).Methods("GET")

	// Start server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createAsset(w http.ResponseWriter, r *http.Request) {
	var req CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Here you would interact with your chaincode
	// For now, we'll just return a success response
	sendResponse(w, http.StatusCreated, "success", "Asset created successfully", nil)
}

func getAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	msisdn := vars["msisdn"]

	// Here you would call your chaincode to get the asset
	// For now, returning dummy data
	asset := map[string]interface{}{
		"msisdn":   msisdn,
		"balance":  1000.0,
		"status":   "ACTIVE",
		"dealerId": "DEALER001",
	}

	sendResponse(w, http.StatusOK, "success", "Asset retrieved successfully", asset)
}

func updateBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	msisdn := vars["msisdn"]

	var req UpdateBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Here you would call your chaincode to update the balance
	// For now, just return success
	sendResponse(w, http.StatusOK, "success", fmt.Sprintf("Balance updated for MSISDN: %s", msisdn), nil)
}

func getAssetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	msisdn := vars["msisdn"]

	// Here you would call your chaincode to get transaction history
	// For now, returning dummy data
	history := []map[string]interface{}{
		{
			"transType": "CREDIT",
			"amount":    100.0,
			"timestamp": "2024-10-25T10:00:00Z",
		},
		{
			"transType": "DEBIT",
			"amount":    50.0,
			"timestamp": "2024-10-25T11:00:00Z",
		},
	}

	sendResponse(w, http.StatusOK, "success", fmt.Sprintf("History retrieved for MSISDN: %s", msisdn), history)
}

func sendResponse(w http.ResponseWriter, statusCode int, status, message string, data interface{}) {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

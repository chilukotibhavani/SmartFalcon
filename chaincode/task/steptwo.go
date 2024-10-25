package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing assets
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up an asset
type Asset struct {
	DocType     string  `json:"docType"` // Used for indexing
	DealerID    string  `json:"dealerId"`
	MSISDN      string  `json:"msisdn"`
	MPIN        string  `json:"mpin"`
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
	TransAmount float64 `json:"transAmount"`
	TransType   string  `json:"transType"`
	Remarks     string  `json:"remarks"`
	LastUpdated string  `json:"lastUpdated"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// CreateAsset issues a new asset to the world state with given details
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerId string, msisdn string, mpin string, balance float64) error {
	exists, err := s.AssetExists(ctx, msisdn)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", msisdn)
	}

	asset := Asset{
		DocType:     "asset",
		DealerID:    dealerId,
		MSISDN:      msisdn,
		MPIN:        mpin,
		Balance:     balance,
		Status:      "ACTIVE",
		TransAmount: 0,
		TransType:   "CREATE",
		Remarks:     "Account Created",
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(msisdn, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given msisdn
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, msisdn string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(msisdn)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", msisdn)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateBalance updates an asset's balance in the world state
func (s *SmartContract) UpdateBalance(ctx contractapi.TransactionContextInterface, msisdn string, mpin string, amount float64, transType string, remarks string) error {
	asset, err := s.ReadAsset(ctx, msisdn)
	if err != nil {
		return err
	}

	if asset.MPIN != mpin {
		return fmt.Errorf("incorrect MPIN")
	}

	if transType == "DEBIT" {
		if asset.Balance < amount {
			return fmt.Errorf("insufficient balance")
		}
		asset.Balance -= amount
	} else if transType == "CREDIT" {
		asset.Balance += amount
	} else {
		return fmt.Errorf("invalid transaction type")
	}

	asset.TransAmount = amount
	asset.TransType = transType
	asset.Remarks = remarks
	asset.LastUpdated = time.Now().Format(time.RFC3339)

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(msisdn, assetJSON)
}

// AssetExists returns true when asset with given MSISDN exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, msisdn string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(msisdn)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAssetHistory returns the chain of custody for an asset since issuance
func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, msisdn string) ([]Asset, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(msisdn)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []Asset
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if err := json.Unmarshal(response.Value, &asset); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating asset-management chaincode: %v", err)
		return
	}

	if err := assetChaincode.Start(); err != nil {
		fmt.Printf("Error starting asset-management chaincode: %v", err)
	}
}

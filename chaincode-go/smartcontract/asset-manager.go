package assetManagement

import (
	"encoding/json"
	"fmt"
	"crypto/sha256"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	DealerID 	string	`json:"DEALERID"`
	MSISDN 		string 	`json:"MSISDN"`
	MPIN 		string	`json:"MPIN"`
	BALANCE		float64	`json:"BALANCE"`
	STATUS		string	`json:"STATUS"`
	TRANSAMOUNT	float64	`json:"TRANSAMOUNT"`
	TRANSTYPE	string	`json:"TRANSTYPE"`
	REMARKS		string	`json:"REMARKS"`
}

type HistoryQueryResult struct{
	Record 		*Asset	`json:"record"`
	TxID 		string	`json:"txId"`
	Timestamp	string	`json:"timestamp"`
	IsDelete	bool	`json:"isDelete"`
}

func (s *SmartContract) InitLedger (ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{{
		DealerID:    "DLR001",
		MSISDN:      "919876543210",
		MPIN:        hashPassword("4321"),
		BALANCE:     10000.00,
		STATUS:      "ACTIVE",
		TRANSAMOUNT: 0,
		TRANSTYPE:   "NA",
		REMARKS:     "Initial account opening",
	},
	{
		DealerID:    "DLR002",
		MSISDN:      "918123456789",
		MPIN:        hashPassword("1234"),
		BALANCE:     5000.25,
		STATUS:      "ACTIVE",
		TRANSAMOUNT: 1000.00,
		TRANSTYPE:   "DEPOSIT",
		REMARKS:     "First top‑up",
	},
	{
		DealerID:    "DLR003",
		MSISDN:      "919912345678",
		MPIN:        hashPassword("9834"),
		BALANCE:     200.50,
		STATUS:      "SUSPENDED",
		TRANSAMOUNT: 150.00,
		TRANSTYPE:   "WITHDRAWAL",
		REMARKS:     "ATM cash‑out — flagged for review",
	},
	{
		DealerID:    "DLR004",
		MSISDN:      "917012345678",
		MPIN:        hashPassword("2468"),
		BALANCE:     0.00,
		STATUS:      "BLOCKED",
		TRANSAMOUNT: 0.00,
		TRANSTYPE:   "NA",
		REMARKS:     "Account blocked due to suspected fraud",
	},
	{
		DealerID:    "DLR005",
		MSISDN:      "919876123450",
		MPIN:        hashPassword("1357"),
		BALANCE:     7500.75,
		STATUS:      "ACTIVE",
		TRANSAMOUNT: 500.75,
		TRANSTYPE:   "TRANSFER",
		REMARKS:     "Transfer to dealer DLR002",
	},
}

	for _,asset := range assets {
		assetJSON,err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.DealerID,assetJSON)
		if err != nil {
			return err
		}
	}
	return nil	
}


func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin string, bal float64, status string, txnAmount float64, txnType string, remarks string) error {
	if err := s.onlyOrg1(ctx); err != nil { return err }
	exists, err := s.Exists(ctx, dealerID)
	if err != nil {
		return fmt.Errorf("Error : %s\n", err)
	}
	if exists  {
		return fmt.Errorf("The DealerId %s already exists\n", dealerID)
	}
	
	asset := Asset{
		DealerID: dealerID,
		MSISDN: msisdn,
		MPIN: hashPassword(mpin),
		BALANCE: bal,
		STATUS: status,
		TRANSAMOUNT: txnAmount,
		TRANSTYPE: txnType,
		REMARKS: remarks,
	}

	assetJson , err := json.Marshal(asset)
	if err != nil {
		return  fmt.Errorf("Error encoding the data : %s \n", err)
	}

	err = ctx.GetStub().PutState(dealerID,assetJson)
	if err != nil {
		return err
	}
	return nil
}


func (s *SmartContract) Exists(ctx contractapi.TransactionContextInterface,DealerID string) (bool,error){
	assetJSON , err := ctx.GetStub().GetState(DealerID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, err
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,dealerID string, msisdn string, mpin string, bal float64, status string, txnAmount float64, txnType string, remarks string) error {
	if err := s.onlyOrg1(ctx); err != nil { return err }
	exists , err := s.Exists(ctx,dealerID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("The Dealer ID Doesnt Exist %s\n", dealerID)
	}

	asset := Asset{
		DealerID: dealerID,
		MSISDN: msisdn,
		MPIN: hashPassword(mpin),
		BALANCE: bal,
		STATUS: status,
		TRANSAMOUNT: txnAmount,
		TRANSTYPE: txnType,
		REMARKS: remarks,
	}

	assetJson , err := json.Marshal(asset)
	if err != nil {
		return  fmt.Errorf("Error encoding the data : %s \n", err)
	}

	return ctx.GetStub().PutState(dealerID,assetJson)
}


func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, dealerID string) (error) {
	if err := s.onlyOrg1(ctx); err != nil { return err }
	exists, err := s.Exists(ctx,dealerID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("The Dealer ID Doesnt Exist %s\n", dealerID)
	}

	return ctx.GetStub().DelState(dealerID)
}


func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface , dealerID string) (*Asset,error) {
	assetJson, err := ctx.GetStub().GetState(dealerID)
	if err != nil {
		return nil,err
	}
	if assetJson == nil {
		return nil,fmt.Errorf("The Dealer ID Doesnt Exist %s\n", dealerID) 
	}

	var asset Asset

	err = json.Unmarshal(assetJson,&asset)
	if err != nil {
		return nil, err
	}
	asset.MPIN = ""
	return &asset,nil
}


func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, dealerID string) ([]HistoryQueryResult,error){
	exists, err := s.Exists(ctx,dealerID)
	if err != nil {
        return nil, fmt.Errorf("failed to check if asset exists: %v", err)
    }
    if !exists {
        return nil, fmt.Errorf("asset with dealerID %s does not exist", dealerID)
    }

	historyIterator, err := ctx.GetStub().GetHistoryForKey(dealerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for dealerID %s: %v", dealerID, err)
	}
	defer historyIterator.Close()

	var results []HistoryQueryResult

	for historyIterator.HasNext(){
		historyData, err := historyIterator.Next()
		if err != nil{ 
			return nil, fmt.Errorf("failed to get next history item: %v", err)
		}

		var asset Asset
		var record *Asset

		if !historyData.IsDelete {
			err := json.Unmarshal(historyData.Value , &asset)
			if err != nil{
				return nil, fmt.Errorf("failed to unmarshal the data: %v", err)
			}
			record = &asset
		}

		timestamp := historyData.Timestamp.AsTime().Format("2006-01-02 15:04:05")

		historyRecord := HistoryQueryResult{
			Record: record,
			TxID: historyData.TxId,
			Timestamp: timestamp,
			IsDelete: historyData.IsDelete,
		}

		results = append(results, historyRecord)
	}

	return results ,nil
}	


func (s * SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error){

	resultsIterator, err := ctx.GetStub().GetStateByRange("","")
	if err != nil {
        return nil, fmt.Errorf("failed to get all assets: %v", err)
    }
    defer resultsIterator.Close()

	var assets []*Asset
	
	for resultsIterator.HasNext(){
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get all assets: %v", err)
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil{
			return nil, fmt.Errorf("failed to unmarshal assets: %v", err)
		}
		asset.MPIN = ""
		assets = append(assets, &asset)

	}
	return assets, nil
}

func (s *SmartContract) VerifyMPIN(ctx contractapi.TransactionContextInterface, dealerID string, inputMPIN string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerID)
	if err != nil {
		return false, fmt.Errorf("failed to read asset: %v", err)
	}
	if assetJSON == nil {
		return false, fmt.Errorf("asset with dealerID %s not found", dealerID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal asset: %v", err)
	}

	if verifyPassword(inputMPIN, asset.MPIN) {
		return true, nil
	}
	return false, nil
}


func hashPassword(password string) string {
    hash := sha256.Sum256([]byte(password))
    return fmt.Sprintf("%x", hash)
}

func verifyPassword(input, stored string) bool {
    return hashPassword(input) == stored
}


func (s *SmartContract) onlyOrg1(ctx contractapi.TransactionContextInterface) error {
    mspID, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get MSP ID: %v", err)
    }

    if mspID != "Org1MSP" {
        return fmt.Errorf("unauthorized: only Org1MSP can invoke this function")
    }
    return nil
}


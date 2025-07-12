package assetManagement

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	DealerID 	string	`json:"DEALERID"`
	MSISDN 		string 	`json:"MSISDN"`
	MPIN 		int		`json:"MPIN"`
	BALANCE		float64	`json:"BALANCE"`
	STATUS		string	`json:"STATUS"`
	TRANSAMOUNT	float64	`json:"TRANSAMOUNT"`
	TRANSTYPE	string	`json:"TRANSTYPE"`
	REMARKS		string	`json:"REMARKS"`
}

func (s *SmartContract) InitLedger (ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{{
		DealerID:    "DLR001",
		MSISDN:      "919876543210",
		MPIN:        4321,
		BALANCE:     10000.00,
		STATUS:      "ACTIVE",
		TRANSAMOUNT: 0,
		TRANSTYPE:   "NA",
		REMARKS:     "Initial account opening",
	},
	{
		DealerID:    "DLR002",
		MSISDN:      "918123456789",
		MPIN:        1234,
		BALANCE:     5000.25,
		STATUS:      "ACTIVE",
		TRANSAMOUNT: 1000.00,
		TRANSTYPE:   "DEPOSIT",
		REMARKS:     "First top‑up",
	},
	{
		DealerID:    "DLR003",
		MSISDN:      "919912345678",
		MPIN:        9876,
		BALANCE:     200.50,
		STATUS:      "SUSPENDED",
		TRANSAMOUNT: 150.00,
		TRANSTYPE:   "WITHDRAWAL",
		REMARKS:     "ATM cash‑out — flagged for review",
	},
	{
		DealerID:    "DLR004",
		MSISDN:      "917012345678",
		MPIN:        2468,
		BALANCE:     0.00,
		STATUS:      "BLOCKED",
		TRANSAMOUNT: 0.00,
		TRANSTYPE:   "NA",
		REMARKS:     "Account blocked due to suspected fraud",
	},
	{
		DealerID:    "DLR005",
		MSISDN:      "919876123450",
		MPIN:        1357,
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


func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin int, bal float64, status string, txnAmount float64, txnType string, remarks string) error {
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
		MPIN: mpin,
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

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,dealerID string, msisdn string, mpin int, bal float64, status string, txnAmount float64, txnType string, remarks string) error {

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
		MPIN: mpin,
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
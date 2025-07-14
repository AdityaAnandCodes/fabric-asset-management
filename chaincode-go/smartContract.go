/*
SPDX-License-Identifier: MIT
*/

package main

import (
	"log"

	assetManagement "github.com/AdityaAnandCodes/fabric-samples/fabric-asset-management/chaincode-go/smartcontract"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func main() {
	assetManagerContract,err := contractapi.NewChaincode(&assetManagement.SmartContract{})
	if err != nil {
		log.Panicf("Error creating the chaincode:  %s\n", err)
	}

	if err := assetManagerContract.Start(); err != nil {
		log.Panicf("Error starting the chaincode:  %s\n", err)
	}
}
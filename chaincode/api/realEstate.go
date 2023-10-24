package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CreateRealEstate creates new real estate (admin)
func CreateRealEstate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Validate parameters
	if len(args) != 4 {
		return shim.Error("Insufficient number of parameters")
	}
	accountId := args[0] // Account ID for verifying admin rights
	proprietor := args[1]
	totalArea := args[2]
	livingSpace := args[3]
	if accountId == "" || proprietor == "" || totalArea == "" || livingSpace == "" {
		return shim.Error("Parameters contain empty values")
	}
	if accountId == proprietor {
		return shim.Error("The operator should be an admin and cannot be the same as the owner")
	}
	// Convert parameter data format
	var formattedTotalArea float64
	if val, err := strconv.ParseFloat(totalArea, 64); err != nil {
		return shim.Error(fmt.Sprintf("Failed to convert 'totalArea' parameter format: %s", err))
	} else {
		formattedTotalArea = val
	}
	var formattedLivingSpace float64
	if val, err := strconv.ParseFloat(livingSpace, 64); err != nil {
		return shim.Error(fmt.Sprintf("Failed to convert 'livingSpace' parameter format: %s", err))
	} else {
		formattedLivingSpace = val
	}
	// Verify if it's an admin operation
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{accountId})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("Operator permission verification failed: %s", err))
	}
	var account model.Account
	if err = json.Unmarshal(resultsAccount[0], &account); err != nil {
		return shim.Error(fmt.Sprintf("Query operator information failed: %s", err))
	}
	if account.UserName != "Admin" {
		return shim.Error("Operator does not have sufficient permissions")
	}
	// Verify the existence of the proprietor
	resultsProprietor, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{proprietor})
	if err != nil || len(resultsProprietor) != 1 {
		return shim.Error(fmt.Sprintf("Owner 'proprietor' information verification failed: %s", err))
	}
	realEstate := &model.RealEstate{
		RealEstateID: stub.GetTxID()[:16],
		Proprietor:   proprietor,
		Encumbrance:  false,
		TotalArea:    formattedTotalArea,
		LivingSpace:  formattedLivingSpace,
	}
	// Write to the ledger
	if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	// Return the information of the successfully created real estate
	realEstateByte, err := json.Marshal(realEstate)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to serialize the information of the created real estate: %s", err))
	}
	// Successful response
	return shim.Success(realEstateByte)
}

// QueryRealEstateList queries real estate (can query all or by owner)
func QueryRealEstateList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var realEstateList []model.RealEstate
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var realEstate model.RealEstate
			err := json.Unmarshal(v, &realEstate)
			if err != nil {
				return shim.Error(fmt.Sprintf("Failed to deserialize QueryRealEstateList: %s", err))
			}
			realEstateList = append(realEstateList, realEstate)
		}
	}
	realEstateListByte, err := json.Marshal(realEstateList)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to serialize QueryRealEstateList: %s", err))
	}
	return shim.Success(realEstateListByte)
}
